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
	// empty is set when the time was parsed from an empty string, or
	// created with EmptyTime, and indicates that the time has been
	// deliberately left undefined so that it can be filled in later.
	// This makes it possible to distinguish between an empty time and
	// midnight (00:00:00).
	empty bool
}

// NewTime provides a pointer to a new time instance.
func NewTime(hour, minute, second int) *Time {
	t := MakeTime(hour, minute, second)
	return &t
}

// MakeTime provides a new time instance.
func MakeTime(hour, minute, second int) Time {
	return Time{
		Time: civil.Time{
			Hour:   hour,
			Minute: minute,
			Second: second,
		},
	}
}

// EmptyTime provides a pointer to a new empty time instance. An empty
// time is deliberately undefined and is expected to be replaced by
// the current time during calculations, as opposed to midnight which
// is a valid time of day. It has the same JSON representation as an
// empty string.
func EmptyTime() *Time {
	return &Time{empty: true}
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
	return Time{Time: ct}
}

// IsZero returns true if the time is the zero value (00:00:00),
// which is also true for empty times. This is used by the `omitzero`
// JSON tag to determine if the time should be omitted from JSON output.
// Use IsEmpty to determine if a time was deliberately left undefined.
func (t Time) IsZero() bool {
	return t.Time.IsZero()
}

// IsEmpty returns true if the time was deliberately left undefined,
// either by parsing an empty string or by using EmptyTime. A time
// of midnight (00:00:00) is not considered empty.
func (t Time) IsEmpty() bool {
	return t.empty
}

// MarshalJSON ensures that empty times are output as an empty string
// so that they can be distinguished from midnight.
func (t Time) MarshalJSON() ([]byte, error) {
	if t.empty {
		return json.Marshal("")
	}
	return json.Marshal(t.Time.String())
}

// UnmarshalJSON is used to parse a time from json and ensures that
// we can handle invalid data reasonably.
func (t *Time) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		*t = Time{empty: true}
		return nil
	}
	dt, err := civil.ParseTime(s)
	if err != nil {
		return err
	}
	*t = Time{Time: dt}
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
