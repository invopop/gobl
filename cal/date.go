package cal

import (
	"errors"
	"time"

	"cloud.google.com/go/civil"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/jsonschema"
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
