// Copyright 2009 The Go9p Authors.  All rights reserved.
// Portions: Copyright 2018 Larry Rau. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package warp9

// Create a Tversion message in the specified Fcall.
func (fc *Fcall) packTversion(msize uint32, version string) error {
	size := 4 + 2 + len(version) /* msize[4] version[s] */
	p, err := fc.packCommon(size, Tversion)
	if err != nil {
		return err
	}

	fc.Msize = msize
	fc.Version = version
	p = pint32(msize, p)
	p = pstr(version, p)

	return nil
}

// Create a Tauth message in the specified Fcall.
func (fc *Fcall) packTauth(fid uint32, uid uint32, aname string) error {
	size := 4 + 4 + 2 + len(aname) /* fid[4] uid[4] aname[s] */

	p, err := fc.packCommon(size, Tauth)
	if err != nil {
		return err
	}

	fc.Fid = fid
	fc.Uid = uid
	fc.Aname = aname
	p = pint32(fid, p)
	p = pint32(uid, p)
	p = pstr(aname, p)

	return nil
}

// Create a Tflush message in the specified Fcall.
func (fc *Fcall) packTflush(oldtag uint16) error {
	p, err := fc.packCommon(2, Tflush)
	if err != nil {
		return err
	}

	fc.Oldtag = oldtag
	p = pint16(oldtag, p)
	return nil
}

// Create a Tattach message in the specified Fcall.
func (fc *Fcall) packTattach(fid uint32, atok uint32, uid uint32, aname string) error {
	size := 4 + 4 + 4 + 2 + len(aname) /* fid[4] afid[4] uid[4] aname[s] */

	p, err := fc.packCommon(size, Tattach)
	if err != nil {
		return err
	}

	fc.Fid = fid
	fc.Atok = atok
	fc.Uid = uid
	fc.Aname = aname
	p = pint32(fid, p)
	p = pint32(atok, p)
	p = pint32(uid, p)
	p = pstr(aname, p)

	return nil
}

// Create a Twalk message in the specified Fcall.
func (fc *Fcall) packTwalk(fid uint32, newfid uint32, wnames []string) error {
	nwname := len(wnames)
	size := 4 + 4 + 2 + nwname*2 /* fid[4] newfid[4] nwname[2] nwname*wname[s] */
	for i := 0; i < nwname; i++ {
		size += len(wnames[i])
	}

	p, err := fc.packCommon(size, Twalk)
	if err != nil {
		return err
	}

	fc.Fid = fid
	fc.Newfid = newfid
	p = pint32(fid, p)
	p = pint32(newfid, p)
	p = pint16(uint16(nwname), p)
	fc.Wname = make([]string, nwname)
	for i := 0; i < nwname; i++ {
		fc.Wname[i] = wnames[i]
		p = pstr(wnames[i], p)
	}

	return nil
}

// Create a Topen message in the specified Fcall.
func (fc *Fcall) packTopen(fid uint32, mode uint8) error {
	size := 4 + 1 /* fid[4] mode[1] */
	p, err := fc.packCommon(size, Topen)
	if err != nil {
		return err
	}

	fc.Fid = fid
	fc.Mode = mode
	p = pint32(fid, p)
	p = pint8(mode, p)
	return nil
}

// Create a Tcreate message in the specified Fcall.
func (fc *Fcall) packTcreate(fid uint32, name string, perm uint32, mode uint8, extattr string) error {
	size := 4 + 2 + len(name) + 4 + 1 /* fid[4] name[s] perm[4] mode[1] */

	//RAU extattr size += 2 + len(extattr)

	p, err := fc.packCommon(size, Tcreate)
	if err != nil {
		return err
	}

	fc.Fid = fid
	fc.Name = name
	fc.Perm = perm
	fc.Mode = mode
	p = pint32(fid, p)
	p = pstr(name, p)
	p = pint32(perm, p)
	p = pint8(mode, p)

	//RAU extended attr ExtAttr
	fc.ExtAttr = extattr
	//p = pstr(extattr, p)

	return nil
}

// Create a Tread message in the specified Fcall.
func (fc *Fcall) packTread(fid uint32, offset uint64, count uint32) error {
	size := 4 + 8 + 4 /* fid[4] offset[8] count[4] */
	p, err := fc.packCommon(size, Tread)
	if err != nil {
		return err
	}

	fc.Fid = fid
	fc.Offset = offset
	fc.Count = count
	p = pint32(fid, p)
	p = pint64(offset, p)
	p = pint32(count, p)
	return nil
}

// Create a Twrite message in the specified Fcall.
func (fc *Fcall) packTwrite(fid uint32, offset uint64, count uint32, data []byte) error {
	c := len(data)
	size := 4 + 8 + 4 + c /* fid[4] offset[8] count[4] data[count] */
	p, err := fc.packCommon(size, Twrite)
	if err != nil {
		return err
	}

	fc.Fid = fid
	fc.Offset = offset
	fc.Count = count
	p = pint32(fid, p)
	p = pint64(offset, p)
	p = pint32(count, p)
	fc.Data = p
	copy(fc.Data, data)
	return nil
}

// Create a Tclunk message in the specified Fcall.
func (fc *Fcall) packTclunk(fid uint32) error {
	p, err := fc.packCommon(4, Tclunk) /* fid[4] */
	if err != nil {
		return err
	}

	fc.Fid = fid
	p = pint32(fid, p)
	return nil
}

// Create a Tremove message in the specified Fcall.
func (fc *Fcall) packTremove(fid uint32) error {
	p, err := fc.packCommon(4, Tremove) /* fid[4] */
	if err != nil {
		return err
	}

	fc.Fid = fid
	p = pint32(fid, p)
	return nil
}

// Create a Tstat message in the specified Fcall.
func (fc *Fcall) packTstat(fid uint32) error {
	p, err := fc.packCommon(4, Tstat) /* fid[4] */
	if err != nil {
		return err
	}

	fc.Fid = fid
	p = pint32(fid, p)
	return nil
}

// Create a Twstat message in the specified Fcall.
func (fc *Fcall) packTwstat(fid uint32, d *Dir) error {
	stsz := 0            // statsz(d, dotu)
	size := 4 + 2 + stsz /* fid[4] stat[n] */
	p, err := fc.packCommon(size, Twstat)
	if err != nil {
		return err
	}

	fc.Fid = fid
	fc.Dir = *d
	p = pint32(fid, p)
	p = pint16(uint16(stsz), p)
	p = pstat(d, p)
	//rau: handle extended attribute
	return nil
}
