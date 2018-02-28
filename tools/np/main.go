// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/lavaorg/warp9/warp9"
)

var verbose = flag.Bool("v", false, "verbose mode")
var dbglev = flag.Int("d", 0, "debuglevel")
var addr = flag.String("a", "127.0.0.1:5640", "network address")
var aname = flag.String("aname", "/", "path on server to use as root")

func main() {

	flag.Parse()

	uid := warp9.OsUsers.Uid2User(os.Geteuid())
	warp9.DefaultDebuglevel = *dbglev

	c9, err := warp9.Mount("tcp", *addr, *aname, 8192, uid)
	if err != nil {
		log.Fatalf("Error:%v\n", err)
	}

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
	case "fcat":
		cmdfcat(c9)
	}
	return
}

func usage() {
	fmt.Println("usage: np [-v][-d dbglev] [-a addr] cmd arg")
	fmt.Println("\tcmd = {ls,stat,cat,echo}")
}

func cmdcat(c9 *warp9.Clnt) {
	f, err := c9.FOpen(flag.Arg(1), warp9.OREAD)
	if err != nil {
		log.Fatalf("Error:%v\n", err)
	}
	defer f.Close()

	buf := make([]byte, 8192)
	for {
		n, err := f.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			log.Fatal("Error reading:%v\n", err)
		}
		os.Stdout.Write(buf[0:n])
		if err == io.EOF {
			break
		}
	}

	if err != nil && err != io.EOF {
		log.Fatalf("Error:%v\n", err)
	}

}

func cmdls(c9 *warp9.Clnt) {

	f, err := c9.FOpen(flag.Arg(1), warp9.OREAD)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for {
		d, err := f.Readdir(0)
		if d == nil || len(d) == 0 || err != nil {
			break
		}
		for i := 0; i < len(d); i++ {
			os.Stdout.WriteString(d[i].Name + "\n")
		}
	}

	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}

func cmdstat(c9 *warp9.Clnt) {
	d, err := c9.FStat(flag.Arg(1))
	if err != nil {
		log.Println("Error", err)
		return
	}

	if !*verbose {
		fmt.Printf("%v\n", d)
	} else {
		fmt.Printf("Name: %s\n", d.Name)
		fmt.Printf("UID/GID: [%x]%s / [%x]%s\n", d.Uidnum, d.Uid, d.Gidnum, d.Gid)
		fmt.Printf("Type: %x\n", d.Type)
		fmt.Printf("Dev: %x\n", d.Dev)
		fmt.Printf("Last Access: %v\n", d.Atime)
		fmt.Printf("Last Modified: %v\n", d.Mtime)
		fmt.Printf("last mod uid: [%x]%s\n", d.Muidnum, d.Muid)
		fmt.Printf("Extensions: %s\n", d.Ext)
	}

}

func cmdfcat(c9 *warp9.Clnt) {

	err := c9.Open(c9.Root, warp9.OREAD)
	if err != nil {
		log.Fatalf("open err:%v\n", err)
	}

	buf, err := c9.Read(c9.Root, uint64(0), uint32(10))
	if err != nil {
		log.Fatalf("Error:%v\n", err)
	}

	os.Stdout.Write(buf[0:])

}
