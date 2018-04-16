// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package osrv

import (
	"github.com/lavaorg/warp/warp9"
)

type WOpener interface {
	WOpen(*warp9.SrvReq) warp9.W9Err
}
type WReader interface {
	WRead(*warp9.SrvReq) warp9.W9Err
}
type WWriter interface {
	WWrite(*warp9.SrvReq) warp9.W9Err
}
type WClunker interface {
	WClunk(*warp9.SrvReq) warp9.W9Err
}
type WStater interface {
	WStat(*warp9.SrvReq) warp9.W9Err
}
type WWstater interface {
	WWstat(*warp9.SrvReq) warp9.W9Err
}
type WWalker interface {
	WWalk(*warp9.SrvReq) warp9.W9Err
}
type WCreater interface {
	WCreate(*warp9.SrvReq) warp9.W9Err
}

type WHandler interface {
}

// A base Warp Base Object definition.
type WOBase struct {
	warp9.Dir
	parent *WODir
	Data   interface{}
}

// A Warp Directory Object defintion.  WODir objects are containers holding references
// to other Warp Objects.
type WODir struct {
	WOBase
	contents []*WOBase //contents of this directory
	stack    []*WOBase //reference to union-peers
}

func (wo *WOBase) Open(req *warp9.SrvReq) warp9.W9Err {

	return warp9.Enotimpl
}

func (wo *WOBase) Read(req *warp9.SrvReq) warp9.W9Err {

	return warp9.Enotimpl
}

func (wo *WOBase) Write(req *warp9.SrvReq) warp9.W9Err {

	return warp9.Enotimpl
}

func (wo *WOBase) Clunk(req *warp9.SrvReq) warp9.W9Err {

	return warp9.Egood
}

func (wo *WOBase) Stat(req *warp9.SrvReq) warp9.W9Err {

	return warp9.Egood
}

func (wo *WOBase) Wstat(req *warp9.SrvReq) warp9.W9Err {

	return warp9.Enotimpl
}

func (wd *WODir) Walk(req *warp9.SrvReq) warp9.W9Err {

	return warp9.Egood
}

func (wd *WODir) Create(req *warp9.SrvReq) warp9.W9Err {

	return warp9.Enotimpl
}

func (wo *WOBase) DirEntry() *warp9.Dir {
	return &(wo.Dir)
}

func (wo *WOBase) Parent() *WODir {
	return wo.parent
}
