// Copyright 2009 The Go9p Authors.  All rights reserved.
// Portions: Copyright 2018 Larry Rau. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package warp9

// Converts a Dir value to its on-the-wire representation and writes it to
// the buf. Returns the number of bytes written, 0 if there is not enough space.
func PackDir(d *Dir) []byte {
	sz := statsz(d)
	buf := make([]byte, sz)
	pstat(d, buf)
	return buf
}

// Converts the on-the-wire representation of a stat to Stat value.
// Returns an error if the conversion is impossible, otherwise
// a pointer to a Stat value.
func UnpackDir(buf []byte) (d *Dir, b []byte, amt int, err error) {
	sz := 2 + 2 + 4 + 13 + 4 + /* size[2] type[2] dev[4] qid[13] mode[4] */
		4 + 4 + 8 + /* atime[4] mtime[4] length[8] */
		2 + 2 + 2 + 2 /* name[s] uid[s] gid[s] muid[s] */

	if len(buf) < sz {
		//log:?? fmt.Sprintf("short buffer: Need %d and have %v", sz, len(buf))
		return nil, nil, 0, &WarpError{Einval, ""}
	}
	d = new(Dir)
	b, err = gstat(buf, d)
	if err != nil {
		return nil, nil, 0, err
	}

	return d, b, len(buf) - len(b), nil

}

//
// private types and functions
//

// fixed msg preample of 7 bytes: msgsize[4] msg-type[1] tag[2]
const (
	fixedFcsz = 7
)

// minimum size of a Warp9 message for a type; not couting fixed preamble
var minFcsize = [...]uint32{
	6,  /* Tversion msize[4] version[s] */
	6,  /* Rversion msize[4] version[s] */
	10, /* Tauth fid[4] uid[4] aname[s] */
	13, /* Rauth aqid[13] */
	14, /* Tattach fid[4] atok[4] uid[4] aname[s] */
	13, /* Rattach qid[13] */
	0,  /* Terror */
	2,  /* Rerror errcode[2] */
	2,  /* Tflush oldtag[2] */
	0,  /* Rflush */
	10, /* Twalk fid[4] newfid[4] nwname[2]... */
	//	2,  /* Rwalk nwqid[2]... */
	13, /* Rwalk wqid[13] */
	5,  /* Topen fid[4] mode[1] */
	17, /* Ropen qid[13] iounit[4] */
	11, /* Tcreate fid[4] name[s] perm[4] mode[1] */
	17, /* Rcreate qid[13] iounit[4] */
	16, /* Tread fid[4] offset[8] count[4] */
	4,  /* Rread count[4]... */
	16, /* Twrite fid[4] offset[8] count[4]... */
	4,  /* Rwrite count[4] */
	4,  /* Tclunk fid[4] */
	0,  /* Rclunk */
	4,  /* Tremove fid[4] */
	0,  /* Rremove */
	4,  /* Tstat fid[4] */
	4,  /* Rstat stat[n] */
	8,  /* Twstat fid[4] stat[n] */
	0,  /* Rwstat */
	0,  /* Tget */
	0,  /* Rget */
	0,  /* Tput */
	0,  /* Rput */
	10, /* Treport atok[4] uid[4] aname[s] */
	2,  /* Rreport count[2]...*/
	13, /* Tstream fid[4] istream[1] offset[8] */
	4,  /* Rstream count[4]... */
}

func gint8(buf []byte) (uint8, []byte) { return buf[0], buf[1:] }

func gint16(buf []byte) (uint16, []byte) {
	return uint16(buf[0]) | (uint16(buf[1]) << 8), buf[2:]
}

func gint32(buf []byte) (uint32, []byte) {
	return uint32(buf[0]) | (uint32(buf[1]) << 8) | (uint32(buf[2]) << 16) |
			(uint32(buf[3]) << 24),
		buf[4:]
}

// signed-int32
func gsint32(buf []byte) (int32, []byte) {
	u32 := uint32(buf[0]) | (uint32(buf[1]) << 8) | (uint32(buf[2]) << 16) | (uint32(buf[3]) << 24)
	return int32(u32), buf[4:]
}

func gint64(buf []byte) (uint64, []byte) {
	return uint64(buf[0]) | (uint64(buf[1]) << 8) | (uint64(buf[2]) << 16) |
			(uint64(buf[3]) << 24) | (uint64(buf[4]) << 32) | (uint64(buf[5]) << 40) |
			(uint64(buf[6]) << 48) | (uint64(buf[7]) << 56),
		buf[8:]
}

