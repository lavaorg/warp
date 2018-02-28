// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

// this file exists only to provide documentation for godoc

package protocol

/*
Create: Client requests server to create a new object.

    size[4] Tcreate tag[2] fid[4] name[s] perm[4] mode[1]

    size[4] Rcreate tag[2] qid[13] iounit[4]

The fid parameter represents a directory and the client should have write
permission in that directory. The owner of the object will be implied by
the current client session and the group will be that of the directory.

The object will have its permissions set based on the supplied perm paramter
along with the current permissions of the directory.  Such as:

    perm & (~0666 | (dir.perm & 0666))   -- regular object

	perm & (~0777 | (dir.perm & 0777))   -- directory object

This means, for example, that if the create allows read permission to others,
but the containing directory does not, then the created object will not allow
others to read the object.

Finally, the newly created object is opened according to mode (see open), and
fid will be associated with the newly opened object. Mode is not checked against
the permissions in perm.

Directories are created by setting the DMDIR bit (0x80000000) in the perm.

Upon a successful create the server returns a qid that represents the opened
object.

The iounit return value, if not zero, represents the server's idea of an atomic unit
of bytes that can be delivered in one message. This only represents what the server
can read or write and not necessarily what the network guarantees in the middle.

If the fid was in an opened state before this message an error will be returned.

The names . and .. are special; it is illegal to create objects with these names.

An attempt to create a object in a directory where the given name already exists
will result in an Error.
*/
func Create() {}
