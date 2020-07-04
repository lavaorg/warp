// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package wkit

import (
	"time"

	"github.com/lavaorg/warp/warp9"
)


// called when SrvFid is destroyed 
func (srv *ServerController) FidDestroy(sfid *warp9.SrvFid) {

	// if an Item is found then invoke clunk on it
	if i, ok := sfid.Aux.(Item); ok {
		err := i.Clunk()
		if err != nil {
			//nothing to do with error; log it.
			warp9.Error("ignoring FidDestroy:Clunk() error: %v", err)
		}
	}
}

// Called when a client attaches to this server.
// the root is always "/"
//
// log the attach
//
func (srv *ServerController) Attach(req *warp9.SrvReq) {
	if req.Afid != nil {
		req.RespondError(warp9.ErrorCode(warp9.Enoauth))
		return
	}

	req.Fid.Aux = srv.root //associate this server's root with client's fid
	warp9.Info("req.Fid:%v, root:%v", req.Fid, srv.root)
	qid := srv.root.GetQid()
	req.RespondRattach(&qid)
}

// Flush has no function
func (*ServerController) Flush(req *warp9.SrvReq) {}

// Ensure fid is a directory and invoke the
// walk method on that directory.
// Promote the fid if successfully moved
func (*ServerController) Walk(req *warp9.SrvReq) {
	d, ok := req.Fid.Aux.(Directory)
	if !ok {
		req.RespondError(warp9.ErrorCode(warp9.Enotdir))
		return
	}
	if d == nil {
		req.RespondError(warp9.ErrorCode(warp9.Ebaduse))
		return
	}

	tc := req.Tc

	item, err := d.Walk(tc.Wname)
	if err != nil {
		req.RespondError(fsRespondError(err, warp9.ErrorCode(warp9.Enoent)))
		return
	}
	req.Newfid.Aux = item
	qid := item.GetQid()
	req.RespondRwalk(&qid)
}

// Invoke the objects SetOpened(true) method.
func (*ServerController) Open(req *warp9.SrvReq) {
	i := req.Fid.Aux.(Item)
	//tc := req.Tc
	reqmode := req.Tc.Mode

	iounit, err := i.Open(reqmode)
	if err != nil {
		req.RespondError(fsRespondError(err, warp9.ErrorCode(warp9.Eio)))
		return
	}

	req.RespondRopen(&i.GetDir().Qid, iounit)
}

// invoke the object's SetOpened(false) method.
func (*ServerController) Clunk(req *warp9.SrvReq) {
	i := req.Fid.Aux.(Item)
	err := i.Clunk()
	if err != nil {
		req.RespondError(fsRespondError(err, warp9.ErrorCode(warp9.Eio)))
		return
	}
	req.Fid.Aux = nil //disassociate Item from fid
	req.RespondRclunk()
}

// Ensure target is a directory and invoke its CreateItem method.
// Promote the fid to new object if successful.
func (*ServerController) Create(req *warp9.SrvReq) {
	d, ok := req.Fid.Aux.(Directory)
	if !ok {
		req.RespondError(warp9.ErrorCode(warp9.Enotdir))
	}
	if d == nil {
		req.RespondError(warp9.ErrorCode(warp9.Ebaduse))
		return
	}

	// tc is the incoming message
	tc := req.Tc

	item := NewDirItem(tc.Name)
	item.SetMode(tc.Perm)

	req.Fid.Aux = item
	req.RespondRcreate(&item.GetDir().Qid, 0)
}

// Invoke the object's Read() method.
func (*ServerController) Read(req *warp9.SrvReq) {
	item := req.Fid.Aux.(Item)
	tc := req.Tc
	rc := req.Rc

	rc.InitRread(tc.Count)

	count, err := item.Read(rc.Data, tc.Offset, tc.Count)
	if err != nil {
		req.RespondError(fsRespondError(err, warp9.ErrorCode(warp9.Eio)))
		return
	}

	// change the a-time
	d := item.GetDir()
	d.Atime = uint32(time.Now().Unix())

	rc.SetRreadCount(count)
	req.Respond()
}

// Invoke the object's Write() method.
func (*ServerController) Write(req *warp9.SrvReq) {
	item := req.Fid.Aux.(Item)
	tc := req.Tc
	warp9.Debug("Write: %T:%v", item, item)
	count, err := item.Write(tc.Data, tc.Offset, tc.Count)
	if err != nil {
		req.RespondError(fsRespondError(err, warp9.ErrorCode(warp9.Eio)))
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
func (*ServerController) Remove(req *warp9.SrvReq) {
	i := req.Fid.Aux.(Item)
	err := i.Remove()
	if err != nil {
		req.RespondError(fsRespondError(err, warp9.ErrorCode(warp9.Eio)))
		return
	}
	req.RespondRremove()
	return
}

// Report the object's current status, reply with meta-data.
func (*ServerController) Stat(req *warp9.SrvReq) {
	i := req.Fid.Aux.(Item)
	if i == nil {
		req.RespondError(warp9.ErrorCode(warp9.Ebaduse))
		return
	}
	dir, err := i.Stat()
	if err != nil {
		req.RespondError(fsRespondError(err, warp9.ErrorCode(warp9.Eio)))
		return
	}
	req.RespondRstat(dir)
	return
}

// not supported
func (u *ServerController) Wstat(req *warp9.SrvReq) {
	req.RespondError(warp9.ErrorCode(warp9.Enotimpl))
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
