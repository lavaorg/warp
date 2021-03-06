// Copyright 2009 The Go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package warp9

// Clunks a fid. Returns nil if successful.
func (clnt *Clnt) Clunk(fid *Fid) error {
	if fid.Fid == NOFID {
		return &WarpError{Efidnil, ""}
	}
	if fid.walked {
		tc := clnt.NewFcall()
		err := tc.packTclunk(fid.Fid)
		if err != nil {
			return err
		}

		_, err = clnt.Rpc(tc)
	}

	clnt.fidpool.putId(fid.Fid)
	fid.walked = false
	fid.Fid = NOFID
	return nil
}

// Closes an object. Returns nil if successful.
func (obj *Object) Close() error {
	// Should we cancel all pending requests for the Object
	return obj.Fid.Clnt.Clunk(obj.Fid)
}
