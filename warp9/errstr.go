// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LIC/*ENS/*E

package warp9

// Errors are int16 all warp errors <0, 0==no err, >0 obj-server errors
type W9Err int16

func (e W9Err) Error() string {
	if e < Egood && e > Emax {
		return ErrStr[e*-1]
	}
	return "w9 error"
}

func (e W9Err) String() string {
	return e.Error()
}

const (
	Egood = iota * -1
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
	/*Eperm      */ "permission denied",
	/*Enotdir    */ "not a directory",
	/*Enoauth    */ "upas/fs: authentication not required",
	/*Enotexist  */ "object does not exist",
	/*Einuse     */ "object in use",
	/*Eexist     */ "object exists",
	/*Enotowner  */ "not owner",
	/*Eisopen    */ "file already open for I/O",
	/*Excl       */ "exclusive use object already open",
	/*Ename      */ "illegal name",
	/*Ebadw9msg  */ "unknown warp9 message",
	/*Eunknownfid*/ "unknown fid",
	/*Ebaduse    */ "bad use of fid",
	/*Eopen      */ "fid already opened",
	/*Etoolarge  */ "i/o count too large",
	/*Ebadoffset */ "bad offset in directory read",
	/*Edirchange */ "cannot convert between files and directories",
	/*Enouser    */ "unknown user",
	/*Enotimpl   */ "not implemented",
	/*Enotempty  */ "directory not empty",
	/*Enoent     */ "no entry found in walk",
	/*Enotopen   */ "file not open",
	/*Ebaduser   */ "bad user",
	/*Emsize     */ "msize too small",
	/*Ebufsz     */ "buffer too small",
	/*Ebadmsgid  */ "bad warp9 msg id",
	/*Ebaduid    */ "bad u/g/m id",
	/*Ebadmsgsz  */ "bad message size",
	/*Eio        */ "IO error",
	/*Einval     */ "unexpected response",
	/*Edial      */ "dial failed",
	/*Econn      */ "connection closed",
	/*Efidnil    */ "nil fid",
	/*Eeof       */ "end of file",
	/*Eauthinit  */ "authentication init failed",
	/*Eauthchk   */ "authentication check failed",
	/*Eauthread  */ "authentication read failed",
	/*Eauthwrite */ "authentication write failed",
}
