// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package warpsrv

import "github.com/lavaorg/warp/warp9"

// If the FCreateOp interface is implemented, the Create operation will be called
// when the client attempts to create a file in the W9File implementing the interface.
// If not implemented, "permission denied" error will be send back. If successful,
// the operation should call (*File)Add() to add the created file to the directory.
// The operation returns the created file, or the error occured while creating it.
type FCreateOp interface {
	Create(fid *W9Fid, name string, perm uint32) (*W9File, error)
}

type FOpenOp interface {
	Open(fid *W9Fid, mode uint8) error
}

func (*W9Srv) Open(req *warp9.SrvReq) {
	fid := req.Fid.Aux.(*W9Fid)
	tc := req.Tc

	if !fid.F.CheckPerm(req.Fid.User, mode2Perm(tc.Mode)) {
		req.RespondError(warp9.Err(warp9.Eperm))
		return
	}

	if op, ok := (fid.F.ops).(FOpenOp); ok {
		err := op.Open(fid, tc.Mode)
		if err != nil {
			req.RespondError(err)
			return
		}
	}
	req.RespondRopen(&fid.F.Qid, 0)
}

func (*W9Srv) Create(req *warp9.SrvReq) {
	fid := req.Fid.Aux.(*W9Fid)
	tc := req.Tc

	dir := fid.F
	if !dir.CheckPerm(req.Fid.User, warp9.DMWRITE) {
		req.RespondError(warp9.Err(warp9.Eperm))
		return
	}

	if cop, ok := (dir.ops).(FCreateOp); ok {
		f, err := cop.Create(fid, tc.Name, tc.Perm)
		if err != nil {
			req.RespondError(err)
		} else {
			fid.F = f
			req.RespondRcreate(&fid.F.Qid, 0)
		}
	} else {
		req.RespondError(warp9.Err(warp9.Eperm))
	}
}
