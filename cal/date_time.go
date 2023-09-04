package cal

import (
	"encoding/json"
	"errors"
	"time"

	"cloud.google.com/go/civil"
	"github.com/invopop/jsonschema"
)

// DateTime represents a date with time but without timezone.
type DateTime struct {
	civil.DateTime
}

// NewDateTime provides a pointer to a new datetime instance.
func NewDateTime(year int, month time.Month, day, hour, min, sec int, nsec int) *DateTime {
	dt := MakeDateTime(year, month, day, hour, min, sec, nsec)
	return &dt
}

// MakeDateTime provides a new datetime instance.
func MakeDateTime(year int, month time.Month, day, hour, min, sec int, nsec int) DateTime {
	return DateTime{
		civil.DateTime{
			Date: civil.Date{
				Year:  year,
				Month: month,
				Day:   day,
			},
			Time: civil.Time{
				Hour:       hour,
				Minute:     min,
				Second:     sec,
				Nanosecond: nsec,
			},
		},
	}
}

// Validate ensures the the date time object looks valid.
func (dt DateTime) Validate() error {
	if dt.IsZero() {
		return nil // there is a specific test for this
	}
	if !dt.DateTime.IsValid() {
		return errors.New("invalid date")
	}
	return nil
}

// UnmarshalJSON is used to parse a date from json and ensures that
// we can handle invalid data reasonably.
func (dt *DateTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	// Zero dates are not great, put pass validation.
	if s == "0000-00-00T00:00:00" {
		*dt = DateTime{}
		return nil
	}
	pdt, err := civil.ParseDateTime(s)
	if err != nil {
		return err
	}
	*dt = DateTime{pdt}
	return nil
}

// JSONSchema returns a custom json schema for the date.
func (DateTime) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Format:      "date-time",
		Title:       "DateTime",
		Description: "Civil datetime in ISO 8601 format, like 2021-05-26T12:20:30.123",
	}
}
