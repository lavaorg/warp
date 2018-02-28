// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

// this file exists only to provide documentation for godoc

package protocol

/*
Read: Request the server to return a sequence of bytes from the object.

    size[4] Tread tag[2] fid[4] offset[8] count[4]

    size[4] Rread tag[2] count[4] data[count]

The read request asks for count bytes of data from the object identified by fid,
starting offset bytes after the beginning of the objects data. The fid must have
been open for read access.

The bytes are returned with the read reply message.

The count field in the reply indicates the number of bytes returned. This may
be less than the requested amount. If the offset field is greater than or equal
to the number of bytes in the object, a count of zero will be returned.

For directory objects, read returns an integral number of directory entries exactly
as in stat (see stat), one for each member of the directory. The read
request message must have offset equal to zero or the value of offset in the
previous read on the directory, plus the number of bytes returned in the
previous read. In other words, seeking other than to the beginning is illegal
in a directory.

Because Warp9 implementations may limit the size of individual messages, more
than one message may be produced by client application level read and write
operations. The iounit field field returned by open, if non-zero, reports the
maximum size that is guaranteed to be transferred atomically.
*/
func Read() {}
