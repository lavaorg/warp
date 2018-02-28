// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

// this file exists only to provide documentation for godoc

package protocol

//
// Version: Client and server agree on a protocol version.
//
//     size[4] Tversion tag[2] msize[4] version[s]
//
//     size[4] Rversion tag[2] msize[4] version[s]
//
/*
The version request negotiates the protocol version and message size to be
used on the connection and initializes the connection for I/O. Tversion
must be the first message sent on the Warp9 connection, and the client cannot
issue any further requests until it has received the Rversion reply. The
tag should be NOTAG (value (ushort)~0) for a version message.

The client suggests a maximum message size, msize, that is the maximum length,
in bytes, it will ever generate or expect to receive in a single Warp9 message.
This count includes all Warp9 protocol data, starting from the size field and
extending through the message, but excludes enveloping transport protocols.
The server responds with its own maximum, msize, which must be less than or
equal to the client’s value. Thenceforth, both sides of the connection must
honor this limit.

The version string identifies the level of the protocol. The string must
always begin with the two characters “W9” and be in the form W9nn.xx. Where
nn is a numeric identifire and xx can be implementation defined. The client
and server will ignore the ".xx"

If the server does not understand the client’s version string, it should
respond with an Rversion message (not Rerror) with the version string
the 7 characters “unknown”.

The server may respond with a version string of the form W9mm.yy where
mm is a numeric identifier of the protocol version it will support. The
server must use the client's version if it can support that version. The
server may add a ".yy" suffix indicating the server's implementation defined
identifying information. (note: both .xx and .yy can allow for logging and
debugging).

The client and server will use the protocol version defined by the server’s
response for all subsequent communication on the connection.

A successful version request initializes the connection. All outstanding
I/O on the connection is aborted; all active fids are freed (‘clunked’)
automatically. The set of messages between version requests is called a
session.
*/
func Version() {}
