// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package wkit

import (
	"net"

	"github.com/lavaorg/warp/warp9"

	"github.com/lavaorg/lrt/mlog"
)

type Mounter interface {
	IsDirectory() bool
}

// A MountPoint attaches a remote object server to the local namespace
// The MountPoint is another directory node locally and will allow the
// Walk operation into this directory. The MoutPoint will then proxy all
// operations to the remote object server.
// A MountPoint acts as a Warp9 Client to the remote object server.
type MountPoint struct {
	warp9.Dir
	parent Directory
	mi     *MountInfo
	fid    *warp9.Fid
}

// holder for the remote mount information
type MountInfo struct {
	Aname  string      //mount point name
	ntype  string      //network type, per net.Dial
	addr   string      //network address, per net.Dial
	msize  uint32      //the message size to negotiate
	user   warp9.User  //warp9 user id
	dialer *net.Dialer //dial options
	conn   net.Conn    //connection, only non-nil if pre-exsited before mount
	clnt   *warp9.Clnt //warp9 remote mount

}

// Attempt to establish a mount of a remote object server,
// upon success return a valid local MountPoint to be placed in the local
// namespace
func MountPointDial(ntype, addr, aname string, msize uint32, user warp9.User) (*MountPoint, error) {
	mi := &MountInfo{aname, ntype, addr, msize, user, &net.Dialer{}, nil, nil}
	var err warp9.W9Err
	if msize == 0 {
		msize = warp9.MSIZE
	} else if msize < warp9.IOHDRSZ {
		msize = warp9.IOHDRSZ
	}

	mi.clnt, err = warp9.Mount(ntype, addr, aname, msize, user)
	if err != warp9.Egood {
		return nil, err
	}
	mi.msize = mi.clnt.Msize
	mt := &MountPoint{warp9.Dir{}, nil, mi, mi.clnt.Root}
	return mt, nil
}

// Attempt to establish a mount of a remote object serer, using the
// Dialer attributes passed, upon success return a valid local MountPoint
// to be placed in the local namespace.
func MountPointDialer(dialer net.Dialer, ntype, addr, aname string, msize uint32, user warp9.User) (*MountPoint, error) {
	mi := &MountInfo{aname, ntype, addr, msize, user, &dialer, nil, nil}
	var err warp9.W9Err

	c, e := mi.dialer.Dial(ntype, addr)
	if e != nil {
		return nil, warp9.Edial
	}

	if msize == 0 {
		msize = warp9.MSIZE
	} else if msize < warp9.IOHDRSZ {
		msize = warp9.IOHDRSZ
	}

	mi.clnt, err = warp9.MountConn(c, aname, msize, user)
	if err != warp9.Egood {
		return nil, err
	}
	mi.msize = mi.clnt.Msize
	mt := &MountPoint{warp9.Dir{}, nil, mi, mi.clnt.Root}

	return mt, nil
}

// Use the provided established net.Conn to mount a remote object server,
// upon success return a valid local MountPoint to be placed in the local
// namespace
func MountPointConn(conn net.Conn, aname string, msize uint32, user warp9.User) (*MountPoint, error) {
	mi := &MountInfo{aname, "", "", msize, user, &net.Dialer{}, conn, nil}
	var err warp9.W9Err

	if msize == 0 {
		msize = warp9.MSIZE
	} else if msize < warp9.IOHDRSZ {
		msize = warp9.IOHDRSZ
	}

	mi.clnt, err = warp9.MountConn(mi.conn, aname, msize, user)
	if err != warp9.Egood {
		return nil, err
	}
	mi.msize = mi.clnt.Msize
	mt := &MountPoint{warp9.Dir{}, nil, mi, mi.clnt.Root}

	return mt, nil
}

func (mt *MountPoint) Unmount() {
	cli := mt.mi.clnt
	mt.mi.clnt = nil
	if cli != nil {
		cli.Unmount()
	}
}

//
// implement the Directory interface
//

func (mt *MountPoint) Name() string {
	return mt.mi.Aname
}

