// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package osrv

// provide a hosted environment for simple object servers.
// the host environment will manage the communication, protocol,
// user management, fid pools, and basic functions

// this is an interface a object server must implement to be hosted
type WOserver interface {
	// Return an instanciated root directory to serve as the root of a mount point
	// the server does not need to instanciate the root until called.
	RootDir() *WODir
	// Host will suggest a min/max size; server can suggest a lower max or refuse
	//SuggestMsize(min, max uint32) (uint32, error)
	// Symbolic name of the service; used for mount requests
	Name() string
}

// map of registered servers
var Servers map[string]WOserver = map[string]WOserver{} //later don't export this

func Register(os WOserver) error {
	n := os.Name()
	Servers[n] = os
	return nil
}

type Mount interface {
	AddObject(name string, handler WHandler)
}

type MountPoint struct {
}

func (mt *MountPoint) AddObject(name string, handler WHandler) {

}

func Handle(name string) {

}
