package wkit

import (
	"encoding/binary"
	"hash"
	"testing"
)

//this implements all interfaces for info; along with testing fields
type infoHash struct {
	hash.Hash
	val  uint32
	rslt uint32
	buf  []byte
}

var (
	verhash    infoHash  = infoHash{buf: make([]byte, 20)}
	statushash infoHash  = infoHash{buf: make([]byte, 20)}
	info       Directory = NewInfoItem(&verhash, &statushash)
	vdigest    *DigestItem
	vdetail    *BytesItem
	sdigest    *DigestItem
	sdetail    *BytesItem
	itemtst    *testing.T
)

// testing  -- test object without remote access
//
// - create a info object
//   - walk to each leaf object
// - test both digest objects
// - test both details objects
//

func TestInfoItem(t *testing.T) {
	//warp9.LogDebug(true)
	itemtst = t

	vdigest = itemWalkTo([]string{"version", "digest"}).(*DigestItem)
	vdetail = itemWalkTo([]string{"version", "details"}).(*BytesItem)
	sdigest = itemWalkTo([]string{"status", "digest"}).(*DigestItem)
	sdetail = itemWalkTo([]string{"status", "details"}).(*BytesItem)

	digestTest(&verhash, 42, vdigest)
	digestTest(&statushash, 43, sdigest)

	detailTest(&verhash, "Hollis", vdetail)
	detailTest(&statushash, "New Hampshire", sdetail)

}

func itemWalkTo(path []string) Item {
	itm, err := info.Walk(path)
	if err != nil {
		return nil
	}
	return itm
}

func digestTest(ih *infoHash, val uint32, d *DigestItem) {
	itemtst.Logf("info-digestTest: val=%d\n", val)
	ih.val = val
	var buf []byte = make([]byte, 4)
	n, err := d.Read(buf, 0, 4)
	if err != nil || n != 4 {
		itemtst.Errorf("Read: failed:%v", err)
		return
	}
	rslt := binary.LittleEndian.Uint32(buf)
	if val != rslt {
		itemtst.Errorf("Read: expected: %d, got: %d\n", val, rslt)
	}
}

func detailTest(ih *infoHash, str string, d *BytesItem) {
	itemtst.Logf("info-detailTest: str=%v\n", str)
	ih.buf = []byte(str)
	var inbuf []byte = make([]byte, 20)
	n, err := d.Read(inbuf, 0, uint32(len(inbuf)))
	if err != nil {
		itemtst.Errorf("Read: failed:%v", err)
		return
	}
	if str != string(inbuf[:n]) {
		itemtst.Errorf("Read: expected: %d, got: %d\n", ih.buf, inbuf)
	}

}

//
// Implement the Digest / BytesSequence interfaces
//

func (h infoHash) Sum32() uint32 {
	return h.val
}

func (h infoHash) Bytes() []byte {
	return h.buf
}
