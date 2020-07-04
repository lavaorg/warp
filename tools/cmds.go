// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package tools

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/lavaorg/warp/warp9"
)


var CmdVerbose bool = false

func Cat(c9 *warp9.Clnt) {
	o, err := c9.Open(flag.Arg(1), warp9.OREAD)
	if err != nil {
		log.Fatalf("Error:%v\n", err)
	}
	defer o.Close()

	cat(o)
}

// cat: copy full contents of object to stdardout
func cat(o *warp9.Object) {

	buf := make([]byte, 8192)
	var err error = nil
	for {
		n, err := o.Read(buf)
		if n == 0 {
			break
		}
		if err != nil && err != warp9.WarpErrorEOF {
			     warp9.Error("Error reading:%v\n", err)
		}
		os.Stdout.Write(buf[0:n])
		if err == warp9.WarpErrorEOF {
			break
		}
	}

	if err != nil && err != warp9.WarpErrorEOF {
		log.Fatalf("Error:%v\n", err)
	}
}

// ls: read contents of diretory.
func Ls(c9 *warp9.Clnt) {

	n := flag.Arg(1)
	if n == "." || n == "/" {
		n = ""
	}

	fid, err := c9.Walk(n)
	if err != nil {
		warp9.Error("error:%v", err)
		return
	}
	defer c9.Clunk(fid)

	if fid.Qid.Type&warp9.QTDIR > 0 {
		// read directory
		f, err := c9.FOpenObject(fid, warp9.OREAD)
		if err != nil {
			warp9.Error("error:%v", err)
			return
		}

		for {
			d, err := f.Readdir(0)
			if d == nil || len(d) == 0 || err != nil {
				break
			}
			for i := 0; i < len(d); i++ {
				//os.Stdout.WriteString(d[i].Name + "\n")
				fmt.Printf("%v\n", d[i])
			}
		}
	} else {
		// stat the file
		d, err := c9.FStat(fid)
		if err != nil {
			log.Println("Error", err)
			return
		}
		fmt.Printf("%v\n", d)
	}
	if err != nil && err != warp9.WarpErrorEOF {
		log.Fatal(err)
	}
}

func Stat(c9 *warp9.Clnt) {
	d, err := c9.Stat(flag.Arg(1))
	if err != nil {
		log.Println("Error", err)
		return
	}

	if !CmdVerbose {
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

// get: perform a "Get" operation. open/read-fully/clunk. like cat.
func Get(c9 *warp9.Clnt) {
	data, qid, err := c9.Get(flag.Arg(1), 0)
	if err != nil {
		log.Fatalf("Error:%v\n", err)
	}
	if CmdVerbose {
		fmt.Printf("qid = %v\nData:\n", qid)
	}
	os.Stdout.Write(data[0:])
}

// write: read from stdin until eof and write to object
func Write(c9 *warp9.Clnt) {

	o, err := c9.Open(flag.Arg(1), warp9.OWRITE)
	if err != nil {
		warp9.Error("Open Error:%v\n", err)
		return
	}
	defer o.Close()

	c, e := io.Copy(o, os.Stdin)
	if e != nil {
		warp9.Error("Copy Error:%v\n", err)
		return
	}

	fmt.Printf("bytes written:%v\n", c)
}

// ctl: write remaining arguments to object. Aguments "Join"-ed with " " between.
func Ctl(c9 *warp9.Clnt) {
	o, err := c9.Open(flag.Arg(1), warp9.ORDWR)
	if err != nil {
		if err != warp9.WarpErrorNOTEXIST {
			warp9.Error("Error:%v\n", err)
		} else {
			fmt.Printf("object does not exist\n")
		}
		return
	}
	defer o.Close()

	//rest of command line to object
	cmd := strings.Join(flag.Args()[2:], " ")
	_, e := o.Write([]byte(cmd))
	if e != nil {
		warp9.Error("Error:%v\n", e)
	}
	// read results back
	cat(o)
}
