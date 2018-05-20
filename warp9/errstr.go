// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 License file

package warp9

import (
	"github.com/lavaorg/lrt/mlog"
)

// Errors are int16 all warp errors <0, 0==no err, >0 obj-server errors
type W9Err int16

func (e W9Err) Error() string {
	if e < Egood && e > Emax {
		return ErrStr[e*-1]
	}
	mlog.Debug("W9Err:%d", e)
	return "w9 error"
}

func (e W9Err) String() string {
	return e.Error()
}

const (
	Egood W9Err = iota * -1
	Ebadver
	Eperm
	Enotdir
	Enoauth
	Enotexist
	Einuse
	Eexist
	Enotowner
	Eisopen
	Excl
	Ename
	Ebadw9msg
	Eunknownfid
	Ebaduse
	Eopen
	Etoolarge
	Ebadoffset
	Edirchange
	Enouser
	Enotimpl
	Enotempty
	Enoent
	Enotopen
	Ebaduser
	Emsize
	Ebufsz
	Ebadmsgid
	Ebaduid
	Ebadmsgsz
	Eio
	Einval
	Edial
	Econn
	Efidnil
	Eeof
	Eauthinit
	Eauthchk
	Eauthread
	Eauthwrite
	Efidnoaux

	Emax

	ENOENT  = -100
	EIO     = -101
	EEXIST  = -102
	ENOTDIR = -103
	EINVAL  = -104
)

var ErrStr []string = []string{
	/*Egood      */ "no err",
	/*Ebadver    */ "bad version",
	/*Eperm      */ "obj: permission denied",
	/*Enotdir    */ "obj: not a directory",
	/*Enoauth    */ "upas/fs: authentication not required",
	/*Enotexist  */ "obj: does not exist",
	/*Einuse     */ "fid: in use",
	/*Eexist     */ "obj: exists",
	/*Enotowner  */ "user: not owner",
	/*Eisopen    */ "obj already open for I/O",
	/*Excl       */ "exclusive use object already open",
	/*Ename      */ "illegal name",
	/*Ebadw9msg  */ "warp9 msg: unknown",
	/*Eunknownfid*/ "fid: unknown",
	/*Ebaduse    */ "fid: bad use",
	/*Eopen      */ "fid: already opened",
	/*Etoolarge  */ "i/o count too large",
	/*Ebadoffset */ "bad offset in directory read",
	/*Edirchange */ "cannot convert between obj and directories",
	/*Enouser    */ "user: unknown",
	/*Enotimpl   */ "not implemented",
	/*Enotempty  */ "directory not empty",
	/*Enoent     */ "no entry found in walk",
	/*Enotopen   */ "obj not open",
	/*Ebaduser   */ "user: bad",
	/*Emsize     */ "msize too small",
	/*Ebufsz     */ "buffer too small",
	/*Ebadmsgid  */ "warp9 msg: bad id",
	/*Ebaduid    */ "id: bad u/g/m",
	/*Ebadmsgsz  */ "bad message size",
	/*Eio        */ "IO error",
	/*Einval     */ "unexpected response",
	/*Edial      */ "dial failed",
	/*Econn      */ "connection closed",
	/*Efidnil    */ "fid: nil",
	/*Eeof       */ "end of obj",
	/*Eauthinit  */ "authentication init failed",
	/*Eauthchk   */ "authentication check failed",
	/*Eauthread  */ "authentication read failed",
	/*Eauthwrite */ "authentication write failed",
	/*Efidnoaux  */ "fid: no aux",
}
