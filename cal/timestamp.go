package cal

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/jsonschema"
)

// RFC3339Milli is the canonical marshaling format used by Timestamp:
// RFC3339 in UTC with a Z suffix and millisecond precision.
const RFC3339Milli = "2006-01-02T15:04:05.000Z"

// Timestamp wraps time.Time with JSON serialization in RFC3339 format, always
// in UTC with a `Z` suffix and millisecond precision. Parsing tolerates
// optional sub-second precision up to nanoseconds and arbitrary offsets,
// always normalizing to UTC.
type Timestamp struct {
	time.Time
}

// TimestampNow returns the current time as a UTC Timestamp.
func TimestampNow() Timestamp {
	return Timestamp{time.Now().UTC()}
}

// TimestampOf returns t as a Timestamp normalized to UTC.
func TimestampOf(t time.Time) Timestamp {
	return Timestamp{t.UTC()}
}

// ParseTimestamp parses an RFC3339 formatted string (sub-second precision
// optional) and returns it as a UTC Timestamp.
func ParseTimestamp(s string) (Timestamp, error) {
	o, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return Timestamp{}, fmt.Errorf("cal: invalid timestamp: %w", err)
	}
	return Timestamp{o.UTC()}, nil
}

// String provides the timestamp in RFC3339 UTC format with millisecond
// precision.
func (t Timestamp) String() string {
	return t.Format(RFC3339Milli)
}

// Clone returns a pointer to a copy of the timestamp.
func (t *Timestamp) Clone() *Timestamp {
	t2 := *t
	return &t2
}

// MarshalJSON emits the timestamp in canonical RFC3339 UTC millisecond format.
// A zero timestamp is emitted as `null`.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + t.String() + `"`), nil
}

// UnmarshalJSON accepts an RFC3339 encoded string (with or without sub-second
// precision, any offset) or a JSON null.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == "null" || s == `""` {
		*t = Timestamp{}
		return nil
	}
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if str == "" {
		*t = Timestamp{}
		return nil
	}
	parsed, err := ParseTimestamp(str)
	if err != nil {
		return err
	}
	*t = parsed
	return nil
}

// JSONSchema returns a custom json schema for the timestamp.
func (Timestamp) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:    "string",
		Format:  "date-time",
		Title:   "Timestamp",
		Pattern: `^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(\.[0-9]+)?Z$`,
		Description: here.Doc(`
			RFC3339 timestamp in UTC with a Z suffix and optional sub-second
			precision, for example: 2021-05-26T13:45:00.000Z
		`),
	}
}

// TimestampTest is used to validate a timestamp according to the provided
// rules.
type TimestampTest struct {
	desc    string
	notZero bool
	after   *Timestamp
	before  *Timestamp
}

// String provides a description of the test.
func (d TimestampTest) String() string { return d.desc }

// Check will perform the defined timestamp test on the provided value.
func (d TimestampTest) Check(value any) bool { return d.Validate(value) == nil }

// Validate is used to check a timestamp's value.
func (d TimestampTest) Validate(value any) error {
	in, ok := value.(Timestamp)
	if !ok {
		ip, ok := value.(*Timestamp)
		if !ok || ip == nil {
			return nil
		}
		in = *ip
	}
	if d.notZero && in.IsZero() {
		return errors.New("required")
	}
	if d.after != nil && !in.After(d.after.Time) {
		return errors.New("too early")
	}
	if d.before != nil && !in.Before(d.before.Time) {
		return errors.New("too late")
	}
	return nil
}

// TimestampNotZero ensures the timestamp is not a zero value.
func TimestampNotZero() TimestampTest {
	return TimestampTest{desc: "not zero", notZero: true}
}

// TimestampAfter returns a validation rule which checks to ensure the
// timestamp is *after* the provided timestamp.
func TimestampAfter(t Timestamp) TimestampTest {
	return TimestampTest{desc: "after " + t.String(), after: &t}
}

// TimestampBefore is used during validation to ensure the timestamp is before
// the value passed in.
func TimestampBefore(t Timestamp) TimestampTest {
	return TimestampTest{desc: "before " + t.String(), before: &t}
}

func timestampRules() *rules.Set {
	return rules.For(new(Timestamp),
		rules.Assert("01", "invalid timestamp", is.Func("valid timestamp", timestampFormatValid)),
	)
}

func timestampFormatValid(val any) bool {
	t, ok := val.(Timestamp)
	if !ok {
		tp, ok := val.(*Timestamp)
		if !ok || tp == nil {
			return false
		}
		t = *tp
	}
	// A time.Time can't be structurally invalid the way a civil date can,
	// so the only way to fail here is to be non-zero yet bogus, which Go's
	// time package does not produce. Always pass.
	_ = t
	return true
}
