package types

import (
	"fmt"
	"time"
)

type Client struct {
	ID           string       `json:"id,omitempty"`
	Name         string       `json:"name"`
	WorkingHours WorkingHours `json:"workingHours"`
	Priority     int          `json:"priority"`
	LeadCount    int          `json:"leadCount"`
	Capacity     int          `json:"capacity"`
}

type WorkingHours struct {
	Start CustomTime `json:"start"`
	End   CustomTime `json:"end"`
}

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	strInput := string(data)
	strInput = strInput[1 : len(strInput)-1]

	if strInput[len(strInput)-1] == 'Z' {
		strInput = strInput[:len(strInput)-1] + "+00:00"
	}

	parsedTime, err := time.Parse(time.RFC3339, strInput)
	if err != nil {
		return fmt.Errorf("error parsing time: %w", err)
	}
	ct.Time = parsedTime
	return nil
}

func (wh WorkingHours) Contains(t CustomTime) bool {
	return !t.Time.Before(wh.Start.Time) && !t.Time.After(wh.End.Time)
}
