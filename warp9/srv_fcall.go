// Copyright 2009 The Go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package warp9

// the set of methods in this file manage the common behavor each of the serving Warp9 message handling.
//
// FCall represents the message (see fcall.go)
//
// A warp9 server can be viewed as a set of layers.
//  1. low level message parsing
//  2. primitive message handling (this file)
//  3. user message handling (the methods here inovke them)
//  4. user objects handling
//
// Many of the message handling methods follow the form of:
//  <preample-handler> --> <user-handler> --> <post-handler>
//
// this allows common book-keeping to be done on behalf of the user's handlers
//

func (srv *Srv) version(req *SrvReq) {
	tc := req.Tc
	conn := req.Conn

	if tc.Version != Warp9Version {
		req.RespondError(&WarpError{Ebadver, ""})
		return
	}

	if tc.Msize < IOHDRSZ {
		req.RespondError(&WarpError{Emsize, ""})
		return
	}

	if tc.Msize < conn.Msize {
		conn.Msize = tc.Msize
	}

	/* make sure that the responses of all current requests will be ignored */
	conn.Lock()
	for tag, r := range conn.reqs {
		if tag == NOTAG {
			continue
		}

		for rr := r; rr != nil; rr = rr.next {
			rr.Lock()
			rr.status |= reqFlush
			rr.Unlock()
		}
	}
	conn.Unlock()

	req.RespondRversion(conn.Msize, Warp9Version)
}

func (srv *Srv) auth(req *SrvReq) {
	tc := req.Tc
	conn := req.Conn
	if tc.Atok == NOTOK {
		req.RespondError(&WarpError{Eunknownfid, ""})
		return
	}

	req.Afid = conn.FidNew(tc.Atok)
	if req.Afid == nil {
		req.RespondError(&WarpError{Einuse, ""})
		return
	}

	var user User = nil
	if tc.Uid != NOUID {
		user = srv.Upool.User(tc.Uid)
	}

	if user == nil {
		req.RespondError(&WarpError{Enouser, ""})
		return
	}

	req.Afid.User = user
	req.Afid.Type = QTAUTH
	if aop, ok := (srv.ops).(AuthOps); ok {
		aqid, err := aop.AuthInit(req.Afid, tc.Aname)
		if err != nil {
			req.RespondError(&WarpError{Eauthinit, ""})
		} else {
			aqid.Type |= QTAUTH // just in case
			req.RespondRauth(aqid)
		}
	} else {
		req.RespondError(&WarpError{Enoauth, ""})
	}

}

func (srv *Srv) authPost(req *SrvReq) {
	if req.Rc != nil && req.Rc.Type == Rauth {
		req.Afid.IncRef()
	}
}

func (srv *Srv) attach(req *SrvReq) {
	tc := req.Tc
	conn := req.Conn
	if tc.Fid == NOFID {
		req.RespondError(&WarpError{Eunknownfid, ""})
		return
	}

	req.Fid = conn.FidNew(tc.Fid)
	if req.Fid == nil {
		req.RespondError(&WarpError{Einuse, ""})
		return
	}

	if tc.Atok != NOTOK {
		req.Afid = conn.FidGet(tc.Atok)
		if req.Afid == nil {
			req.RespondError(&WarpError{Eunknownfid, ""})
		}
	}

	var user User = nil
	if tc.Uid != NOUID {
		user = srv.Upool.User(tc.Uid)
	}

	if user == nil {
		req.RespondError(&WarpError{Enouser, ""})
		return
	}

	req.Fid.User = user
	if aop, ok := (srv.ops).(AuthOps); ok {
		err := aop.AuthCheck(req.Fid, req.Afid, tc.Aname)
		if err != nil {
			req.RespondError(&WarpError{Eauthchk, ""})
			return
		}
	}

	(srv.ops).(SrvReqOps).Attach(req)
}

func (srv *Srv) attachPost(req *SrvReq) {
	if req.Rc != nil && req.Rc.Type == Rattach {
		req.Fid.Type = req.Rc.Qid.Type
		req.Fid.IncRef()
	}
}

