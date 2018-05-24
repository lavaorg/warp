// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

// this file exists only to provide documentation for godoc

package protocol

/*
Error: A reply only message for server to return an error for a request.

	**size[4] Rerror tag[2] ename[s]

	size[4] Rerror tag[2] errcode[2] ename[s]

The Rerror message is used to return an error indicating the failure of a
transaction. It replaces the Rmsg associated with the requestoers Tmsg.

The errcode is a 16bit signed integer where <0 indicates a warp9 framework error code
and >0 inidcates a user-server error code.

Optionally the UTF-8 string 'ename' can be provided.
If ename is present it will follow the UTF-8 Warp9 representiation. E.g.

The tag parameter is that of the failing Tmsg request.

There is no Terror message.

By convention, clients may truncate error messages after ERRMAX−1 bytes.
*/
func Error() {}
