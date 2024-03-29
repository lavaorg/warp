// Copyright 2009 The Go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE object.

// The srv* files provides definitions and functions used to implement
// a Warp9 object server.
package warp9

import (
	"sync"
)

type reqStatus int

const (
	reqFlush     reqStatus = (1 << iota) /* request is flushed (no response will be sent) */
	reqWork                              /* goroutine is currently working on it */
	reqResponded                         /* response is already produced */
	reqSaved                             /* no response was produced after the request is worked on */
)

// Authentication operations. The object server should implement them if
// it requires user authentication. The authentication in Warp9 is
// done by creating special authentication fids and performing I/O
// operations on them. Once the authentication is done, the authentication
// fid can be used by the user to get access to the actual objects.
type AuthOps interface {
	// AuthInit is called when the user starts the authentication
	// process on SrvFid afid. The user that is being authenticated
	// is referred by afid.User. The function should return the Qid
	// for the authentication object, or an Error if the user can't be
	// authenticated
	AuthInit(afid *SrvFid, aname string) (*Qid, error)

	// AuthDestroy is called when an authentication fid is destroyed.
	AuthDestroy(afid *SrvFid)

	// AuthCheck is called after the authentication process is finished
	// when the user tries to attach to the object server. If the function
	// returns nil, the authentication was successful and the user has
	// permission to access the objects.
	AuthCheck(fid *SrvFid, afid *SrvFid, aname string) error

	// AuthRead is called when the user attempts to read data from an
	// authentication fid.
	AuthRead(afid *SrvFid, offset uint64, data []byte) (count int, err error)

	// AuthWrite is called when the user attempts to write data to an
	// authentication fid.
	AuthWrite(afid *SrvFid, offset uint64, data []byte) (count int, err error)
}

// Connection operations. These should be implemented if the object server
// needs to be called when a connection is opened or closed.
type ConnOps interface {
	ConnOpened(*Conn)
	ConnClosed(*Conn)
}

// SrvFid operations. This interface should be implemented if the object server
// needs to be called when a SrvFid is destroyed.
type SrvFidOps interface {
	FidDestroy(*SrvFid)
}

// Request operations. This interface should be implemented if the object server
// needs to bypass the default request process, or needs to perform certain
// operations before the (any) request is processed, or before (any) response
// sent back to the client.
type SrvReqProcessOps interface {
	// Called when a new request is received from the client. If the
	// interface is not implemented, (req *SrvReq) srv.Process() method is
	// called. If the interface is implemented, it is the user's
	// responsibility to call srv.Process. If srv.Process isn't called,
	// SrvFid, Afid and Newfid fields in SrvReq are not set, and the SrvReqOps
	// methods are not called.
	SrvReqProcess(*SrvReq)

	// Called when a request is responded, i.e. when (req *SrvReq)srv.Respond()
	// is called and before the response is sent. If the interface is not
	// implemented, (req *SrvReq) srv.PostProcess() method is called to finalize
	// the request. If the interface is implemented and SrvReqProcess calls
	// the srv.Process method, SrvReqRespond should call the srv.PostProcess
	// method.
	SrvReqRespond(*SrvReq)
}

// Flush operation. This interface should be implemented if the object server
// can flush pending requests. If the interface is not implemented, requests
// that were passed to the object server implementation won't be flushed.
// The flush method should call the (req *SrvReq) srv.Flush() method if the flush
// was successful so the request can be marked appropriately.
type FlushOp interface {
	Flush(*SrvReq)
}

// The Srv type contains the basic fields used to control the Warp9
// object server. Each server implementation should create a value
// of Srv type, initialize the values it cares about and pass the
// struct to the (Srv *) srv.Start(ops) method together with the object
// that implements the object server operations.
type Srv struct {
	sync.Mutex
	Id         string // Used for debugging and stats
	Msize      uint32 // Maximum size of the Warp9 messages supported by the server
	Debuglevel int    // debug level
	Upool      Users  // Interface for finding users and groups known to the object server
	Maxpend    int    // Maximum pending outgoing requests

	ops   interface{}     // operations
	conns map[*Conn]*Conn // List of connections
}

