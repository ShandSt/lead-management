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
	// Convert data to a string and remove surrounding quotes.
	strInput := string(data)
	strInput = strInput[1 : len(strInput)-1] // strip the surrounding quotes

	// If time ends with 'Z', it's in UTC; replace 'Z' with "+00:00" for Go's parsing compatibility
	if strInput[len(strInput)-1] == 'Z' {
		strInput = strInput[:len(strInput)-1] + "+00:00"
	}

	// Parse the time string assuming RFC3339 format, which includes the timezone
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
