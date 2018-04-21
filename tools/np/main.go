// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lavaorg/lrt/mlog"
	"github.com/lavaorg/warp/warp9"
)

var verbose = flag.Bool("v", false, "verbose mode")
var dbglev = flag.Int("d", 0, "debuglevel")
var addr = flag.String("a", "127.0.0.1:5640", "network address")
var aname = flag.String("aname", "/", "path on server to use as root")

func main() {

	flag.Parse()

	uid := uint32(0xFFFFFFFF & uint32(os.Getuid()))
	user := warp9.Identity.User(uid)
	warp9.DefaultDebuglevel = *dbglev

	c9, err := warp9.Mount("tcp", *addr, *aname, 8192, user)
	if err != warp9.Egood {
		log.Fatalf("Error:%v:%v\n", err, err.String())
	}
	defer c9.Clunk(c9.Root)

	if flag.NArg() != 2 {
		usage()
		log.Fatal("expected an argument")
	}

	cmd := flag.Arg(0)

	switch cmd {

	default:
		{
			usage()
			log.Fatal("unkonwn cmd")
		}

	case "cat":
		cmdcat(c9)
	case "ls":
		cmdls(c9)

	case "stat":
		cmdstat(c9)
	case "get":
		cmdget(c9)
	}
	return
}

func usage() {
	fmt.Println("usage: np [-v][-d dbglev] [-a addr] cmd arg")
	fmt.Println("\tcmd = {ls,stat,cat,echo}")
}

func cmdcat(c9 *warp9.Clnt) {
	f, err := c9.Open(flag.Arg(1), warp9.OREAD)
	if err != warp9.Egood {
		log.Fatalf("Error:%v\n", err)
	}
	defer f.Close()

	buf := make([]byte, 8192)
	for {
		n, err := f.Read(buf)
		if n == 0 {
			break
		}
		if err != warp9.Egood {
			log.Fatal("Error reading:%v\n", err)
		}
		os.Stdout.Write(buf[0:n])
		if err == warp9.Eeof {
			break
		}
	}

	if err != warp9.Egood && err != warp9.Eeof {
		log.Fatalf("Error:%v\n", err)
	}

}

func cmdls(c9 *warp9.Clnt) {

	fid, err := c9.Walk(flag.Arg(1))
	if err != warp9.Egood {
		mlog.Error("error:%v", err)
		return
	}
	defer c9.Clunk(fid)

	mlog.Debug("ls: fid: %v", fid)

	if fid.Qid.Type&warp9.QTDIR > 0 {
		// read directory
		f, err := c9.FOpenObject(fid, warp9.OREAD)
		if err != warp9.Egood {
			mlog.Error("error:%v", err)
			return
		}

		mlog.Debug("fid opened; readdir:%v", fid)
		for {
			d, err := f.Readdir(0)
			if d == nil || len(d) == 0 || err != warp9.Egood {
				break
			}
			for i := 0; i < len(d); i++ {
				os.Stdout.WriteString(d[i].Name + "\n")
			}
		}
	} else {
		// stat the file
		d, err := c9.FStat(fid)
		if err != warp9.Egood {
			log.Println("Error", err)
			return
		}
		fmt.Printf("%v\n", d)
	}
	if err != warp9.Egood && err != warp9.Eeof {
		log.Fatal(err)
	}
}

func cmdstat(c9 *warp9.Clnt) {
	d, err := c9.Stat(flag.Arg(1))
	if err != warp9.Egood {
		log.Println("Error", err)
		return
	}

	if !*verbose {
		fmt.Printf("%v\n", d)
	} else {
		fmt.Printf("    Name: %s\n", d.Name)
		fmt.Printf("    Size: %d\n", d.Length)
		fmt.Printf("    Mode: %s\n", warp9.PermToString(d.Mode))
		fmt.Printf(" UID:GID: %d:%d\n", d.Uid, d.Gid)
		fmt.Printf("     Qid: %s\n", d.Qid.String())
		fmt.Printf("  Access: %v\n", d.Atime)
		fmt.Printf("  Modify: %v\n", time.Unix(int64(d.Mtime), 0))
		fmt.Printf("Last uid: %d\n", d.Muid)
		fmt.Printf("     Ext: %s\n", d.ExtAttr)
	}

}

func cmdget(c9 *warp9.Clnt) {
	data, qid, err := c9.Get(flag.Arg(1), 0)
	if err != warp9.Egood {
		log.Fatalf("Error:%v\n", err)
	}
	if *verbose {
		fmt.Printf("qid = %v\nData:\n", qid)
	}
	os.Stdout.Write(data[0:])
}
