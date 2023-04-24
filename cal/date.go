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

// Validate ensures the the date object looks valid.
func (d Date) Validate() error {
	if d.IsZero() {
		return nil // there is a specific test for this
	}
	if !d.Date.IsValid() {
		return errors.New("invalid date")
	}
	return nil
}

// Clone returns a new pointer to a copy of the date.
func (d *Date) Clone() *Date {
	d2 := *d
	return &d2
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
		if in.Date.IsZero() {
			return errors.New("required")
		}
	}
	if d.after != nil {
		if in.Date.DaysSince(d.after.Date) < 0 {
			return errors.New("too early")
		}
	}
	if d.before != nil {
		if in.Date.DaysSince(d.before.Date) > 0 {
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
