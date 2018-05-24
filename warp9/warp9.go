// Copyright 2009 The Go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE object.

package warp9

// definitions for the Warp9 protocol

// Version
const Warp9Version = "Warp9.0"

// Warp9 message types
const (
	Tversion = 100 + iota
	Rversion
	Tauth
	Rauth
	Tattach
	Rattach
	Terror
	Rerror
	Tflush
	Rflush
	Twalk
	Rwalk
	Topen
	Ropen
	Tcreate
	Rcreate
	Tread
	Rread
	Twrite
	Rwrite
	Tclunk
	Rclunk
	Tremove
	Rremove
	Tstat
	Rstat
	Twstat
	Rwstat
	Tget
	Rget
	Tput
	Rput
	Treport
	Rreport
	Tstream
	Rstream
	Tlast
)

const (
	MSIZE   = 8192 + IOHDRSZ // default message size (data-size+IOHdrSz)
	IOHDRSZ = 24             // the non-data size of the Twrite messages
	PORT    = 9090           // default port for Warp9 object servers
)

// Qid types
const (
	QTDIR    = 0x80 // directories
	QTAPPEND = 0x40 // append only objects
	QTEXCL   = 0x20 // exclusive use objects
	QTMOUNT  = 0x10 // mounted channel
	QTAUTH   = 0x08 // authentication object
	QTTMP    = 0x04 // non-backed-up object
	QTOBJ    = 0x00
)

// Flags for the mode field in Topen and Tcreate messages
const (
	OREAD   = 0  // open read-only
	OWRITE  = 1  // open write-only
	ORDWR   = 2  // open read-write
	OUSE    = 3  // USE; directories can be used for walk without Read perm
	OTRUNC  = 16 // or'ed in truncate object first
	ORCLOSE = 64 // or'ed in, remove on close
)

// Object modes
const (
	// object types -- high order 8bits
	DMMASK   = 0xFF000000 // masks for top 8bits
	DMDIR    = 0x80000000 // mode bit for directories
	DMAPPEND = 0x40000000 // mode bit for append only objects, offset ignored in write
	DMEXCL   = 0x20000000 // mode bit for exclusive use objects; one client open at a time
	DMAUTH   = 0x08000000 // mode bit for authentication object
	DMTMP    = 0x04000000 // mode bit for non-backed-up object

	// permissions
	DMREAD  = 0x4 // mode bit for read permission
	DMWRITE = 0x2 // mode bit for write permission
	DMUSE   = 0x1 // mode bit for execute permission

)

const (
	NOTAG uint16 = 0xFFFF     // no tag specified
	NOFID uint32 = 0xFFFFFFFF // no fid specified
	NOUID uint32 = 0xFFFFFFFF // no uid specified
	NOTOK uint32 = 0xFFFFFFFF // no auth token
)

// Server's descriptor of an Object
type Qid struct {
	Type    uint8  // type of the object (high 8 bits of the mode)
	Version uint32 // version number for the path
	Path    uint64 // server's unique identification of the object
}

// Dir describes a directory object
type Dir struct {
	DirSize uint16 // size-2 of the Dir on the wire
	//Type    uint16
	//Dev     uint32
	Qid            // object's Qid
	Mode    uint32 // permissions and flags
	Atime   uint32 // last access time in seconds
	Mtime   uint32 // last modified time in seconds
	Length  uint64 // object length in bytes
	Name    string // object name
	Uid     uint32 // owner id
	Gid     uint32 // group id
	Muid    uint32 // id of the last user that modified the object
	ExtAttr string // extended attributes
}

// stats callbacks
type StatsOps interface {
	statsRegister()
	statsUnregister()
}

// Debug flags
const (
	DbgPrintFcalls  = (1 << iota) // print all 9P messages on stderr
	DbgPrintPackets               // print the raw packets on stderr
	DbgLogFcalls                  // keep the last N 9P messages (can be accessed over http)
	DbgLogPackets                 // keep the last N 9P messages (can be accessed over http)
)
