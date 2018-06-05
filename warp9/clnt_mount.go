// Copyright 2009 The Go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package warp9

import (
	"net"
)

// Creates an authentication fid for the specified user. Returns the fid, if
// successful, or an Error.
func (clnt *Clnt) Auth(user User, aname string) (*Fid, error) {
	fid := clnt.FidAlloc()
	tc := clnt.NewFcall()
	err := tc.packTauth(fid.Fid, user.Id(), aname)
	if err != nil {
		return nil, clnt.Perr(err)
	}

	_, err = clnt.Rpc(tc)
	if err != nil {
		return nil, clnt.Perr(err)
	}

	fid.User = user
	fid.walked = true
	return fid, nil
}

// Creates a fid for the specified user that points to the root
// of the file server's file tree. Returns a Fid pointing to the root,
// if successful, or an Error.
func (clnt *Clnt) Attach(afid *Fid, user User, aname string) (*Fid, error) {
	var afno uint32

	if afid != nil {
		afno = afid.Fid
	} else {
		afno = NOFID
	}

	fid := clnt.FidAlloc()
	tc := clnt.NewFcall()
	err := tc.packTattach(fid.Fid, afno, user.Id(), aname)
	if err != nil {
		return nil, clnt.Perr(err)
	}

	rc, err := clnt.Rpc(tc)
	if err != nil {
		return nil, clnt.Perr(err)
	}

	fid.Qid = rc.Qid
	fid.User = user
	fid.walked = true
	return fid, nil
}

// Connects to a file server and attaches to it as the specified user.
func Mount(ntype, addr, aname string, msize uint32, user User) (*Clnt, error) {
	c, e := net.Dial(ntype, addr)
	if e != nil {
		return nil, &WarpError{Edial, ""}
	}

	return MountConn(c, aname, msize, user)
}

func MountConn(c net.Conn, aname string, msize uint32, user User) (*Clnt, error) {
	clnt, err := Connect(c, msize+IOHDRSZ)
	if err != nil {
		return nil, err
	}

	fid, err := clnt.Attach(nil, user, aname)
	if err != nil {
		clnt.Unmount()
		return nil, err
	}

	clnt.Root = fid
	return clnt, nil
}

// Closes the connection to the file sever.
func (clnt *Clnt) Unmount() {
	clnt.Lock()
	clnt.err = &WarpError{Econn, ""}
	clnt.conn.Close()
	clnt.Unlock()
}
