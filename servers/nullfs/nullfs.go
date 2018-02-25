// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

// Nullfs is primarily a template to start building multi-level FS from.
// This is a read only file system with one file /info.  This means you can
// open/read/stat/ the files "/" and "/info"
// "/" is a directory
// "/info" is a read only file
//
package nullfs

import (
	"log"

	"github.com/lavaorg/warp9/ninep"
)

type nullfsFid struct {
	entry       *NullfsDir
	direntrybuf []byte
}

type Nullfs struct {
	ninep.Srv
	ninep.StatsOps
}

// NullDir represents an entry in the NullFS. It will contain a ninep dir
// but can carry additional data/state needed as necessary
type NullfsDir struct {
	ninep.Dir
}

var root *NullfsDir = newNullfsDir("/")

func newNullfsDir(n string) *NullfsDir {
	var d NullfsDir

	d.Name = n
	d.Uid = "nobody"
	d.Gid = "nobody"
	d.Muid = "nobody"

	d.Mode = ninep.DMDIR | uint32(perms(ninep.DMREAD, ninep.DMREAD, ninep.DMREAD))
	d.Atime = 0
	d.Mtime = 0

	d.Qid = ninep.Qid{ninep.QTDIR, 0, 9999}

	d.Type = 0
	d.Dev = 0

	return &d
}

func perms(u, g, o byte) uint16 {
	return uint16(uint16(u)<<6 | uint16(g)<<3 | uint16(o))
}

func (*Nullfs) ConnOpened(conn *ninep.Conn) {
	if conn.Srv.Debuglevel > 0 {
		log.Println("connected")
	}
}

func (*Nullfs) ConnClosed(conn *ninep.Conn) {
	if conn.Srv.Debuglevel > 0 {
		log.Println("disconnected")
	}
}

func (*Nullfs) FidDestroy(sfid *ninep.SrvFid) {
	var fid *nullfsFid

	if sfid.Aux == nil {
		return
	}

	fid = sfid.Aux.(*nullfsFid)
	if sfid.Fconn.Debuglevel > 0 {
		log.Printf("fid destroy:%v\n", fid)
	}
	//cleanup fid
}

func (ufs *Nullfs) Attach(req *ninep.SrvReq) {
	if req.Afid != nil {
		req.RespondError(ninep.Enoauth)
		return
	}
	//tc := req.Tc
	// ignore the aname; just mount "/"
	fid := new(nullfsFid)
	fid.entry = root
	req.Fid.Aux = fid
	req.RespondRattach(&root.Qid)
}

func (*Nullfs) Flush(req *ninep.SrvReq) {}

func (*Nullfs) Walk(req *ninep.SrvReq) {
	fid := req.Fid.Aux.(*nullfsFid)
	tc := req.Tc

	if fid == nil {
		req.RespondError("bad fid")
		return
	}

	if req.Newfid.Aux == nil {
		req.Newfid.Aux = new(nullfsFid)
	}

	// there are no entries so if path is not "." or ".." return an error
	// "." and ".." by definition are alias for the current node, so valid.
	if len(tc.Wname) != 1 {
		req.RespondError(ninep.Enoent)
		return
	}
	p := tc.Wname[0]
	if p != "." && p != ".." {
		req.RespondError(ninep.Enoent)
		return
	}

	req.Newfid.Aux = req.Fid.Aux
	wqids := make([]ninep.Qid, 1)
	wqids[0] = fid.entry.Qid

	req.RespondRwalk(wqids[0:])
}

func (*Nullfs) Open(req *ninep.SrvReq) {

	tc := req.Tc
	mode := tc.Mode
	if mode != ninep.OREAD {
		req.RespondError(ninep.EIO)
		return
	}

	req.RespondRopen(&root.Qid, 0)
}

func (*Nullfs) Create(req *ninep.SrvReq) {
	// no creation
	req.RespondError(&ninep.Error{"not supported", ninep.EIO})
}

func (*Nullfs) Read(req *ninep.SrvReq) {
	tc := req.Tc
	rc := req.Rc

	rc.InitRread(tc.Count)

	// convert our directory to byte buffer; we aren't caching
	b := ninep.PackDir(&root.Dir, req.Conn.Dotu)

	// determine which and how many bytes to return
	var count int
	switch {
	case tc.Offset > uint64(len(b)):
		count = 0
	case len(b[tc.Offset:]) > int(tc.Count):
		count = int(tc.Count)
	default:
		count = len(b[tc.Offset:])
	}
	copy(rc.Data, b[tc.Offset:int(tc.Offset)+count])
	log.Printf("buf:%v, rc.Data: %v, off:%v,  count:%v\n", len(b), len(rc.Data), tc.Offset, count)
	rc.SetRreadCount(uint32(count))
	req.Respond()
}

func (*Nullfs) Write(req *ninep.SrvReq) {
	req.RespondError(&ninep.Error{"write not supported", ninep.EIO})
	return
}

func (*Nullfs) Clunk(req *ninep.SrvReq) { req.RespondRclunk() }

func (*Nullfs) Remove(req *ninep.SrvReq) {
	req.RespondError(&ninep.Error{"remove not supported", ninep.EIO})
	return
}

func (*Nullfs) Stat(req *ninep.SrvReq) {
	req.RespondError(&ninep.Error{"stat not supported", ninep.EIO})
	return
}
func (u *Nullfs) Wstat(req *ninep.SrvReq) {
	req.RespondError(&ninep.Error{"wstat not supported", ninep.EIO})
	return
}
