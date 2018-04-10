// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

// sensim provides a simple single level object server.
// The root directory will have no directory objects; only entry objects.
//
// The two objects provider are:
//    ctl -- a control file
//    sensors -- a read only file to serve sensor readings
//
// The simulator will periodically update the values of 'sensors'
// If the ctl file is read without previously writing a commmand a simple
// info message is returned.
// if the ctl file is open rdwr and a valid command is writen the results
// of the command can be immediately read. If multiple commands are writen
// the result of each command can be read with line breaks in between.
//
package sensim

import (
	"log"

	"github.com/lavaorg/lrt/mlog"
	"github.com/lavaorg/warp/warp9"
)

type senFid struct {
	entry       *SenDir
	direntrybuf []byte
}

type SenSrv struct {
	warp9.Srv
	warp9.StatsOps
}

// SenDir represents an entry in the SenSim object server.
type SenDir struct {
	warp9.Dir
	item Item
}

type Item interface {
	Stat(dir *SenDir) error
	Read() ([]byte, error)
}

var root *SenDir = newSenDir(".", true)

var objects map[string]*SenDir //all our objects
var qidp uint64 = uint64(0xF0) //unique qid-path counter

func init() {
	objects = make(map[string]*SenDir)
	objects[root.Name] = root

	objects["ctl"], _ = newSensorItem(newSenDir("ctl", false))
	objects["sensors"], _ = newSensorItem(newSenDir("sensors", false))
}

func newSenDir(n string, dir bool) *SenDir {
	var d SenDir

	d.Name = n
	d.Uid = "nobody"
	d.Gid = "nobody"
	d.Muid = "nobody"

	if n == "/" {
		d.Mode = warp9.DMDIR | uint32(perms(warp9.DMREAD, warp9.DMREAD, warp9.DMREAD))
	} else {
		d.Mode = uint32(perms(warp9.DMREAD, warp9.DMREAD, warp9.DMREAD))
	}

	d.Atime = 0
	d.Mtime = 0
	typ := uint8(warp9.QTFILE)
	if dir {
		typ = warp9.QTDIR
	}
	d.Qid = warp9.Qid{typ, 0, qidp}
	qidp++

	d.Type = 0
	d.Dev = 0
	mlog.Debug("new SenDir:%v", d)
	return &d
}

func perms(u, g, o byte) uint16 {
	return uint16(uint16(u)<<6 | uint16(g)<<3 | uint16(o))
}

func (*SenSrv) ConnOpened(conn *warp9.Conn) {
	if conn.Srv.Debuglevel > 0 {
		log.Println("connected")
	}
}

func (*SenSrv) ConnClosed(conn *warp9.Conn) {
	if conn.Srv.Debuglevel > 0 {
		log.Println("disconnected")
	}
}

func (*SenSrv) FidDestroy(sfid *warp9.SrvFid) {
	var fid *senFid

	if sfid.Aux == nil {
		return
	}

	fid = sfid.Aux.(*senFid)
	if sfid.Fconn.Debuglevel > 0 {
		log.Printf("fid destroy:%v\n", fid)
	}
	//cleanup fid
}

func (ufs *SenSrv) Attach(req *warp9.SrvReq) {
	if req.Afid != nil {
		req.RespondError(warp9.Enoauth)
		return
	}
	//tc := req.Tc
	// ignore the aname; just mount "/"
	fid := new(senFid)
	fid.entry = root
	req.Fid.Aux = fid
	req.RespondRattach(&root.Qid)
}

func (*SenSrv) Flush(req *warp9.SrvReq) {}

func (*SenSrv) Walk(req *warp9.SrvReq) {
	fid := req.Fid.Aux.(*senFid)
	tc := req.Tc
	if fid == nil {
		req.RespondError(warp9.Ebaduse)
		return
	}

	if req.Newfid.Aux == nil {
		req.Newfid.Aux = new(senFid)
	}

	p := "."
	if len(tc.Wname) > 0 {
		p = tc.Wname[0]
		if p == "." || p == ".." {
			p = "."
		}
	}
	o := objects[p]
	if o == nil {
		req.RespondError(warp9.Enotexist)
		log.Printf("obj not found: %v\n", p)
		return
	}
	senfid := req.Newfid.Aux.(*senFid)
	senfid.entry = o
	wqids := make([]warp9.Qid, 1)
	wqids[0] = senfid.entry.Qid

	req.RespondRwalk(wqids[0:])
}

func (*SenSrv) Open(req *warp9.SrvReq) {

	tc := req.Tc
	mode := tc.Mode
	if mode != warp9.OREAD {
		req.RespondError(warp9.Eperm)
		return
	}
	sfid := req.Fid.Aux.(*senFid)
	req.RespondRopen(&sfid.entry.Qid, 0)
}

func (*SenSrv) Create(req *warp9.SrvReq) {
	// no creation
	req.RespondError(warp9.Enotimpl)
}

func (*SenSrv) Read(req *warp9.SrvReq) {
	tc := req.Tc
	fid := req.Fid

	rc := req.Rc
	rc.InitRread(tc.Count)

	// should check for having been opened

	// convert our directory to byte buffer; we aren't caching
	//b := warp9.PackDir(&root.Dir, req.Conn.Dotu)
	var b []byte
	var err warp9.W9Err
	if fid.Type&warp9.QTDIR > 0 {
		b, err = readdir(req)
	} else {
		b, err = readobj(req)
	}
	if err != warp9.Egood {
		req.RespondError(err)
		return
	}

	// determine which and how many bytes to return
	var count int
	switch {
	case tc.Offset > uint64(len(b)):
		count = 0
	case len(b[tc.Offset:]) > int(tc.Count):
		count = int(tc.Count)
	default:
		count = len(b[tc.Offset:])
	}
	copy(rc.Data, b[tc.Offset:int(tc.Offset)+count])
	log.Printf("buf:%v, rc.Data: %v, off:%v,  count:%v\n", len(b), len(rc.Data), tc.Offset, count)
	rc.SetRreadCount(uint32(count))
	req.Respond()
}

func readdir(req *warp9.SrvReq) ([]byte, warp9.W9Err) {
	buf := make([]byte, 100)
	for _, o := range objects {
		b := warp9.PackDir(&o.Dir)
		buf = append(buf, b...)
	}
	return buf, warp9.Egood
}

func readobj(req *warp9.SrvReq) ([]byte, warp9.W9Err) {
	sfid := req.Fid.Aux.(*senFid)
	sdir := sfid.entry

	if sdir.item != nil {
		b, e := sdir.item.Read()
		if e != nil {
			return nil, warp9.Eio
		}
		return b, warp9.Egood
	}
	return nil, warp9.Eio
}

func (*SenSrv) Write(req *warp9.SrvReq) {
	req.RespondError(warp9.Enotimpl)
	return
}

func (*SenSrv) Clunk(req *warp9.SrvReq) { req.RespondRclunk() }

func (*SenSrv) Remove(req *warp9.SrvReq) {
	req.RespondError(warp9.Enotimpl)
	return
}

func (*SenSrv) Stat(req *warp9.SrvReq) {
	mlog.Debug("Stat: %v", req)
	fid := req.Fid.Aux.(*senFid)
	sdir := fid.entry
	if sdir.item != nil {
		sdir.item.Stat(sdir)
	}
	req.RespondRstat(&fid.entry.Dir)
	return
}
func (u *SenSrv) Wstat(req *warp9.SrvReq) {
	req.RespondError(warp9.Enotimpl)
	return
}