func (srv *Srv) flush(req *SrvReq) {
	conn := req.Conn
	tag := req.Tc.Oldtag
	req.Rc.packRflush()
	conn.Lock()
	r := conn.reqs[tag]
	if r != nil {
		req.flushreq = r.flushreq
		r.flushreq = req
	}
	conn.Unlock()

	if r == nil {
		// there are no requests with that tag
		req.Respond()
		return
	}

	r.Lock()
	status := r.status
	if (status & (reqWork | reqSaved)) == 0 {
		/* the request is not worked on yet */
		r.status |= reqFlush
	}
	r.Unlock()

	if (status & (reqWork | reqSaved)) == 0 {
		r.Respond()
	} else {
		if op, ok := (srv.ops).(FlushOp); ok {
			op.Flush(r)
		}
	}
}

func (srv *Srv) walk(req *SrvReq) {
	conn := req.Conn
	tc := req.Tc
	fid := req.Fid

	if fid == nil {
		req.RespondError(&WarpError{Efidnil, ""})
		return
	}

	// we can't walk regular objects, only clone them
	if len(tc.Wname) > 0 && (fid.Type&QTDIR) == 0 {
		req.RespondError(&WarpError{Enotdir, ""})
		return
	}

	// some common bad names -- client should not pass
	if len(tc.Wname) == 1 {
		p := tc.Wname[0]
		if p == "." || p == ".." || p == "/" {
			req.RespondError(&WarpError{Ename, ""})
			return
		}
	}

	//we can't walk open objects
	if fid.opened {
		req.RespondError(&WarpError{Eopen, ""})
		return
	}

	if req.Fid.Aux == nil {
		req.RespondError(&WarpError{Efidnoaux, ""})
	}
	if tc.Fid != tc.Newfid {
		req.Newfid = conn.FidNew(tc.Newfid)
		if req.Newfid == nil {
			req.RespondError(&WarpError{Einuse, ""})
			return
		}

		req.Newfid.User = fid.User
		req.Newfid.Type = fid.Type
	} else {
		req.Newfid = req.Fid
		req.Newfid.IncRef()
	}

	(req.Conn.Srv.ops).(SrvReqOps).Walk(req)
}

func (srv *Srv) walkPost(req *SrvReq) {
	rc := req.Rc
	if rc == nil || rc.Type != Rwalk || req.Newfid == nil {
		return
	}

	req.Newfid.Type = rc.Qid.Type

	if req.Newfid.fid != req.Fid.fid {
		req.Newfid.IncRef()
	}
}

func (srv *Srv) open(req *SrvReq) {
	fid := req.Fid
	tc := req.Tc

	if fid == nil {
		req.RespondError(&WarpError{Efidnil, ""})
		return
	}

	if fid.opened {
		req.RespondError(&WarpError{Eopen, ""})
		return
	}

	if (fid.Type&QTDIR) != 0 && tc.Mode != OREAD {
		req.RespondError(&WarpError{Eperm, ""})
		return
	}

	fid.Omode = tc.Mode
	(req.Conn.Srv.ops).(SrvReqOps).Open(req)
}

func (srv *Srv) openPost(req *SrvReq) {
	if req.Fid != nil {
		req.Fid.opened = req.Rc != nil && req.Rc.Type == Ropen
	}
}

func (srv *Srv) create(req *SrvReq) {
	fid := req.Fid
	tc := req.Tc

	if fid == nil {
		req.RespondError(&WarpError{Efidnil, ""})
		return
	}
	if fid.opened {
		req.RespondError(&WarpError{Eopen, ""})
		return
	}

	if (fid.Type & QTDIR) == 0 {
		req.RespondError(&WarpError{Enotdir, ""})
		return
	}

	/* can't open directories for other than reading */
	if (tc.Perm&DMDIR) != 0 && tc.Mode != OREAD {
		req.RespondError(&WarpError{Eperm, ""})
		return
	}

	fid.Omode = tc.Mode
	(req.Conn.Srv.ops).(SrvReqOps).Create(req)
}

func (srv *Srv) createPost(req *SrvReq) {
	if req.Rc != nil && req.Rc.Type == Rcreate && req.Fid != nil {
		req.Fid.Type = req.Rc.Qid.Type
		req.Fid.opened = true
	}
}

