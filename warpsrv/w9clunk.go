// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package warpsrv

import "github.com/lavaorg/warp9/warp9"

type FClunkOp interface {
	Clunk(fid *W9Fid) error
}

func (*W9Srv) Clunk(req *warp9.SrvReq) {
	fid := req.Fid.Aux.(*W9Fid)

	if op, ok := (fid.F.ops).(FClunkOp); ok {
		err := op.Clunk(fid)
		if err != nil {
			req.RespondError(err)
		}
	}
	req.RespondRclunk()
}
