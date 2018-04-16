// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

// this file exists only to provide documentation for godoc

package protocol

/*
Stat: Obtain the attributes for a given object.

    size[4] Tstat tag[2] fid[4]

    size[4] Rstat tag[2] stat[n]

The stat transaction inquires about the object identified by fid. The reply will
contain a directory entry, stat, laid out as follows:

    size[2]     -- total byte count of the following data

    qid.type[1] -- the object's type (directory, etc.), represented as a bit vector
                   corresponding to the high 8 bits of the object's mode word.

    qid.vers[4] -- version number for given path.

    qid.path[8] -- the server’s unique identification for the object.

    mode[4]     -- permissions and flags

    atime[4]    -- last access time

    mtime[4]    -- last modification time

    osize[8]    -- size of the object in bytes

    name[s]     -- object name; must be / if the object is the root directory of
                   the server.

    uid[s]      -- owner name

    gid[s]      -- group name

    muid[s]     -- name of the user who last modified the object.

Integers in this encoding are in little-endian order (least significant byte
first).

The mode contains permission bits as described in intro and the following:

    DMDIR (0x80000000)   -- The object is a Directory object.

    DMAPPED (0x40000000) -- The object can only be appended to.

    DMEXCL (0x20000000)  -- The object will be accessed exclusivly.

    DMTMP (0x04000000)   -- The object is tempoary

these attributes are echoed in Qid.type.

Writes to append-only objects always place their data at the end of the object;
the offset in the write message is ignored, as is the OTRUNC bit in an open.

Exclusive use objects may be open for I/O by only one fid at a time across all
clients of the server. If a second open is attempted, it draws an error. Servers
may implement a timeout on the lock on an exclusive use object: if the fid holding
the object open has been unused for an extended period (of order at least minutes),
it is reasonable to break the lock and deny the initial fid further I/O.

Temporary objects may be removed by the server at its discretion when the objects
are not being used by any client. Clients cannot depend on the object removal.

The two time fields are measured in seconds since the epoch (Jan 1 00:00 1970 GMT).

The mtime field reflects the time of the last change of content (except when
later changed by wstat). For a plain object, mtime is the time of the most recent
create, open with truncation, or write; for a directory object it is the time of
the most recent remove, create, or wstat of a object in the directory.

Similarly, the atime field records the last read of the contents; also it is set
whenever mtime is set. In addition, for a directory, it is set by an attach,
walk, or create, all whether successful or not.

The muid field names the user whose actions most recently changed the mtime of the
object.

The osize is the size of the object in bytes. Directories and most synthetic
objects have a conventional length of 0.

The stat request requires no special permissions.

A read of a directory yields an integral number of directory entries in the
encoding given above (see read.

Note that since the stat information is sent as a Warp9 variable-length datum,
it is limited to a maximum of 65535 bytes.

*/
func Stat() {}
