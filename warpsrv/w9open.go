// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package warpsrv

// If the FCreateOp interface is implemented, the Create operation will be called
// when the client attempts to create a file in the srvFile implementing the interface.
// If not implemented, "permission denied" error will be send back. If successful,
// the operation should call (*File)Add() to add the created file to the directory.
// The operation returns the created file, or the error occured while creating it.
type FCreateOp interface {
	Create(fid *FFid, name string, perm uint32) (*srvFile, error)
}

type FOpenOp interface {
	Open(fid *FFid, mode uint8) error
}

func (*Fsrv) Open(req *SrvReq) {
	fid := req.Fid.Aux.(*FFid)
	tc := req.Tc

	if !fid.F.CheckPerm(req.Fid.User, mode2Perm(tc.Mode)) {
		req.RespondError(Err(Eperm))
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

func (*Fsrv) Create(req *SrvReq) {
	fid := req.Fid.Aux.(*FFid)
	tc := req.Tc

	dir := fid.F
	if !dir.CheckPerm(req.Fid.User, DMWRITE) {
		req.RespondError(Err(Eperm))
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
		req.RespondError(Err(Eperm))
	}
}
