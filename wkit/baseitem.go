// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package wkit

import (
	"time"

	"github.com/lavaorg/warp/warp9"
)

type (

	// BaseItem is a generic in-memory blob object. The contents of the
	// object are arbitrary bytes and can be written/read.
	BaseItem struct {
		warp9.Dir
		parent Directory
		opened bool
	}
)

func NewBaseItem(name string, mkdir bool) *BaseItem {
	var otyp uint8 = warp9.QTOBJ
	if mkdir {
		otyp = warp9.QTDIR
	}
	return &BaseItem{
		Dir: warp9.Dir{
			Name:  name,
			Qid:   warp9.Qid{otyp, 0, NextQid()},
			Mode:  warp9.DMAPPEND | uint32(Perms(warp9.ORDWR, warp9.ORDWR, warp9.ORDWR)),
			Atime: uint32(time.Now().Unix()),
			Mtime: uint32(time.Now().Unix()),
		},
		parent: nil,
		opened: false,
	}
}

// not sure this should be allowed
//func (o *BaseItem) SetDirectory(dir warp9.Dir) Item {
//	o.Dir = dir
//	return o
//}

func (o *BaseItem) SetOpen(isOpen bool) {
	o.opened = isOpen
}

func (o *BaseItem) SetMode(mode uint32) {
	o.Dir.Mode = mode
}

// Return the object's Dir structure
func (o *BaseItem) GetDir() *warp9.Dir {
	return &o.Dir
}

// Return the object as the interface type Item.
// This method should be provided by the concrete type
func (o *BaseItem) GetItem() Item {
	return o
}

// Indicate that the object is not a directory, returns nil.
func (o *BaseItem) IsDirectory() Directory {
	return nil
}

// Return the object's parent (eg container).
func (o *BaseItem) Parent() Directory {
	return o.parent
}

// Set the parent of the object.
func (o *BaseItem) SetParent(d Directory) error {
	o.parent = d
	return nil
}

// GetQid returns
func (o *BaseItem) GetQid() warp9.Qid {
	return o.Dir.Qid
}

// Walked performs no action.
func (o *BaseItem) Walked() (Item, error) {
	return o, nil
}

// Read is a no-op. If opened will return (0,nil). If not opened will return error
// indicating the object is not open (0,Enotopen).
//
func (o *BaseItem) Read(obuf []byte, off uint64, rcount uint32) (uint32, error) {
	if o.opened {
		return 0, nil
	}
	// otherwise return an error
	return 0, warp9.ErrorCode(warp9.Enotopen)
}

// Write is a no-op. If opened write will return (0,nil). If not opened will
// return error indicating the object is not open. (0.Enotopen)
func (o *BaseItem) Write(ibuf []byte, off uint64, count uint32) (uint32, error) {
	if o.opened {
		return 0, nil
	}
	// otherwise return an error
	return 0, warp9.ErrorCode(warp9.Enotopen)
}

// set item to open status
func (o *BaseItem) Open(mode byte) (uint32, error) {
	o.opened = true
	return 0, nil
}

// Clunk will mark the object closed and return nil. If the item
// is not in the open state an error is returned (Enotopen).
//
func (o *BaseItem) Clunk() error {
	if o.opened {
		o.opened = false
		return nil
	}
	//return warp9.ErrorCode(warp9.Enotopen)  //implement a ref count; diff fids can open object (walked should clone)
	return nil
}

// Remove is not implemented. Returns Enotimpl.
func (o *BaseItem) Remove() error {
	return warp9.ErrorCode(warp9.Enotimpl)
}

// Stat returns the Warp Stat object (Dir). and nil.
// A created object always has a Dir structure.
func (o *BaseItem) Stat() (*warp9.Dir, error) {
	return &o.Dir, nil
}

// WStat is not implemented. Returns Enotimpl
func (o *BaseItem) WStat(dir *warp9.Dir) error {
	return warp9.ErrorCode(warp9.Enotimpl)
}
