# Warp9

This package implements the Warp9 protocol.

The Warp9 protocol allows for the remote access to a set of hierarchical named 
objects. The protocol provides a small set of primitives with a simple, but powerful,set of behaviors enabled by these primitives. 

The protocol is simple to use and simple to build resulting in reliable applications that can span from small simple devices to micro-servcies running in large-scale cloud hosted clusters.

The basic primitives:

- Object primitives:
  - TWalk   -- traversal of a hierarchical namespace.
  - TOpen   -- switch state of an object for allowing access.
  - TClunk  -- switch state of an object to stop access.
  - TRead   -- Read bytes from an object.
  - TWrite  -- Write bytes to an object. 
  - TCreate -- Create a new objects.
  - TRemove -- delete an object.
  - TStat   -- Obtain object status/meta-data.
  - TWStat  -- Modify an object's meta-data.
- Object Server primitives
  - TVersion -- Handshake and exchange protocol information.
  - TAuth   -- establish an authenticated connection.
  - TAttach -- Initiate access to an object provider's objects.
  - TFlush  -- Flush/disregard outstanding I/O messages.
 
