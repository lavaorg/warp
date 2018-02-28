// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

// this file exists only to provide documentation for godoc

/*
Warp9: A remote resource access protocol.

A Warp9 server is an entity that provides a set of resources in the form of one
or more collections of objects arranged in a hierarchical namespace.
A server responds to requests by clients to navigate the hierarchy, and to access,
create, remove objects.

The collection of objects can be highly dyanamic.  They can represent existing
fixed resources or by synthesized on the fly based on when/how clients are
traversing the namespace hierarchy or how they are interacting with existing objects.

Accesing a remote set of resources will occurr over a bidirectional, reliable,
communication channel. A resource server should generally assume multiple clients
however, this does not preclude a low-end resource server from being single threaded

   Note: Warp9 is a direct decendent of the Plan 9 9P2000 protocol

The Warp9 protocol, following its 9P ancessor, will follow a request/reply protocol.
The 9P notation of describing requests (T-messags) and replies (R-messages) will
be used. Each corresponing T-message and R-message pair comprises a transaction.

Each message consists of a sequence of bytes. Two-, four-, and eight-byte fields
hold unsigned integers represented in little-endian order (least significant byte
first). Data items of larger or variable lengths are represented by a two-byte
field specifying a count, n, followed by n bytes of data. Text strings are
represented this way, with the text itself stored as a UTF-8 encoded sequence
of Unicode characters (see utf(7)). Text strings in Warp9 messages are not
NUL-terminated: n counts the bytes of UTF-8 data, which include no final zero
byte. The NUL character is illegal in all text strings in Warp9, and is therefore
excluded from object names, user names, and so on.

Each Warp9 message follows the structure of:

    size[4] msg[1] tag[2] data[n]

size is a four-byte size field specifying the length in bytes
of the complete message including the four bytes of the size field itself.

msg is a single byte that enumerates the protocol message type.

tag is a two byte identifier, described below.

data is the remaining parameters of different sizes.

In the message descriptions, the number of bytes in a field is given in brackets
after the field name. The notation parameter[n] where n is not a constant
represents a variable-length parameter: n[2] followed by n bytes of data forming
the parameter. The notation string[s] (using a literal s character) is
shorthand for s[2] followed by s bytes
of UTF-8 text.

(Systems may choose to reduce the set of legal characters to reduce
syntactic problems, for example to remove slashes from name components, but
the protocol has no such restriction. By convention servers should support
names with any printable character (that is, any character outide hexadecimal
00-1F and 80-9F) except slash, which is used as a name component separator.

Messages are transported in byte form to allow for machine independence;

MESSAGES

    size[4] Tversion tag[2] msize[4] version[s]

    size[4] Rversion tag[2] msize[4] version[s]

    size[4] Rerror tag[2] ename[s]

    size[4] Tflush tag[2] oldtag[2]

    size[4] Rflush tag[2]

    size[4] Tattach tag[2] fid[4] atok[4] uname[s] aname[s]

    size[4] Rattach tag[2] qid[13]

    size[4] Twalk tag[2] fid[4] newfid[4] nwname[2] nwname*(wname[s])

    size[4] Rwalk tag[2] nwqid[2] nwqid*(wqid[13])

    size[4] Topen tag[2] fid[4] mode[1]

    size[4] Ropen tag[2] qid[13] iounit[4]

    size[4] Tcreate tag[2] fid[4] name[s] perm[4] mode[1]

    size[4] Rcreate tag[2] qid[13] iounit[4]

    size[4] Tread tag[2] fid[4] offset[8] count[4]

    size[4] Rread tag[2] count[4] data[count]

    size[4] Twrite tag[2] fid[4] offset[8] count[4] data[count]

    size[4] Rwrite tag[2] count[4]

    size[4] Tclunk tag[2] fid[4]

    size[4] Rclunk tag[2]

    size[4] Tremove tag[2] fid[4]

    size[4] Rremove tag[2]

    size[4] Tstat tag[2] fid[4]

    size[4] Rstat tag[2] stat[n]

    size[4] Twstat tag[2] fid[4] stat[n]

    size[4] Rwstat tag[2]

BRIEF TERMINILOGY

The following provides brief definitions for the rest of the doc:

	fid    -- identifier provided by the client
	qid    -- identifier provided by the server
	tag    -- a client chosen id for a message transaction
	object -- a sequence of bytes whoses structre and meaning defined by object
			  An object's bytes can be considered the serialized form of its data.
			  Objects are organized as a tree.
	perm   -- permissions
	mode   -- a given objects mode of access
    stat   -- the object status, or object's meta-data

MESSAGE HANDLING

Each T-message has a tag field, chosen and used by the client to identify the
message. The reply to the message will have the same tag. Clients must arrange
that no two outstanding messages on the same connection have the same tag. An
exception is the tag NOTAG, defined as (ushort)~0: the client can
use it, when establishing a connection, to override tag matching in version
messages.

The type of an R-message will either be one greater than the type of the
corresponding T-message or Rerror, indicating that the request failed. In the
latter case, the ename field contains a string describing the reason for failure.
The version message identifies the version of the protocol and indicates the
maximum message size the system is prepared to handle. It also initializes the
connection and aborts all outstanding I/O on the connection. The set of messages
between version requests is called a session.

Most T-messages contain a fid, a 32-bit unsigned integer that the client uses to
identify a “current object" on the server. Fids can be seen as client opaque
handles to a remote object, regardless of that objects current disposition. That is
a Fid does not imply the remote object is open for I/O operations. Fids can be used,
for example, to examine directory objects, perform a stat on an object, etc. Fids
are chosent by the client. Each connection has a single fid-space for the duration
of that connection.

The fid supplied in an attach message will be taken by the server to refer to the
root of the served object tree. The attach identifies the user to the server and may
specify a particular object tree served by the server (for those that supply more
than one).

Permission to attach to the service is proven by providing a special authentication
token -- atok -- in the attach message. The atok is established outside the definition
of Warp9.  A Warp9 server and client will have its own method to obtain and verify
the validity of the atok for the given connection. The atok allows a server, at its
choice, to verify the passed uname in the attach message and the association with the
current connection. For the duration of the session the uname will be used for
establishing ownership, group memeber ship and access permissions.

A walk message causes the server to change the current object associated with a fid
to be a object in the directory that is the old current object, or one of its
subdirectories. Walk returns a new fid that refers to the resulting object. Usually,
a client maintains a fid for the root, and navigates by walks from the root fid.

A client can send multiple T-messages without waiting for the corresponding
R-messages, but all outstanding T-messages must specify different tags. The server
may delay the response to a request and respond to later ones; this is sometimes
necessary, for example when the client reads from an object that the server synthesizes
from external events such as from a sensor.

Replies (R-messages) to attach, walk, open, and create requests convey a qid
field back to the client. The qid represents the server’s unique identification for
the object being accessed: two objects on the same server hierarchy are the same if and
only if their qids are the same. (The client may have multiple fids pointing to a
single object on a server and hence having a single qid.) The thirteen-byte qid
fields hold a one-byte type, specifying whether the object is a directory, append-only
object, etc., and two unsigned integers: first the four-byte qid version, then the
eight-byte qid path. The path is an integer unique among all objects in the hierarchy.
If an object is deleted and recreated with the same name in the same directory, the old
and new path components of the qids should be different. The version is a version
number for a object; typically, it is incremented every time the object is modified.
An existing object can be opened, or a new object may be created in the current (directory)
object. I/O of a given number of bytes at a given offset on an open object is done by read
and write.

A client should clunk any fid that is no longer needed. The remove transaction deletes
objects.

The stat transaction retrieves information about the object. The stat field in the reply
includes the object's name, access permissions (read, write and execute for owner, group
and public), access and modification times, and owner and group identifications
(see stat. The owner and group identifications are textual names. The wstat
transaction allows some of an object's properties to be changed.

A request can be aborted with a flush request. When a server receives a Tflush,
it should not reply to the message with tag oldtag (unless it has already replied),
and it should immediately send an Rflush. The client must wait until it gets the
Rflush (even if the reply to the original message arrives in the interim), at which
point oldtag may be reused.

Because the message size is negotiable and some elements of the protocol are variable
length, it is possible (although unlikely) to have a situation where a valid
message is too large to fit within the negotiated size. For example, a very long
object name may cause a Rstat of the object or Rread of its directory entry to be too
large to send. In most such cases, the server should generate an error rather than
modify the data to fit, such as by truncating the object name. The exception is that
a long error string in an Rerror message should be truncated if necessary, since
the string is only advisory and in some sense arbitrary.

The expectation is most client programs don't use the Warp9 protcol directly, but
rather use library routines appropriate for their language and performing actions
on remote objects will be translated int 1 or more Warp9 messages.

OBJECT TYPES

Each object will have a specified, and well known, type. Currently there are
two types of objects:

   Directory -- a collection of objects. Can have decendents.
   Entry -- a leaf object. Can not have decendents

Each object, regardless of type, can have the same base set of operations, while
directories have a few extra operations that can be performed:
   Base Operations: Open, Clunk, Read, Write, Stat, WStat, Remove
   Dir Operations: Walk, Create

An Object, regardless of type, is created via the Create operation. The type
of object is specified in the mode parmeter (see Create)

OBJECT METADATA

Each object carries with it a set of meta data. This meta data contains the
objects modes, permissions, ownership, etc. When a Stat operatoins is performed
on an object the meta data is well defined structure returned in a byte-serialized
encoding. This structure and encoding is defined in the Stat command.

DIRECTORY OBJECT

When directory objects are read they behave as any other object. The Read command
returns a set of bytes and they are returned as a stream. The contents are the
meta data for each object in the directory in sequence. Specifically, its the
Stat structure, one for each object in the directory. See the Stat command.

Directory objects are created, as with any object, via the Create command. The
object type is encoded in the mode parmeter of that command.

All Directory objects have implicit named aliases for the parent and itself:
    ".." is the name alias for the Directory objects parent directory.
    "." is the name alias for it's self.
The root Directory for an exported set of resources has a parent of it's self.
That is when a root: self == "." == ".."

ACCESS PERMISSIONS

Access permissions are not enforced by the protocol.  Each resource server
will manage how it interprets permissions.  However, the following conventions
are expected:

This section describes conventional access permissions implemented by most
servers. These conventions are not enforced by the protocol and may
differ between servers. These conventions are similar to UNIX or Plan 9
conventions.

Each server manages a set of user and group names. Each user can be a member
of any number of groups. Each group has a group leader who has special
privileges (see stat). Every object request has an implicit user id (copied
from the original attach) and an implicit set of groups (every group of which
the user is a member).

Each object has an associated owner and group id and three sets of permissions:
those of the owner, those of the group, and those of “other” users. When the
owner attempts to do something to an object, the owner, group, and other permissions
are consulted, and if any of them grant the requested permission, the operation
is allowed. For someone who is not the owner, but is a member of the object's group,
the group and other permissions are consulted. For everyone else, the other
permissions are used. Each set of permissions says whether reading is allowed,
whether writing is allowed, and whether using is allowed. A walk in a
directory is regarded as using the directory, not reading it. Permissions
are kept in the low-order bits of the object mode: owner read/write/use
permission represented as 1 in bits 8, 7, and 6 respectively (using 0 to number
the low order). The group permissions are in bits 5, 4, and 3, and the other
permissions are in bits 2, 1, and 0. The object mode contains some additional
attributes besides the permissions. If bit 31 (DMDIR) is set, the object is a
directory; if bit 30 (DMAPPEND) is set, the object is append-only (offset is
ignored in writes); if bit 29 (DMEXCL) is set, the object is exclusive-use
(only one client may have it open at a time); if bit 26 (DMTMP) is set, the
object may be removed at the servers discetion, when the server concludes the
object is not in use. Other bits should be left as zero. These bits
are reproduced, from the top bit down, in the type byte of the Qid: QTDIR,
QTAPPEND, QTEXCL, (skipping two bits), and QTTMP. The name QTOBJ,
defined to be zero, identifies the value of the type for a plain object.

This specification is derrived from the Plan 9 Documentation and Manual Pages.
The source material is Copyright 2003, Lucent Technologies Inc. and others.
All rights reserved.

Extensions to the original source material are
Copyright 2005, Larry Rau. All rights reserved.
*/
package protocol

//Warp9 Protocol constants.
//
const (
	ERRMAX   = 28 // Clients can truncate error strings after ERRMAX-1
	MAXWELEM = 16

	NOTAG = -1
)
