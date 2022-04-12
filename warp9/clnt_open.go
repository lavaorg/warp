// Copyright 2009 The Go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Copyright 2018,2022 Larry Rau.

package warp9

import (
	"strings"
)

// Create creates and opens a named object.
// Returns the object if the operation is successful, or an Error.
func (clnt *Clnt) Create(path string, perm uint32, mode uint8) (*Object, error) {
	n := strings.LastIndex(path, "/")
	if n < 0 {
		n = 0
	}

	fid, err := clnt.Walk(path[0:n])
	if err != nil {
		return nil, err
	}

	if path[n] == '/' {
		n++
	}

	err = clnt.FCreate(fid, path[n:], perm, mode, "")
	if err != nil {
		clnt.Clunk(fid)
		return nil, err
	}

	return &Object{fid, 0}, nil
}

// Open opens a named object. Returns the opened object, or an Error.
func (clnt *Clnt) Open(path string, mode uint8) (*Object, error) {
	fid, err := clnt.Walk(path)
	if err != nil {
		return nil, err
	}

	err = clnt.FOpen(fid, mode)
	if err != nil {
		clnt.Clunk(fid)
		return nil, err
	}

	return &Object{fid, 0}, nil
}

// FOpen opens the object currently associated with the fid. Returns nil if
// the operation is successful.
func (clnt *Clnt) FOpen(fid *Fid, mode uint8) error {
	tc := clnt.NewFcall()
	err := tc.packTopen(fid.Fid, mode)
	if err != nil {
		return clnt.Perr(err)
	}

	rc, err := clnt.Rpc(tc)
	if err != nil {
		werr, ok := err.(*WarpError)
		if ok && werr.errcode == Enotexist {
			return clnt.Perr(WarpErrorNOTEXIST)
		}
		return clnt.Perr(err)
	}

	fid.Qid = rc.Qid
	fid.Iounit = rc.Iounit
	if fid.Iounit == 0 || fid.Iounit > clnt.Msize-IOHDRSZ {
		fid.Iounit = clnt.Msize - IOHDRSZ
	}
	fid.Mode = mode
	return nil
}

// FOpenObject like FOpen but returns an Object in the open state. The expectation is
// the fid was already the result of a walk.
func (clnt *Clnt) FOpenObject(fid *Fid, mode uint8) (*Object, error) {
	if fid == nil {
		return nil, &WarpError{Efidnil, ""}
	}
	if !fid.walked {
		return nil, &WarpError{Ebaduse, ""}
	}

	err := clnt.FOpen(fid, mode)
	if err != nil {
		clnt.Clunk(fid)
		return nil, err
	}

	return &Object{fid, 0}, nil

}

// FCreate creates an object in the directory associated with the fid. Returns nil
// if the operation is successful.
func (clnt *Clnt) FCreate(fid *Fid, name string, perm uint32, mode uint8, extattr string) error {
	tc := clnt.NewFcall()
	err := tc.packTcreate(fid.Fid, name, perm, mode, extattr)
	if err != nil {
		return clnt.Perr(err)
	}

	rc, err := clnt.Rpc(tc)
	if err != nil {
		return clnt.Perr(err)
	}

	fid.Qid = rc.Qid
	fid.Iounit = rc.Iounit
	if fid.Iounit == 0 || fid.Iounit > clnt.Msize-IOHDRSZ {
		fid.Iounit = clnt.Msize - IOHDRSZ
	}
	fid.Mode = mode
	return nil
}
