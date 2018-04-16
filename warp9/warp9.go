// Copyright 2009 The Go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE object.

// The definitions and functions used to implement the Warp9 protocol
// TODO.
// All the packet conversion code in this file is crap and needs a rewrite.
package warp9

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
	MSIZE   = 1048576 + IOHDRSZ // default message size (1048576+IOHdrSz)
	IOHDRSZ = 24                // the non-data size of the Twrite messages
	PORT    = 564               // default port for Warp9 object servers
)

// Qid types
const (
	QTDIR    = 0x80 // directories
	QTAPPEND = 0x40 // append only objects
	QTEXCL   = 0x20 // exclusive use objects
	QTMOUNT  = 0x10 // mounted channel
	QTAUTH   = 0x08 // authentication object
	QTTMP    = 0x04 // non-backed-up object
	QTFILE   = 0x00
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

// File modes
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

// File identifier
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
	2,  /* Rwalk nwqid[2]... */
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

func gerr(buf []byte) (W9Err, []byte) {
	return W9Err(uint16(buf[0]) | (uint16(buf[1]) << 8)), buf[2:]
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

func gstat(buf []byte, d *Dir) ([]byte, W9Err) {
	d.DirSize, buf = gint16(buf)
	buf = gqid(buf, &d.Qid)
	d.Mode, buf = gint32(buf)
	d.Atime, buf = gint32(buf)
	d.Mtime, buf = gint32(buf)
	d.Length, buf = gint64(buf)
	d.Name, buf = gstr(buf)
	if buf == nil {
		//s := fmt.Sprintf("Buffer too short for basic 9p: need %d, have %d", 49, sz)
		return nil, Ebufsz
	}

	d.Uid, buf = gint32(buf)
	if buf == nil {
		return nil, Ebaduid
	}
	d.Gid, buf = gint32(buf)
	if buf == nil {
		return nil, Ebaduid
	}

	d.Muid, buf = gint32(buf)
	if buf == nil {
		return nil, Ebaduid
	}

	return buf, Egood
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

func perr(val W9Err, buf []byte) []byte {
	buf[0] = uint8(val)
	buf[1] = uint8(val >> 8)
	return buf[2:]
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
func UnpackDir(buf []byte) (d *Dir, b []byte, amt int, err W9Err) {
	sz := 2 + 2 + 4 + 13 + 4 + /* size[2] type[2] dev[4] qid[13] mode[4] */
		4 + 4 + 8 + /* atime[4] mtime[4] length[8] */
		2 + 2 + 2 + 2 /* name[s] uid[s] gid[s] muid[s] */

	if len(buf) < sz {
		//log:?? fmt.Sprintf("short buffer: Need %d and have %v", sz, len(buf))
		return nil, nil, 0, Einval
	}
	d = new(Dir)
	b, err = gstat(buf, d)
	if err != Egood {
		return nil, nil, 0, err
	}

	return d, b, len(buf) - len(b), Egood

}
