// Copyright 2009 The Go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package warp9

// Reads count bytes starting from offset from the object associated with the fid.
// Returns a slice with the data read, if the operation was successful, or an
// Error.
func (clnt *Clnt) Read(fid *Fid, offset uint64, count uint32) ([]byte, W9Err) {
	if count > fid.Iounit {
		count = fid.Iounit
	}
	tc := clnt.NewFcall()
	err := tc.packTread(fid.Fid, offset, count)
	if err != Egood {
		return nil, err
	}

	rc, err := clnt.Rpc(tc)
	if err != Egood {
		return nil, err
	}

	return rc.Data, Egood
}

// Reads up to len(buf) bytes from the Object. Returns the number
// of bytes read, or an Error.
func (obj *Object) Read(buf []byte) (int, W9Err) {
	n, err := obj.ReadAt(buf, int64(obj.offset))
	if err == Egood {
		obj.offset += uint64(n)
	}

	return n, Egood
}

// Reads up to len(buf) bytes from the object starting from offset.
// Returns the number of bytes read, or an Error.
func (obj *Object) ReadAt(buf []byte, offset int64) (int, W9Err) {
	b, err := obj.Fid.Clnt.Read(obj.Fid, uint64(offset), uint32(len(buf)))
	if err != Egood {
		return 0, err
	}

	if len(b) == 0 {
		return 0, Eeof
	}

	copy(buf, b)
	return len(b), Egood
}

// Reads exactly len(buf) bytes from the Object starting from offset.
// Returns the number of bytes read (could be less than len(buf) if
// end-of-data of the object is reached), or an Error.
func (obj *Object) Readn(buf []byte, offset uint64) (int, W9Err) {
	ret := 0
	for len(buf) > 0 {
		n, err := obj.ReadAt(buf, int64(offset))
		if err != Egood {
			return 0, err
		}

		if n == 0 {
			break
		}

		buf = buf[n:]
		offset += uint64(n)
		ret += n
	}

	return ret, Egood
}

// Reads the content of the directory associated with the Object.
// Returns an array of maximum num entries (if num is 0, returns
// all entries from the directory). If the operation fails, returns
// an Error.
func (obj *Object) Readdir(num int) ([]*Dir, W9Err) {
	buf := make([]byte, obj.Fid.Clnt.Msize-IOHDRSZ)
	dirs := make([]*Dir, 32)
	pos := 0
	offset := obj.offset
	defer func() {
		obj.offset = offset
	}()
	for {
		n, err := obj.Read(buf)
		if err != Egood && err != Eeof {
			return nil, err
		}

		if n == 0 {
			break
		}

		for b := buf[0:n]; len(b) > 0; {
			d, _, _, perr := UnpackDir(b)
			if perr != Egood {
				// If we have unpacked anything, it is almost certainly
				// a too-short buffer. So return what we got.
				if pos > 0 {
					return dirs[0:pos], Egood
				}
				return nil, perr
			}
			b = b[d.DirSize+2:]
			offset += uint64(d.DirSize + 2)
			if pos >= len(dirs) {
				s := make([]*Dir, len(dirs)+32)
				copy(s, dirs)
				dirs = s
			}

			dirs[pos] = d
			pos++
			if num != 0 && pos >= num {
				break
			}
		}
	}

	return dirs[0:pos], Egood
}
