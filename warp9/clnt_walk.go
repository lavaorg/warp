// Copyright 2009 The Go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package warp9

import (
	"strings"

	"github.com/lavaorg/lrt/mlog"
)

// Starting from the object associated with fid, walks all wnames in
// sequence and associates the resulting object with newfid. If no wnames
// were walked successfully, an Error is returned. Otherwise a Qid for the
// object represented by the last wname is returned.
func (clnt *Clnt) FWalk(fid *Fid, newfid *Fid, wnames []string) (*Qid, error) {
	tc := clnt.NewFcall()
	err := tc.packTwalk(fid.Fid, newfid.Fid, wnames)
	if err != nil {
		return nil, clnt.Perr(err)
	}

	rc, err := clnt.Rpc(tc)
	if err != nil {
		return nil, clnt.Perr(err)
	}

	newfid.walked = true
	return &rc.Qid, nil
}

// Walks to a named object. Returns a Fid associated with the object,
// or an Error.
func (clnt *Clnt) Walk(path string) (*Fid, error) {
	var err error = nil

	var i, m int
	for i = 0; i < len(path); i++ {
		if path[i] != '/' {
			break
		}
	}

	if i > 0 {
		path = path[i:]
	}

	wnames := strings.Split(path, "/")
	newfid := clnt.FidAlloc()
	fid := clnt.Root
	newfid.User = fid.User

	// get rid of the empty names
	for i, m = 0, 0; i < len(wnames); i++ {
		if wnames[i] != "" {
			wnames[m] = wnames[i]
			m++
		}
	}

	wnames = wnames[0:m]
	for {
		n := len(wnames)
		if n > 16 {
			n = 16
		}

		tc := clnt.NewFcall()
		err = tc.packTwalk(fid.Fid, newfid.Fid, wnames[0:n])
		if err != nil {
			mlog.Debug("err=%T,%v,%p,%t,%t", err, err, err, (err != nil), (err == nil))
			goto error
		}

		var rc *Fcall
		rc, err = clnt.Rpc(tc)
		if err != nil {
			mlog.Debug("(rpc)err=%T,%v", err, err)
			goto error
		}

		newfid.walked = true
		newfid.Qid = rc.Qid

		wnames = wnames[n:]
		fid = newfid
		if len(wnames) == 0 {
			break
		}
	}

	return newfid, nil

error:
	clnt.Clunk(newfid)
	return nil, err
}
