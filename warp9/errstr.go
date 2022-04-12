// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 License file

package warp9

import (
	"fmt"
	"io"
)

// The Warp9 protocol expresses an error as a int16 value with negative values representing framework
// internal errors and positive values representing user level object server errors unique to the
// specific object server. Framework errors will have well known ID's that are associated with well
// known ascii strings. The protocol optionally allows for optional messages to be tranmissted as well
// please see the Rerror message type.
// This error type is compatible with 'error', but all framework calls will use this type.
// Object server implemenations generating errors must use this type.
// Framework errors <0, 0==no err, >0 obj-server errors.
type WarpError struct {
	errcode int16
	optmsg  string
}

// Create a new W9Err instance.
func ErrorCode(code int16) *WarpError {
	return &WarpError{code, ""}
}

// Create a new W9Err instance with an optional msg.
func ErrorMsg(code int16, msg string) *WarpError {
	return &WarpError{code, msg}
}

// Implement the Error interface and return a string representation.
// For framework errors the protocol errorcode is converted to string representation.
// if an optional message was transmitted its appended to the standard message.
func (e *WarpError) Error() string {
	if e.errcode < Egood && e.errcode > Emax {
		ec := e.errcode * -1
		if e.optmsg != "" {
			return fmt.Sprintf("%s:%s", ErrStr[ec], e.optmsg)
		}
		return ErrStr[ec]
	}
	return fmt.Sprintf("%d:%s", e.errcode, e.optmsg)
}

// Implement the standard convertion string interface. Call Error().
func (e *WarpError) String() string {
	return e.Error()
}

// End of Object Data. alias of io.EOF
var WarpErrorEOF = io.EOF

// Object not exist
var WarpErrorNOTEXIST = &WarpError{Enotexist, ""}

const (
	Egood int16 = iota * -1
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
)

var ErrStr [Emax * -1]string

func init() {
	ErrStr[Egood] = "no err"
	ErrStr[Ebadver*-1] = "bad version"
	ErrStr[Enotimpl*-1] = "not implemented"
	ErrStr[Emsize*-1] = "msize too small"
	ErrStr[Ebufsz*-1] = "buffer too small"
	ErrStr[Ebadmsgid*-1] = "warp9 msg: bad id"
	ErrStr[Ebadw9msg*-1] = "warp9 msg: unknown"
	ErrStr[Eio*-1] = "IO error"
	ErrStr[Ebadmsgsz*-1] = "bad message size"
	ErrStr[Einval*-1] = "unexpected response"
	ErrStr[Edial*-1] = "conn: dial failed"
	ErrStr[Econn*-1] = "conn: connection closed"
	ErrStr[Efidnil*-1] = "fid: nil"
	ErrStr[Einuse*-1] = "fid: in use"
	ErrStr[Ebaduse*-1] = "fid: bad use"
	ErrStr[Eopen*-1] = "fid: already opened"
	ErrStr[Efidnoaux*-1] = "fid: no aux"
	ErrStr[Eunknownfid*-1] = "fid: unknown"
	ErrStr[Enoent*-1] = "fid: no entry found in walk"
	ErrStr[Ename*-1] = "fid: illegal name"
	ErrStr[Eperm*-1] = "obj: permission denied"
	ErrStr[Enotopen*-1] = "obj: not open"
	ErrStr[Eeof*-1] = "obj: end of data"
	ErrStr[Enotdir*-1] = "obj: not a directory"
	ErrStr[Eisopen*-1] = "obj already open for I/O"
	ErrStr[Enotexist*-1] = "obj: does not exist"
	ErrStr[Eexist*-1] = "obj: exists"
	ErrStr[Excl*-1] = "obj: exclusive use object already open"
	ErrStr[Edirchange*-1] = "obj: cannot convert between obj and directories"
	ErrStr[Enotempty*-1] = "obj: directory not empty"
	ErrStr[Ebadoffset*-1] = "obj: bad offset in directory read"
	ErrStr[Etoolarge*-1] = "obj: i/o count too large"
	ErrStr[Enotowner*-1] = "user: not owner"
	ErrStr[Enouser*-1] = "user: unknown"
	ErrStr[Ebaduid*-1] = "user: bad u/g/m"
	ErrStr[Ebaduser*-1] = "user: bad"
	ErrStr[Enoauth*-1] = "auth: authentication not required"
	ErrStr[Eauthinit*-1] = "auth: init failed"
	ErrStr[Eauthchk*-1] = "auth: check failed"
	ErrStr[Eauthread*-1] = "auth: read failed"
	ErrStr[Eauthwrite*-1] = "auth: write failed"
}
