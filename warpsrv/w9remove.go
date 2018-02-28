// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package warpsrv

import (
	"log"

	"github.com/lavaorg/warp/warp9"
)

// If the FRemoveOp interface is implemented, the Remove operation will be called
// when the client attempts to create a file in the srvFile implementing the interface.
// If not implemented, "permission denied" error will be send back.
// The operation returns nil if successful, or the error that occured while removing
// the file.
type W9RemoveCmd interface {
	W9Remove(*W9Fid) error
}

func (*W9Srv) Remove(req *warp9.SrvReq) {
	fid := req.Fid.Aux.(*W9Fid)
	f := fid.F
	f.Lock()
	if f.cfirst != nil {
		f.Unlock()
		req.RespondError(warp9.Err(warp9.Enotempty))
		return
	}
	f.Unlock()

	if rop, ok := (f.ops).(W9RemoveCmd); ok {
		err := rop.W9Remove(fid)
		if err != nil {
			req.RespondError(err)
		} else {
			f.Remove()
			req.RespondRremove()
		}
	} else {
		log.Println("remove not implemented")
		req.RespondError(warp9.Err(warp9.Eperm))
	}
}
