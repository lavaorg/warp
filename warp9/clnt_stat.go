// Copyright 2009 The Go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package warp9

// Returns the metadata for a named object, or an Error.
func (clnt *Clnt) Stat(path string) (*Dir, error) {
	fid, err := clnt.Walk(path)
	if err != nil {
		return nil, err
	}

	d, err := clnt.FStat(fid)
	clnt.Clunk(fid)
	return d, err
}

// Returns the metadata for the object associated with the Fid, or an Error.
func (clnt *Clnt) FStat(fid *Fid) (*Dir, error) {
	tc := clnt.NewFcall()
	err := tc.packTstat(fid.Fid)
	if err != nil {
		return nil, err
	}

	rc, err := clnt.Rpc(tc)
	if err != nil {
		return nil, err
	}

	return &rc.Dir, nil
}

// Modifies the data of the object associated with the Fid, or an Error.
func (clnt *Clnt) FWstat(fid *Fid, dir *Dir) error {
	tc := clnt.NewFcall()
	err := tc.packTwstat(fid.Fid, dir)
	if err != nil {
		return err
	}

	_, err = clnt.Rpc(tc)
	return err
}
