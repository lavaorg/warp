// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package warpsrv

import "github.com/lavaorg/warp9/ninep"

func (*W9Srv) Walk(req *ninep.SrvReq) {
	fid := req.Fid.Aux.(*W9Fid)
	tc := req.Tc

	if req.Newfid.Aux == nil {
		nfid := new(W9Fid)
		nfid.Fid = req.Newfid
		req.Newfid.Aux = nfid
	}

	nfid := req.Newfid.Aux.(*W9Fid)
	wqids := make([]ninep.Qid, len(tc.Wname))
	i := 0
	f := fid.F
	for ; i < len(tc.Wname); i++ {
		if tc.Wname[i] == ".." {
			// handle dotdot
			f = f.Parent
			wqids[i] = f.Qid
			continue
		}
		if (wqids[i].Type & ninep.QTDIR) > 0 {
			if !f.CheckPerm(req.Fid.User, ninep.DMEXEC) {
				break
			}
		}

		p := f.Find(tc.Wname[i])
		if p == nil {
			break
		}

		f = p
		wqids[i] = f.Qid
	}

	if len(tc.Wname) > 0 && i == 0 {
		req.RespondError(ninep.Err(ninep.Enoent))
		return
	}

	nfid.F = f
	req.RespondRwalk(wqids[0:i])
}
