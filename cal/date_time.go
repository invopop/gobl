package cal

import (
	"encoding/json"
	"errors"
	"time"

	"cloud.google.com/go/civil"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// DateTime represents a combination of date and time without location
// specific information nor support for millisecond precision.
type DateTime struct {
	civil.DateTime
}

// NewDateTime provides a pointer to a new date time instance.
func NewDateTime(year int, month time.Month, day, hour, minute, second int) *DateTime {
	dt := MakeDateTime(year, month, day, hour, minute, second)
	return &dt
}

// MakeDateTime provides a new date time instance.
func MakeDateTime(year int, month time.Month, day, hour, minute, second int) DateTime {
	return DateTime{
		civil.DateTime{
			Date: civil.Date{
				Year:  year,
				Month: month,
				Day:   day,
			},
			Time: civil.Time{
				Hour:   hour,
				Minute: minute,
				Second: second,
			},
		},
	}
}

// ThisSecond produces a new date time instance for the current UTC time
// to the nearest second.
func ThisSecond() DateTime {
	t := time.Now().UTC()
	return DateTime{
		civil.DateTimeOf(t),
	}
}

// ThisSecondIn provides a new date time using the current time from the provided
// location as a reference.
func ThisSecondIn(loc *time.Location) DateTime {
	t := time.Now().In(loc)
	return DateTime{
		civil.DateTimeOf(t),
	}
}

// DateTimeOf returns the DateTime from the provided time.
func DateTimeOf(t time.Time) DateTime {
	return DateTime{
		civil.DateTimeOf(t),
	}
}

// Clone returns a new pointer to a copy of the date time.
func (dt *DateTime) Clone() *DateTime {
	dt2 := *dt
	return &dt2
}

// Validate ensures the date time object looks valid
func (dt DateTime) Validate() error {
	if dt.IsZero() {
		return nil
	}
	if !dt.DateTime.IsValid() {
		return errors.New("invalid date time")
	}
	return nil
}

// In returns a new time.Time instance with the provided location.
func (dt DateTime) In(loc *time.Location) time.Time {
	return dt.DateTime.In(loc)
}

// TimeZ returns a new time.Time instance with the UTC location.
func (dt DateTime) TimeZ() time.Time {
	return dt.In(time.UTC)
}

// JSONSchema returns a custom json schema for the date time.
func (DateTime) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:   "string",
		Format: "date-time",
		Title:  "Date Time",
		Description: here.Doc(`
			Civil date time in simplified ISO format with no time zone
			information, for example: 2021-05-26T13:45:00
		`),
	}
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
	dtn, err := civil.ParseDateTime(s)
	if err != nil {
		return err
	}
	*dt = DateTime{dtn}
	return nil
}

type dateTimeValidationRule struct {
	notZero bool
	after   *DateTime
	before  *DateTime
}

// Validate is used to check a date time's value.
func (d *dateTimeValidationRule) Validate(value interface{}) error {
	in, ok := value.(DateTime)
	if !ok {
		inp, ok := value.(*DateTime)
		if !ok {
			return nil
		}
		if inp == nil {
			return nil
		}
		in = *inp
	}
	if d.notZero {
		if in.IsZero() {
			return errors.New("required")
		}
	}
	if d.after != nil {
		if !in.DateTime.After(d.after.DateTime) {
			return errors.New("too early")
		}
	}
	if d.before != nil {
		if !in.DateTime.Before(d.before.DateTime) {
			return errors.New("too late")
		}
	}
	return nil
}

// DateTimeNotZero ensures the date is not a zero value.
func DateTimeNotZero() validation.Rule {
	return &dateTimeValidationRule{
		notZero: true,
	}
}

// DateTimeAfter returns a validation rule which checks to ensure the date
// is *after* the provided date.
func DateTimeAfter(dt DateTime) validation.Rule {
	return &dateTimeValidationRule{
		after: &dt,
	}
}

// DateTimeBefore is used during validation to ensure the date is before
// the value passed in.
func DateTimeBefore(dt DateTime) validation.Rule {
	return &dateTimeValidationRule{
		before: &dt,
	}
}
