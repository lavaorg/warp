// Copyright 2009 The Go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package warp9

// Removes the file associated with the Fid. Returns nil if the
// operation is successful.
func (clnt *Clnt) Remove(fid *Fid) W9Err {
	tc := clnt.NewFcall()
	err := tc.packTremove(fid.Fid)
	if err != Egood {
		return err
	}

	_, err = clnt.Rpc(tc)
	clnt.fidpool.putId(fid.Fid)
	fid.Fid = NOFID

	return err
}

// Removes the named file. Returns nil if the operation is successful.
func (clnt *Clnt) FRemove(path string) W9Err {
	var err W9Err
	fid, err := clnt.FWalk(path)
	if err != Egood {
		return err
	}

	err = clnt.Remove(fid)
	return err
}
