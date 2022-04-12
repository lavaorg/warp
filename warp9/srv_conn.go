// Copyright 2009 The Go9p Authors.  All rights reserved.
// Copyright 2019 RMG Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package warp9

import (
	"fmt"
	"net"
	"sync"
)

// The Conn type represents a connection from a client to the object server
type Conn struct {
	sync.Mutex
	Srv        *Srv
	Msize      uint32 // maximum size of Warp9 messages for the connection
	Id         string // used for debugging and stats
	Debuglevel int

	conn    net.Conn
	fidpool map[uint32]*SrvFid
	reqs    map[uint16]*SrvReq // all outstanding requests

	reqout chan *SrvReq
	rchan  chan *Fcall
	done   chan bool

	// stats -- only for logging and debugging
	nreqs   int    // number of requests processed by the server
	tsz     uint64 // total size of the T messages received
	rsz     uint64 // total size of the R messages sent
	npend   int    // number of currently pending messages
	maxpend int    // maximum number of pending messages
	nreads  int    // number of reads
	nwrites int    // number of writes
}

func (conn *Conn) String() string {
	return conn.Srv.Id + "/" + conn.Id
}

func (srv *Srv) newConnSetup(c net.Conn) *Conn {

	conn := new(Conn)
	conn.Srv = srv
	conn.Msize = srv.Msize
	conn.Debuglevel = srv.Debuglevel
	conn.conn = c
	conn.fidpool = make(map[uint32]*SrvFid)
	conn.reqs = make(map[uint16]*SrvReq)
	conn.reqout = make(chan *SrvReq, srv.Maxpend)
	conn.done = make(chan bool)
	conn.rchan = make(chan *Fcall, 64)

	srv.Lock()
	if srv.conns == nil {
		srv.conns = make(map[*Conn]*Conn)
	}
	srv.conns[conn] = conn
	srv.Unlock()

	conn.Id = c.RemoteAddr().String()
	if op, ok := (conn.Srv.ops).(ConnOps); ok {
		op.ConnOpened(conn)
	}

	if sop, ok := (interface{}(conn)).(StatsOps); ok {
		sop.statsRegister()
	}
	return conn
}

// Async serve connection.
func (srv *Srv) NewConn(c net.Conn) {
	conn := srv.newConnSetup(c)
	go conn.recv()
	go conn.send()
}

// Block and serve the connection.
func (srv *Srv) NewConnWait(c net.Conn) {
	conn := srv.newConnSetup(c)
	go conn.recv()
	conn.send()
}

// Initiate a network connection than serve the connection.
// net,addr are same as used in net.Dial()
// wait==true block and serve (see NewConnWait()), else return see NewConn())
// return nil on success; else err is a Dial() error
//
func (srv *Srv) InitiateConn(nettyp, addr string, wait bool) error {
	c, e := net.Dial(nettyp, addr)
	if e != nil {
		return e
	}
	if wait {
		srv.NewConnWait(c)
	} else {
		srv.NewConn(c)
	}
	return nil
}

func (conn *Conn) close() {
	conn.done <- true
	conn.Srv.Lock()
	delete(conn.Srv.conns, conn)
	conn.Srv.Unlock()

	if sop, ok := (interface{}(conn)).(StatsOps); ok {
		sop.statsUnregister()
	}
	if op, ok := (conn.Srv.ops).(ConnOps); ok {
		op.ConnClosed(conn)
	}

	/* call FidDestroy for all remaining fids */
	if op, ok := (conn.Srv.ops).(SrvFidOps); ok {
		for _, fid := range conn.fidpool {
			op.FidDestroy(fid)
		}
	}
}

