// Copyright 2009 The Go9p Authors.  All rights reserved.
// Portions: Copyright 2018 Larry Rau. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Fcall is the structure to hold contents for an on-the-wire message to/from a Warp9 server.
package warp9

// Fcall represents a Warp9 message. Not all fields are used in all messages.
type Fcall struct {
	FcSize  uint32   // size of the message
	Type    uint8    // message type
	Fid     uint32   // file identifier
	Tag     uint16   // message tag
	Msize   uint32   // maximum message size (used by Tversion, Rversion)
	Version string   // protocol version (used by Tversion, Rversion)
	Oldtag  uint16   // tag of the message to flush (used by Tflush)
	Error   W9Err    // error (used by Rerror)
	Qid              // file Qid (used by Rauth, Rattach, Ropen, Rcreate)
	Iounit  uint32   // maximum bytes read without breaking in multiple messages (used by Ropen, Rcreate)
	Atok    uint32   // authentication fid (used by Tauth, Tattach)
	Uid     uint32   // user uid (used by Tauth, Tattach)
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
	ExtAttr string   // used by Tcreate

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

func (fc *Fcall) packCommon(size int, id uint8) ([]byte, W9Err) {
	size += fixedFcsz /* size[4] id[1] tag[2] */
	if len(fc.Buf) < int(size) {
		return nil, Ebufsz
	}

	fc.FcSize = uint32(size)
	fc.Type = id
	fc.Tag = NOTAG
	p := fc.Buf
	p = pint32(uint32(size), p)
	p = pint8(id, p)
	p = pint16(NOTAG, p)
	fc.Pkt = fc.Buf[0:size]

	return p, Egood
}

// Creates a Fcall value from the on-the-wire representation.
// Returns the unpacked message,
// error and how many bytes from the buffer were used by the message.
func Unpack(buf []byte) (fc *Fcall, err W9Err, fcsz int) {
	var m uint16

	if len(buf) < fixedFcsz {
		return nil, Ebufsz, 0
	}

	fc = new(Fcall)
	fc.Fid = NOFID
	fc.Atok = NOTOK
	fc.Newfid = NOFID

	p := buf
	fc.FcSize, p = gint32(p)
	fc.Type, p = gint8(p)
	fc.Tag, p = gint16(p)

	if int(fc.FcSize) > len(buf) || fc.FcSize < fixedFcsz {
		return nil,
			Ebufsz,
			0
	}

	p = p[0 : fc.FcSize-7]
	fc.Pkt = buf[0:fc.FcSize]
	fcsz = int(fc.FcSize) //lr this seems to be truncating...
	if fc.Type < Tversion || fc.Type >= Tlast {
		return nil, Ebadmsgid, 0
	}

	var sz uint32
	sz = minFcsize[fc.Type-Tversion]

	if fc.FcSize < sz {
		goto szerror
	}

	err = Egood
	switch fc.Type {
	default:
		return nil, Ebadmsgid, 0

	case Tversion, Rversion:
		fc.Msize, p = gint32(p)
		fc.Version, p = gstr(p)
		if p == nil {
			goto szerror
		}

	case Tauth:
		fc.Atok, p = gint32(p)
		fc.Uid, p = gint32(p)
		if p == nil {
			goto szerror
		}

		fc.Aname, p = gstr(p)
		if p == nil {
			goto szerror
		}

	case Rauth, Rattach:
		p = gqid(p, &fc.Qid)

	case Tflush:
		fc.Oldtag, p = gint16(p)

	case Tattach:
		fc.Fid, p = gint32(p)
		fc.Atok, p = gint32(p)
		fc.Uid, p = gint32(p)
		if p == nil {
			goto szerror
		}

		fc.Aname, p = gstr(p)
		if p == nil {
			goto szerror
		}

	case Rerror:
		var ecode W9Err
		ecode, p = gerr(p) //gstr(p)
		fc.Error = ecode
		if p == nil {
			goto szerror
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
		p, err = gstat(p, &fc.Dir)
		if err != Egood {
			return nil, err, 0
		}

	case Twstat:
		fc.Fid, p = gint32(p)
		m, p = gint16(p)
		p, _ = gstat(p, &fc.Dir)

	case Rflush, Rclunk, Rremove, Rwstat:
	}

	if len(p) > 0 {
		goto szerror
	}

	return

szerror:
	return nil, Ebadmsgsz, 0
}
