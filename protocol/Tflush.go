// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

// this file exists only to provide documentation for godoc

package protocol

//
// Flush: Client sends to indicate a reply to an outstanding request is no longer
// needed.
//
//     size[4] Tflush tag[2] oldtag[2]
//
//     size[4] Rflush tag[2]
//
// The message being flushed is identified by oldtag. The semantics of flush
// depends on messages arriving in order. The server should answer the flush message
// immediately. If it recognizes oldtag as the tag of a pending transaction, it
// should abort any pending response and discard that tag. In either case, it should
// respond with an Rflush echoing the tag (not oldtag) of the Tflush message. A T
// flush can never be responded to by an Rerror message. The server may respond to
// the pending request before responding to the Tflush. It is possible for a client
// to send multiple Tflush messages for a particular pending request. Each subsequent
// Tflush must contain as oldtag the tag of the pending request (not a previous Tflush).
// Should multiple Tflushes be received for a pending request, they must be answered
// n order. A Rflush for any of the multiple Tflushes implies an answer for all
// previous ones. Therefore, should a server receive a request and then multiple
// flushes for that request, it need respond only to the last flush. When the client
// sends a Tflush, it must wait to receive the corresponding Rflush before reusing
// oldtag for subsequent messages. If a response to the flushed request is received
// before the Rflush, the client must honor the response as if it had not been flushed,
// since the completed request may signify a state change in the server. For instance,
// Tcreate may have created an object and Twalk may have allocated a fid. If no response
// is received before the Rflush, the flushed transaction is considered to have been
// canceled, and should be treated as though it had never been sent. Several
// exceptional conditions are handled correctly by the above specification:
// sending multiple flushes for a single tag, flushing after a transaction is
// completed, flushing a Tflush, and flushing an invalid tag.
//
func Flush() {}
