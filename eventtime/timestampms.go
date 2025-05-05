package eventtime

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type TimestampMS struct {
	time.Duration
}

func (s *TimestampMS) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	return s.FromString(str)
}

func (s *TimestampMS) FromString(str string) error {
	parts := strings.Split(str, ".")
	if len(parts) != 2 {
		return fmt.Errorf("wrong format: wrong number of dots")
	}

	noMS, err := NewTimestamp(parts[0])
	if err != nil {
		return err
	}

	if len(parts[1]) != 3 {
		return fmt.Errorf("wrong ms format: invalid number of digits")
	}

	ms, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("wrong ms format: non-numeric")
	}

	s.Duration = noMS.Duration + time.Duration(ms)*time.Millisecond

	return nil
}

func (s TimestampMS) String() string {
	hours := int(s.Hours())
	minutes := int(s.Minutes()) - 60*hours
	seconds := int(s.Seconds()) - 60*(minutes+60*hours)
	ms := int(s.Milliseconds()) - 1000*(seconds+60*(minutes+60*hours))
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, ms)
}

func NewTimestampMS(str string) (TimestampMS, error) {
	var s TimestampMS
	err := s.FromString(str)
	return s, err
}
