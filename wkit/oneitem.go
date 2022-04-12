// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package wkit

import (
	"github.com/lavaorg/warp/warp9"
)

type (

	// OneItem is a generic in-memory blob object. The contents of the
	// object are arbitrary bytes and can be written/read.
	OneItem struct {
		*BaseItem
		buffer []byte
	}
)

// NewItem returns a new OneItem instance.
//
func NewItem(name string) *OneItem {
	return &OneItem{
		BaseItem: NewBaseItem(name, false),
		buffer:   make([]byte, 0),
	}
}

// Buffer returns the current object's byte buffer.
//
func (o *OneItem) Buffer() []byte {
	return o.buffer
}

// SetBuffer replaces the current object's byte buffer with buf.
//
func (o *OneItem) SetBuffer(buf []byte) Item {
	o.buffer = buf
	return o
}

// Return the object as the interface type Item.
func (o *OneItem) GetItem() Item {
	return o
}

// Walked performs no action.
func (o *OneItem) Walked() (Item, error) {
	return o, nil
}

// Return the requested set of bytes from the object's byte buffer.
func (o *OneItem) Read(obuf []byte, off uint64, rcount uint32) (uint32, error) {
	// determine which and how many bytes to return
	var count uint32
	switch {
	case off > uint64(len(o.buffer)):
		count = 0
	case uint32(len(o.buffer[off:])) > rcount:
		count = rcount
	default:
		count = uint32(len(o.buffer[off:]))
	}
	n := copy(obuf, o.buffer[off:uint32(off)+count])
	if uint32(n) != count {
		return 0, warp9.ErrorCode(warp9.Eio)
	}
	warp9.Debug("o.Read:buffer:%v, obuf: %v, off:%v, rcount:%v\n", len(o.buffer), len(obuf), off, count)
	warp9.Debug("o.Read:%T %p %v", o, o, o.Qid)
	return count, nil

}

// Write bytes into the object's byte buffer.
// The object's byte buffer is smaller than len(max(int))
func (o *OneItem) Write(ibuf []byte, off uint64, count uint32) (uint32, error) {

	// our file will not be super large;  convert everything to int
	ioff := int(off)
	icnt := int(count)
	if uint64(ioff) != off || uint32(icnt) != count {
		return 0, warp9.ErrorCode(warp9.Etoolarge)
	}

	// if append file always just append
	// if offset is the current len; just append
	if check(o.Mode, warp9.DMAPPEND) || ioff == len(o.buffer) {
		o.buffer = append(o.buffer, ibuf[:icnt]...)
		return count, nil
	}
	// if offset < cur len; truncate current and append
	if ioff < len(o.buffer) {
		o.buffer = append(o.buffer[:off], ibuf[:icnt]...)
		return count, nil
	}

	// if we are seeking past eof then add 0's first
	if ioff >= len(o.buffer) {
		zsz := ioff - len(o.buffer) - 1
		z := make([]byte, zsz, zsz+icnt)
		z = append(z, ibuf[:icnt]...)
		o.buffer = append(o.buffer, z...)
		return count, nil
	}

	return 0, warp9.ErrorCode(warp9.Eio)
}
