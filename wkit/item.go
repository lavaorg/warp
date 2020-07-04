// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package wkit

import (
	"github.com/lavaorg/warp/warp9"
	"sync/atomic"
	"time"
)

type (

	// Item represents a generic object in a hierarchical tree.  This interface
	// will allow the object server to perform generic operations on any object.
	// Where necessary it can learn more details of what the object is (e.g. a directory,
	// or a bind-point, etc)
	Item interface {
		GetDir() *warp9.Dir
		GetItem() Item
		IsDirectory() Directory
		Parent() Directory
		SetParent(d Directory) error
		GetQid() warp9.Qid
		Walked() (Item, error)
		SetMode(uint32) Item
		Read(obuf []byte, off uint64, rcount uint32) (uint32, error)
		Write(ibuf []byte, off uint64, count uint32) (uint32, error)
		Open(mode byte) (uint32, error)
		Clunk() error
		Remove() error
		Stat() (*warp9.Dir, error)
		WStat(dir *warp9.Dir) error
	}

	// OneItem is a generic in-memory blob object. The contents of the
	// object are arbitrary bytes and can be written/read.
	OneItem struct {
		warp9.Dir
		parent Directory
		opened bool
		buffer []byte
	}
)

func NewItem(name string) *OneItem {
	return &OneItem{
		Dir: warp9.Dir{
			Name:  name,
			Qid:   warp9.Qid{warp9.QTOBJ, 0, NextQid()},
			Mode:  warp9.DMAPPEND | uint32(Perms(warp9.ORDWR, warp9.ORDWR, warp9.ORDWR)),
			Atime: uint32(time.Now().Unix()),
			Mtime: uint32(time.Now().Unix()),
		},
		parent: nil,
		opened: false,
		buffer: make([]byte, 0),
	}
}

func (o *OneItem) Buffer() []byte {
	return o.buffer
}

func (o *OneItem) SetBuffer(buf []byte) Item {
	o.buffer = buf
	return o
}

func (o *OneItem) SetDirectory(dir warp9.Dir) Item {
	o.Dir = dir
	return o
}

func (o *OneItem) SetOpen(isOpen bool) Item {
	o.opened = isOpen
	return o
}

func (o *OneItem) SetMode(mode uint32) Item {
	o.Dir.Mode = mode
	return o
}

// Return the object's Dir structure
func (o *OneItem) GetDir() *warp9.Dir {
	return &o.Dir
}

// Return the object as the interface type Item.
func (o *OneItem) GetItem() Item {
	return o
}

// Indicate that the object is not a directory, returns nil.
func (o *OneItem) IsDirectory() Directory {
	return nil
}

// Return the object's parent (eg container).
func (o *OneItem) Parent() Directory {
	return o.parent
}

// Set the parent of the object.
func (o *OneItem) SetParent(d Directory) error {
	o.parent = d
	return nil
}

// GetQid returns
func (o *OneItem) GetQid() warp9.Qid {
	return o.Dir.Qid
}

// Walked prints a debug message.
func (o *OneItem) Walked() (Item, error) {
	warp9.Debug("o.Walked: %T.%v", o, o)
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

// set item to open status
func (o *OneItem) Open(mode byte) (uint32, error) {
	o.opened = true
	return 0, nil
}

func (o *OneItem) Clunk() error {
	o.opened = false
	return nil
}

func (o *OneItem) Remove() error {
	return warp9.ErrorCode(warp9.Enotimpl)
}

func (o *OneItem) Stat() (*warp9.Dir, error) {
	return &o.Dir, nil
}

func (o *OneItem) WStat(dir *warp9.Dir) error {
	return warp9.ErrorCode(warp9.Enotimpl)
}

// General Item functions

// Create a Permissions word out of user/group/other components
// The lower 3 bytes set in User/Group/Other order
func Perms(u, g, o byte) uint16 {
	return uint16(uint16(u)<<6 | uint16(g)<<3 | uint16(o))
}

// Helper for creation of Qid-ID. Simple incremented counter.
func NextQid() uint64 {
	return atomic.AddUint64(&qidpGlob, 1)
}

//
// private helper functions
//

func check(mode, kind uint32) bool {
	return mode&kind > 0
}
