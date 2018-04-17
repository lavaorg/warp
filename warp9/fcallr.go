// Copyright 2009 The Go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package warp9

// Create a Rversion message in the specified Fcall.
func (fc *Fcall) packRversion(msize uint32, version string) W9Err {
	size := 4 + 2 + len(version) /* msize[4] version[s] */
	p, err := fc.packCommon(size, Rversion)
	if err != Egood {
		return err
	}

	fc.Msize = msize
	fc.Version = version
	p = pint32(msize, p)
	p = pstr(version, p)

	return Egood
}

// Create a Rauth message in the specified Fcall.
func (fc *Fcall) packRauth(aqid *Qid) W9Err {
	size := 13 /* aqid[13] */
	p, err := fc.packCommon(size, Rauth)
	if err != Egood {
		return err
	}

	fc.Qid = *aqid
	p = pqid(aqid, p)
	return Egood
}

// Create a Rerror message in the specified Fcall.
func (fc *Fcall) packRerror(err W9Err) W9Err {

	size := 2 /* error[2] */

	p, locerr := fc.packCommon(size, Rerror)
	if locerr != Egood {
		return locerr
	}

	fc.Error = err
	p = perr(err, p)
	return Egood
}

// Create a Rflush message in the specified Fcall.
func (fc *Fcall) packRflush() W9Err {
	_, err := fc.packCommon(0, Rflush)

	return err
}

// Create a Rattach message in the specified Fcall.
func (fc *Fcall) packRattach(aqid *Qid) W9Err {
	size := 13 /* aqid[13] */
	p, err := fc.packCommon(size, Rattach)
	if err != Egood {
		return err
	}

	fc.Qid = *aqid
	p = pqid(aqid, p)
	return Egood
}

// Create a Rwalk message in the specified Fcall.
func (fc *Fcall) packRwalk(wqid *Qid) W9Err {
	size := 13 //wqid
	p, err := fc.packCommon(size, Rwalk)
	if err != Egood {
		return err
	}

	fc.Qid = *wqid
	p = pqid(wqid, p)
	return Egood
}

/* -- replaced
func (fc *Fcall) packRwalk(wqids []Qid) W9Err {
	nwqid := len(wqids)
	size := 2 + nwqid*13 // nwqid[2] nwname*wqid[13]
	p, err := fc.packCommon(size, Rwalk)
	if err != Egood {
		return err
	}

	p = pint16(uint16(nwqid), p)
	fc.Wqid = make([]Qid, nwqid)
	for i := 0; i < nwqid; i++ {
		fc.Wqid[i] = wqids[i]
		p = pqid(&wqids[i], p)
	}

	return Egood
}
*/
// Create a Ropen message in the specified Fcall.
func (fc *Fcall) packRopen(qid *Qid, iounit uint32) W9Err {
	size := 13 + 4 /* qid[13] iounit[4] */
	p, err := fc.packCommon(size, Ropen)
	if err != Egood {
		return err
	}

	fc.Qid = *qid
	fc.Iounit = iounit
	p = pqid(qid, p)
	p = pint32(iounit, p)
	return Egood
}

// Create a Rcreate message in the specified Fcall.
func (fc *Fcall) packRcreate(qid *Qid, iounit uint32) W9Err {
	size := 13 + 4 /* qid[13] iounit[4] */
	p, err := fc.packCommon(size, Rcreate)
	if err != Egood {
		return err
	}

	fc.Qid = *qid
	fc.Iounit = iounit
	p = pqid(qid, p)
	p = pint32(iounit, p)
	return Egood
}

// Initializes the specified Fcall value to contain Rread message.
// The user should copy the returned data to the slice pointed by
// fc.Data and call SetRreadCount to update the data size to the
// actual value.
func (fc *Fcall) InitRread(count uint32) W9Err {
	size := int(4 + count) /* count[4] data[count] */
	p, err := fc.packCommon(size, Rread)
	if err != Egood {
		return err
	}

	fc.Count = count
	fc.Data = p[4 : fc.Count+4]
	p = pint32(count, p)
	return Egood
}

// Updates the size of the data returned by Rread. Expects that
// the Fcall value is already initialized by InitRread.
func (fc *Fcall) SetRreadCount(count uint32) {
	/* we need to update both the packet size as well as the data count */
	size := 4 + 1 + 2 + 4 + count /* size[4] id[1] tag[2] count[4] data[count] */
	pint32(size, fc.Pkt)
	pint32(count, fc.Pkt[7:])
	fc.FcSize = size
	fc.Count = count
	fc.Pkt = fc.Pkt[0:size]
	fc.Data = fc.Data[0:count]
}

// Create a Rread message in the specified Fcall.
func (fc *Fcall) packRread(data []byte) W9Err {
	count := uint32(len(data))
	err := fc.InitRread(count)
	if err != Egood {
		return err
	}

	copy(fc.Data, data)
	return Egood
}

// Create a Rwrite message in the specified Fcall.
func (fc *Fcall) packRwrite(count uint32) W9Err {
	p, err := fc.packCommon(4, Rwrite) /* count[4] */
	if err != Egood {
		return err
	}

	fc.Count = count

	p = pint32(count, p)
	return Egood
}

// Create a Rclunk message in the specified Fcall.
func (fc *Fcall) packRclunk() W9Err {
	_, err := fc.packCommon(0, Rclunk)
	return err
}

// Create a Rremove message in the specified Fcall.
func (fc *Fcall) packRremove() W9Err {
	_, err := fc.packCommon(0, Rremove)
	return err
}

// Create a Rstat message in the specified Fcall.
func (fc *Fcall) packRstat(d *Dir) W9Err {
	stsz := statsz(d)
	size := 2 + stsz /* stat[n] */
	p, err := fc.packCommon(size, Rstat)
	if err != Egood {
		return err
	}

	p = pint16(uint16(stsz), p)
	p = pstat(d, p)
	fc.Dir = *d
	return Egood
}

// Create a Rwstat message in the specified Fcall.
func (fc *Fcall) packRwstat() W9Err {
	_, err := fc.packCommon(0, Rwstat)
	return err
}
