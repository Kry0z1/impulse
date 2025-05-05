package lib

import (
	"fmt"
	"github.com/Kry0z1/impulse/eventtime"
)

type LapInfo struct {
	entered eventtime.TimestampMS
	left    eventtime.TimestampMS
	lapLen  int
}

func (li *LapInfo) Entered(tm eventtime.TimestampMS) {
	li.entered = tm
}

func (li *LapInfo) Left(tm eventtime.TimestampMS) {
	li.left = tm
}

func (li *LapInfo) Time() eventtime.TimestampMS {
	return eventtime.TimestampMS{Duration: li.left.Duration - li.entered.Duration}
}

func (li *LapInfo) AverageSpeed() float64 {
	return float64(li.lapLen) / li.Time().Seconds()
}

func (li LapInfo) String() string {
	if li.entered.Duration == 0 || li.left.Duration == 0 {
		return "{,}"
	}

	return fmt.Sprintf("{%v, %.3f}", li.Time(), li.AverageSpeed())
}

func NewLapInfo(lapLen int) LapInfo {
	return LapInfo{lapLen: lapLen}
}
