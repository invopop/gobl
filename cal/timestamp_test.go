package cal_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimestampJSONUnmarshal(t *testing.T) {
	cases := []struct {
		name     string
		given    string
		expected time.Time
		wantErr  bool
	}{
		{"null", `null`, time.Time{}, false},
		{"empty string", `""`, time.Time{}, false},
		{"short form Z", `"2025-03-15T10:00:00Z"`,
			time.Date(2025, time.March, 15, 10, 0, 0, 0, time.UTC), false},
		{"milli form Z", `"2025-03-15T10:00:00.123Z"`,
			time.Date(2025, time.March, 15, 10, 0, 0, 123000000, time.UTC), false},
		{"nano form Z", `"2025-03-15T10:00:00.123456789Z"`,
			time.Date(2025, time.March, 15, 10, 0, 0, 123456789, time.UTC), false},
		{"with offset normalized to UTC", `"2025-03-15T12:00:00+02:00"`,
			time.Date(2025, time.March, 15, 10, 0, 0, 0, time.UTC), false},
		{"with offset and nanos", `"2025-03-15T12:00:00.5-02:00"`,
			time.Date(2025, time.March, 15, 14, 0, 0, 500000000, time.UTC), false},
		{"bad form", `"not a timestamp"`, time.Time{}, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got cal.Timestamp
			err := json.Unmarshal([]byte(c.given), &got)
			if c.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.True(t, got.Time.Equal(c.expected), "got %s want %s", got.Time, c.expected)
			if !got.IsZero() {
				z, _ := got.Zone()
				assert.Equal(t, "UTC", z, "timestamp should be in UTC")
			}
		})
	}
}

func TestTimestampJSONMarshal(t *testing.T) {
	cases := []struct {
		name     string
		given    time.Time
		expected string
	}{
		{"zero", time.Time{}, `null`},
		{"utc integer seconds", time.Date(2025, time.March, 15, 10, 0, 0, 0, time.UTC),
			`"2025-03-15T10:00:00.000Z"`},
		{"utc with millis", time.Date(2025, time.March, 15, 10, 0, 0, 123000000, time.UTC),
			`"2025-03-15T10:00:00.123Z"`},
		{"nanos truncate to millis", time.Date(2025, time.March, 15, 10, 0, 0, 123456789, time.UTC),
			`"2025-03-15T10:00:00.123Z"`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ts := cal.Timestamp{Time: c.given}
			out, err := json.Marshal(ts)
			require.NoError(t, err)
			assert.Equal(t, c.expected, string(out))
		})
	}
}

func TestTimestampMarshalOfNonUTCInput(t *testing.T) {
	// Timestamp holds whatever time is given to it; marshaling formats
	// that value against the fixed UTC layout. Callers should normalize
	// via TimestampOf or ParseTimestamp before trusting the output.
	loc, err := time.LoadLocation("Europe/Madrid")
	require.NoError(t, err)
	raw := time.Date(2025, time.March, 15, 12, 0, 0, 0, loc)

	// TimestampOf normalizes to UTC, so the wall-clock shifts by the offset.
	ts := cal.TimestampOf(raw)
	out, err := json.Marshal(ts)
	require.NoError(t, err)
	assert.Equal(t, `"2025-03-15T11:00:00.000Z"`, string(out))
}

func TestTimestampRoundTrip(t *testing.T) {
	data := `"2025-03-15T10:00:00.000Z"`
	var ts cal.Timestamp
	require.NoError(t, json.Unmarshal([]byte(data), &ts))
	out, err := json.Marshal(ts)
	require.NoError(t, err)
	assert.Equal(t, data, string(out))
}

func TestTimestampNow(t *testing.T) {
	ts := cal.TimestampNow()
	z, _ := ts.Zone()
	assert.Equal(t, "UTC", z)
	assert.WithinDuration(t, time.Now(), ts.Time, time.Second)
}

func TestTimestampOf(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)
	raw := time.Date(2025, time.March, 15, 6, 0, 0, 0, loc)
	ts := cal.TimestampOf(raw)
	z, _ := ts.Zone()
	assert.Equal(t, "UTC", z)
	// 06:00 New York (EDT, UTC-4) = 10:00 UTC
	assert.Equal(t, 10, ts.Hour())
}

