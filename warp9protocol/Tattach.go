// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

// this file exists only to provide documentation for godoc

package protocol

/*
Attach: A client initiation of access to a resource server.

    size[4] Tattach tag[2] fid[4] atok[4] uid[4] aname[s]

    size[4] Rattach tag[2] qid[13]

The attach message serves as a fresh introduction from a client to the
server. The clients user id (uid) is provided and this uid will be
used to determine what objects the client has permission to use. The aname
parameter is optional, if provided it will be the name of a specific top
level object the server is willing to provide.  If not present it defaults
to an object of the server's choosing.  Note "/", "." and ".." are the equivalent
of a blank aname.

The fid parameter, upon success, will represent the root object exported
by the server. If the fid is in use an error is returned.

Upon success the server will return a qid representing the root object.

The uid is a int32 -- signed 32bit integer, little-endian, 2's compliment

The atok field represents an authentication token for the session.  If the
client does not wish to authenticate it should send a NOFID (uint32(~0) value.
If the server requires authenticatio it will return an error and close the
connection.

The server's use of the authentication token is implementation specific. The
authetnication token should be representativ of (connection, time, uname) at
the server's discrtion.

How the client obtains a token is not specified in this protocol.
*/
func Attach() {}
