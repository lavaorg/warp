// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package wkit

import (
	"sync"

	"github.com/lavaorg/warp/warp9"
)

type Binder interface {
	Lookup(item string) Item
	IsDirectory() bool
}

// a BindPoint is both a binder and a DirItem.  As a DirItem it will have a directory
// entry of its contents Item.  When its data is read the data is the concatenation
// of all of its bind-points chain with early entries obscuring the later entries (e.g.
// a uniion directory)
type BindPoint struct {
	sync.Mutex
	DirItem
	dir    bool
	parent *BindPoint
	next   *BindPoint //reference to union-peers (we only ever look back)
}

func NewBindPoint(i Item) *BindPoint {
	var b BindPoint
	b.dir = i.IsDirectory() != nil
	return &b
}

func (bp *BindPoint) GetDir() *warp9.Dir {
	return bp.GetDir()
}

func (bp *BindPoint) GetData() []byte {
	return bp.DirItem.GetData()
}

func (d *BindPoint) Lookup(n string) Item {
	return nil
}

func (bp *BindPoint) IsDirectory() bool {
	return bp.dir
}