// Create a new object in the directory associated with mt. Open the object
// according to mode and associate the new object with the current fid of mt.
func (mt *MountPoint) Create(name string, perm uint32, mode uint8) (Item, error) {

	err := mt.mi.clnt.FCreate(mt.fid, name, perm, mode, "")
	if err != warp9.Egood {
		return nil, err
	}

	return mt, nil
}

func (mt *MountPoint) Walk(path []string) (Item, error) {
	mlog.Info("mt.Walk:%v", path)
	newfid := mt.mi.clnt.FidAlloc()
	newfid.User = mt.fid.User
	qid, err := mt.mi.clnt.FWalk(mt.fid, newfid, path)
	if err != warp9.Egood {
		return nil, err
	}

	newfid.Qid = *qid
	newfid.Iounit = mt.fid.Iounit

	var newmt MountPoint
	newmt = *mt
	newmt.Qid = *qid
	newmt.fid = newfid
	mlog.Info("mt.Walk:newfid:%v", newfid)
	return &newmt, nil
}

//
// implement the Item interface
//
func (mt *MountPoint) GetDir() *warp9.Dir {
	return &mt.Dir
}
func (mt *MountPoint) GetItem() Item {
	return mt
}

func (mt *MountPoint) IsDirectory() Directory {

	if mt.fid != nil && mt.fid.Qid.Type&warp9.QTDIR == 0 {
		return nil
	}
	return mt
}

func (mt *MountPoint) Parent() Directory {
	return mt.parent
}

func (mt *MountPoint) SetParent(d Directory) error {
	if mt.parent != nil {
		return warp9.Einval
	}
	mt.parent = d
	return nil
}

func (mt *MountPoint) GetQid() warp9.Qid {
	return mt.fid.Qid
}

func (mt *MountPoint) Walked() (Item, error) {
	return mt.Walk([]string(nil))
}

func (mt *MountPoint) Read(obuf []byte, off uint64, rcount uint32) (uint32, error) {
	mlog.Debug("mt.Read:off:%v, rcount:%v", off, rcount)
	buf, err := mt.mi.clnt.Read(mt.fid, off, rcount)
	if err != warp9.Egood {
		return 0, err
	}
	copy(obuf, buf)
	return uint32(len(buf)), nil
}

func (mt *MountPoint) Write(ibuf []byte, off uint64, count uint32) (uint32, error) {
	i, err := mt.mi.clnt.Write(mt.fid, ibuf[:count], off)
	if err != nil {
		return 0, err
	}
	return uint32(i), nil
}

func (mt *MountPoint) Open(mode byte) (uint32, error) {
	mlog.Debug("mt.Open:")
	err := mt.mi.clnt.FOpen(mt.fid, mode)
	if err != warp9.Egood {
		return 0, err
	}
	return mt.fid.Iounit, nil
}

func (mt *MountPoint) Clunk() error {
	mlog.Debug("mt.Clunk:%v, fid#:%v", mt.fid, mt.fid.Fid)
	err := mt.mi.clnt.Clunk(mt.fid)
	if err != warp9.Egood {
		return err
	}
	return nil
}

func (mt *MountPoint) Remove() error {
	return mt.mi.clnt.FRemove(mt.fid)
}

func (mt *MountPoint) Stat() (*warp9.Dir, error) {
	mlog.Info("mt.Stat:fid:%v", mt.fid)
	d, e := mt.mi.clnt.FStat(mt.fid)
	if e != warp9.Egood {
		return nil, e
	}
	return d, nil
}

func (mt *MountPoint) WStat(dir *warp9.Dir) error {
	return mt.mi.clnt.FWstat(mt.fid, dir)
}

// set the debug level for client actions to target object servers
// -1=don't change, 0=off, >0=fcall, >1=raw msg bytes
// return the previous state
func (mt *MountPoint) Debug(level int) int {
	past := mt.mi.clnt.Debuglevel
	if level >= 0 {
		mt.mi.clnt.Debuglevel = level
	}
	return past
}