// The SrvFid type references an object on the object server.
// A new SrvFid is created when the user attaches to the object server (the Attach
// operation), or when Walk-ing to an object. The SrvFid values are created
// automatically by the srv implementation.
// There can be multiple SrvFid's refering to the same Object, both within a
// Conn or from multple Conn's
// A fid, thus a SrvFid, must be unique within a Conn (e.g. Conn's are a fid space)
// The SrvFidDestroy operation is called when a SrvFid is destroyed.
type SrvFid struct {
	sync.Mutex
	fid       uint32
	refcount  int
	opened    bool        // True if the SrvFid is opened
	Fconn     *Conn       // Connection the SrvFid belongs to
	Omode     uint8       // Open mode (O* flags), if the fid is opened
	Type      uint8       // SrvFid type (QT* flags)
	Diroffset uint64      // If directory, the next valid read position
	Dirents   []byte      // If directory, the serialized dirents
	User      User        // The SrvFid's user
	Aux       interface{} // Can be used by the object server implementation for per-SrvFid data
}

// The SrvReq type represents a Warp9 request. Each request has a
// T-message (Tc) and a R-message (Rc). If the SrvReqProcessOps don't
// override the default behavior, the implementation initializes SrvFid,
// Afid and Newfid values and automatically keeps track on when the SrvFids
// should be destroyed.
type SrvReq struct {
	sync.Mutex
	Tc     *Fcall  // Incoming message
	Rc     *Fcall  // Outgoing response
	Fid    *SrvFid // The SrvFid value for all messages that contain fid[4]
	Afid   *SrvFid // The SrvFid value for the messages that contain afid[4] (Tauth and Tattach)
	Newfid *SrvFid // The SrvFid value for the messages that contain newfid[4] (Twalk)
	Conn   *Conn   // Connection that the request belongs to

	status     reqStatus
	flushreq   *SrvReq
	prev, next *SrvReq
}

// The Start method should be called once the object server implementor
// initializes the Srv struct with the preferred values. It sets default
// values to the fields that are not initialized and creates the goroutines
// required for the server's operation. The method receives an empty
// interface value, ops, that should implement the interfaces the object server is
// interested in. Ops must implement the SrvReqOps interface.
func (srv *Srv) Start(ops interface{}) bool {
	if _, ok := (ops).(SrvReqOps); !ok {
		return false
	}

	srv.ops = ops
	if srv.Upool == nil {
		srv.Upool = Identity
	}

	if srv.Msize < IOHDRSZ {
		srv.Msize = MSIZE
	}

	if sop, ok := (interface{}(srv)).(StatsOps); ok {
		Info("start stats")
		sop.statsRegister()
	}

	return true
}

func (srv *Srv) String() string {
	return srv.Id
}

func (req *SrvReq) process() {
	req.Lock()
	flushed := (req.status & reqFlush) != 0
	if !flushed {
		req.status |= reqWork
	}
	req.Unlock()

	if flushed {
		req.Respond()
	}

	if rop, ok := (req.Conn.Srv.ops).(SrvReqProcessOps); ok {
		rop.SrvReqProcess(req)
	} else {
		req.Process()
	}

	req.Lock()
	req.status &= ^reqWork
	if !(req.status&reqResponded != 0) {
		req.status |= reqSaved
	}
	req.Unlock()
}

// Performs the default processing of a request. Initializes
// the SrvFid, Afid and Newfid fields and calls the appropriate
// SrvReqOps operation for the message. The object server implementer
// should call it only if the object server implements the SrvReqProcessOps
// within the SrvReqProcess operation.
func (req *SrvReq) Process() {
	conn := req.Conn
	srv := conn.Srv
	tc := req.Tc

	if tc.Fid != NOFID && tc.Type != Tattach {
		srv.Lock()
		req.Fid = conn.FidGet(tc.Fid)
		srv.Unlock()
		if req.Fid == nil {
			req.RespondError(&WarpError{Eunknownfid, ""})
			return
		}
	}

	switch req.Tc.Type {
	default:
		req.RespondError(&WarpError{Ebadw9msg, ""})

	case Tversion:
		srv.version(req)

	case Tauth:
		srv.auth(req)

	case Tattach:
		srv.attach(req)

	case Tflush:
		srv.flush(req)

	case Twalk:
		srv.walk(req)

	case Topen:
		srv.open(req)

	case Tcreate:
		srv.create(req)

	case Tread:
		srv.read(req)

	case Twrite:
		srv.write(req)

	case Tclunk:
		srv.clunk(req)

	case Tremove:
		srv.remove(req)

	case Tstat:
		srv.stat(req)

	case Twstat:
		srv.wstat(req)
	}
}

