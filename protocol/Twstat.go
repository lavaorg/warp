// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

// this file exists only to provide documentation for godoc

package protocol

/*
Wstat: Modify the attributes for a given object.

    size[4] Twstat tag[2] fid[4] stat[n]

    size[4] Rwstat tag[2]

The wstat request can change some of the object's attributes.

    see the Tstat message for the structure of an objects attributes.

The name can be changed by anyone with write permission in the parent
directory; it is an error to change the name to that of an existing object.

The length can be changed (affecting the actual length of the object) by
anyone with write permission on the object. It is an error to attempt to
set the length of a directory to a non-zero value, and servers may
decide to reject length changes for other reasons.

The mode and mtime can be changed by the owner of the object or the group
leader of the object's current group.

The directory bit cannot be changed by a wstat; the other defined
permission and mode bits can.

The gid can be changed: by the owner if also a member of the new group;
or by the group leader of the object's current group if also leader of the
new group (see intro for more information about permissions, users, and
groups).

None of the other data can be altered by a wstat and attempts to change them
will trigger an error. In particular, it is illegal to attempt to change the
owner of a object.

Either all the changes in wstat request happen, or none of them does: if the
request succeeds, all changes were made; if it fails, none were.

A wstat request can avoid modifying some properties of the object by providing
explicit “don’t touch” values in the stat data that is sent: zero-length strings
for text values and the maximum unsigned value of appropriate size for integral
values. As a special case, if all the elements of the directory entry in a Twstat
message are “don’t touch” values, the server may interpret it as a request to
guarantee that the contents of the associated object are committed to stable storage
before the Rwstat message is returned. (Consider the message to mean, “make the
state of the object exactly what it claims to be.”)

Note that since the stat information is sent as a Warp9 variable-length datum,
it is limited to a maximum of 65535 bytes.
*/
func Wstat() {}
