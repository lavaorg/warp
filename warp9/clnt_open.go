// Copyright 2009 The Go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Copyright 2018 Larry Rau.

package warp9

import (
	"strings"

	"github.com/lavaorg/lrt/mlog"
)

// Creates and opens a named object.
// Returns the object if the operation is successful, or an Error.
func (clnt *Clnt) Create(path string, perm uint32, mode uint8) (*Object, W9Err) {
	n := strings.LastIndex(path, "/")
	if n < 0 {
		n = 0
	}

	fid, err := clnt.Walk(path[0:n])
	if err != Egood {
		return nil, err
	}

	if path[n] == '/' {
		n++
	}

	err = clnt.FCreate(fid, path[n:], perm, mode, "")
	if err != Egood {
		clnt.Clunk(fid)
		return nil, err
	}

	return &Object{fid, 0}, Egood
}

// Opens a named object. Returns the opened object, or an Error.
func (clnt *Clnt) Open(path string, mode uint8) (*Object, W9Err) {
	fid, err := clnt.Walk(path)
	if err != Egood {
		return nil, err
	}

	err = clnt.FOpen(fid, mode)
	if err != Egood {
		clnt.Clunk(fid)
		return nil, err
	}

	return &Object{fid, 0}, Egood
}

// Opens the object currently associated with the fid. Returns nil if
// the operation is successful.
func (clnt *Clnt) FOpen(fid *Fid, mode uint8) W9Err {
	tc := clnt.NewFcall()
	err := tc.packTopen(fid.Fid, mode)
	if err != Egood {
		return err
	}

	rc, err := clnt.Rpc(tc)
	if err != Egood {
		return err
	}

	fid.Qid = rc.Qid
	fid.Iounit = rc.Iounit
	if fid.Iounit == 0 || fid.Iounit > clnt.Msize-IOHDRSZ {
		fid.Iounit = clnt.Msize - IOHDRSZ
	}
	mlog.Debug("FOpen:msize:%v rc.Iounit:%v %v", clnt.Msize, rc.Iounit, fid)
	fid.Mode = mode
	return Egood
}

// like FOpen but returns an Object in the open state. The expectation is
// the fid was already the result of a walk.
func (clnt *Clnt) FOpenObject(fid *Fid, mode uint8) (*Object, W9Err) {
	if fid == nil {
		return nil, Efidnil
	}
	if !fid.walked {
		return nil, Ebaduse
	}

	err := clnt.FOpen(fid, mode)
	if err != Egood {
		clnt.Clunk(fid)
		return nil, err
	}

	return &Object{fid, 0}, Egood

}

// Creates an object in the directory associated with the fid. Returns nil
// if the operation is successful.
func (clnt *Clnt) FCreate(fid *Fid, name string, perm uint32, mode uint8, extattr string) W9Err {
	tc := clnt.NewFcall()
	err := tc.packTcreate(fid.Fid, name, perm, mode, extattr)
	if err != Egood {
		return err
	}

	rc, err := clnt.Rpc(tc)
	if err != Egood {
		return err
	}

	fid.Qid = rc.Qid
	fid.Iounit = rc.Iounit
	if fid.Iounit == 0 || fid.Iounit > clnt.Msize-IOHDRSZ {
		fid.Iounit = clnt.Msize - IOHDRSZ
	}
	fid.Mode = mode
	return Egood
}
