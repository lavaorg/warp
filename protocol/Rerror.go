// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

// this file exists only to provide documentation for godoc

package protocol

/*
Error: A reply only message for server to return an error for a request.

	size[4] Rerror tag[2] ename[s]

The Rerror message is used to return an error string describing the failure
of a transaction. It replacess the Rmsg associated with the requesters Tmsg.

The tag parameter is that of the failing Tmsg request.

There is no Terror message.

By convention, clients may truncate error messages after ERRMAX−1 bytes.
*/
func Error() {}
