// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

// this file exists only to provide documentation for godoc

package protocol

/*
Write: Request the server to change the contents of an object.

    size[4] Twrite tag[2] fid[4] offset[8] count[4] data[count]

    size[4] Rwrite tag[2] count[4]

The write request asks that count bytes of data be recorded in the object
identified by fid, which must be opened for writing, starting offset bytes
after the beginning of the object. If the object is append-only, the data will
be placed at the end of the object regardless of offset.

Directory Objects may not be written.

The write reply records the number of bytes actually written. It is usually
an error if this is not the same as requested.

Because Warp9 implementations may limit the size of individual messages, more
than one message may be produced by client application level read and write
operations. The iounit field field returned by open, if non-zero, reports the
maximum size that is guaranteed to be transferred atomically.
*/
func Write() {}
