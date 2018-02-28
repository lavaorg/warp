// Copyright 2009 The Go9p Authors.  All rights reserved.
// Portions: Copyright 2018 Larry Rau. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Fcall is the structure to hold contents for an on-the-wire message to/from a 9p server.
package warp9

import (
	"fmt"
)

// Fcall represents a 9P2000 message
type Fcall struct {
	FcSize  uint32   // size of the message
	Type    uint8    // message type
	Fid     uint32   // file identifier
	Tag     uint16   // message tag
	Msize   uint32   // maximum message size (used by Tversion, Rversion)
	Version string   // protocol version (used by Tversion, Rversion)
	Oldtag  uint16   // tag of the message to flush (used by Tflush)
	Error   string   // error (used by Rerror)
	Qid              // file Qid (used by Rauth, Rattach, Ropen, Rcreate)
	Iounit  uint32   // maximum bytes read without breaking in multiple messages (used by Ropen, Rcreate)
	Afid    uint32   // authentication fid (used by Tauth, Tattach)
	Uname   string   // user name (used by Tauth, Tattach)
	Aname   string   // attach name (used by Tauth, Tattach)
	Perm    uint32   // file permission (mode) (used by Tcreate)
	Name    string   // file name (used by Tcreate)
	Mode    uint8    // open mode (used by Topen, Tcreate)
	Newfid  uint32   // the fid that represents the file walked to (used by Twalk)
	Wname   []string // list of names to walk (used by Twalk)
	Wqid    []Qid    // list of Qids for the walked files (used by Rwalk)
	Offset  uint64   // offset in the file to read/write from/to (used by Tread, Twrite)
	Count   uint32   // number of bytes read/written (used by Tread, Rread, Twrite, Rwrite)
	Data    []uint8  // data read/to-write (used by Rread, Twrite)
	Dir              // file description (used by Rstat, Twstat)

	/* 9P2000.u extensions */
	Errornum uint32 // error code, 9P2000.u only (used by Rerror)
	Ext      string // special file description, 9P2000.u only (used by Tcreate)
	Unamenum uint32 // user ID, 9P2000.u only (used by Tauth, Tattach)

	Pkt []uint8 // raw packet data
	Buf []uint8 // buffer to put the raw data in
}

// Allocates a new Fcall.
func NewFcall(sz uint32) *Fcall {
	fc := new(Fcall)
	fc.Buf = make([]byte, sz)

	return fc
}

// Sets the tag of a Fcall.
func (fc *Fcall) SetTag(tag uint16) {
	fc.Tag = tag
	pint16(tag, fc.Pkt[5:])
}

func (fc *Fcall) packCommon(size int, id uint8) ([]byte, error) {
	size += 4 + 1 + 2 /* size[4] id[1] tag[2] */
	if len(fc.Buf) < int(size) {
		return nil, &Error{"buffer too small", EINVAL}
	}

	fc.FcSize = uint32(size)
	fc.Type = id
	fc.Tag = NOTAG
	p := fc.Buf
	p = pint32(uint32(size), p)
	p = pint8(id, p)
	p = pint16(NOTAG, p)
	fc.Pkt = fc.Buf[0:size]

	return p, nil
}

