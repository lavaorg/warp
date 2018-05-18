// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package wkit

import (
	"bytes"
	"strings"

	"github.com/lavaorg/lrt/mlog"
	"github.com/lavaorg/warp/warp9"
)

// A Command is a type of object that allows simple commands
// to be invoked.  The command follows the form of:
//  <name><space><args>  where:
//  name  == [string]single world command
//  space == single ascii space ' ' char
//  args  == arbitrary sequence of bytes
//
// the command is invoked by writing the byte representation
// of the command to the object.
//
// Results of the command can be provided as a sequence of bytes.
// The byte-sequence of the result can be read from the object.
//
// The commands acceppted can be provided as a map to the object
// following the form specified below (see .Fcts)
//
// An optional context, Ctx, can be provided and will be passed back
// to any command upon execution.
//
// The last results of a command will be available for as long as the
// object is available or until a new command is invoked. The results
// do not persist beyond the life of the server. The contents may be
// removed upon the last Clunk operatin on the file.
//
// Note: A Command object fulfills the Item interface
// and is a supertype of a OneItem
//
type Command struct {
	OneItem
	Buf  string
	Fcts map[string]CommandFct
	Ctx  CmdCtx
}

type CmdCtx interface {
}

type CommandFct func(ctxt CmdCtx, cmd *Command, cmdname string, args []byte) error

// create a new Command object.
// fcts == pre-created map of methods or
// fcts == nil will cause an empty mapp to be created
func NewCommand(name string, fcts map[string]CommandFct, ctx CmdCtx) *Command {
	var cmd Command
	cmd.Name = name
	if fcts == nil {
		fcts = make(map[string]CommandFct)
	}
	cmd.Fcts = fcts
	cmd.Ctx = ctx
	cmd.Buffer = make([]byte, 20, 80)
	return &cmd
}

// Add/replace a command.
// if method is nil the command is deleted.
func (cmd *Command) AddMethod(name string, meth CommandFct) {
	if meth == nil {
		delete(cmd.Fcts, name)
		return
	}
	cmd.Fcts[name] = meth
	return
}

// command file is always append only. complete command
// must be provided in a single write message.
// command is executed before write is complete
// successful write == command completed normally
// and contents will be the results of the command
// if command fails to complete; write fails
// and contents will be 0
//
func (o *Command) Write(ibuf []byte, off uint64, count uint32) (uint32, error) {

	// our file will not be super large;  convert everything to int
	ioff := int(off)
	icnt := int(count)
	if uint64(ioff) != off || uint32(icnt) != count {
		return 0, warp9.Etoolarge
	}

	// split buffer into the command and the args
	parts := bytes.SplitAfterN(ibuf, []byte{' '}, 2)

	if len(parts) < 1 {
		return 0, warp9.Eio
	}
	cmd := strings.TrimSpace(string(parts[0]))
	args := []byte(nil)
	if len(parts) > 1 {
		args = parts[1]
	}

	// invoke the command if found
	fct := o.Fcts[cmd]
	if fct == nil {
		mlog.Error("bad fct: [%v]", cmd)
		return 0, warp9.Eio
	}
	e := fct(o.Ctx, o, cmd, args)
	if e != nil {
		return 0, e
	}

	return count, nil
}

func (o *Command) Walked() (Item, error) {
	return o, nil
}
