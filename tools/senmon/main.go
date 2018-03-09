// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/lavaorg/lrt/mlog"
	"github.com/lavaorg/warp/warp9"
)

var verbose = flag.Bool("v", false, "verbose mode")
var dbglev = flag.Int("d", 0, "debuglevel")
var addr = flag.String("a", "127.0.0.1:9901", "network address")
var aname = flag.String("aname", ".", "path on server to use as root")

func main() {

	flag.Parse()

	// accept a network connection
	listenForServers("tcp", *addr)

	return
}

func usage() {
	fmt.Println("usage: np [-v][-d dbglev] [-a addr] cmd arg")
	fmt.Println("\tcmd = {ls,stat,cat,echo}")
}

// list on the indicated network and address and then mount the calling server.
// This is a reverse mount; the server iniitates a connection and then this
// client performs the mount on that conneciton.
func listenForServers(ntype, addr string) {

	l, err := net.Listen(ntype, addr)
	if err != nil {
		mlog.Error("listen failed:%v", err)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			mlog.Error("accept fail:v", err)
			return
		}
		mlog.Debug("accepted connection: %v", c)
		go handleConnection(c)
	}
}

// mount the server that just called on the Conn
func handleConnection(c net.Conn) {

	uid := warp9.OsUsers.Uid2User(os.Geteuid())
	warp9.DefaultDebuglevel = *dbglev

	c9, err := warp9.MountConn(c, *aname, 500, uid)
	if err != nil {
		mlog.Error("Error:%v\n", err)
		return // end thread
	}

	// read target sensor
	readSensor(c9)

	// close connection
	c9.Unmount()
}

func readSensor0(c9 *warp9.Clnt) {
	f, err := c9.FOpen("sensors", warp9.OREAD)
	if err != nil {
		log.Fatalf("Error:%v\n", err)
	}
	defer f.Close()

	buf := make([]byte, 8192)
	for {
		n, err := f.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			log.Fatalf("Error reading:%v\n", err)
		}
		mlog.Info("%v", string(buf))
		if err == io.EOF {
			break
		}
	}

	if err != nil && err != io.EOF {
		mlog.Error("error:%v", err)
		return
	}

}

// this shortens the number of requests due to avoiding a
// last read that just looks for EOF. We have knowledge that
// the sensors being read is a small number of bytes.
func readSensor(c9 *warp9.Clnt) {

	fid, err := c9.FWalk("sensors")
	if err != nil {
		mlog.Error("could not Walk:%v", err)
		return
	}
	err = c9.Open(fid, warp9.OREAD)
	if err != nil {
		c9.Clunk(fid)
		mlog.Error("open failed:%v", err)
		return
	}

	buf, err := c9.Read(fid, uint64(0), uint32(100))
	if err != nil {
		log.Fatalf("Error:%v\n", err)
	}
	mlog.Info("%v", string(buf))
}
