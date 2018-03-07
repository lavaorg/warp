// Copyright 2009 The ninep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/lavaorg/lrt/mlog"
	"github.com/lavaorg/warp/servers/sensim"
)

var addr = flag.String("a", "127.0.0.1:9901", "network address")
var debug = flag.Int("d", 0, "print debug messages")
var listen = flag.Bool("listen", false, "listen for connections")
var count = flag.Int("count", 100, "sensor count")
var life = flag.Int("life", 100, "life span")
var sleep = flag.Int("sleep", 10, "sleep in sec")

type sensor struct {
	wg    *sync.WaitGroup
	srv   *sensim.SenSrv
	addr  string
	life  int
	sleep int
}

func main() {
	flag.Parse()
	sensrv := new(sensim.SenSrv)

	sensrv.Dotu = false
	sensrv.Id = "sensrv"
	sensrv.Debuglevel = *debug
	sensrv.Start(sensrv)
	fmt.Print("sensrv starting\n")

	// serve our object tree
	if *listen {
		err := sensrv.StartNetListener("tcp", *addr)
		if err != nil {
			log.Println(err)
		}
	} else {
		runSensors(sensrv, *count, *life, *sleep)
	}

}

func runSensors(sensrv *sensim.SenSrv, count, life, sleep int) {

	var wg sync.WaitGroup

	wg.Add(count)
	for ; count > 0; count-- {
		go sensorMain(sensor{&wg, sensrv, *addr, life, sleep})
	}
	wg.Wait()

}

func sensorMain(s sensor) {

	defer s.wg.Done()
	for ; s.life > 0; s.life-- {
		c, e := net.Dial("tcp", s.addr)
		if e != nil {
			mlog.Error("Dial error:%v", e)
			return
		}
		s.srv.NewConnWait(c)
		time.Sleep(time.Duration(s.sleep) * time.Second)
	}
}
