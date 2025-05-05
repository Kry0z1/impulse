package lib

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Kry0z1/impulse/config"
	"github.com/Kry0z1/impulse/eventtime"
)

var (
	ErrInvalidNumberOfArguments = errors.New("invalid number of arguments")
	ErrInvalidTimeStamp         = errors.New("invalid timestamp")
	ErrInvalidEventID           = errors.New("invalid event id")
	ErrInvalidCompetitorID      = errors.New("invalid competitor id")
	ErrUnexpectedExtraParam     = errors.New("unexpected extra param")
	ErrNoExtraParam             = errors.New("extra param is not found")
	ErrDoubleRegister           = errors.New("competitor is already registered")
	ErrInvalidExtraParam        = errors.New("invalid extra param")
	ErrCompetitorNotFound       = errors.New("competitor is not found")
)

var extraParamsID = []int{2, 5, 6, 11}

type Orchestrator struct {
	competitors     map[int]*Competitor
	competitorOrder []int

	cfg *config.Config

	append string
}

func (o *Orchestrator) Result() string {
	var sb strings.Builder

	for _, comp := range o.competitorOrder {
		sb.WriteString(o.competitors[comp].String())
		sb.WriteString("\n")
	}

	return sb.String()
}

func (o *Orchestrator) ParseLine(line string) (string, error) {
	p, err := parseParamsFromLine(line)
	if err != nil {
		return "", fmt.Errorf("failed to parse line: %w", err)
	}

	if err := o.LoadEvent(p.time, p.eventID, p.competitorID, p.extraParam); err != nil {
		return "", fmt.Errorf("failed to load event: %w", err)
	}

	if o.append == "" {
		return outputString(p.time, p.eventID, p.competitorID, p.extraParam), nil
	}
	output := outputString(p.time, p.eventID, p.competitorID, p.extraParam) + "\n" + o.append
	o.append = ""
	return output, nil
}

func (o *Orchestrator) LoadEvent(
	time eventtime.TimestampMS,
	eventID int,
	competitorID int,
	extraParam any,
) error {
	if eventID == 1 {
		return o.register(competitorID)
	}

	comp, ok := o.competitors[competitorID]
	if !ok {
		return ErrCompetitorNotFound
	}

	switch eventID {
	case 2:
		comp.ScheduleStartTime(extraParam.(eventtime.TimestampMS))
	case 3:
		return nil
	case 4:
		if comp.Start(time) {
			o.append = outputString(time, 32, competitorID, extraParam)
		}
	case 5:
		comp.EnterFiringRange()
	case 6:
		comp.HitTarget()
	case 7:
		return nil
	case 8:
		comp.EnterPenaltyLaps(time)
	case 9:
		comp.LeavePenaltyLaps(time)
	case 10:
		if comp.EndMainLap(time) {
			o.append = outputString(time, 33, competitorID, extraParam)
		}
	case 11:
		comp.CantContinue()
	default:
		return ErrInvalidEventID
	}

	return nil
}

func (o *Orchestrator) register(competitorID int) error {
	_, ok := o.competitors[competitorID]
	if ok {
		return ErrDoubleRegister
	}

	o.competitorOrder = append(o.competitorOrder, competitorID)
	comp := NewCompetitor(o.cfg.StartDelta, o.cfg.Laps, o.cfg.LapLen, o.cfg.PenaltyLen, competitorID)
	o.competitors[competitorID] = &comp

	return nil
}

func parseExtraParam(eventID int, param string) (any, error) {
	switch eventID {
	case 2:
		return eventtime.NewTimestampMS(param)
	case 5:
		return strconv.Atoi(param)
	case 6:
		return strconv.Atoi(param)
	case 11:
		return param, nil
	}

	return nil, nil
}

type lineParams struct {
	time         eventtime.TimestampMS
	eventID      int
	competitorID int
	extraParam   any
}

func parseParamsFromLine(line string) (lineParams, error) {
	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		return lineParams{}, ErrInvalidNumberOfArguments
	}
	if len(parts[0]) < 2 {
		return lineParams{}, ErrInvalidTimeStamp
	}

	time, err := eventtime.NewTimestampMS(parts[0][1 : len(parts[0])-1])
	if err != nil {
		return lineParams{}, fmt.Errorf("%w: %w", ErrInvalidTimeStamp, err)
	}

	eventID, err := strconv.Atoi(parts[1])
	if err != nil {
		return lineParams{}, fmt.Errorf("%w: %w", ErrInvalidEventID, err)
	}

	competitorID, err := strconv.Atoi(parts[2])
	if err != nil {
		return lineParams{}, fmt.Errorf("%w: %w", ErrInvalidCompetitorID, err)
	}

	var extraParam string

	if len(parts) >= 4 {
		if !slices.Contains(extraParamsID, eventID) {
			return lineParams{}, ErrUnexpectedExtraParam
		}
		extraParam = strings.Join(parts[3:], " ")
	} else if slices.Contains(extraParamsID, eventID) {
		return lineParams{}, ErrNoExtraParam
	}

	parsedParam, err := parseExtraParam(eventID, extraParam)
	if err != nil {
		return lineParams{}, fmt.Errorf("%w: %w", ErrInvalidExtraParam, err)
	}

	return lineParams{
		time:         time,
		eventID:      eventID,
		competitorID: competitorID,
		extraParam:   parsedParam,
	}, nil
}

// assuming eventID and extraParam are already validated
func outputString(
	time eventtime.TimestampMS,
	eventID int,
	competitorID int,
	extraParam any,
) string {
	switch eventID {
	case 1:
		return fmt.Sprintf("[%v] The competitor(%d) registered", time, competitorID)
	case 2:
		return fmt.Sprintf("[%v] The start time for the competitor(%d) was set by a draw to %v", time, competitorID, extraParam.(eventtime.TimestampMS))
	case 3:
		return fmt.Sprintf("[%v] The competitor(%d) is on the start line", time, competitorID)
	case 4:
		return fmt.Sprintf("[%v] The competitor(%d) has started", time, competitorID)
	case 5:
		return fmt.Sprintf("[%v] The competitor(%d) is on the firing range(%d)", time, competitorID, extraParam.(int))
	case 6:
		return fmt.Sprintf("[%v] The target(%d) has been hit by competitor(%d)", time, extraParam.(int), competitorID)
	case 7:
		return fmt.Sprintf("[%v] The competitor(%d) left the firing range", time, competitorID)
	case 8:
		return fmt.Sprintf("[%v] The competitor(%d) entered the penalty laps", time, competitorID)
	case 9:
		return fmt.Sprintf("[%v] The competitor(%d) left the penalty laps", time, competitorID)
	case 10:
		return fmt.Sprintf("[%v] The competitor(%d) ended the main lap", time, competitorID)
	case 11:
		return fmt.Sprintf("[%v] The competitor(%d) can`t continue: %s", time, competitorID, extraParam.(string))

	case 32:
		return fmt.Sprintf("[%v] The competitor(%d) is disqualified", time, competitorID)
	case 33:
		return fmt.Sprintf("[%v] The competitor(%d) has finished", time, competitorID)
	}
	return ""
}

func NewOrchestrator(cfg *config.Config) *Orchestrator {
	return &Orchestrator{
		competitors:     make(map[int]*Competitor),
		competitorOrder: make([]int, 0),
		cfg:             cfg,
	}
}