func gerr(buf []byte) (*WarpError, []byte) {
	var e WarpError
	var b []byte
	e.errcode = (int16(buf[0]) | (int16(buf[1]) << 8))
	e.optmsg, b = gstr(buf[2:])

	return &e, b
}

func gstr(buf []byte) (string, []byte) {
	var n uint16

	if buf == nil || len(buf) < 2 {
		return "", nil
	}

	n, buf = gint16(buf)

	if int(n) > len(buf) {
		return "", nil
	}

	return string(buf[0:n]), buf[n:]
}

func gqid(buf []byte, qid *Qid) []byte {
	qid.Type, buf = gint8(buf)
	qid.Version, buf = gint32(buf)
	qid.Path, buf = gint64(buf)

	return buf
}

func gstat(buf []byte, d *Dir) ([]byte, error) {
	d.DirSize, buf = gint16(buf)
	buf = gqid(buf, &d.Qid)
	d.Mode, buf = gint32(buf)
	d.Atime, buf = gint32(buf)
	d.Mtime, buf = gint32(buf)
	d.Length, buf = gint64(buf)
	d.Name, buf = gstr(buf)
	if buf == nil {
		//s := fmt.Sprintf("Buffer too short for basic 9p: need %d, have %d", 49, sz)
		return nil, &WarpError{Ebufsz, ""}
	}

	d.Uid, buf = gint32(buf)
	if buf == nil {
		return nil, &WarpError{Ebaduid, "buf"}
	}
	d.Gid, buf = gint32(buf)
	if buf == nil {
		return nil, &WarpError{Ebaduid, "gid"}
	}

	d.Muid, buf = gint32(buf)
	if buf == nil {
		return nil, &WarpError{Ebaduid, "muid"}
	}

	return buf, nil
}

func pint8(val uint8, buf []byte) []byte {
	buf[0] = val
	return buf[1:]
}

func pint16(val uint16, buf []byte) []byte {
	buf[0] = uint8(val)
	buf[1] = uint8(val >> 8)
	return buf[2:]
}

func pint32(val uint32, buf []byte) []byte {
	buf[0] = uint8(val)
	buf[1] = uint8(val >> 8)
	buf[2] = uint8(val >> 16)
	buf[3] = uint8(val >> 24)
	return buf[4:]
}

func psint32(val int32, buf []byte) []byte {
	buf[0] = uint8(val)
	buf[1] = uint8(val >> 8)
	buf[2] = uint8(val >> 16)
	buf[3] = uint8(val >> 24)
	return buf[4:]
}

func pint64(val uint64, buf []byte) []byte {
	buf[0] = uint8(val)
	buf[1] = uint8(val >> 8)
	buf[2] = uint8(val >> 16)
	buf[3] = uint8(val >> 24)
	buf[4] = uint8(val >> 32)
	buf[5] = uint8(val >> 40)
	buf[6] = uint8(val >> 48)
	buf[7] = uint8(val >> 56)
	return buf[8:]
}

//TODO: if empty optmsg then don't send the size uint16 either
func perr(val *WarpError, buf []byte) []byte {
	buf[0] = uint8(val.errcode)
	buf[1] = uint8(val.errcode >> 8)
	buf = pstr(val.optmsg, buf[2:])
	return buf
}

func pstr(val string, buf []byte) []byte {
	n := uint16(len(val))
	buf = pint16(n, buf)
	b := []byte(val)
	copy(buf, b)
	return buf[n:]
}

func pqid(val *Qid, buf []byte) []byte {
	buf = pint8(val.Type, buf)
	buf = pint32(val.Version, buf)
	buf = pint64(val.Path, buf)

	return buf
}

func statsz(d *Dir) int {
	//sz := 2 + 2 + 4 + 13 + 4 + 4 + 4 + 8 + 2 + len(d.Name) + 4 + 4 + 4
	sz := 2 + 13 + 4 + 4 + 4 + 8 + 2 + len(d.Name) + 4 + 4 + 4
	return sz
}

func pstat(d *Dir, buf []byte) []byte {
	sz := statsz(d)
	buf = pint16(uint16(sz-2), buf)
	buf = pqid(&d.Qid, buf)
	buf = pint32(d.Mode, buf)
	buf = pint32(d.Atime, buf)
	buf = pint32(d.Mtime, buf)
	buf = pint64(d.Length, buf)
	buf = pstr(d.Name, buf)
	buf = pint32(d.Uid, buf)
	buf = pint32(d.Gid, buf)
	buf = pint32(d.Muid, buf)
	return buf
}