func TestParseTimestamp(t *testing.T) {
	t.Run("short Z form", func(t *testing.T) {
		ts, err := cal.ParseTimestamp("2025-03-15T10:00:00Z")
		require.NoError(t, err)
		assert.Equal(t, "2025-03-15T10:00:00.000Z", ts.String())
	})
	t.Run("milli Z form", func(t *testing.T) {
		ts, err := cal.ParseTimestamp("2025-03-15T10:00:00.5Z")
		require.NoError(t, err)
		assert.Equal(t, "2025-03-15T10:00:00.500Z", ts.String())
	})
	t.Run("offset form normalizes to UTC", func(t *testing.T) {
		ts, err := cal.ParseTimestamp("2025-03-15T12:00:00+02:00")
		require.NoError(t, err)
		assert.Equal(t, "2025-03-15T10:00:00.000Z", ts.String())
	})
	t.Run("invalid", func(t *testing.T) {
		_, err := cal.ParseTimestamp("not a timestamp")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid timestamp")
	})
}

func TestTimestampString(t *testing.T) {
	ts := cal.TimestampOf(time.Date(2025, time.March, 15, 10, 0, 0, 0, time.UTC))
	assert.Equal(t, "2025-03-15T10:00:00.000Z", ts.String())

	ts = cal.TimestampOf(time.Date(2025, time.March, 15, 10, 0, 0, 123456789, time.UTC))
	assert.Equal(t, "2025-03-15T10:00:00.123Z", ts.String())
}

func TestTimestampClone(t *testing.T) {
	orig := cal.TimestampOf(time.Date(2025, time.March, 15, 10, 0, 0, 0, time.UTC))
	cpy := orig.Clone()
	assert.Equal(t, orig.String(), cpy.String())

	orig = cal.TimestampOf(time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC))
	assert.NotEqual(t, orig.String(), cpy.String())
}

func TestTimestampValidation(t *testing.T) {
	t.Run("rules.Validate passes on zero and non-zero", func(t *testing.T) {
		assert.NoError(t, rules.Validate(cal.Timestamp{}))
		assert.NoError(t, rules.Validate(cal.TimestampNow()))

		var nilp *cal.Timestamp
		assert.NoError(t, rules.Validate(nilp))
	})

	t.Run("not zero", func(t *testing.T) {
		assert.False(t, cal.TimestampNotZero().Check(cal.Timestamp{}))
		assert.True(t, cal.TimestampNotZero().Check(cal.TimestampNow()))

		// nil pointer returns nil error per the embedded convention
		var nilp *cal.Timestamp
		assert.True(t, cal.TimestampNotZero().Check(nilp))
	})

	t.Run("after", func(t *testing.T) {
		base := cal.TimestampOf(time.Date(2025, time.March, 15, 10, 0, 0, 0, time.UTC))
		later := cal.TimestampOf(time.Date(2025, time.March, 15, 11, 0, 0, 0, time.UTC))
		earlier := cal.TimestampOf(time.Date(2025, time.March, 15, 9, 0, 0, 0, time.UTC))

		assert.True(t, cal.TimestampAfter(base).Check(later))
		assert.False(t, cal.TimestampAfter(base).Check(earlier))
		assert.False(t, cal.TimestampAfter(base).Check(base))
	})

	t.Run("before", func(t *testing.T) {
		base := cal.TimestampOf(time.Date(2025, time.March, 15, 10, 0, 0, 0, time.UTC))
		earlier := cal.TimestampOf(time.Date(2025, time.March, 15, 9, 0, 0, 0, time.UTC))
		later := cal.TimestampOf(time.Date(2025, time.March, 15, 11, 0, 0, 0, time.UTC))

		assert.True(t, cal.TimestampBefore(base).Check(earlier))
		assert.False(t, cal.TimestampBefore(base).Check(later))
		assert.False(t, cal.TimestampBefore(base).Check(base))
	})
}

func TestTimestampJSONSchema(t *testing.T) {
	schema := cal.Timestamp{}.JSONSchema()
	out, err := json.Marshal(schema)
	require.NoError(t, err)
	s := string(out)
	assert.Contains(t, s, `"type":"string"`)
	assert.Contains(t, s, `"format":"date-time"`)
	assert.Contains(t, s, `"title":"Timestamp"`)
	assert.True(t, strings.Contains(s, `Z$`), "schema pattern should require Z suffix")
}

func TestTimestampInModel(t *testing.T) {
	// Mirrors lib/at: zero timestamp marshals as null; a pointer field with
	// omitempty is omitted entirely when nil.
	type model struct {
		At     cal.Timestamp  `json:"at"`
		AtPtr  *cal.Timestamp `json:"at_ptr,omitempty"`
	}

	data, err := json.Marshal(&model{})
	require.NoError(t, err)
	assert.Contains(t, string(data), `"at":null`)
	assert.NotContains(t, string(data), `at_ptr`)
}
