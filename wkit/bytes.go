// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package wkit

import (
	"github.com/lavaorg/warp/warp9"
)

// BytesSequence provides a sequence of bytes
type BytesSequence interface {
	Bytes() []byte
}

type (
	BytesItem struct {
		*BaseItem
		bytes BytesSequence
	}
)

// NewBytesItem constructs a new DigestItem object whtat will invoke the
// provided hash function each time read is called.
//
func NewBytesItem(name string, b BytesSequence) *BytesItem {
	return &BytesItem{
		BaseItem: NewBaseItem(name, false),
		bytes:    b,
	}
}

// GetItem return the object as the interface type Item.
func (o *BytesItem) GetItem() Item {
	return o
}

// Walked performs no action.
func (o *BytesItem) Walked() (Item, error) {
	return o, nil
}

// Read will get the bytes from the service and return back.
// obuf: must be large enough to hold all bytes
// off: should be 0;  >0 == EOF (this handles reading until EOF)
// rcount: should be size of obuf
//
func (o *BytesItem) Read(obuf []byte, off uint64, rcount uint32) (uint32, error) {
	if off > 0 {
		return 0, warp9.WarpErrorEOF
	}
	buf := o.bytes.Bytes()
	n := copy(obuf, buf)
	if n != len(buf) {
		return 0, warp9.ErrorCode(warp9.Ebufsmall)
	}
	return uint32(n), nil
}
