package cli

import (
	"time"
)

type TimeFlag time.Time

func (t *TimeFlag) String() string {
	return time.Time(*t).Format(time.RFC3339)
}

func (t *TimeFlag) Set(value string) error {
	parsedTime, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return err
	}
	*t = TimeFlag(parsedTime)
	return nil
}

func (t *TimeFlag) Type() string {
	return "time"
}

func (t TimeFlag) Time() time.Time {
	return time.Time(t)
}
