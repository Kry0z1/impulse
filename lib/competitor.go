package lib

import (
	"fmt"
	"github.com/Kry0z1/impulse/eventtime"
	"strings"
)

type Competitor struct {
	scheduledStartTime eventtime.TimestampMS

	notStarted  bool
	notFinished bool

	lapsInfo        []LapInfo
	penaltyLapsInfo PenaltyLapInfo

	currentLap int

	startDelta eventtime.Timestamp

	currentHitTargets int
	totalHitTargets   int
	totalTargets      int

	id int
}

func (c *Competitor) ScheduleStartTime(timestamp eventtime.TimestampMS) {
	c.scheduledStartTime = timestamp
}

func (c *Competitor) Start(timestamp eventtime.TimestampMS) bool {
	if timestamp.Duration-c.scheduledStartTime.Duration > c.startDelta.Duration {
		c.notStarted = true
		return true
	} else {
		c.EnterMainLap(c.scheduledStartTime)
		return false
	}
}

func (c *Competitor) EnterMainLap(timestamp eventtime.TimestampMS) {
	c.lapsInfo[c.currentLap].Entered(timestamp)
}

func (c *Competitor) EnterFiringRange() {
	c.totalTargets += 5
	c.currentHitTargets = 0
}

func (c *Competitor) HitTarget() {
	c.totalHitTargets++
	c.currentHitTargets++
}

func (c *Competitor) EnterPenaltyLaps(timestamp eventtime.TimestampMS) {
	c.penaltyLapsInfo.Entered(timestamp, 5-c.currentHitTargets)
}

func (c *Competitor) LeavePenaltyLaps(timestamp eventtime.TimestampMS) {
	c.penaltyLapsInfo.Left(timestamp)
}

func (c *Competitor) EndMainLap(timestamp eventtime.TimestampMS) bool {
	c.lapsInfo[c.currentLap].Left(timestamp)
	c.currentLap++

	if c.currentLap == len(c.lapsInfo) {
		return true
	}

	c.EnterMainLap(timestamp)
	return false
}

func (c *Competitor) CantContinue() {
	c.notFinished = true
}

func (c Competitor) String() string {
	var totalTimeStr string
	if c.notStarted {
		totalTimeStr = "NotStarted"
	} else if c.notFinished {
		totalTimeStr = "NotFinished"
	} else {
		var totalTime eventtime.TimestampMS
		for _, info := range c.lapsInfo {
			totalTime.Duration = totalTime.Duration + info.Time().Duration
		}
		totalTime.Duration = totalTime.Duration + c.penaltyLapsInfo.Time().Duration
		totalTimeStr = totalTime.String()
	}

	var lapsString strings.Builder
	lapsString.WriteString("[")
	lapsString.WriteString(c.lapsInfo[0].String())
	for _, lap := range c.lapsInfo[1:] {
		lapsString.WriteString(", ")
		lapsString.WriteString(lap.String())
	}
	lapsString.WriteString("]")

	return fmt.Sprintf(
		"[%s] %d %s %v %d/%d",
		totalTimeStr, c.id, lapsString.String(), c.penaltyLapsInfo, c.totalHitTargets, c.totalTargets,
	)
}

func NewCompetitor(
	startDelta eventtime.Timestamp,
	laps int,
	lapLen int,
	penaltyLapLen int,
	id int,
) Competitor {
	lapsInfo := make([]LapInfo, laps)
	for i := range laps {
		lapsInfo[i] = NewLapInfo(lapLen)
	}

	return Competitor{
		startDelta:      startDelta,
		lapsInfo:        lapsInfo,
		penaltyLapsInfo: NewPenaltyLapInfo(penaltyLapLen),
		id:              id,
	}
}
