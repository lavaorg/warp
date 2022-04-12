// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

/*
Wkit provides a toolkit for the creation of a Warp9 object server that allows
for the creation of a dynamic adjustable namespace of a hierarchical set of objects.

The namespace allows for bind points to be added into the hierarchy that can
support union object trees.  Each bind point can contain a stackable set of
object trees with the option of presening a union view of the set.

A set of specific **Object Typs** are provided to either build more complex
objects or to provide some __specific__ objects.

Objects will fulfill two basic primary interface types:
- Item       -- a leaf in the namespace tree
- Directory  -- a node in the namespace tree that is a root of a subtree

The objects exist in the wkit package so clients can use an object by
using the import

   import "github.com/lavaorg/warp/wkit/<object-name>"


Example Objects:

**__BaseItem__**: is a concrete type that fulfills the Item interface and
provides a base upon which other specific objects can be constructed.

**__OneItem__**: is a concrete class implementting a simple in-memory
buffer where the contents can be maintained by the server and accessed
by remote clients.

**__DirItem__**: is a concrete directory implementation that allows for
other objects fulfilling the Item or Directory interfaces to be added as
sub-components of its tree.  This allows for building an object hierarchy.

**__MountPoint__**: is a concrete directory implementation that allwos for
mounting remote servers and placing them into the current namespace tree.

Please see documentation generated for each Object type.

*/
package wkit
