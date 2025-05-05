package lib

import (
	"fmt"
	"github.com/Kry0z1/impulse/eventtime"
)

type PenaltyLapInfo struct {
	totalTime eventtime.TimestampMS

	entered eventtime.TimestampMS
	curLaps int

	lapLen    int
	totalLaps int
}

func (p *PenaltyLapInfo) Entered(tm eventtime.TimestampMS, laps int) {
	p.entered = tm
	p.curLaps = laps
}

func (p *PenaltyLapInfo) Left(tm eventtime.TimestampMS) {
	p.totalTime.Duration = p.totalTime.Duration + tm.Duration - p.entered.Duration
	p.totalLaps += p.curLaps
}

func (p *PenaltyLapInfo) Time() eventtime.TimestampMS {
	return p.totalTime
}

func (p *PenaltyLapInfo) AverageSpeed() float64 {
	return float64(p.totalLaps*p.lapLen) / p.Time().Seconds()
}

func (p PenaltyLapInfo) String() string {
	if p.totalLaps == 0 {
		return "{,}"
	}

	return fmt.Sprintf("{%v, %.3f}", p.Time(), p.AverageSpeed())
}

func NewPenaltyLapInfo(lapLen int) PenaltyLapInfo {
	return PenaltyLapInfo{lapLen: lapLen}
}
