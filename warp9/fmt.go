// Copyright 2009 The Go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package warp9

import "fmt"

func PermToString(perm uint32) string {
	ret := ""

	if perm&DMMASK == 0 {
		ret = "o" //regular object
	} else {
		if perm&DMDIR != 0 {
			ret += "d"
		}

		if perm&DMAPPEND != 0 {
			ret += "a"
		}

		if perm&DMAUTH != 0 {
			ret += "A"
		}

		if perm&DMEXCL != 0 {
			ret += "l"
		}

		if perm&DMTMP != 0 {
			ret += "t"
		}
	}

	ret += fmt.Sprintf("%s", p2str(perm&0777))
	return ret
}

func p2str(a uint32) string {
	buf := []byte{'r', 'w', 'x', 'r', 'w', 'x', 'r', 'w', 'x'}

	for i, _ := range buf {
		if a&0400 == 0 {
			buf[i] = '-'
		}
		a = a << 1
	}
	return string(buf)
}

func (qid *Qid) String() string {
	b := ""
	if qid.Type == 0 {
		b = "o" // regular object
	} else {

		if qid.Type&QTDIR != 0 {
			b += "d"
		}
		if qid.Type&QTAPPEND != 0 {
			b += "a"
		}
		if qid.Type&QTAUTH != 0 {
			b += "A"
		}
		if qid.Type&QTEXCL != 0 {
			b += "l"
		}
		if qid.Type&QTTMP != 0 {
			b += "t"
		}
	}

	return fmt.Sprintf("(%s.%x.%x)", b, qid.Version, qid.Path)
}

func (d *Dir) String() string {
	ret := fmt.Sprintf("%s %s [%d:%d:%d] %d Q", PermToString(d.Mode), d.Name, d.Uid, d.Gid, d.Muid, d.Length)
	ret += d.Qid.String() + " "
	ret += fmt.Sprintf("at[%d] mt[%d]", d.Atime, d.Mtime)

	return ret
}

func (fc *Fcall) String() string {
	ret := ""

	switch fc.Type {
	default:
		ret = fmt.Sprintf("invalid call: %d", fc.Type)
	case Tversion:
		ret = fmt.Sprintf("Tversion tag %d msize %d version '%s'", fc.Tag, fc.Msize, fc.Version)
	case Rversion:
		ret = fmt.Sprintf("Rversion tag %d msize %d version '%s'", fc.Tag, fc.Msize, fc.Version)
	case Tauth:
		ret = fmt.Sprintf("Tauth tag %d afid %d uname '%s' aname '%s'",
			fc.Tag, fc.Atok, fc.Uid, fc.Aname)
	case Rauth:
		ret = fmt.Sprintf("Rauth tag %d aqid %v", fc.Tag, &fc.Qid)
	case Rattach:
		ret = fmt.Sprintf("Rattach tag %d aqid %v", fc.Tag, &fc.Qid)
	case Tattach:
		ret = fmt.Sprintf("Tattach tag %d fid %d atok %d uid '%d' aname '%s'",
			fc.Tag, fc.Fid, fc.Atok, fc.Uid, fc.Aname)
	case Tflush:
		ret = fmt.Sprintf("Tflush tag %d oldtag %d", fc.Tag, fc.Oldtag)
	case Rerror:
		ret = fmt.Sprintf("Rerror tag %d err %d '%s'", fc.Tag, fc.Error, fc.Error)
	case Twalk:
		ret = fmt.Sprintf("Twalk tag %d fid %d newfid %d [", fc.Tag, fc.Fid, fc.Newfid)
		for i := 0; i < len(fc.Wname); i++ {
			ret += fmt.Sprintf("'%s',", fc.Wname[i])
		}
		ret += "]"
	case Rwalk:
		ret = fmt.Sprintf("Rwalk tag %d wqid %v", fc.Tag, &fc.Qid)
	case Topen:
		ret = fmt.Sprintf("Topen tag %d fid %d mode %x", fc.Tag, fc.Fid, fc.Mode)
	case Ropen:
		ret = fmt.Sprintf("Ropen tag %d qid %v iounit %d", fc.Tag, &fc.Qid, fc.Iounit)
	case Rcreate:
		ret = fmt.Sprintf("Rcreate tag %d qid %v iounit %d", fc.Tag, &fc.Qid, fc.Iounit)
	case Tcreate:
		ret = fmt.Sprintf("Tcreate tag %d fid %d name '%s' perm ", fc.Tag, fc.Fid, fc.Name)
		ret += PermToString(fc.Perm)
		ret += fmt.Sprintf(" mode %x ", fc.Mode)
	case Tread:
		ret = fmt.Sprintf("Tread tag %d fid %d offset %d count %d", fc.Tag, fc.Fid, fc.Offset, fc.Count)
	case Rread:
		ret = fmt.Sprintf("Rread tag %d count %d", fc.Tag, fc.Count)
	case Twrite:
		ret = fmt.Sprintf("Twrite tag %d fid %d offset %d count %d", fc.Tag, fc.Fid, fc.Offset, fc.Count)
	case Rwrite:
		ret = fmt.Sprintf("Rwrite tag %d count %d", fc.Tag, fc.Count)
	case Tclunk:
		ret = fmt.Sprintf("Tclunk tag %d fid %d", fc.Tag, fc.Fid)
	case Rclunk:
		ret = fmt.Sprintf("Rclunk tag %d", fc.Tag)
	case Tremove:
		ret = fmt.Sprintf("Tremove tag %d fid %d", fc.Tag, fc.Fid)
	case Tstat:
		ret = fmt.Sprintf("Tstat tag %d fid %d", fc.Tag, fc.Fid)
	case Rstat:
		ret = fmt.Sprintf("Rstat tag %d st (%v)", fc.Tag, &fc.Dir)
	case Twstat:
		ret = fmt.Sprintf("Twstat tag %d fid %d st (%v)", fc.Tag, fc.Fid, &fc.Dir)
	case Rflush:
		ret = fmt.Sprintf("Rflush tag %d", fc.Tag)
	case Rremove:
		ret = fmt.Sprintf("Rremove tag %d", fc.Tag)
	case Rwstat:
		ret = fmt.Sprintf("Rwstat tag %d", fc.Tag)
	}

	return ret
}
