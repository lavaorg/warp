// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

// this file exists only to provide documentation for godoc

package protocol

/*
Open: Client requests access to an object.

    size[4] Topen tag[2] fid[4] mode[1]

    size[4] Ropen tag[2] qid[13] iounit[4]

The fid parmeter represents the object the client is requesting access to.
The object must not already be open.

The mode parameter indicates the access the client is requesting and allows
the client to indicate special conditions to open the object under.  The mode
is a bit mask with the following potential settings:

The low order 3 bits should be one of the following values:

	READ -- read contents from an object

	OWRITE -- write contents to object

	ORDWR -- read from or write to object

	OUSE -- use the object (example: walk allowed)

In addition the mode parameter can have the following bits set. All other bits
should be 0:

	OTRUNC -- object is truncated. Requires Write permission. Ignored if object is
	append only.

	ORCLOSE -- Object is removed when fid clunked. Requires permission to remove
	a object from the containing directory.


Upon this request the server will verify the implied user (uname established
during attach) has the permissions to perform the operation.

Upon success the server prepares the object (associated with fid) for subsequent
I/O operations via read/write messages.

It is illegal to write a directory, truncate it, or attempt to remove it
on close.

If the object is marked for exclusive use, only one client can have the
object open at any time. That is, after such a object has been opened, further
opens will fail until fid has been clunked. All these permissions are
checked at the time of the open request; subsequent changes to the
permissions of objects do not affect the ability to read, write, or remove
an open object.

Upon a successful create the server returns a qid that represents the opened
object.

The iounit return value, if not zero, represents the server's idea of an atomic unit
of bytes that can be delivered in one message. This only represents what the server
can read or write and not necessarily what the network guarantees in the middle.

If the fid was in an opened state before this message an error will be returned.
*/
func Open() {}
