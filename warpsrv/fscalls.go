// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package warpsrv

type FDestroyOp interface {
	FidDestroy(fid *W9Fid)
}
