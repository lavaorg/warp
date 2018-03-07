// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package sensim

import (
	"fmt"

	"github.com/lavaorg/lrt/mlog"
)

type sensorItem struct {
	lasttemp float32
}

func newSensorItem(sdir *SenDir) (*SenDir, error) {
	var sensor sensorItem = sensorItem{32.8}
	sdir.item = &sensor
	return sdir, nil
}

func (s *sensorItem) Read() ([]byte, error) {
	if s.lasttemp > 50 {
		s.lasttemp = 28
	}
	s.lasttemp += 5
	return []byte(fmt.Sprintf("%f", s.lasttemp)), nil
}

func (s *sensorItem) Stat(dir *SenDir) error {
	mlog.Debug("s:%v; sd:%v", s, dir)
	return nil
}

func (s *sensorItem) String() string {
	return fmt.Sprintf("(temp:%v)", s.lasttemp)
}
