// Copyright 2009 The ninep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/lavaorg/warp9/ninep"
	"github.com/lavaorg/warp9/servers/nullfs"
)

var addr = flag.String("addr", ":5640", "network address")
var debug = flag.Int("debug", 0, "print debug messages")

func main() {
	flag.Parse()
	nullfs := new(nullfs.Nullfs)
	showInterfaces(nullfs)

	nullfs.Dotu = true
	nullfs.Id = "nullfs"
	nullfs.Debuglevel = *debug
	nullfs.Start(nullfs)
	fmt.Print("nullfs starting\n")

	err := nullfs.StartNetListener("tcp", *addr)
	if err != nil {
		log.Println(err)
	}
}

func showInterfaces(ifaces interface{}) {
	if _, ok := (ifaces).(ninep.SrvReqOps); ok {
		fmt.Println("implements: SrvReqOps")
	}
	if _, ok := (ifaces).(ninep.StatsOps); ok {
		fmt.Println("implements: StatsOps")
	}

}
