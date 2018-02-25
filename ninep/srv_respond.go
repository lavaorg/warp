// Copyright 2009 The Go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ninep

import "fmt"

// SrvRequest operations. This interface should be implemented by all file servers.
// The operations correspond directly to most of the 9P2000 message types.
type SrvReqOps interface {
	Attach(*SrvReq)
	Walk(*SrvReq)
	Open(*SrvReq)
	Create(*SrvReq)
	Read(*SrvReq)
	Write(*SrvReq)
	Clunk(*SrvReq)
	Remove(*SrvReq)
	Stat(*SrvReq)
	Wstat(*SrvReq)
}

// Respond to the request with Rerror message
func (req *SrvReq) RespondError(err interface{}) {
	switch e := err.(type) {
	case *Error:
		req.Rc.packRerror(e.Error(), uint32(e.Errornum), req.Conn.Dotu)
	case error:
		req.Rc.packRerror(e.Error(), uint32(EIO), req.Conn.Dotu)
	default:
		req.Rc.packRerror(fmt.Sprintf("%v", e), uint32(EIO), req.Conn.Dotu)
	}

	req.Respond()
}

// Respond to the request with Rversion message
func (req *SrvReq) RespondRversion(msize uint32, version string) {
	err := req.Rc.packRversion(msize, version)
	if err != nil {
		req.RespondError(err)
	} else {
		req.Respond()
	}
}

// Respond to the request with Rauth message
func (req *SrvReq) RespondRauth(aqid *Qid) {
	err := req.Rc.packRauth(aqid)
	if err != nil {
		req.RespondError(err)
	} else {
		req.Respond()
	}
}

// Respond to the request with Rflush message
func (req *SrvReq) RespondRflush() {
	err := req.Rc.packRflush()
	if err != nil {
		req.RespondError(err)
	} else {
		req.Respond()
	}
}

// Respond to the request with Rattach message
func (req *SrvReq) RespondRattach(aqid *Qid) {
	err := req.Rc.packRattach(aqid)
	if err != nil {
		req.RespondError(err)
	} else {
		req.Respond()
	}
}

// Respond to the request with Rwalk message
func (req *SrvReq) RespondRwalk(wqids []Qid) {
	err := req.Rc.packRwalk(wqids)
	if err != nil {
		req.RespondError(err)
	} else {
		req.Respond()
	}
}

// Respond to the request with Ropen message
func (req *SrvReq) RespondRopen(qid *Qid, iounit uint32) {
	err := req.Rc.packRopen(qid, iounit)
	if err != nil {
		req.RespondError(err)
	} else {
		req.Respond()
	}
}

// Respond to the request with Rcreate message
func (req *SrvReq) RespondRcreate(qid *Qid, iounit uint32) {
	err := req.Rc.packRcreate(qid, iounit)
	if err != nil {
		req.RespondError(err)
	} else {
		req.Respond()
	}
}

// Respond to the request with Rread message
func (req *SrvReq) RespondRread(data []byte) {
	err := req.Rc.packRread(data)
	if err != nil {
		req.RespondError(err)
	} else {
		req.Respond()
	}
}

// Respond to the request with Rwrite message
func (req *SrvReq) RespondRwrite(count uint32) {
	err := req.Rc.packRwrite(count)
	if err != nil {
		req.RespondError(err)
	} else {
		req.Respond()
	}
}

// Respond to the request with Rclunk message
func (req *SrvReq) RespondRclunk() {
	err := req.Rc.packRclunk()
	if err != nil {
		req.RespondError(err)
	} else {
		req.Respond()
	}
}

// Respond to the request with Rremove message
func (req *SrvReq) RespondRremove() {
	err := req.Rc.packRremove()
	if err != nil {
		req.RespondError(err)
	} else {
		req.Respond()
	}
}

// Respond to the request with Rstat message
func (req *SrvReq) RespondRstat(st *Dir) {
	err := req.Rc.packRstat(st, req.Conn.Dotu)
	if err != nil {
		req.RespondError(err)
	} else {
		req.Respond()
	}
}

// Respond to the request with Rwstat message
func (req *SrvReq) RespondRwstat() {
	err := req.Rc.packRwstat()
	if err != nil {
		req.RespondError(err)
	} else {
		req.Respond()
	}
}
