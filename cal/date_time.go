package cal

import (
	"encoding/json"
	"errors"
	"time"

	"cloud.google.com/go/civil"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/jsonschema"
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
	return ThisSecondIn(time.UTC)
}

// ThisSecondIn provides a new date time using the current time from the provided
// location as a reference.
func ThisSecondIn(loc *time.Location) DateTime {
	t := time.Now().In(loc)
	ct := civil.DateTimeOf(t)
	ct.Time.Nanosecond = 0 // ignore nanoseconds
	return DateTime{ct}
}

// DateTimeOf returns the DateTime from the provided time.
func DateTimeOf(t time.Time) DateTime {
	ct := civil.DateTimeOf(t)
	ct.Time.Nanosecond = 0 // ignore nanoseconds
	return DateTime{ct}
}

// Clone returns a new pointer to a copy of the date time.
func (dt *DateTime) Clone() *DateTime {
	dt2 := *dt
	return &dt2
}

// In returns a new time.Time instance with the provided location.
func (dt DateTime) In(loc *time.Location) time.Time {
	return dt.DateTime.In(loc)
}

// TimeZ returns a new time.Time instance with the UTC location.
func (dt DateTime) TimeZ() time.Time {
	return dt.In(time.UTC)
}

// Date returns the date component of the date time.
func (dt DateTime) Date() Date {
	return Date{
		Date: dt.DateTime.Date,
	}
}

// Time returns the time component of the date time.
func (dt DateTime) Time() Time {
	return Time{
		Time: dt.DateTime.Time,
	}
}

// JSONSchema returns a custom json schema for the date time.
func (DateTime) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:    "string",
		Title:   "Date Time",
		Pattern: "^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}$",
		Description: here.Doc(`
			Civil date time in simplified ISO format with no time zone
			nor location information, for example: 2021-05-26T13:45:00
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

// DateTimeTest is used to validate a date time according to the provided rules.
type DateTimeTest struct {
	desc    string
	notZero bool
	after   *DateTime
	before  *DateTime
}

// String provides a description of the test.
func (d DateTimeTest) String() string { return d.desc }

// Check will perform the defined date time test on the provided value.
func (d DateTimeTest) Check(value any) bool { return d.Validate(value) == nil }

// Validate is used to check a date time's value.
func (d DateTimeTest) Validate(value any) error {
	in, ok := value.(DateTime)
	if !ok {
		inp, ok := value.(*DateTime)
		if !ok || inp == nil {
			return nil
		}
		in = *inp
	}
	if d.notZero && in.IsZero() {
		return errors.New("required")
	}
	if d.after != nil && !in.After(d.after.DateTime) {
		return errors.New("too early")
	}
	if d.before != nil && !in.Before(d.before.DateTime) {
		return errors.New("too late")
	}
	return nil
}

func dateTimeRules() *rules.Set {
	return rules.For(new(DateTime),
		rules.Assert("01", "invalid date time", is.Func("valid date time", dateTimeFormatValid)),
	)
}

func dateTimeFormatValid(val any) bool {
	d, ok := val.(DateTime)
	if !ok {
		dp, ok := val.(*DateTime)
		if !ok || dp == nil {
			return false
		}
		d = *dp
	}
	return d.IsZero() || d.IsValid()
}

// DateTimeNotZero ensures the date time is not a zero value.
func DateTimeNotZero() DateTimeTest {
	return DateTimeTest{desc: "not zero", notZero: true}
}

// DateTimeAfter returns a validation rule which checks to ensure the date time
// is *after* the provided date time.
func DateTimeAfter(dt DateTime) DateTimeTest {
	return DateTimeTest{desc: "after " + dt.String(), after: &dt}
}

// DateTimeBefore is used during validation to ensure the date time is before
// the value passed in.
func DateTimeBefore(dt DateTime) DateTimeTest {
	return DateTimeTest{desc: "before " + dt.String(), before: &dt}
}
