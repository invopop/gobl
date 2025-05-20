package cal

import (
	"encoding/json"
	"errors"
	"time"

	"cloud.google.com/go/civil"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Date represents a simple date without time used most frequently
// with business documents.
type Date struct {
	civil.Date
}

// NewDate provides a pointer to a new date instance.
func NewDate(year int, month time.Month, day int) *Date {
	d := MakeDate(year, month, day)
	return &d
}

// MakeDate provides a new date instance.
func MakeDate(year int, month time.Month, day int) Date {
	return Date{
		civil.Date{
			Year:  year,
			Month: month,
			Day:   day,
		},
	}
}

// Today generates a new date instance for today.
func Today() Date {
	t := time.Now().UTC()
	return Date{
		civil.DateOf(t),
	}
}

// TodayIn generates a new date instance for today in the given location.
func TodayIn(loc *time.Location) Date {
	t := time.Now().In(loc)
	return Date{
		civil.DateOf(t),
	}
}

// DateOf returns the Date in which a time occurs in the time's location.
func DateOf(t time.Time) Date {
	return Date{
		civil.DateOf(t),
	}
}

// Validate ensures the the date object looks valid.
func (d Date) Validate() error {
	if d.IsZero() {
		return nil // there is a specific test for this
	}
	if !d.IsValid() {
		return errors.New("invalid date")
	}
	return nil
}

// Clone returns a new pointer to a copy of the date.
func (d *Date) Clone() *Date {
	d2 := *d
	return &d2
}

// Time returns a time object for the date.
func (d Date) Time() time.Time {
	return d.TimeIn(time.UTC)
}

// TimeIn returns a time object for the date in the given location which
// may be important when considering timezones and arithmetic.
func (d Date) TimeIn(loc *time.Location) time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, loc)
}

// Add returns a new date with the given number of years, months and days.
// This uses the time package to do the arithmetic.
func (d Date) Add(years, months, days int) Date {
	t := d.Time()
	t = t.AddDate(years, months, days)
	return DateOf(t)
}

// WithTime appends the time to the date to create a DateTime object.
func (d Date) WithTime(t Time) DateTime {
	return MakeDateTime(d.Year, d.Month, d.Day, t.Hour, t.Minute, t.Second)
}

// UnmarshalJSON is used to parse a date from json and ensures that
// we can handle invalid data reasonably.
func (d *Date) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	// Zero dates are not great, put pass validation.
	if s == "0000-00-00" {
		*d = Date{}
		return nil
	}
	dt, err := civil.ParseDate(s)
	if err != nil {
		return err
	}
	*d = Date{dt}
	return nil
}

// JSONSchema returns a custom json schema for the date.
func (Date) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Format:      "date",
		Title:       "Date",
		Description: "Civil date in simplified ISO format, like 2021-05-26",
	}
}

type dateValidationRule struct {
	notZero bool
	after   *Date
	before  *Date
}

// Validate is used to check a dates value.
func (d *dateValidationRule) Validate(value interface{}) error {
	in, ok := value.(Date)
	if !ok {
		inp, ok := value.(*Date)
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
		if in.DaysSince(d.after.Date) < 0 {
			return errors.New("too early")
		}
	}
	if d.before != nil {
		if in.DaysSince(d.before.Date) > 0 {
			return errors.New("too late")
		}
	}
	return nil
}

// DateNotZero ensures the date is not a zero value.
func DateNotZero() validation.Rule {
	return &dateValidationRule{
		notZero: true,
	}
}

// DateAfter returns a validation rule which checks to ensure the date
// is *after* the provided date.
func DateAfter(date Date) validation.Rule {
	return &dateValidationRule{
		after: &date,
	}
}

// DateBefore is used during validation to ensure the date is before
// the value passed in.
func DateBefore(date Date) validation.Rule {
	return &dateValidationRule{
		before: &date,
	}
}
