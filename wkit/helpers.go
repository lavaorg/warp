// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package wkit

import (
	"sync/atomic"
)

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
