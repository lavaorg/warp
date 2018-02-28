// Copyright 2009 The ninep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/lavaorg/warp9/servers/ufs"
	"github.com/lavaorg/warp9/warp9"
)

var addr = flag.String("addr", ":5640", "network address")
var debug = flag.Int("debug", 0, "print debug messages")
var root = flag.String("root", "/", "root filesystem")

func main() {
	flag.Parse()
	ufs := new(ufs.Ufs)
	showInterfaces(ufs)

	ufs.Dotu = true
	ufs.Id = "ufs"
	ufs.Root = *root
	ufs.Debuglevel = *debug
	ufs.Start(ufs)
	fmt.Print("ufs starting\n")
	// determined by build tags
	//extraFuncs()
	err := ufs.StartNetListener("tcp", *addr)
	if err != nil {
		log.Println(err)
	}
}

func showInterfaces(ifaces interface{}) {
	if _, ok := (ifaces).(warp9.SrvReqOps); ok {
		fmt.Println("implements: SrvReqOps")
	}
	if _, ok := (ifaces).(warp9.StatsOps); ok {
		fmt.Println("implements: StatsOps")
	}

}
