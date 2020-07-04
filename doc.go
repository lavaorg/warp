// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

/*
Warp -- a framework for constructing and using remote resource object
servers and clients.

The Warp framework uses a few key concepts to allow the construction of
distributed systems built around loosly connected remote resources.

A resource is provide as a hierarchy of named objects.

Object hierarchy follow a simple set of semantics for accessing and using
contents and functions.

Locally, clients can organize their view of remote object hierarchies into
a dynamically adjustable namespace of those object hierarchies.

This allows clients to be developed with an assumed layout of object hierarchies 
and at deployment time attache to remote object hierarchies and then bind
them into its own namespace in an expected arrangement.

Using a consistent manner for accessing remote named object hierarchies allows
for unique sets of deployment time orchestration.  A client and a server can
have interposers between them that can perform operations in a manner that
does not violate the client or servers expectations.  These interposers may
provide transformations, filters, redundancies, caches, load-balancers, etc.

The Warp framework will provide:

   Warp9 -- a remote object access protocol
   nspace -- an intelligent gateway dynamic namespace for mounting remote servers
   client frameworks
   simple single level servers
   multi-level servers
*/
package warp
