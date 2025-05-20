package cal

import (
	"encoding/json"
	"time"

	"cloud.google.com/go/civil"
	"github.com/invopop/jsonschema"
)

// Time represents a simple time of day without a date component.
type Time struct {
	civil.Time
}

// NewTime provides a pointer to a new time instance.
func NewTime(hour, minute, second int) *Time {
	t := MakeTime(hour, minute, second)
	return &t
}

// MakeTime provides a new time instance.
func MakeTime(hour, minute, second int) Time {
	return Time{
		civil.Time{
			Hour:   hour,
			Minute: minute,
			Second: second,
		},
	}
}

// TimeNow generates a new time instance for now.
func TimeNow() Time {
	return TimeNowIn(time.UTC)
}

// TimeNowIn provides the current time of the day in the provided
// location.
func TimeNowIn(loc *time.Location) Time {
	t := time.Now().In(loc)
	ct := civil.TimeOf(t)
	ct.Nanosecond = 0 // ignore nanoseconds
	return Time{ct}
}

// UnmarshalJSON is used to parse a time from json and ensures that
// we can handle invalid data reasonably.
func (t *Time) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		*t = Time{}
		return nil
	}
	dt, err := civil.ParseTime(s)
	if err != nil {
		return err
	}
	*t = Time{dt}
	return nil
}

// JSONSchema returns a custom json schema for the current hour, minute,
// and seconds of the day.
func (Time) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Title:       "Time",
		Pattern:     `^([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$`,
		Description: "Civil time in simplified ISO format, like 13:45:30",
	}
}
