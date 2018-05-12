// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package wkit

import (
	"time"

	"github.com/lavaorg/lrt/mlog"
	"github.com/lavaorg/warp/warp9"
)

// item represents a generic object in a hierarchical tree.  This interface
// will allow the object server to perform generic operations on any object.
// Where necessary it can learn more details of what the object is (e.g. a directory,
// or a bind-point, etc)
type Item interface {
	GetDir() *warp9.Dir
	GetItem() Item
	Data() []byte
	IsDirectory() Directory
	Opened() bool
	SetOpened(b bool) bool
	Parent() Directory
	SetParent(d Directory) error
	Read(obuf []byte, off uint64, rcount uint32) (uint32, error)
	Write(ibuf []byte, off uint64, count uint32) (uint32, error)
}

// A special container object that allows objects to be attached and removed.
type Directory interface {
	CreateItem(item Item, name string, perm uint32, mode uint8) (Item, error)
	Walk(path []string) (Item, error)
}

// OneItem is a generic in-memory blob object. The contents of the
// object are arbitrary bytes and can be written/read.
type OneItem struct {
	warp9.Dir
	parent Directory
	opened bool
	Buffer []byte
}

// Basic Directory object to hold other Items.
type DirItem struct {
	OneItem
	Content map[string]Item
}

// Return the object's Dir structure
func (o *OneItem) GetDir() *warp9.Dir {
	return &o.Dir
}

// Return the object as the interface type Item.
func (o *OneItem) GetItem() Item {
	return o
}

// Return the byte data of the object
func (o *OneItem) Data() []byte {
	return o.Buffer
}

// Indicate that the object is not a directory, returns nil.
func (o *OneItem) IsDirectory() Directory {
	return nil
}

// set the open state to true or false; return previous state.
func (o *OneItem) SetOpened(v bool) bool {
	b := o.opened
	o.opened = v
	return b
}

// Return open status.
func (o *OneItem) Opened() bool {
	return o.opened
}

// Return the object's parent (eg container).
func (o *OneItem) Parent() Directory {
	return o.parent
}

// Set the parent of the object.
func (o *OneItem) SetParent(d Directory) error {
	o.parent = d
	return nil
}

// Return the requested set of bytes from the object's byte buffer.
func (o *OneItem) Read(obuf []byte, off uint64, rcount uint32) (uint32, error) {

	// determine which and how many bytes to return
	var count uint32
	switch {
	case off > uint64(len(o.Buffer)):
		count = 0
	case uint32(len(o.Buffer[off:])) > rcount:
		count = rcount
	default:
		count = uint32(len(o.Buffer[off:]))
	}
	n := copy(obuf, o.Buffer[off:uint32(off)+count])
	if uint32(n) != count {
		return 0, warp9.Eio
	}
	mlog.Debug("d.Buffer:%v, obuf: %v, off:%v, rcount:%v\n", len(o.Buffer), len(obuf), off, count)
	mlog.Debug("o:%T %p", o, o)
	return count, nil

}

// Write bytes into the object's byte buffer.
// The object's byte buffer is smaller than len(max(int))
func (o *OneItem) Write(ibuf []byte, off uint64, count uint32) (uint32, error) {

	// our file will not be super large;  convert everything to int
	ioff := int(off)
	icnt := int(count)
	if uint64(ioff) != off || uint32(icnt) != count {
		return 0, warp9.Etoolarge
	}

	// if append file always just append
	// if offset is the current len; just append
	if check(o.Mode, warp9.DMAPPEND) || ioff == len(o.Buffer) {
		o.Buffer = append(o.Buffer, ibuf[:icnt]...)
		return count, nil
	}
	// if offset < cur len; truncate current and append
	if ioff < len(o.Buffer) {
		o.Buffer = append(o.Buffer[:off], ibuf[:icnt]...)
		return count, nil
	}

	// if we are seeking past eof then add 0's first
	if ioff >= len(o.Buffer) {
		zsz := ioff - len(o.Buffer) - 1
		z := make([]byte, zsz, zsz+icnt)
		z = append(z, ibuf[:icnt]...)
		o.Buffer = append(o.Buffer, z...)
		return count, nil
	}

	return 0, warp9.Eio
}

