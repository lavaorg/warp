// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 License file

package warp9

// A helper function that will open the named object, perform a read operations
// then close the object. The read operation will begin at a given offset and will
// return at max the number of bytes the server is willing to return in one message.
// (this is the size of the iounit as would be returned by Open)
// The associated Qid is also returned.
// A suitable error code is returned if any error.
func (clnt *Clnt) Get(path string, offset uint64) ([]byte, *Qid, W9Err) {

	obj, err := clnt.Open(path, OREAD)
	if err != Egood {
		return nil, nil, err
	}
	qid := &obj.Fid.Qid //remember or qid
	data, err := clnt.Read(obj.Fid, offset, obj.Fid.Iounit)
	if err != Egood {
		return nil, qid, err
	}
	err = clnt.Clunk(obj.Fid)
	return data, qid, err
}
