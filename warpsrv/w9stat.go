// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package warpsrv

// The FStatOp interface provides a single operation (Stat) that will be
// called before a file stat is sent back to the client. If implemented,
// the operation should update the data in the srvFile struct.
type FStatOp interface {
	Stat(fid *FFid) error
}

// The FWstatOp interface provides a single operation (Wstat) that will be
// called when the client requests the srvFile metadata to be modified. If
// implemented, the operation will be called when Twstat message is received.
// If not implemented, "permission denied" error will be sent back. If the
// operation returns an Error, the error is send back to the client.
type FWstatOp interface {
	Wstat(*FFid, *Dir) error
}

func (*Fsrv) Stat(req *SrvReq) {
	fid := req.Fid.Aux.(*FFid)
	f := fid.F

	if sop, ok := (f.ops).(FStatOp); ok {
		err := sop.Stat(fid)
		if err != nil {
			req.RespondError(err)
		} else {
			req.RespondRstat(&f.Dir)
		}
	} else {
		req.RespondRstat(&f.Dir)
	}
}

func (*Fsrv) Wstat(req *SrvReq) {
	tc := req.Tc
	fid := req.Fid.Aux.(*FFid)
	f := fid.F

	if wop, ok := (f.ops).(FWstatOp); ok {
		err := wop.Wstat(fid, &tc.Dir)
		if err != nil {
			req.RespondError(err)
		} else {
			req.RespondRwstat()
		}
	} else {
		req.RespondError(Err(Eperm))
	}
}
