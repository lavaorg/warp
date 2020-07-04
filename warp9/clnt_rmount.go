// Copyright 2019 RMG Technologies, inc.  All rights reserved.

package warp9

import (
	"net"
)

// Callback from ReverseMountListener with the mount client context.
// This callback is made on a thread context specific for this mount attempt, thread will exit after this function returns.
//   (c9) the Context representing the remote server
//   (err) nil if success, else err returned from MountConn
//
type RMountConn func(c9 *Clnt)

// Callback from ReverseMountListener upon an error.  
// This callback is made on a thread context specific for this callback, thread will exit after this function returns.
//   (kind) 0==error during listen; 1==error during mount
//   (err) the causing error
type RMountError func (kind int, err error)

// allows client to Close the RMountListener
type RMountCloser interface {
	Close() error
} 

// Start listening on the specified symbolic network type/address and then mount the calling server.
// This is a reverse mount: the server initiates a connection and then this
// client performs the mount on that connection.
// This function spawns a thread to continuosly listen for new connections, until closed.
// Each accepted connection spawns a thread that mounts the object server.
// This function uses net.Listen(ntype,addr).
//   (ntype,addr) type and address to listen on - similar to net.Listen()
//   (aname,msize,user) parameters to perform mount
//	 (mounted) called when a valid mount is performed
//	 (rmounterr) called upon an error during 
// Return (Closer,nil) on success; else (nil,error)
//
func ReverseMountListener(ntype, addr string, aname string, msize uint32, user User, mounted RMountConn, rmounterr RMountError ) (RMountCloser, error) {

	l, err := net.Listen(ntype, addr)
	if err != nil {
		return nil,err
	}

	go handleListen(l, aname,msize,user,mounted, rmounterr) 

	return l,nil
}
	
	
// perform Accept on listener in a loop; until closed or error.
func handleListen(l net.Listener, aname string, msize uint32, user User, mounted RMountConn, rmounterr RMountError) {

	for {
		c, err := l.Accept()
		if err != nil {
			rmounterr(0,err)
			return //exit thread
		} else {
			go handleConnection(c, aname, msize, user, mounted, rmounterr)
		}
	}
}

// mount the server that just called on the Conn
func handleConnection(c net.Conn, aname string, msize uint32, user User, mounted RMountConn, rmounterr RMountError) {

	c9, err := MountConn(c, aname, msize, user)
	if err != nil {
		c.Close()
		rmounterr(1,err)
	} else {
		mounted(c9)
	}
	//exit thread
}

