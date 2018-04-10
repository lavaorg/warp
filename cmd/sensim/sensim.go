// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

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

// simulate 'count' (configurable) number of independent sensors
// make a thread for each sensor; wait for them all to complete
func runSensors(sensrv *sensim.SenSrv, count, life, sleep int) {

	var wg sync.WaitGroup

	wg.Add(count)
	for ; count > 0; count-- {
		go sensorMain(sensor{&wg, sensrv, *addr, life, sleep})
	}
	wg.Wait()

}

// initiate a connection and then serve our object server on that connection
// e.g. we will be expecting the server contacted to be a client of our object tree
// we will then sleep for a configured period and repeast the process
// we repeast this process for a configured life; then thread ends
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
