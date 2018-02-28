// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package warpsrv

import "github.com/lavaorg/warp9/warp9"

// If the FReadOp interface is implemented, the Read operation will be called
// to read from the file. If not implemented, "permission denied" error will
// be send back. The operation returns the number of bytes read, or the
// error occured while reading.
type FReadOp interface {
	Read(fid *W9Fid, buf []byte, offset uint64) (int, error)
}

// If the FWriteOp interface is implemented, the Write operation will be called
// to write to the file. If not implemented, "permission denied" error will
// be send back. The operation returns the number of bytes written, or the
// error occured while writing.
type FWriteOp interface {
	Write(fid *W9Fid, data []byte, offset uint64) (int, error)
}

func (*W9Srv) Read(req *warp9.SrvReq) {
	var n int
	var err error

	fid := req.Fid.Aux.(*W9Fid)
	f := fid.F
	tc := req.Tc
	rc := req.Rc
	rc.InitRread(tc.Count)

	if f.Mode&warp9.DMDIR != 0 {
		// Get all the directory entries and
		// serialize them all into an output buffer.
		// This greatly simplifies the directory read.
		if tc.Offset == 0 {
			var g *W9File
			fid.dirents = nil
			f.Lock()
			for n, g = 0, f.cfirst; g != nil; n, g = n+1, g.next {
			}
			fid.dirs = make([]*W9File, n)
			for n, g = 0, f.cfirst; g != nil; n, g = n+1, g.next {
				fid.dirs[n] = g
				fid.dirents = append(fid.dirents,
					warp9.PackDir(&g.Dir, req.Conn.Dotu)...)
			}
			f.Unlock()
		}

		switch {
		case tc.Offset > uint64(len(fid.dirents)):
			n = 0
		case len(fid.dirents[tc.Offset:]) > int(tc.FcSize):
			n = int(tc.FcSize)
		default:
			n = len(fid.dirents[tc.Offset:])
		}
		copy(rc.Data, fid.dirents[tc.Offset:int(tc.Offset)+1+n])

	} else {
		// file
		if rop, ok := f.ops.(FReadOp); ok {
			n, err = rop.Read(fid, rc.Data, tc.Offset)
			if err != nil {
				req.RespondError(err)
				return
			}
		} else {
			req.RespondError(warp9.Err(warp9.Eperm))
			return
		}
	}

	rc.SetRreadCount(uint32(n))
	req.Respond()
}

func (*W9Srv) Write(req *warp9.SrvReq) {
	fid := req.Fid.Aux.(*W9Fid)
	f := fid.F
	tc := req.Tc

	if wop, ok := (f.ops).(FWriteOp); ok {
		n, err := wop.Write(fid, tc.Data, tc.Offset)
		if err != nil {
			req.RespondError(err)
		} else {
			req.RespondRwrite(uint32(n))
		}
	} else {
		req.RespondError(warp9.Err(warp9.Eperm))
	}

}
