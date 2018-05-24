// Copyright 2009 The Go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package warp9

// Removes the object associated with the Fid. Returns nil if the
// operation is successful.
func (clnt *Clnt) FRemove(fid *Fid) error {
	tc := clnt.NewFcall()
	err := tc.packTremove(fid.Fid)
	if err != nil {
		return err
	}

	_, err = clnt.Rpc(tc)
	clnt.fidpool.putId(fid.Fid)
	fid.Fid = NOFID

	return err
}

// Removes the named object. Returns nil if the operation is successful.
func (clnt *Clnt) Remove(path string) error {
	fid, err := clnt.Walk(path)
	if err != nil {
		return err
	}

	err = clnt.FRemove(fid)
	return err
}