// Creates a Fcall value from the on-the-wire representation. If
// dotu is true, reads 9P2000.u messages. Returns the unpacked message,
// error and how many bytes from the buffer were used by the message.
func Unpack(buf []byte, dotu bool) (fc *Fcall, err error, fcsz int) {
	var m uint16

	if len(buf) < 7 {
		return nil, &Error{"buffer too short", EINVAL}, 0
	}

	fc = new(Fcall)
	fc.Fid = NOFID
	fc.Afid = NOFID
	fc.Newfid = NOFID

	p := buf
	fc.FcSize, p = gint32(p)
	fc.Type, p = gint8(p)
	fc.Tag, p = gint16(p)

	if int(fc.FcSize) > len(buf) || fc.FcSize < 7 {
		return nil,
			&Error{fmt.Sprintf("buffer too short: %d expected %d", len(buf), fc.FcSize), EINVAL},
			0
	}

	p = p[0 : fc.FcSize-7]
	fc.Pkt = buf[0:fc.FcSize]
	fcsz = int(fc.FcSize)
	if fc.Type < Tversion || fc.Type >= Tlast {
		return nil, &Error{"invalid id", EINVAL}, 0
	}

	var sz uint32
	if dotu {
		sz = minFcsize[fc.Type-Tversion]
	} else {
		sz = minFcusize[fc.Type-Tversion]
	}

	if fc.FcSize < sz {
		goto szerror
	}

	err = nil
	switch fc.Type {
	default:
		return nil, &Error{"invalid message id", EINVAL}, 0

	case Tversion, Rversion:
		fc.Msize, p = gint32(p)
		fc.Version, p = gstr(p)
		if p == nil {
			goto szerror
		}

	case Tauth:
		fc.Afid, p = gint32(p)
		fc.Uname, p = gstr(p)
		if p == nil {
			goto szerror
		}

		fc.Aname, p = gstr(p)
		if p == nil {
			goto szerror
		}

		if dotu {
			if len(p) > 0 {
				fc.Unamenum, p = gint32(p)
			} else {
				fc.Unamenum = NOUID
			}
		} else {
			fc.Unamenum = NOUID
		}

	case Rauth, Rattach:
		p = gqid(p, &fc.Qid)

	case Tflush:
		fc.Oldtag, p = gint16(p)

	case Tattach:
		fc.Fid, p = gint32(p)
		fc.Afid, p = gint32(p)
		fc.Uname, p = gstr(p)
		if p == nil {
			goto szerror
		}

		fc.Aname, p = gstr(p)
		if p == nil {
			goto szerror
		}

		if dotu {
			if len(p) > 0 {
				fc.Unamenum, p = gint32(p)
			} else {
				fc.Unamenum = NOUID
			}
		}

	case Rerror:
		fc.Error, p = gstr(p)
		if p == nil {
			goto szerror
		}
		if dotu {
			fc.Errornum, p = gint32(p)
		} else {
			fc.Errornum = 0
		}

	case Twalk:
		fc.Fid, p = gint32(p)
		fc.Newfid, p = gint32(p)
		m, p = gint16(p)
		fc.Wname = make([]string, m)
		for i := 0; i < int(m); i++ {
			fc.Wname[i], p = gstr(p)
			if p == nil {
				goto szerror
			}
		}

	case Rwalk:
		m, p = gint16(p)
		fc.Wqid = make([]Qid, m)
		for i := 0; i < int(m); i++ {
			p = gqid(p, &fc.Wqid[i])
		}

	case Topen:
		fc.Fid, p = gint32(p)
		fc.Mode, p = gint8(p)

	case Ropen, Rcreate:
		p = gqid(p, &fc.Qid)
		fc.Iounit, p = gint32(p)

	case Tcreate:
		fc.Fid, p = gint32(p)
		fc.Name, p = gstr(p)
		if p == nil {
			goto szerror
		}
		fc.Perm, p = gint32(p)
		fc.Mode, p = gint8(p)
		if dotu {
			fc.Ext, p = gstr(p)
			if p == nil {
				goto szerror
			}
		}

	case Tread:
		fc.Fid, p = gint32(p)
		fc.Offset, p = gint64(p)
		fc.Count, p = gint32(p)

	case Rread:
		fc.Count, p = gint32(p)
		if len(p) < int(fc.Count) {
			goto szerror
		}
		fc.Data = p
		p = p[fc.Count:]

	case Twrite:
		fc.Fid, p = gint32(p)
		fc.Offset, p = gint64(p)
		fc.Count, p = gint32(p)
		if len(p) != int(fc.Count) {
			fc.Data = make([]byte, fc.Count)
			copy(fc.Data, p)
			p = p[len(p):]
		} else {
			fc.Data = p
			p = p[fc.Count:]
		}

	case Rwrite:
		fc.Count, p = gint32(p)

	case Tclunk, Tremove, Tstat:
		fc.Fid, p = gint32(p)

	case Rstat:
		m, p = gint16(p)
		p, err = gstat(p, &fc.Dir, dotu)
		if err != nil {
			return nil, err, 0
		}

	case Twstat:
		fc.Fid, p = gint32(p)
		m, p = gint16(p)
		p, _ = gstat(p, &fc.Dir, dotu)

	case Rflush, Rclunk, Rremove, Rwstat:
	}

	if len(p) > 0 {
		goto szerror
	}

	return

szerror:
	return nil, &Error{"invalid size", EINVAL}, 0
}