func (conn *Conn) recv() {
	var err error
	var n int

	buf := make([]byte, conn.Msize*8)
	pos := 0
	for {
		if len(buf) < int(conn.Msize) {
			b := make([]byte, conn.Msize*8)
			copy(b, buf[0:pos])
			buf = b
			b = nil
		}

		n, err = conn.conn.Read(buf[pos:])
		if err != nil || n == 0 {
			conn.close()
			return
		}

		pos += n
		for pos > 4 {
			sz, _ := gint32(buf)
			if sz > conn.Msize {
				Error("bad client connection: ", conn.conn.RemoteAddr())
				conn.conn.Close()
				conn.close()
				return
			}
			if pos < int(sz) {
				if len(buf) < int(sz) {
					b := make([]byte, conn.Msize*8)
					copy(b, buf[0:pos])
					buf = b
					b = nil
				}

				break
			}
			fc, err, fcsize := Unpack(buf)
			if err != nil {
				Error(fmt.Sprintf("invalid packet : %v %v", err, buf))
				conn.conn.Close()
				conn.close()
				return
			}

			tag := fc.Tag
			req := new(SrvReq)
			select {
			case req.Rc = <-conn.rchan:
				break
			default:
				req.Rc = NewFcall(conn.Msize)
			}

			req.Conn = conn
			req.Tc = fc
			//			req.Rc = rc
			if conn.Debuglevel > 0 {
				if conn.Debuglevel&DbgPrintPackets != 0 {
					Info(">-> %v %v", conn.Id, req.Tc.Pkt)
				}

				if conn.Debuglevel&DbgPrintFcalls != 0 {
					Info(">>> %v %v", conn.Id, req.Tc.String())
				}
			}

			conn.Lock()
			conn.nreqs++
			conn.tsz += uint64(fc.FcSize)
			conn.npend++
			if conn.npend > conn.maxpend {
				conn.maxpend = conn.npend
			}

			req.next = conn.reqs[tag]
			conn.reqs[tag] = req
			process := req.next == nil
			if req.next != nil {
				req.next.prev = req
			}
			conn.Unlock()
			if process {
				// Tversion may change some attributes of the
				// connection, so we block on it. Otherwise,
				// we may loop back to reading and that is a race.
				// This fix brought to you by the race detector.
				if req.Tc.Type == Tversion {
					req.process()
				} else {
					go req.process()
				}
			}

			buf = buf[fcsize:]
			pos -= fcsize
		}
	}

}

func (conn *Conn) send() {
	for {
		select {
		case <-conn.done:
			return

		case req := <-conn.reqout:
			req.Rc.SetTag(req.Tc.Tag)
			conn.Lock()
			conn.rsz += uint64(req.Rc.FcSize)
			conn.npend--
			conn.Unlock()
			if conn.Debuglevel > 0 {
				if conn.Debuglevel&DbgPrintPackets != 0 {
					Info("<-< %v %v", conn.Id, req.Rc.Pkt)
				}

				if conn.Debuglevel&DbgPrintFcalls != 0 {
					Info("<<< %v %v", conn.Id, req.Rc.String())
				}
			}

			for buf := req.Rc.Pkt; len(buf) > 0; {
				n, err := conn.conn.Write(buf)
				if err != nil {
					/* just close the socket, will get signal on conn.done */
					Error("error while writing")
					conn.conn.Close()
					break
				}

				buf = buf[n:]
			}

			select {
			case conn.rchan <- req.Rc:
				break
			default:
			}
		}
	}

	//panic("unreached")
}

// Return the remote address of the connection.
func (conn *Conn) RemoteAddr() net.Addr {
	return conn.conn.RemoteAddr()
}

// Return the local address of the connection.
func (conn *Conn) LocalAddr() net.Addr {
	return conn.conn.LocalAddr()
}

// Start listening on the specified symbolic network type/address. This function
// creates a net.Listen(ntype,addr) object and invokes StartListener() with it.
//
func (srv *Srv) StartNetListener(ntype, addr string) error {
	l, err := net.Listen(ntype, addr)
	if err != nil {
		return err
	}

	return srv.StartListener(l)
}

// Start listening on the specified network and address for incoming
// connections. Once a connection is established, create a new Conn
// value, read messages from the socket, send them to the specified
// server, and send back responses received from the server.
func (srv *Srv) StartListener(l net.Listener) error {
	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}

		srv.NewConn(c)
	}
}
