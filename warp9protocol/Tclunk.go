// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

// this file exists only to provide documentation for godoc

package protocol

/*
Clunk: Client indicates to server it is no longer interested an object

    size[4] Tclunk tag[2] fid[4]

    size[4] Rclunk tag[2]

The clunk request informs the server that the current object represented
by fid is no longer needed by the client. The actual object is not removed on
the server unless the fid had been opened with ORCLOSE.

Once a fid has been clunked, the same fid can be reused in a new walk or attach
request. Even if the clunk returns an error, the fid is no longer valid.
*/
func Clunk() {}
