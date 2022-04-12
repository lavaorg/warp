package wkit

import (
	"errors"

	"github.com/lavaorg/warp/warp9"
)

const (
	infoName    = "info"
	versionName = "version"
	statusName  = "status"
)

type (

	// InfoItem provides a common method to present information about a service.
	// Each service can simply publish it's version information and it's current status indicators.
	// In order to make it easy for client to detect changes to either a simple 32-but unsigned digest
	// is presented for a quick check before querying details.  The structure if InfoItem is as follows:
	//
	//    /info
	//        /version
	//            /digest
	//            /details
	//        /status
	//            /digest
	//            /details
	//
	// digest -> 32bit unsigned value indicative of the contents (example CRC32)
	//
	// details -> sequence of bytes, encoding determimed by service.
	//
	// The warp tree maintained by InfoItem are all read-only elements.
	//
	InfoItem struct {
		*DirItem //promote the method-set
		version  InfoData
		status   InfoData
	}
)

// NewDirItem returns a pointer to a new Directory named 'info' that
// contains the /info tree of objects.
func NewInfoItem(verdata InfoData, statusdata InfoData) Directory {

	info := &InfoItem{
		DirItem: NewDirItem(infoName).(*DirItem),
	}
	// create version objects
	vdir := NewDirItem(versionName)
	vdir.AddItem(NewDigestItem("digest", verdata))
	vdir.AddItem(NewBytesItem("details", verdata))

	// create status objects
	sdir := NewDirItem(statusName)
	sdir.AddItem(NewDigestItem("digest", statusdata))
	sdir.AddItem(NewBytesItem("details", statusdata))

	info.DirItem.AddDirectory(vdir)
	info.DirItem.AddDirectory(sdir)

	return info
}

type InfoData interface {
	Digest32
	BytesSequence
}

// promoted from DirItem
// Directory Interface
//   Name()
//   Walk()
// Item Interface
//   GetDir
//   Parent
//   SetParent
//   GetQid
//   SetMode
//   Write
//   Open
//   Clunk
//   Remove
//   Stat
//   WStat
//   Read()
// DirItem methods
//   SetOpened()
//   SetUGMId()
//   SetAtime()
//   SetMTime()
//   SetQid()

//
// overrides of DirItem
//

// ResetBuffer has no side effects.
func (o *InfoItem) ResetBuffer() {
	return
}

// SetContent sets the map of strings -> Item mappings.
func (o *InfoItem) SetContent(content map[string]Item) {
	return
}

//
// Directory Interface
//

func (o *InfoItem) Name() string {
	return infoName
}

//
// Directory interface implementations
//

// AddDirectory returns with no side effects.  InfoItem is a fixed tree.
func (o *InfoItem) AddDirectory(newDir Directory) {
	return
}

// AddItem returns with no side effects. InfoItem is a fixed tree.
func (o *InfoItem) AddItem(item Item) {
	return
}

func (o *InfoItem) Children() map[string]Item {
	return o.Content
}

//RemoveItem returns with an Eperm error. Tree is fixed.
func (o *InfoItem) RemoveItem(item Item) error {
	return warp9.ErrorCode(warp9.Eperm)
}

//
// Item Interface implementations
//

// Return the object as the interface type Item.
func (o *InfoItem) GetItem() Item {
	return o
}

// IsDirectory returns itself.
func (o *InfoItem) IsDirectory() Directory {
	return o
}

// SetParent sets the back reference to the directory this object is in.
// setting to nil is not allowed.
func (o *InfoItem) SetParent(dir Directory) error {
	if dir == nil {
		return errors.New("cannot set nil directory as a parent")
	}
	o.parent = dir
	return nil
}

func (o *InfoItem) Walked() (Item, error) {
	return o, nil
}
