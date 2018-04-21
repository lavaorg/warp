// Copyright 2009 The Go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package warp9

// Returns the metadata for a named object, or an Error.
func (clnt *Clnt) Stat(path string) (*Dir, W9Err) {
	fid, err := clnt.Walk(path)
	if err != Egood {
		return nil, err
	}

	d, err := clnt.FStat(fid)
	clnt.Clunk(fid)
	return d, err
}

// Returns the metadata for the object associated with the Fid, or an Error.
func (clnt *Clnt) FStat(fid *Fid) (*Dir, W9Err) {
	tc := clnt.NewFcall()
	err := tc.packTstat(fid.Fid)
	if err != Egood {
		return nil, err
	}

	rc, err := clnt.Rpc(tc)
	if err != Egood {
		return nil, err
	}

	return &rc.Dir, Egood
}

// Modifies the data of the object associated with the Fid, or an Error.
func (clnt *Clnt) FWstat(fid *Fid, dir *Dir) W9Err {
	tc := clnt.NewFcall()
	err := tc.packTwstat(fid.Fid, dir)
	if err != Egood {
		return err
	}

	_, err = clnt.Rpc(tc)
	return err
}
