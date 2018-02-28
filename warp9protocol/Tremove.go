// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

// this file exists only to provide documentation for godoc

package protocol

/*
Remove: Request the server to remove the object.

    size[4] Tremove tag[2] fid[4]

    size[4] Rremove tag[2]

The remove request asks the server both to remove the object represented by fid
and to clunk the fid, even if the remove fails.

This request will fail if the client does not have write permission in the
object's parent directory object.

It is correct to consider remove to be a clunk with the side effect of
removing the object if permissions allow.

If an object has been opened as multiple fids, possibly on different connections,
and one fid is used to remove the object, whether the other fids continue to
provide access to the object is implementation-defined.
*/
func Remove() {}
