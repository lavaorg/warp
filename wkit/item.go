// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package wkit

import (
	"github.com/lavaorg/warp/warp9"
)

type (

	// Item represents a generic object in a hierarchical tree.  This interface
	// will allow the object server to perform generic operations on any object.
	// Where necessary it can learn more details of what the object is (e.g. a directory,
	// or a bind-point, etc)
	Item interface {
		GetDir() *warp9.Dir
		GetItem() Item //a sub-type should implement to return actual instance
		IsDirectory() Directory
		Parent() Directory
		SetParent(d Directory) error
		GetQid() warp9.Qid
		Walked() (Item, error) //a sub-type should directly implement to return actual instance
		SetMode(uint32)
		Read(obuf []byte, off uint64, rcount uint32) (uint32, error)
		Write(ibuf []byte, off uint64, count uint32) (uint32, error)
		Open(mode byte) (uint32, error)
		Clunk() error
		Remove() error
		Stat() (*warp9.Dir, error)
		WStat(dir *warp9.Dir) error
	}
)
