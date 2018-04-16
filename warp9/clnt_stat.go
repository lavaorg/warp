// Copyright 2009 The Go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package warp9

// Returns the metadata for the file associated with the Fid, or an Error.
func (clnt *Clnt) Stat(fid *Fid) (*Dir, W9Err) {
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

// Returns the metadata for a named file, or an Error.
func (clnt *Clnt) FStat(path string) (*Dir, W9Err) {
	fid, err := clnt.FWalk(path)
	if err != Egood {
		return nil, err
	}

	d, err := clnt.Stat(fid)
	clnt.Clunk(fid)
	return d, err
}

// Modifies the data of the file associated with the Fid, or an Error.
func (clnt *Clnt) Wstat(fid *Fid, dir *Dir) W9Err {
	tc := clnt.NewFcall()
	err := tc.packTwstat(fid.Fid, dir)
	if err != Egood {
		return err
	}

	_, err = clnt.Rpc(tc)
	return err
}
