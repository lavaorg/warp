// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

// this file exists only to provide documentation for godoc

package protocol

/*
Walk: Traverse the given directory associated with fid, using the named elements.

   size[4] Twalk tag[2] fid[4] newfid[4] nwname[2] nwname*(wname[s])

   size[4] Rwalk tag[2] nwqid[2] nwqid*(wqid[13])

The fid argument must associated with a directory and newfid can be the same as fid
otherwise it must not be in use. The newfid will be associated with the result of
walking the directory hierarchy using the path name elements stored in wname.
Each path name element must be a successive directory, except for the last element.

The fid must be valid for the current session and not have been opened for I/O
via an open or create message. If the full sequence of nwname elements is
walked successfully, newfid will represent the object that results. If not,
newfid (and fid) will be unaffected. However, if newfid is in use or
otherwise illegal, an Rerror is returned.

The name “..” (dot-dot) represents the parent directory. The name “.” (dot),
meaning the current directory, is not used in the protocol. It is legal for
nwname to be zero, in which case newfid will represent the same object as fid
and the walk will usually succeed; this is equivalent to walking to dot. The
rest of this discussion assumes nwname is greater than zero.

The nwname path name elements wname are walked in order, “elementwise”.
For the first elementwise walk to succeed, the object identified by fid must
be a directory, and the implied user of the request must have permission
to search the directory (package intro). Subsequent elementwise walks have
equivalent restrictions applied to the implicit fid that results from the
preceding elementwise walk.

If the first element cannot be walked for any reason, Rerror is returned.
Otherwise, the walk will return an Rwalk message containing nwqid qids
corresponding, in order, to the objects that are visited by the nwqid successful
elementwise walks; nwqid is therefore either nwname or the index of the first
elementwise walk that failed. The value of nwqid cannot be zero unless nwname
is zero. Also, nwqid will always be less than or equal to nwname. Only if it
is equal, however, will newfid be affected, in which case newfid will
represent the object reached by the final elementwise walk requested in the
message. A walk of the name “..” in the root directory of a server is
equivalent to a walk with no name elements. If newfid is the same as fid, the
above discussion applies, with the obvious difference that if the walk changes
the state of newfid, it also changes the state of fid; and if newfid is
unaffected, then fid is also unaffected.

To simplify the implementation of the servers, a maximum of sixteen name
elements or qids may be packed in a single message. This constant is called
MAXWELEM in fcall(3). Despite this restriction, the system imposes no limit
on the number of elements in an object name, only the number that may be
transmitted in a single message.
*/
func Walk() {}
