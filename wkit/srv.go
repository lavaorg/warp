// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package wkit

import (
	"time"

	"github.com/lavaorg/lrt/mlog"
	"github.com/lavaorg/warp/warp9"
)

type Server interface {
	StartNetListener(ntype, addr string) error
	GetRoot() Directory
}

type Serv struct {
	warp9.Srv
	warp9.StatsOps
	qidp uint64
	root *DirItem
}

var srv *Serv

func init() {
	srv = new(Serv)
	srv.qidp = uint64(1000)
}

//
// host application invokes this when its ready to start serving objects.
//
// id is a symbolic name used for the root mount name
// dbgLevel sets level of debug to emit to standard out log
//
func StartServer(id string, dbgLevel int) Server {
	srv.Id = id
	srv.Debuglevel = dbgLevel
	srv.setRoot()
	srv.Start(srv)
	return srv
}

func (srv *Serv) GetRoot() Directory {
	return srv.root
}

func (srv *Serv) setRoot() {

	d := NewDirItem() //root is not a bind location; but is a diritem
	d.Name = srv.Id
	d.Uid = 1
	d.Gid = 1
	d.Muid = 1

	d.Atime = uint32(time.Now().Unix())
	d.Mtime = d.Atime

	// mark as a directory that allows READ operation for all users
	d.Mode = warp9.DMDIR | uint32(Perms(warp9.DMREAD, warp9.DMREAD, warp9.DMREAD))

	d.Qid = warp9.Qid{warp9.QTDIR, 0, srv.qidp}
	srv.qidp++

	mlog.Debug("getRoot:%v", d)
	srv.root = d
	return
}

func Perms(u, g, o byte) uint16 {
	return uint16(uint16(u)<<6 | uint16(g)<<3 | uint16(o))
}

func NextQid() uint64 {
	srv.qidp++
	return srv.qidp
}
