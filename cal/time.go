package cal

import (
	"fmt"
	"time"
)

// Millisecond time formats to comply with W3C datetime format that
// contains rules for local timezones so that they always include a ":".
const (
	RFC3339Milli         string = "2006-01-02T15:04:05.000Z"
	RFC3339MilliWithZone string = "2006-01-02T15:04:05.000-07:00"
)

const (
	nullString = "null"
)

// Time stores when something happened and allows us to
// provide standard formatting.
type Time struct {
	time.Time
}

// Now provides a timestamp for the current UTC system Time
func Now() Time {
	return Time{time.Now().UTC()}
}

// NowIn provides a localised time
func NowIn(loc *time.Location) Time {
	return Time{time.Now().In(loc)}
}

// ParseTime attempts to parse the time string.
func ParseTime(str string) (Time, error) {
	// parse with generic RFC3339 precision, which supports milliseconds
	// and helps us get around issues around timestamps that don't
	// include the milliseconds for whatever reason.
	o, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return Time{}, fmt.Errorf("at: unable to parse timestamp: %w", err)
	}
	return Time{o}, nil
}

// String provides the timestamp in RFC3339 format including milliseconds.
func (t *Time) String() string {
	if t.Time.Location() != time.UTC {
		return t.Time.Format(RFC3339MilliWithZone)
	}
	return t.Time.Format(RFC3339Milli)
}

// UnmarshalJSON uses our time parser.
func (t *Time) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == nullString {
		return nil
	}
	s = s[1 : len(s)-1] // no quotes
	var err error
	*t, err = ParseTime(s)
	return err
}

// MarshalJSON provides the time in JSON format
func (t Time) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte(nullString), nil
	}
	return []byte(`"` + t.String() + `"`), nil
}
