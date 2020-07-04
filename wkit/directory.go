package wkit

import (
	"errors"
	"log"
	"time"

	"github.com/lavaorg/warp/warp9"
)

type (
	// A special container object that allows objects to be attached and removed.
	Directory interface {
		Item
		Name() string
		Walk([]string) (Item, error)
		AddDirectory(Directory) Directory
		AddItem(Item) Directory
		Children() map[string]Item
		RemoveItem(Item) error
	}

	// Basic Directory object to hold other Items.
	DirItem struct {
		OneItem
		Content map[string]Item
		root    Directory
	}
)

// NewDirItem returns a pointer to a new DirItem
func NewDirItem(name string) Directory {
	return &DirItem{
		OneItem: OneItem{
			Dir: warp9.Dir{
				Qid:   warp9.Qid{warp9.QTDIR, 0, NextQid()},
				Mode:  (warp9.DMDIR | uint32(Perms(warp9.DMREAD, warp9.DMREAD, warp9.DMREAD))),
				Atime: uint32(time.Now().Unix()),
				Mtime: uint32(time.Now().Unix()),
				Name:  name,
				Uid:   1,
				Gid:   1,
				Muid:  1,
			},
			opened: false,
		},
		Content: make(map[string]Item, 0),
	}

}

func (d *DirItem) Children() map[string]Item {
	return d.Content
}

func (d *DirItem) GetDir() *warp9.Dir {
	return &d.Dir
}

func (d *DirItem) GetItem() Item {
	return &d.OneItem
}

func (d *DirItem) Parent() Directory {
	return d.parent
}

func (d *DirItem) SetParent(dir Directory) error {
	if dir == nil {
		return errors.New("cannot set nil directory as a parent")
	}
	d.parent = dir
	return nil
}

func (d *DirItem) GetQid() warp9.Qid {
	return d.Qid
}

func (d *DirItem) Write(ibuf []byte, off uint64, count uint32) (uint32, error) {
	return 0, nil
}

func (d *DirItem) Open(mode byte) (uint32, error) {
	d.opened = true
	return 0, nil
}

func (d *DirItem) Clunk() error {
	d.opened = false
	return nil
}

func (d *DirItem) Remove() error {
	return nil
}

func (d *DirItem) Stat() (*warp9.Dir, error) {
	return &d.Dir, nil
}

func (d *DirItem) WStat(dir *warp9.Dir) error {
	d.ResetBuffer()
	return nil
}

func (d *DirItem) SetName(name string) {
	d.Dir.Name = name
}

// SetContent sets the map of strings -> Item mappings.
func (d *DirItem) SetContent(content map[string]Item) {
	if content != nil {
		d.Content = content
	} else {
		d.Content = make(map[string]Item, 0)
	}
}

// SetUGMId sets the user, group and modified bits.
func (d *DirItem) SetUGMId(user, group, modified uint32) {
	d.Uid, d.Gid, d.Muid = user, group, modified
}

// SetAtime sets the access time.
func (d *DirItem) SetAtime(atime uint32) {
	d.Atime = atime
}

// SetMtime sets the last modified time.
func (d *DirItem) SetMTime(mtime uint32) {
	d.Mtime = mtime
}

// SetQid sets the QID for this item.
func (d *DirItem) SetQid(qid warp9.Qid) {
	d.Qid = qid
}

// SetMode sets the mode bits for this item.
func (d *DirItem) SetMode(mode uint32) Item {
	d.Mode = mode
	return d
}

// IsDirectory returns itself.
func (d *DirItem) IsDirectory() Directory {
	return d
}

func (d *DirItem) Walked() (Item, error) {
	return d, nil
}

func (d *DirItem) Name() string {
	return d.Dir.Name
}

// AddDirectory sets the passed in directory as a child of this one.
func (d *DirItem) AddDirectory(newDir Directory) Directory {
	if newDir == nil {
		log.Println("newDir is nil. Directory not set as child.")
		return d
	}

	if err := newDir.SetParent(d); err != nil {
		log.Println(err.Error())
	}

	return d.AddItem(newDir)
}

