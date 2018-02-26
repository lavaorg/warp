// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package warpsrv

var Eexist = &Error{"file already exists", EEXIST}
var Enoent = &Error{"file not found", ENOENT}
var Enotempty = &Error{"directory not empty", EPERM}
