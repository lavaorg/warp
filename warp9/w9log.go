// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package warp9

// provide a logging interface that clients to the warp9 library can provide

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type LogFct func(template string, args ...interface{})
type LogDbg func(on bool)

var Debug LogFct = w9debug
var Event LogFct = w9event
var Info LogFct = w9info
var Stat LogFct = w9stat
var Error LogFct = w9error
var Alarm LogFct = w9alarm

var LogDebug LogDbg = logdebug

func SetLogging(d, i, e, s, err, a LogFct, ld LogDbg) {
	Debug = d
	Info = i
	Event = e
	Stat = s
	Error = err
	Alarm = a
	LogDebug = ld
}

func (clnt *Clnt) Perr(err error) error {
	if clnt.Debuglevel&DbgPrintAtErrMsg != 0 {
		// get caller statistics
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "???"
			line = 0
		}
		emit(DEBUG, file, line, err.Error())
	}
	return err
}

var gdbg bool = false

func logdebug(on bool) {
	gdbg = on
}

const (
	cversion     = "1"
	cmarker      = "*" + cversion
	cseparator   = "|"
	ctmformat    = "2006/01/02 15:04:05.999999"
	ctmformatlen = len(ctmformat)
)

// Severity Enumeration -- <INFO emit to stderr; else stdout
const (
	ALARM uint8 = iota
	ERROR
	STAT
	EVENT
	INFO
	DEBUG
	UNKNOWN
)

var (
	name string // process name
	pid  string // os pid
	// these must be initialized early
	stderr = os.Stderr
	stdout = os.Stdout

	// map sev enum to strings
	sevstr []string = []string{
		"ALARM",
		"ERROR",
		"STAT",
		"EVENT",
		"INFO",
		"DEBUG",
		"UNKNOWN",
	}
)

func init() {

	// determine basic application information
	pathx := strings.Split(os.Args[0], "/")
	name = pathx[len(pathx)-1]
	pid = strconv.Itoa(os.Getpid())
}

// Emit debug message if global debug flag set
func w9debug(template string, args ...interface{}) {
	if gdbg {
		// get caller statistics
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "???"
			line = 0
		}
		emit(DEBUG, file, line, fmt.Sprintf(template, args...))
	}
}

// Emit an Event message
func w9event(template string, args ...interface{}) {
	// get caller statistics
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	emit(EVENT, file, line, fmt.Sprintf(template, args...))
}

// Emit an Info message
func w9info(template string, args ...interface{}) {
	// get caller statistics
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	emit(INFO, file, line, fmt.Sprintf(template, args...))
}

// Emit a Stat message
func w9stat(template string, args ...interface{}) {
	// get caller statistics
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	emit(STAT, file, line, fmt.Sprintf(template, args...))
}

// Emit an Error message
func w9error(template string, args ...interface{}) {
	// get caller statistics
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	emit(ERROR, file, line, fmt.Sprintf(template, args...))
}

// Emit using the alarm severity level
func w9alarm(template string, args ...interface{}) {
	// get caller statistics
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	emit(ALARM, file, line, fmt.Sprintf(template, args...))
}

func emit(sev uint8, file string, line int, m string) {

	// determine the shorted-version of the filename
	// and avoid the func call of strings.SplitAfter
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short

	linx := strconv.Itoa(line)
	fileAndLine := strings.Join([]string{file, linx}, ":")

	// then split into individual lines (by CR)
	lines := strings.Split(m, "\n")

	// choose output stream to use
	stream := stderr
	if sev > EVENT {
		stream = stdout
	}

	// create a structured log message to emit
	var message string

	timestamp := time.Now().UTC().Format(ctmformat)
	message = strings.Join([]string{
		cmarker, sevstr[sev], pid, name, fileAndLine, timestamp,
	}, cseparator)

	for _, line := range lines {
		if line == "" {
			continue
		}
		// Write message to stdout or stderr
		_, _ = fmt.Fprintln(stream, message+cseparator+line)
	}
}