func (d *DirItem) AddItem(item Item) Directory {
	if err := item.SetParent(d); err != nil {
		log.Println(err.Error())
		return d
	}

	ndir := item.GetDir()
	// Inherit directory user, group and modified bits.
	ndir.Uid, ndir.Gid, ndir.Muid = d.Uid, d.Gid, d.Muid
	ndir.Atime = uint32(time.Now().Unix())
	ndir.Mtime = d.Atime

	// Set this item to a string in the content map.
	d.Content[ndir.Name] = item
	d.ResetBuffer()
	return d
}

//TODO handle concurrent access
func (d *DirItem) RemoveItem(item Item) error {
	if d.Mode&uint32(Perms(warp9.DMWRITE, 0, 0)) == 0 {
		warp9.Error("No permission to remove item. %#v", d)
		return warp9.ErrorCode(warp9.Eperm)
	}
	ndir := item.GetDir()

	if _, found := d.Content[ndir.Name]; !found {
		return errors.New("item not found")
	}
	delete(d.Content, ndir.Name)
	d.ResetBuffer()
	warp9.Debug("Item removed:%s", ndir.Name)
	return nil
}

// walk the content items looking for a match for path. each element in path
// except the last must be a directory.
// return the found Item or an error
func (d *DirItem) Walk(path []string) (Item, error) {

	warp9.Debug("d.Walk:%T.%v, path:%v", d, d.Dir.Name, path)
	if len(path) < 1 {
		// empty path succeeds in finding self
		return d, nil
	}

	if len(path) == 1 {
		// leaf item return if found
		n := path[0]
		var item Item
		if n == ".." {
			if d.Parent() != nil {
				item = d.parent.(Item)
			} else {
				return nil, warp9.ErrorCode(warp9.Enotexist)
			}
		} else {
			item = d.Content[n]
		}
		if item == nil {
			return nil, warp9.ErrorCode(warp9.Enotexist)
		}
		// let item know its walked to, it may return a clone
		warp9.Debug("d.Walk: %T.%v", item, item)
		return item.Walked()
	}

	// element must be a directory to further walk
	elem := path[0]
	path = path[1:]
	var item Item
	var dir Directory
	if elem == ".." {
		dir = d.Parent()
	} else {
		item := d.Content[elem]
		if item == nil {
			return nil, warp9.ErrorCode(warp9.Enotexist)
		}
		dir = item.IsDirectory()
	}
	if dir == nil {
		return item, warp9.ErrorCode(warp9.Enotdir)
	}
	// walk to next dir
	warp9.Debug("d.Walk:next dir:%T.%v, path:%v", dir, dir.Name(), path)
	return dir.Walk(path)
}

func (d *DirItem) ResetBuffer() {
	//TODO Lock
	d.buffer = nil
	return
}

// Set the open state of the directory
func (d *DirItem) SetOpened(o bool) bool {
	b := d.opened
	d.opened = o
	if !o {
		d.buffer = nil
	}
	return b
}

// Return the requested byte sequence of the Directory contents.
// The byte-buffer is the byte representation of all the current
// object's Dir entry. (see Warp9.Stat for representation)
func (d *DirItem) Read(obuf []byte, off uint64, rcount uint32) (uint32, error) {
	// walk all contents; get Dir structure; pack as bytes
	if d.buffer == nil {
		d.buffer = make([]byte, 0, 300)
		for _, item := range d.Content {
			buf := warp9.PackDir(item.GetDir())
			d.buffer = append(d.buffer, buf...)
			warp9.Debug("d.Read: dir item:%v, len(buf):%v, len(buffer):%v", item, len(buf), len(d.buffer))
		}
	}

	// determine which and how many bytes to return
	var count uint32
	switch {
	case off > uint64(len(d.buffer)):
		count = 0
	case uint32(len(d.buffer[off:])) > rcount:
		count = rcount
	default:
		count = uint32(len(d.buffer[off:]))
	}
	copy(obuf, d.buffer[off:uint32(off)+count])
	warp9.Debug("d.Read:buffer:%v, obuf: %v, off:%v, rcount:%v\n", len(d.buffer), len(obuf), off, count)

	return count, nil
}
