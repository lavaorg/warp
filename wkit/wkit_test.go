// Copyright 2009 The go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wkit

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/lavaorg/warp/tools"
	"github.com/lavaorg/warp/warp9"
)

const TESTDATA string = "Test Data"

var srvport int = 0
var gMount *warp9.Clnt
var tracelevel = 0 //0,1,2,3
var root Directory

func init() {
	//warp9.LogDebug(true)

	root = NewDirItem("/")
	startServer()

	// Test the mount
	var err error
	gMount, err = mountServer()
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to server: %s\n", err))
	}
}

// return the single root; each test can add to root
func getRoot() Directory {
	return root
}

func startServer() {
	// Get a random port between 1025 and 9999
	srvport = rand.Intn(9999-1025) + 1025
	go serverThread(srvport)
	// Wait for the server to start
	time.Sleep(2 * time.Second)
	return
}

func mountServer() (*warp9.Clnt, error) {
	// Test the mount
	user := warp9.Identity.User(1)
	addr := "127.0.0.1:" + strconv.Itoa(srvport)
	fmt.Printf("Mounting to %s\n", addr)
	c9, err := warp9.Mount("tcp", addr, "/", 8192, user)
	if c9 != nil {
		c9.Debuglevel = tracelevel
	}
	return c9, err
}

func serverThread(port int) {

	srv := NewServer("test server", tracelevel, root)
	if srv.Start(srv) == false {
		panic(fmt.Errorf("Unable to start server\n"))
	}
	fmt.Printf("Starting Server on 127.0.0.1:%d\n", port)
	srv.StartNetListener("tcp", "127.0.0.1:"+strconv.Itoa(port))
}

func testWriteOp(mount *warp9.Clnt) error {
	input := TESTDATA
	ir := strings.NewReader(input)

	o, err := mount.Open("/testdir/testfile", warp9.OWRITE)
	if err != nil {
		return err
	}
	defer o.Close()

	_, e := io.Copy(o, ir)
	if e != nil {
		return e
	}

	return nil
}

func testReadOp(mount *warp9.Clnt) error {
	obj, err := mount.Open("/testdir/testfile", warp9.OREAD)
	if err != nil {
		fmt.Errorf("Failed to open file: %s\n", err)
		return err
	}

	data, err2 := mount.Read(obj.Fid, 0, obj.Fid.Iounit)
	if err2 != nil {
		fmt.Errorf("Failed to Read file: %s\n", err2)
		return err2
	}
	err = mount.Clunk(obj.Fid)

	sdata := string(data)
	fmt.Printf("%s\n", sdata)
	if sdata != TESTDATA {
		return fmt.Errorf("Data mismatch, Read fails: %s != %s\n", sdata, TESTDATA)
	}
	return nil
}

func TestBaseItem(t *testing.T) {

	root.AddItem(NewBaseItem("base", false))

	obj, err := gMount.Open("/base", warp9.OREAD)
	if err != nil {
		t.Error("Open base failed")
		return
	}
	err = gMount.Clunk(obj.Fid)
	if err != nil {
		t.Error("Clunk failed:")
	}
}

// TestServerOps will setup a tree like:
//   /events        (EventItem)
//   /base          (BaseItem)
//   /testdir       (DirItem)
//      /testfile   (OneItem)
//   /info          (InfoItem)
//
func TestServerOps(t *testing.T) {
	fmt.Printf("Testing server Operations\n")
	defer gMount.Clunk(gMount.Root)

	testfile := NewItem("testfile")
	testdir := NewDirItem("testdir")
	testdir.AddItem(testfile)

	events := NewEventItem("events")
	bi := NewBaseItem("base", false)

	root := getRoot()

	root.AddItem(events)
	root.AddItem(bi)
	root.AddDirectory(testdir)
	root.AddDirectory(NewInfoItem(&wkitInfoHash, &wkitInfoStatus))

	t.Logf("test:write")
	err := testWriteOp(gMount)
	if err != nil {
		t.Error("Failed to write object")
		return
	}
	t.Logf("test:read")
	err = testReadOp(gMount)
	if err != nil {
		t.Error("Failed to read object")
		return
	}

	t.Logf("test:/info/version")
	err = testInfoVer(gMount)
	if err != nil {
		t.Errorf("Failed to correctly read version: %v\n", err)
		return
	}
	t.Logf("test:/info/status")
	err = testInfoStatus(gMount)
	if err != nil {
		t.Errorf("Failed to correctly read status:%v\n", err)
		return
	}

	testLs(gMount)
	testCat(gMount)
}

//
// test for inf/version/
//

var wkitInfoHash infoHash = infoHash{val: 42, buf: []byte("mile")}
var wkitInfoStatus infoHash = infoHash{val: 43, buf: []byte("marathon")}

func testInfoVer(mount *warp9.Clnt) error {
	obj, err := mount.Open("/info/version/digest", warp9.OREAD)
	if err != nil {
		return err
	}
	var buf []byte = make([]byte, 4)
	n, err := obj.Read(buf[:])
	if err != nil {
		return err
	}
	if n != 4 {
		return fmt.Errorf("info/version/digest: bad size:%d\n", n)
	}
	val := binary.LittleEndian.Uint32(buf)
	if val != 42 {
		return fmt.Errorf("inf/version/digest: bad value want 42 got:%d\n", val)
	}

	return nil
}

func testInfoStatus(mount *warp9.Clnt) error {
	obj, err := mount.Open("/info/status/digest", warp9.OREAD)
	if err != nil {
		return err
	}
	var buf []byte = make([]byte, 4)
	n, err := obj.Read(buf[:])
	if err != nil {
		return err
	}
	if n != 4 {
		return fmt.Errorf("info/status/digest: bad size:%d\n", n)
	}
	val := binary.LittleEndian.Uint32(buf)
	if val != 43 {
		return fmt.Errorf("info/status/digest: bad value want 43 got:%d\n", val)
	}

	return nil
}

func testLs(mount *warp9.Clnt) error {
	//warp9.LogDebug(true)
	fmt.Println("---- testLs ----")
	fmt.Println("expect: /version <cr> /status")
	tools.Ls(mount, "/info")
	fmt.Println("\nexpect: /digest <cr> /details")
	tools.Ls(mount, "/info/version")
	fmt.Println("\nexpect: a-w--w--w- digest [0:0:0] 0 Q(o.0.13) at[1612787555] mt[1612787555]")
	tools.Ls(mount, "/info/version/digest")
	fmt.Println("\nexpect: a-w--w--w- details [0:0:0] 0 Q(o.0.14) at[1612787555] mt[1612787555]")
	tools.Ls(mount, "/info/version/details")
	fmt.Println("\n---- done ----")
	return nil
}

func testCat(mount *warp9.Clnt) error {
	fmt.Println("---- testCat ----")
	tools.Cat(mount, "/info/version/digest")
	fmt.Println("\nexpect: mile")
	tools.Cat(mount, "/info/version/details")
	fmt.Println("\n---- done ----")
	return nil
}
