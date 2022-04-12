package wkit

import (
	"errors"
	"time"

	"github.com/lavaorg/warp/warp9"
)

type (

	// Basic Directory object to hold other Items.  Create a buffer to use
	// during a read of the directory (e.g. the ls function).
	DirItem struct {
		*BaseItem
		Content map[string]Item
		root    Directory
		buffer  []byte
	}
)

// NewDirItem returns a pointer to a new DirItem
func NewDirItem(name string) Directory {

	dir := &DirItem{
		BaseItem: NewBaseItem(name, true),
		Content:  make(map[string]Item, 0),
		root:     nil,
		buffer:   make([]byte, 0),
	}
	return dir
}

// SetParent sets the back reference to the directory this object is in.
// setting to nil is not allowed.
func (d *DirItem) SetParent(dir Directory) error {
	if dir == nil {
		return errors.New("cannot set nil directory as a parent")
	}
	d.parent = dir
	return nil
}

//
// specific to DirItem
//

func (d *DirItem) ResetBuffer() {
	//TODO Lock
	d.buffer = nil
	return
}

// Set the open state of the directory, return prev state.
func (d *DirItem) SetOpened(o bool) bool {
	b := d.opened
	d.opened = o
	if !o {
		d.buffer = nil
	}
	return b
}

// don't think we should do this.  remove in future; could not find any uses
//func (d *DirItem) SetName(name string) {
//	d.Dir.Name = name
//}

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

//
// Directory Interface
//

func (d *DirItem) Name() string {
	return d.Dir.Name
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

// AddDirectory sets the passed in directory as a child of this one.
func (d *DirItem) AddDirectory(newDir Directory) {
	if newDir == nil {
		warp9.Info("newDir is nil. Directory not set as child.")
		return
	}

	if err := newDir.SetParent(d); err != nil {
		warp9.Error(err.Error())
	}

	d.AddItem(newDir)
}

func (d *DirItem) AddItem(item Item) {
	if err := item.SetParent(d); err != nil {
		warp9.Error(err.Error())
		return
	}

	ndir := item.GetDir()
	// Inherit directory user, group and modified bits.
	ndir.Uid, ndir.Gid, ndir.Muid = d.Uid, d.Gid, d.Muid
	ndir.Atime = uint32(time.Now().Unix())
	ndir.Mtime = d.Atime

	// Set this item to a string in the content map.
	d.Content[ndir.Name] = item
	d.ResetBuffer()
	return
}

func (d *DirItem) Children() map[string]Item {
	return d.Content
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

//
// Item Interface
//

// pick up from BaseItem
// GetDir
// Parent
// SetParent
// GetQid
// SetMode
// Write
// Open
// Clunk
// Remove
// Stat
// WStat

// Return the object as the interface type Item.
func (d *DirItem) GetItem() Item {
	return d
}

// IsDirectory returns itself.
func (d *DirItem) IsDirectory() Directory {
	return d
}

func (d *DirItem) Walked() (Item, error) {
	return d, nil
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
