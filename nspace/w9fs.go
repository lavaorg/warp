// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

// warp9 package manages a dynamic namespace of remote resources.
// The namespace is represnted by a special FS that can be altered by issuing commands
// to the ctl file of this FS.  The Warp9FS consistes of the following top level directories:
// /warp9/ctl
// /svc
// /mnt
//
// the /warp9 directory contains a single /ctl file which commands can be issued to in
// order to manipulate the warp9 environment.
//
// /svc is the top level directory for organizing servcies.
//
// /mnt is the top level directory where remote FS will be mounted.  After mouunting
// bind operations can be performed to move these resources around under the /svc
//
// the operations permitted on the /warp9/ctl file are:
//    mount
//    unmount
//    bind
//
// mount operations will only allow mounting to the /mnt tree.
// bind opeations can only happen within the /svc tree
//
package warp9

import (
	_ "github.com/lavaorg/warp9/warp9"
)

// the client's fid information
type w9Fid struct {
	path string
}
