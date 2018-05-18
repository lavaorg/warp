// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package wkit

import (
	"sync/atomic"
	"time"

	"github.com/lavaorg/lrt/mlog"
	"github.com/lavaorg/warp/warp9"
)

// Interface represetning an object server.
type Server interface {
	StartNetListener(ntype, addr string) error
	GetRoot() *DirItem
}

// A Simple object server containing a tree of objects to serve.
type Serv struct {
	warp9.Srv
	warp9.StatsOps
	root *DirItem
}

var qidpGlob uint64 = 0

//
// host application invokes this when its ready to start serving objects.
//
// id is a symbolic name used for the root mount name
// dbgLevel sets level of debug to emit to standard out log
//
func StartServer(id string, dbgLevel int) Server {
	srv := &Serv{}
	srv.Id = id
	srv.Debuglevel = dbgLevel
	srv.setRoot()
	srv.Start(srv)
	return srv
}

// Return the root object for the server.
func (srv *Serv) GetRoot() *DirItem {
	return srv.root
}

func (srv *Serv) setRoot() {

	d := NewDirItem() //root is not a bind location; but is a diritem
	d.Dir.Name = srv.Id
	d.Uid = 1
	d.Gid = 1
	d.Muid = 1

	d.Atime = uint32(time.Now().Unix())
	d.Mtime = d.Atime

	// mark as a directory that allows READ operation for all users
	d.Mode = warp9.DMDIR | uint32(Perms(warp9.DMREAD, warp9.DMREAD, warp9.DMREAD))

	d.Qid = warp9.Qid{warp9.QTDIR, 0, NextQid()}

	mlog.Debug("getRoot:%v", d)
	srv.root = d
	return
}

// Create a Permissions word out of user/group/other components
// The lower 3 bytes set in User/Group/Other order
func Perms(u, g, o byte) uint16 {
	return uint16(uint16(u)<<6 | uint16(g)<<3 | uint16(o))
}

// Helper for creation of Qid-ID. Simple incremented counter.
func NextQid() uint64 {
	return atomic.AddUint64(&qidpGlob, 1)
}