func (srv *Srv) read(req *SrvReq) {
	tc := req.Tc
	fid := req.Fid

	if fid == nil {
		req.RespondError(&WarpError{Efidnil, ""})
		return
	}

	if tc.Count+IOHDRSZ > req.Conn.Msize {
		req.RespondError(&WarpError{Etoolarge, ""})
		return
	}

	if (fid.Type & QTAUTH) != 0 {

		rc := req.Rc
		err := rc.InitRread(tc.Count)
		if err != nil {
			req.RespondError(err.(*WarpError))
			return
		}

		if op, ok := (req.Conn.Srv.ops).(AuthOps); ok {
			n, e := op.AuthRead(fid, tc.Offset, rc.Data)
			if e != nil {
				req.RespondError(&WarpError{Eauthread, ""})
				return
			}
			rc.SetRreadCount(uint32(n))
			req.Respond()
		} else {
			req.RespondError(&WarpError{Enotimpl, ""})
		}

		return
	}

	if (fid.Type & QTDIR) != 0 {
		if tc.Offset == 0 {
			fid.Diroffset = 0
		} else if tc.Offset != fid.Diroffset {
			fid.Diroffset = tc.Offset
		}
	}

	(req.Conn.Srv.ops).(SrvReqOps).Read(req)
}

func (srv *Srv) readPost(req *SrvReq) {
	if req.Rc != nil && req.Rc.Type == Rread && (req.Fid.Type&QTDIR) != 0 {
		req.Fid.Diroffset += uint64(req.Rc.Count)
	}
}

func (srv *Srv) write(req *SrvReq) {
	fid := req.Fid
	tc := req.Tc

	if fid == nil {
		req.RespondError(&WarpError{Efidnil, ""})
		return
	}

	if (fid.Type & QTAUTH) != 0 {
		tc := req.Tc
		if op, ok := (req.Conn.Srv.ops).(AuthOps); ok {
			n, err := op.AuthWrite(req.Fid, tc.Offset, tc.Data)
			if err != nil {
				//log err??
				req.RespondError(&WarpError{Eauthwrite, ""})
			} else {
				req.RespondRwrite(uint32(n))
			}
		} else {
			req.RespondError(&WarpError{Enotimpl, ""})
		}

		return
	}

	if !fid.opened || (fid.Type&QTDIR) != 0 || (fid.Omode&3) == OREAD {
		req.RespondError(&WarpError{Ebaduse, ""})
		return
	}

	if tc.Count+IOHDRSZ > req.Conn.Msize {
		req.RespondError(&WarpError{Etoolarge, ""})
		return
	}

	(req.Conn.Srv.ops).(SrvReqOps).Write(req)
}

func (srv *Srv) clunk(req *SrvReq) {
	fid := req.Fid

	if fid == nil {
		req.RespondError(&WarpError{Efidnil, ""})
		return
	}

	if (fid.Type & QTAUTH) != 0 {
		if op, ok := (req.Conn.Srv.ops).(AuthOps); ok {
			op.AuthDestroy(fid)
			req.RespondRclunk()
		} else {
			req.RespondError(&WarpError{Enotimpl, ""})
		}

		return
	}

	(req.Conn.Srv.ops).(SrvReqOps).Clunk(req)
}

func (srv *Srv) clunkPost(req *SrvReq) {
	if req.Rc != nil && req.Rc.Type == Rclunk && req.Fid != nil {
		req.Fid.DecRef()
	}
}

func (srv *Srv) remove(req *SrvReq) {

	if req.Fid == nil {
		req.RespondError(&WarpError{Efidnil, ""})
		return
	}

	(req.Conn.Srv.ops).(SrvReqOps).Remove(req)
}

func (srv *Srv) removePost(req *SrvReq) {
	if req.Rc != nil && req.Fid != nil {
		req.Fid.DecRef()
	}
}

func (srv *Srv) stat(req *SrvReq) {

	if req.Fid == nil {
		req.RespondError(&WarpError{Efidnil, ""})
		return
	}

	(req.Conn.Srv.ops).(SrvReqOps).Stat(req)
}

func (srv *Srv) wstat(req *SrvReq) {
	/*
		fid := req.Fid
		d := &req.Tc.Dir
		if d.Type != uint16(0xFFFF) || d.Dev != uint32(0xFFFFFFFF) || d.Version != uint32(0xFFFFFFFF) ||
			d.Path != uint64(0xFFFFFFFFFFFFFFFF) {
			req.RespondError(Eperm))
			return
		}

		if (d.Mode != 0xFFFFFFFF) && (((fid.Type&QTDIR) != 0 && (d.Mode&DMDIR) == 0) ||
			((d.Type&QTDIR) == 0 && (d.Mode&DMDIR) != 0)) {
			req.RespondError(Edirchange)
			return
		}
	*/
	if req.Fid == nil {
		req.RespondError(&WarpError{Efidnil, ""})
		return
	}

	(req.Conn.Srv.ops).(SrvReqOps).Wstat(req)
}
