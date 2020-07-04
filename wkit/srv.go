// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package wkit

import (
	"github.com/lavaorg/warp/warp9"
)

// Interface representing an object server.
type (
	Server interface {
		StartNetListener(ntype, addr string) error
		GetRoot() *DirItem
	}

	// A server controller holds a warp9 server along with
	// associated metadata and configuration data.
	ServerController struct {
		warp9.Srv
		stats warp9.StatsOps
		root  Directory
	}
)

var qidpGlob uint64 = 0

func NewServer(id string, debugLevel int, root Directory) *ServerController {
	server := &ServerController{
		Srv: warp9.Srv{
			Id:         id,
			Debuglevel: debugLevel,
			Msize:      8192,
		},
		root: root,
	}
	return server
}

// Return the root object for the server.
func (srv *ServerController) GetRoot() Directory {
	return srv.root
}
