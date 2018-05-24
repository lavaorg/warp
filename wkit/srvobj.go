// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package wkit

import (
	"time"

	"github.com/lavaorg/lrt/mlog"
	"github.com/lavaorg/warp/warp9"
)

// Called when a client attaches to this server.
// the root is always "/"
//
// log the attach
//
func (srv *Serv) Attach(req *warp9.SrvReq) {
	if req.Afid != nil {
		req.RespondError(warp9.Error(warp9.Enoauth))
		return
	}
	//tc := req.Tc
	// ignore the aname; just mount "/"
	//rm fid := new(nullfsFid)
	//rm fid.entry = root
	req.Fid.Aux = srv.root
	mlog.Info("req.Fid:%v, root:%v", req.Fid, srv.root)
	req.RespondRattach(&srv.root.Qid)
}

// Flush has not function
func (*Serv) Flush(req *warp9.SrvReq) {}

// Ensure fid is a directory and invoke the
// walk method on that diretory.
// Promote the fid if succsfully moved
func (*Serv) Walk(req *warp9.SrvReq) {
	d, ok := req.Fid.Aux.(Directory)
	if !ok {
		req.RespondError(warp9.Error(warp9.Enotdir))
		return
	}
	if d == nil {
		req.RespondError(warp9.Error(warp9.Ebaduse))
		return
	}

	tc := req.Tc

	item, err := d.Walk(tc.Wname)
	if err != nil {
		req.RespondError(fsRespondError(err, warp9.Error(warp9.Enoent)))
		return
	}
	req.Newfid.Aux = item
	qid := item.GetQid()
	req.RespondRwalk(&qid)
}

// Invoke the objects SetOpened(true) method.
func (*Serv) Open(req *warp9.SrvReq) {
	i := req.Fid.Aux.(Item)
	//tc := req.Tc
	reqmode := req.Tc.Mode

	iounit, err := i.Open(reqmode)
	if err != nil {
		req.RespondError(fsRespondError(err, warp9.Error(warp9.Eio)))
		return
	}

	req.RespondRopen(&i.GetDir().Qid, iounit)
}

// invoke the object's SetOpened(false) method.
func (*Serv) Clunk(req *warp9.SrvReq) {
	i := req.Fid.Aux.(Item)
	err := i.Clunk()
	if err != nil {
		req.RespondError(fsRespondError(err, warp9.Error(warp9.Eio)))
		return
	}
	req.RespondRclunk()
}

// Ensure target is a directory and invoke its CreateItem method.
// Promote the fid to new object if successful.
func (*Serv) Create(req *warp9.SrvReq) {
	d, ok := req.Fid.Aux.(Directory)
	if !ok {
		req.RespondError(warp9.Error(warp9.Enotdir))
	}
	if d == nil {
		req.RespondError(warp9.Error(warp9.Ebaduse))
		return
	}

	tc := req.Tc

	item, err := d.Create(tc.Name, tc.Perm, tc.Mode)
	if err != nil {
		req.RespondError(fsRespondError(err, warp9.Error(warp9.Eio)))
		return
	}

	req.Fid.Aux = item
	req.RespondRcreate(&item.GetDir().Qid, 0)
}

// Invoke the object's Read() method.
func (*Serv) Read(req *warp9.SrvReq) {
	item := req.Fid.Aux.(Item)
	tc := req.Tc
	rc := req.Rc

	rc.InitRread(tc.Count)

	count, err := item.Read(rc.Data, tc.Offset, tc.Count)
	if err != nil {
		req.RespondError(fsRespondError(err, warp9.Error(warp9.Eio)))
		return
	}

	// change the a-time
	d := item.GetDir()
	d.Atime = uint32(time.Now().Unix())

	rc.SetRreadCount(count)
	req.Respond()
}

// Invoke the object's Write() method.
func (*Serv) Write(req *warp9.SrvReq) {
	item := req.Fid.Aux.(Item)
	tc := req.Tc

	count, err := item.Write(tc.Data, tc.Offset, tc.Count)
	if err != nil {
		req.RespondError(fsRespondError(err, warp9.Error(warp9.Eio)))
		return
	}

	// change the m-time and a-time
	d := item.GetDir()
	d.Atime = uint32(time.Now().Unix())
	d.Mtime = d.Atime

	req.RespondRwrite(count)
	return
}

// Not supported
func (*Serv) Remove(req *warp9.SrvReq) {
	i := req.Fid.Aux.(Item)
	err := i.Remove()
	if err != nil {
		req.RespondError(fsRespondError(err, warp9.Error(warp9.Eio)))
		return
	}
	req.RespondRremove()
	return
}

// Report the object's current status, reply with meta-data.
func (*Serv) Stat(req *warp9.SrvReq) {
	i := req.Fid.Aux.(Item)
	if i == nil {
		req.RespondError(warp9.Error(warp9.Ebaduse))
		return
	}
	dir, err := i.Stat()
	if err != nil {
		req.RespondError(fsRespondError(err, warp9.Error(warp9.Eio)))
		return
	}
	req.RespondRstat(dir)
	return
}

// not supported
func (u *Serv) Wstat(req *warp9.SrvReq) {
	req.RespondError(warp9.Error(warp9.Enotimpl))
	return
}

// helper functions

// error helper
func fsRespondError(err error, alterr *warp9.WarpError) *warp9.WarpError {
	werr, ok := err.(*warp9.WarpError)
	if !ok {
		werr = alterr
	}
	return werr
}
