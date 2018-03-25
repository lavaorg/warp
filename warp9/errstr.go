// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package warp9

const (
	Eperm       = "permission denied"
	Enotdir     = "not a directory"
	Enoauth     = "upas/fs: authentication not required"
	Enotexist   = "file does not exist"
	Einuse      = "file in use"
	Eexist      = "file exists"
	Enotowner   = "not owner"
	Eisopen     = "file already open for I/O"
	Excl        = "exclusive use file already open"
	Ename       = "illegal name"
	Ebadctl     = "unknown control message"
	Eunknownfid = "unknown fid"
	Ebaduse     = "bad use of fid"
	Eopen       = "fid already opened"
	Etoolarge   = "i/o count too large"
	Ebadoffset  = "bad offset in directory read"
	Edirchange  = "cannot convert between files and directories"
	Enouser     = "unknown user"
	Enotimpl    = "not implemented"
	Enotempty   = "directory not empty"
	Enoent      = "no entry found in walk"
	Enotopen    = "file not open"
)
