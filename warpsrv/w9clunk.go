// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package warpsrv

type FClunkOp interface {
	Clunk(fid *FFid) error
}

func (*Fsrv) Clunk(req *SrvReq) {
	fid := req.Fid.Aux.(*FFid)

	if op, ok := (fid.F.ops).(FClunkOp); ok {
		err := op.Clunk(fid)
		if err != nil {
			req.RespondError(err)
		}
	}
	req.RespondRclunk()
}
