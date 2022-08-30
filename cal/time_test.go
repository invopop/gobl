package cal_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/invopop/gobl/cal"
)

func TestTimeUnmarshal(t *testing.T) {
	tl, _ := time.LoadLocation("America/Lima") // always -5 (no DST)
	tl2, _ := time.LoadLocation("Asia/Dubai")  // always +4 (no DST)
	var cases = []struct {
		Given    string
		Expected time.Time
		Error    bool
	}{
		// long form
		{`"2009-11-10T23:19:45.123Z"`, time.Date(2009, time.November, 10, 23, 19, 45, 123000000, time.UTC), false},
		// short form
		{`"2009-11-10T13:19:04Z"`, time.Date(2009, time.November, 10, 13, 19, 4, 0, time.UTC), false},
		// local form
		{`"2009-11-10T13:19:04+02:00"`, time.Date(2009, time.November, 10, 11, 19, 4, 0, time.UTC), false},
		// bad form
		{`"Z2009-11-10T13:19:04Z"`, time.Time{}, true},
		// nil string
		{`null`, time.Time{}, false},
		// local long form
		{`"2009-11-10T23:19:45.123-05:00"`, time.Date(2009, time.November, 10, 23, 19, 45, 123000000, tl), false},
		// local long form 2
		{`"2009-11-10T23:19:45.123+04:00"`, time.Date(2009, time.November, 10, 23, 19, 45, 123000000, tl2), false},
		// local short form
		{`"2009-11-10T23:19:45-05:00"`, time.Date(2009, time.November, 10, 23, 19, 45, 0, tl), false},
		// local bad long form
		{`"2009-11-10T23:19:45.123+0400"`, time.Time{}, true},
	}

	for _, c := range cases {
		payload := []byte(c.Given)
		var output cal.Time
		if err := json.Unmarshal(payload, &output); err != nil && !c.Error {
			t.Error(err)
			continue
		}
		if !output.Equal(c.Expected) {
			t.Errorf("Expected: %q, Given: %q", c.Expected, output)
		}
	}
}

func TestTimeMarshal(t *testing.T) {
	var cases = []struct {
		Given    time.Time
		Expected string
	}{
		{time.Date(2009, time.November, 10, 23, 19, 45, 0, time.UTC), `"2009-11-10T23:19:45.000Z"`},
		{time.Date(2009, time.November, 10, 13, 19, 4, 0, time.UTC), `"2009-11-10T13:19:04.000Z"`},
		{time.Date(2009, time.November, 10, 23, 19, 45, 123456000, time.UTC), `"2009-11-10T23:19:45.123Z"`},
		{time.Time{}, `null`},
	}

	for _, c := range cases {
		ct := cal.Time{c.Given}
		output, err := json.Marshal(ct)
		if err != nil {
			t.Error(err)
		}
		if string(output) != c.Expected {
			t.Errorf("Expected: %q, Got: %q", c.Expected, output)
		}
	}
}

func TestNow(t *testing.T) {
	ct := cal.Now()
	if z, _ := ct.Zone(); z != "UTC" {
		t.Errorf("Failed to get current time in UTC, got: %v", z)
	}
}

func TestLocalTimeNow(t *testing.T) {
	loc, _ := time.LoadLocation("America/Lima")
	ct := cal.NowIn(loc)
	if z, _ := ct.Zone(); z != "-05" {
		t.Errorf("Failed to get expected time zone, got: %v", z)
	}
}

func TestLocalTimeMarshal(t *testing.T) {
	tl, _ := time.LoadLocation("America/Lima") // always -5 (no DST)
	tl2, _ := time.LoadLocation("Asia/Dubai")  // always +4 (no DST)
	var cases = []struct {
		Given    time.Time
		Expected string
	}{
		{time.Date(2009, time.November, 10, 23, 19, 30, 0, tl), `"2009-11-10T23:19:30.000-05:00"`},
		{time.Date(2009, time.November, 10, 13, 19, 4, 123000000, tl), `"2009-11-10T13:19:04.123-05:00"`},
		{time.Date(2009, time.November, 10, 13, 19, 4, 123000000, tl2), `"2009-11-10T13:19:04.123+04:00"`},
		{time.Time{}, `null`},
	}

	for _, c := range cases {
		ct := cal.Time{c.Given}
		output, err := json.Marshal(ct)
		if err != nil {
			t.Error(err)
		}
		if string(output) != c.Expected {
			t.Errorf("Expected: %q, Given: %q", c.Expected, output)
		}
	}
}

func TestTimestampInModel(t *testing.T) {
	type tmodel struct {
		Value     string    `json:"v"`
		ExampleAt cal.Time  `json:"example_at"`
		EmptyAt   *cal.Time `json:"empty_at,omitempty"`
	}
	x := new(tmodel)
	x.Value = "bar"

	data, err := json.Marshal(x)
	if err != nil {
		t.Error(err)
		return
	}
	if !strings.Contains(string(data), `"example_at":null`) {
		t.Errorf("Expected output to contain value example_at, got: %v", string(data))
	}
	if strings.Contains(string(data), "empty_at") {
		t.Errorf("Did not expect output to contain value, got: %v", string(data))
	}
}
