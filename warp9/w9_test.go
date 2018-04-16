// Copyright 2009 The go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package warp9

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"testing"
)

const numDir = 20 //16384

var addr = flag.String("addr", ":5640", "network address")
var pipefsaddr = flag.String("pipefsaddr", ":5641", "pipefs network address")
var debug = flag.Int("debug", 0, "print debug messages")
var root = flag.String("root", "/", "root filesystem")

// Two files, dotu was true.
var testunpackbytes = []byte{
	79, 0, 0, 0, 0, 0, 0, 0, 0, 228, 193, 233, 248, 44, 145, 3, 0, 0, 0, 0, 0, 164, 1, 0, 0, 0, 0, 0, 0, 47, 117, 180, 83, 102, 3, 0, 0, 0, 0, 0, 0, 6, 0, 112, 97, 115, 115, 119, 100, 4, 0, 110, 111, 110, 101, 4, 0, 110, 111, 110, 101, 4, 0, 110, 111, 110, 101, 0, 0, 232, 3, 0, 0, 232, 3, 0, 0, 255, 255, 255, 255, 78, 0, 0, 0, 0, 0, 0, 0, 0, 123, 171, 233, 248, 42, 145, 3, 0, 0, 0, 0, 0, 164, 1, 0, 0, 0, 0, 0, 0, 41, 117, 180, 83, 195, 0, 0, 0, 0, 0, 0, 0, 5, 0, 104, 111, 115, 116, 115, 4, 0, 110, 111, 110, 101, 4, 0, 110, 111, 110, 101, 4, 0, 110, 111, 110, 101, 0, 0, 232, 3, 0, 0, 232, 3, 0, 0, 255, 255, 255, 255,
}

func TestUnpackDir(t *testing.T) {
	b := testunpackbytes
	for len(b) > 0 {
		var err error
		if _, b, _, err = UnpackDir(b, true); err != nil {
			t.Fatalf("Unpackdir: %v", err)
		}
	}
}

func TestAttachOpenReaddir1(t *testing.T) {
	var err error
	flag.Parse()

	// assume a ufs server is running

	/* this may take a few tries ... */
	var conn net.Conn
	for i := 0; i < 16; i++ {
		if conn, err = net.Dial("tcp", *addr); err != nil {
			t.Logf("Try go connect, %d'th try, %v", i, err)
		} else {
			t.Logf("Got a conn %d'th try: %v\n", i, conn)
			break
		}
	}
	if err != nil {
		t.Fatalf("Connect failed after many tries ...")
	}

	root := OsUsers.Uid2User(0)

	dir, err := ioutil.TempDir("/tmp/test9p/", "go9p")
	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	defer os.RemoveAll(dir)

	// Now create files in directory to test readdir
	t.Logf("host-dirname:%v", dir)
	for i := 0; i < 10; i++ {
		f := fmt.Sprintf(path.Join(dir, fmt.Sprintf("%d", i)))
		if err := ioutil.WriteFile(f, []byte(f), 0600); err != nil {
			t.Fatalf("Create %v: got %v, want nil", f, err)
		}
	}

	// mount the remote fs
	clnt, err := Mount("tcp", *addr, "", 8192, root)
	if err != nil {
		t.Fatalf("Connect failed: %v\n", err)
	}
	defer clnt.Unmount()
	t.Logf("attached, rootfid %v\n", clnt.Root)

	// readdir and File
	var dirfile *File
	if dirfile, err = clnt.FOpen(path.Base(dir), OREAD); err != nil {
		t.Fatalf("%v", err)
	}
	i := 0
	var d []*Dir
	for i < 10 {
		if d, err = dirfile.Readdir(11); err != nil {
			t.Fatalf("%v", err)
		}
		i += len(d)
	}
	if i != 10 {
		t.Fatalf("Readdir %v: got %d,%v entries, wanted %d", dir, i, len(d), 10)
	}

}

