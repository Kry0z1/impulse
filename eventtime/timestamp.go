package eventtime

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Timestamp struct {
	time.Duration
}

func (s *Timestamp) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	return s.FromString(str)
}

func (s *Timestamp) FromString(str string) error {
	parts := strings.Split(str, ":")
	if len(parts) != 3 {
		return fmt.Errorf("wrong format: wrong number of colons")
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("wrong format: noninteger hours")
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("wrong format: noninteger minutes")
	}

	seconds, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("wrong format: noninteger seconds")
	}

	if hours < 0 {
		return fmt.Errorf("invalid hours: negative hours")
	}

	if minutes < 0 || minutes > 59 {
		return fmt.Errorf("invalid minutes: minutes out of range")
	}

	if seconds < 0 || seconds > 59 {
		return fmt.Errorf("invalid seconds: seconds out of range")
	}

	s.Duration = time.Duration((hours*60+minutes)*60+seconds) * time.Second

	return nil
}

func (s Timestamp) String() string {
	hours := int(s.Duration.Hours())
	minutes := int(s.Duration.Minutes()) - 60*hours
	seconds := int(s.Duration.Seconds()) - 60*(minutes+60*hours)
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func NewTimestamp(str string) (Timestamp, error) {
	var sd Timestamp
	err := sd.FromString(str)
	return sd, err
}
