package types

import (
	"fmt"
	"strings"
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
	strInput = strings.Trim(strInput, "\"")

	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"15:04:05",
	}

	if strInput[len(strInput)-1] == 'Z' {
		strInput = strInput[:len(strInput)-1] + "+00:00"
	}

	var err error
	var parsedTime time.Time

	for _, format := range formats {
		parsedTime, err = time.Parse(format, strInput)
		if err == nil {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("error parsing time '%s': %w", strInput, err)
	}

	ct.Time = parsedTime
	return nil
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(time.RFC3339))), nil
}

func (wh WorkingHours) Contains(t CustomTime) bool {
	return !t.Time.Before(wh.Start.Time) && !t.Time.After(wh.End.Time)
}
