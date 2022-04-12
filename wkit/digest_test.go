package wkit

import (
	"encoding/binary"
	"hash"
	"testing"
)

type MyHash struct {
	hash.Hash
	val  uint32
	rslt uint32
}

var digest *DigestItem = NewDigestItem("MyDigest", &myhash)
var tst *testing.T
var myhash MyHash

// testing  -- test object without remote access
//
// - create a digest object
//   - pass in a hash function
// - call read several times
//   - verify results
//

func TestDigestItem(t *testing.T) {
	tst = t

	test(42)
	test(43)
	test(0xFFFFFFFF)
	test(0x7FFFFFFF)
	test(0)

	testbigbuf(42)
}

func test(val uint32) {
	tst.Logf("test: val=%d\n", val)
	myhash.val = val
	var buf []byte = make([]byte, 4)
	n, err := digest.Read(buf, 0, 4)
	if err != nil || n != 4 {
		tst.Errorf("Read: failed:%v", err)
		return
	}
	if !myhash.compare(buf) {
		tst.Errorf("Read: expected: %d, got: %d\n", myhash.val, myhash.rslt)
	}
}

func testbigbuf(val uint32) {
	tst.Logf("testbigbuf: val=%d\n", val)
	myhash.val = val
	var buf []byte = make([]byte, 1000)
	n, err := digest.Read(buf, 0, 4)
	if err != nil || n != 4 {
		tst.Errorf("Read: failed:%v", err)
		return
	}
	if !myhash.compare(buf[:n]) {
		tst.Errorf("Read: expected: %d, got: %d\n", myhash.val, myhash.rslt)
	}
}

func (h MyHash) Sum32() uint32 {
	return h.val
}

func (h MyHash) compare(buf []byte) bool {
	h.rslt = binary.LittleEndian.Uint32(buf)
	return h.val == h.rslt
}
