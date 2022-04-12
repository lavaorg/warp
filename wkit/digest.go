// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package wkit

import (
	"encoding/binary"

	"github.com/lavaorg/warp/warp9"
)

// Digest32 is a subset of hash.Hash32
type Digest32 interface {
	Sum32() uint32
}

type (

	// DigestItem allows a service to present an object that is a digest of content controlled
	// by the service. the content it represents will depend on the service.  The digest is
	// not intended to be cryptographically secure.
	// The digest will be a 32bit unsigned entity in little-endian format.
	DigestItem struct {
		*BaseItem
		hash    Digest32
		lastval uint32 //cache last value
		buffer  []byte //cache conversion
	}
)

// NewDigestItem constructs a new DigestItem object whtat will invoke the
// provided hash function each time read is called.
//
func NewDigestItem(name string, h Digest32) *DigestItem {
	return &DigestItem{
		BaseItem: NewBaseItem(name, false),
		hash:     h,
		lastval:  0,
		buffer:   make([]byte, 4),
	}
}

// GetItem return the object as the interface type Item.
func (o *DigestItem) GetItem() Item {
	return o
}

// Walked performs no action.
func (o *DigestItem) Walked() (Item, error) {
	return o, nil
}

// Read will invoke the associated hash function and return the uint32 as
// a sequence of bytes in little endian form.
// off should be 0.  >0 will result in a eof return
// len(obuf) and rcount must be >=4
//
// using the "encoding/binary" package, this method does:
//
//     binary.LittleEndian.PutUin32(obuf,value)
// the client should be able to perform:
//     val := binary.LittleEndian.Uint32(buf)
//
// note: handling off >0 will allow for reading to EOF
//
func (o *DigestItem) Read(obuf []byte, off uint64, rcount uint32) (uint32, error) {
	if off > 0 {
		return 0, warp9.WarpErrorEOF
	}
	if len(obuf) < 4 || rcount < 4 {
		return 0, warp9.ErrorCode(warp9.Ebufsmall)
	}
	warp9.Debug("o.hash=%p, %T\n", o.hash, o.hash)
	// inovke the hash function
	digest := o.hash.Sum32()

	if digest != o.lastval {
		//convert to byte sequence
		binary.LittleEndian.PutUint32(o.buffer, digest)
		o.lastval = digest
	}
	copy(obuf, o.buffer)
	return 4, nil
}