// Performs the post processing required if the (*SrvReq) Process() method
// is called for a request. The object server implementer should call it
// only if the object server implements the SrvReqProcessOps within the
// SrvReqRespond operation.
func (req *SrvReq) PostProcess() {
	srv := req.Conn.Srv

	/* call the post-handlers (if needed) */
	switch req.Tc.Type {
	case Tauth:
		srv.authPost(req)

	case Tattach:
		srv.attachPost(req)

	case Twalk:
		srv.walkPost(req)

	case Topen:
		srv.openPost(req)

	case Tcreate:
		srv.createPost(req)

	case Tread:
		srv.readPost(req)

	case Tclunk:
		srv.clunkPost(req)

	case Tremove:
		srv.removePost(req)
	}

	if req.Fid != nil {
		req.Fid.DecRef()
		req.Fid = nil
	}

	if req.Afid != nil {
		req.Afid.DecRef()
		req.Afid = nil
	}

	if req.Newfid != nil {
		req.Newfid.DecRef()
		req.Newfid = nil
	}
}

// The Respond method sends response back to the client. The req.Rc value
// should be initialized and contain valid Warp9 message. In most cases
// the object server implementer shouldn't call this method directly. Instead
// one of the RespondR* methods should be used.
func (req *SrvReq) Respond() {
	var flushreqs *SrvReq

	conn := req.Conn
	req.Lock()
	status := req.status
	req.status |= reqResponded
	req.status &= ^reqWork
	req.Unlock()

	if (status & reqResponded) != 0 {
		return
	}

	/* remove the request and all requests flushing it */
	conn.Lock()
	nextreq := req.prev
	if nextreq != nil {
		nextreq.next = nil
		// if there are flush requests, move them to the next request
		if req.flushreq != nil {
			var p *SrvReq = nil
			r := nextreq.flushreq
			for ; r != nil; p, r = r, r.flushreq {
			}

			if p == nil {
				nextreq.flushreq = req.flushreq
			} else {
				nextreq = req.flushreq
			}
		}

		flushreqs = nil
	} else {
		delete(conn.reqs, req.Tc.Tag)
		flushreqs = req.flushreq
	}
	conn.Unlock()

	if rop, ok := (req.Conn.Srv.ops).(SrvReqProcessOps); ok {
		rop.SrvReqRespond(req)
	} else {
		req.PostProcess()
	}

	if (status & reqFlush) == 0 {
		conn.reqout <- req
	}

	// process the next request with the same tag (if available)
	if nextreq != nil {
		go nextreq.process()
	}

	// respond to the flush messages
	// can't send the responses directly to conn.reqout, because the
	// the flushes may be in a tag group too
	for freq := flushreqs; freq != nil; freq = freq.flushreq {
		freq.Respond()
	}
}

// Should be called to cancel a request. Should only be called
// from the Flush operation if the FlushOp is implemented.
func (req *SrvReq) Flush() {
	req.Lock()
	req.status |= reqFlush
	req.Unlock()
	req.Respond()
}

// Lookup a SrvFid struct based on the 32-bit identifier sent over the wire.
// Returns nil if the fid is not found. Increases the reference count of
// the returned fid. The user is responsible to call DecRef once it no
// longer needs it.
func (conn *Conn) FidGet(fidno uint32) *SrvFid {
	conn.Lock()
	fid, present := conn.fidpool[fidno]
	conn.Unlock()
	if present {
		fid.IncRef()
	}

	return fid
}

// Creates a new SrvFid struct for the fidno integer. Returns nil
// if the SrvFid for that number already exists. The returned fid
// has reference count set to 1.
func (conn *Conn) FidNew(fidno uint32) *SrvFid {
	conn.Lock()
	_, present := conn.fidpool[fidno]
	if present {
		conn.Unlock()
		return nil
	}

	fid := new(SrvFid)
	fid.fid = fidno
	fid.refcount = 1
	fid.Fconn = conn
	conn.fidpool[fidno] = fid
	conn.Unlock()

	return fid
}

// Increase the reference count for the fid.
func (fid *SrvFid) IncRef() {
	fid.Lock()
	fid.refcount++
	fid.Unlock()
}

// Decrease the reference count for the fid. When the
// reference count reaches 0, the fid is no longer valid.
func (fid *SrvFid) DecRef() {
	fid.Lock()
	fid.refcount--
	n := fid.refcount
	fid.Unlock()

	if n > 0 {
		return
	}

	conn := fid.Fconn
	conn.Lock()
	delete(conn.fidpool, fid.fid)
	conn.Unlock()

	if fop, ok := (conn.Srv.ops).(SrvFidOps); ok {
		fop.FidDestroy(fid)
	}
}