//
// DirItem methods
//

//
// create a new empty directory
//
func NewDirItem() *DirItem {
	d := new(DirItem)
	d.Content = make(map[string]Item, 0)
	return d
}

// Not used for DirItem. Returns nil
func (d *DirItem) GetData() []byte {
	return nil
}

// Always returns self. Is a Directory.
func (d *DirItem) IsDirectory() Directory {
	return d
}

// Places the Item object into the directory with the given attributes.
// Return the item,nil for success
// Return nil,error otherwise
func (d *DirItem) CreateItem(item Item, name string, perm uint32, mode uint8) (Item, error) {

	var i Item
	nqid := warp9.Qid{warp9.QTOBJ, 0, NextQid()}

	if perm&warp9.DMDIR > 0 {
		i = NewDirItem()
		nqid.Type = warp9.QTDIR
	} else {
		if item == nil {
			i = new(OneItem)
		} else {
			i = item
		}
	}
	i.SetParent(d)
	ndir := i.GetDir()
	ndir.Name = name
	ndir.Qid = nqid
	ndir.Uid = d.Uid
	ndir.Gid = d.Gid
	ndir.Muid = ndir.Uid

	ndir.Atime = uint32(time.Now().Unix())
	ndir.Mtime = d.Atime

	ndir.Mode = perm

	d.Content[name] = i
	return i, nil
}

// walk the content items looking for a match for path. each element in path
// except the last must be a directory.
// return the found Item or an error
func (d *DirItem) Walk(path []string) (Item, error) {

	if len(path) < 1 {
		// empty path succeeds in finding self
		return d, nil
	}

	if len(path) == 1 {
		// leaf item return if found
		n := path[0]
		var item Item
		if n == ".." {
			//ToDo: should check for current Conn's root!!
			item = d.parent.(Item)
		} else {
			item = d.Content[n]
		}
		if item == nil {
			return nil, warp9.Enotexist
		}
		return item, nil
	}

	// element must be a diretory to further walk
	elem := path[0]
	path = path[1:]
	var item Item
	var dir Directory
	if elem == ".." {
		dir = d.Parent()
	} else {
		item := d.Content[elem]
		if item == nil {
			return nil, warp9.Enotexist
		}
		dir = item.IsDirectory()
	}
	if dir == nil {
		return item, warp9.Enotdir
	}
	// walk to next dir
	return dir.Walk(path)
}

// Set the open state of the directory
func (d *DirItem) SetOpened(o bool) bool {
	b := d.opened
	d.opened = o
	if !o {
		d.Buffer = nil
	}
	return b
}

// Return the requested byte seqence of the Directory contents.
// The bytebuffer is the byte representation of all the current
// object's Dir entry. (see Warp9.Stat for representation)
func (d *DirItem) Read(obuf []byte, off uint64, rcount uint32) (uint32, error) {

	// walk all contents; get Dir structure; pack as bytes
	if d.Buffer == nil {
		d.Buffer = make([]byte, 0, 300)
		for _, item := range d.Content {
			buf := warp9.PackDir(item.GetDir())
			d.Buffer = append(d.Buffer, buf...)
			mlog.Debug("dir item:%v, len(buf):%v, len(Buffer):%v", item, len(buf), len(d.Buffer))
		}
	}

	// determine which and how many bytes to return
	var count uint32
	switch {
	case off > uint64(len(d.Buffer)):
		count = 0
	case uint32(len(d.Buffer[off:])) > rcount:
		count = rcount
	default:
		count = uint32(len(d.Buffer[off:]))
	}
	copy(obuf, d.Buffer[off:uint32(off)+count])
	mlog.Debug("d.BUffer:%v, obuf: %v, off:%v, rcount:%v\n", len(d.Buffer), len(obuf), off, count)

	return count, nil
}

//
// private helper functions
//

func check(mode, kind uint32) bool {
	return (mode&kind > 0)
}
