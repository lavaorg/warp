package wkit

import (
	"testing"
)

var baseItem *BaseItem = NewBaseItem("MyBase", false)

func TestBaseItemCreate(t *testing.T) {

	g := baseItem != nil
	g = g && (baseItem.opened == false)

	if !g {
		t.Error("bad baseItem")
	}
}

func TestBaseItemDir(t *testing.T) {
	d := baseItem.Dir
	d2 := baseItem.GetDir()
	g := (d == *d2)
	g = g && (d2.Name == "MyBase")

	if !g {
		t.Error("GetDir faulty")
	}
}

func TestBaseItemSet(t *testing.T) {
	m := baseItem.Dir.Mode
	baseItem.SetMode(0x777)
	g := (m != baseItem.Dir.Mode)
	g = g && (baseItem.Dir.Mode == 0x777)

	if !g {
		t.Error("Set failed")
	}

	o := baseItem.opened
	baseItem.SetOpen(!o)
	g = o != baseItem.opened
	if !g {
		t.Error("SetOpen faulty")
	}
	baseItem.opened = false
}

func TestBaseItemRead(t *testing.T) {
	//test read fails
	buf := []byte{' '}
	n, e := baseItem.Read(buf, 0, 1)
	if e == nil {
		t.Error("Read should fail")
	}
	// test read works
	baseItem.opened = true
	n, e = baseItem.Read(buf, 0, 1)
	if e != nil || n != 0 {
		t.Error("Read expected to return 0")
	}
	baseItem.opened = false
}

func TestBaseItemWrite(t *testing.T) {
	//test read fails
	buf := []byte{' '}
	n, e := baseItem.Write(buf, 0, 1)
	if e == nil {
		t.Error("Write should fail")
	}
	// test write return 0
	baseItem.opened = true
	n, e = baseItem.Write(buf, 0, 1)
	if e != nil || n != 0 {
		t.Error("write expected to return 0")
	}
	baseItem.opened = false
}

func TestBaseItemClunk(t *testing.T) {
	baseItem.opened = false
	e := baseItem.Clunk()
	if e != nil {
		t.Error("Clunk should not fail")
	}

	baseItem.opened = true
	e = baseItem.Clunk()
	if e != nil {
		t.Error("Clunk expected to succeed but failed")
	}
	if baseItem.opened {
		t.Error("Clunk should have set closed flag")
	}
	baseItem.opened = false
}