func TestAttachOpenReaddir2(t *testing.T) {
	var err error
	flag.Parse()

	// assume a ufs server is running

	/* this may take a few tries ... */
	var conn net.Conn
	for i := 0; i < 16; i++ {
		if conn, err = net.Dial("tcp", *addr); err != nil {
			t.Logf("Try go connect, %d'th try, %v", i, err)
		} else {
			t.Logf("Got a conn, %v\n", conn)
			break
		}
	}
	if err != nil {
		t.Fatalf("Connect failed after many tries ...")
	}

	root := OsUsers.Uid2User(0)

	dir, err := ioutil.TempDir("/Users/larry/foo/", "go9p")
	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	defer os.RemoveAll(dir)

	// Now create a whole bunch of files to test readdir
	for i := 0; i < numDir; i++ {
		f := fmt.Sprintf(path.Join(dir, fmt.Sprintf("%d", i)))
		if err := ioutil.WriteFile(f, []byte(f), 0600); err != nil {
			t.Fatalf("Create %v: got %v, want nil", f, err)
		}
	}

	var clnt *Clnt
	//for i := 0; i < 16; i++ {
	//	clnt, err = Mount("tcp", *addr, "", 8192, root)
	//}
	clnt, err = Mount("tcp", *addr, "", 8192, root)
	if err != nil {
		t.Fatalf("Connect failed: %v\n", err)
	}
	defer clnt.Unmount()
	t.Logf("attached, rootfid %v\n", clnt.Root)

	dirfid := clnt.FidAlloc()
	if _, err = clnt.Walk(clnt.Root, dirfid, []string{"."}); err != nil {
		t.Fatalf("%v", err)
	}
	if err = clnt.Open(dirfid, 0); err != nil {
		t.Fatalf("%v", err)
	}
	var b []byte
	var i, amt int
	var offset uint64
	for i < numDir {
		if b, err = clnt.Read(dirfid, offset, 64*1024); err != nil {
			t.Fatalf("%v", err)
		}
		for b != nil && len(b) > 0 {
			if _, b, amt, err = UnpackDir(b, false); err != nil {
				break
			} else {
				i++
				offset += uint64(amt)
			}
		}
	}
	if i != numDir {
		t.Fatalf("Reading %v: got %d entries, wanted %d", dir, i, numDir)
	}
	fmt.Println("Created directories")
	// Alternate form, using readdir and File
	var dirfile *File
	if dirfile, err = clnt.FOpen(".", OREAD); err != nil {
		t.Fatalf("%v", err)
	}
	i, amt, offset = 0, 0, 0
	for i < numDir {
		if d, err := dirfile.Readdir(numDir); err != nil {
			t.Fatalf("%v", err)
		} else {
			i += len(d)
		}
	}
	if i != numDir {
		t.Fatalf("Readdir %v: got %d entries, wanted %d", dir, i, numDir)
	}

	// now test partial reads.
	// Read 128 bytes at a time. Remember the last successful offset.
	// if UnpackDir fails, read again from that offset
	t.Logf("NOW TRY PARTIAL")
	i, amt, offset = 0, 0, 0
	for {
		var b []byte
		var d *Dir
		if b, err = clnt.Read(dirfid, offset, 128); err != nil {
			t.Fatalf("%v", err)
		}
		if len(b) == 0 {
			break
		}
		t.Logf("b %v\n", b)
		for b != nil && len(b) > 0 {
			t.Logf("len(b) %v\n", len(b))
			if d, b, amt, err = UnpackDir(b, true); err != nil {
				// this error is expected ...
				t.Logf("unpack failed (it's ok!). retry at offset %v\n",
					offset)
				break
			} else {
				t.Logf("d %v\n", d)
				offset += uint64(amt)
			}
		}
	}

	t.Logf("NOW TRY WAY TOO SMALL")
	i, amt, offset = 0, 0, 0
	for {
		var b []byte
		if b, err = clnt.Read(dirfid, offset, 32); err != nil {
			t.Logf("dirread fails as expected: %v\n", err)
			break
		}
		if offset == 0 && len(b) == 0 {
			t.Fatalf("too short dirread returns 0 (no error)")
		}
		if len(b) == 0 {
			break
		}
		// todo: add entry accumulation and validation here..
		offset += uint64(len(b))
	}
}

var f *File
var b = make([]byte, 1048576/8)

// Not sure we want this, and the test has issues. Revive it if we ever find a use for it.
