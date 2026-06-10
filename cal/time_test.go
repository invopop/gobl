package cal_test

import (
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/invopop/gobl/cal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeMake(t *testing.T) {
	t.Run("valid time", func(t *testing.T) {
		tm := cal.MakeTime(12, 34, 56)
		assert.Equal(t, tm.Hour, 12)
		assert.Equal(t, tm.Minute, 34)
		assert.Equal(t, tm.Second, 56)
	})
}

func TestTimeNew(t *testing.T) {
	t.Run("valid time", func(t *testing.T) {
		tm := cal.NewTime(12, 34, 56)
		assert.Equal(t, tm.Hour, 12)
		assert.Equal(t, tm.Minute, 34)
		assert.Equal(t, tm.Second, 56)
	})
}

func TestTimeNow(t *testing.T) {
	t.Run("matches current UTC wall clock", func(t *testing.T) {
		before := time.Now().UTC()
		tm := cal.TimeNow()
		after := time.Now().UTC()
		assertTimeWithin(t, tm, before, after)
	})
}

func TestTimeNowIn(t *testing.T) {
	t.Run("matches current wall clock in location", func(t *testing.T) {
		loc, err := time.LoadLocation("America/New_York")
		require.NoError(t, err)
		before := time.Now().In(loc)
		tm := cal.TimeNowIn(loc)
		after := time.Now().In(loc)
		assertTimeWithin(t, tm, before, after)
	})
}

// assertTimeWithin checks that tm is a valid time-of-day falling between the
// wall-clock samples taken immediately before and after the call, tolerating
// a midnight rollover.
func assertTimeWithin(t *testing.T, tm cal.Time, before, after time.Time) {
	t.Helper()
	assert.GreaterOrEqual(t, tm.Hour, 0)
	assert.Less(t, tm.Hour, 24)
	assert.GreaterOrEqual(t, tm.Minute, 0)
	assert.Less(t, tm.Minute, 60)
	assert.GreaterOrEqual(t, tm.Second, 0)
	assert.Less(t, tm.Second, 60)
	assert.Zero(t, tm.Nanosecond, "nanoseconds should be stripped")

	secOfDay := func(h, m, s int) int { return h*3600 + m*60 + s }
	got := secOfDay(tm.Hour, tm.Minute, tm.Second)
	lo := secOfDay(before.Hour(), before.Minute(), before.Second())
	hi := secOfDay(after.Hour(), after.Minute(), after.Second())
	if hi < lo {
		assert.True(t, got >= lo || got <= hi, "time %v outside window %v..%v across midnight", tm, before, after)
	} else {
		assert.True(t, got >= lo && got <= hi, "time %v outside window %v..%v", tm, before, after)
	}
}

func TestTimeString(t *testing.T) {
	t.Run("valid time", func(t *testing.T) {
		tm := cal.MakeTime(12, 34, 56)
		assert.Equal(t, tm.String(), "12:34:56")
	})
	t.Run("with location", func(t *testing.T) {
		tl, err := time.LoadLocation("America/New_York")
		require.NoError(t, err)
		tm := cal.TimeNowIn(tl)
		assert.Regexp(t, regexp.MustCompile(`^\d{2}:\d{2}:\d{2}$`), tm.String())
	})
}

func TestTimeParsing(t *testing.T) {
	t.Run("zero time", func(t *testing.T) {
		var tm cal.Time
		err := json.Unmarshal([]byte(`"00:00:00"`), &tm)
		require.NoError(t, err)
		assert.Equal(t, tm.Hour, 0)
		assert.Equal(t, tm.Minute, 0)
		assert.Equal(t, tm.Second, 0)
		assert.True(t, tm.IsZero())
	})
	t.Run("empty time", func(t *testing.T) {
		var tm cal.Time
		err := json.Unmarshal([]byte(`""`), &tm)
		require.NoError(t, err)
		assert.Equal(t, tm.Hour, 0)
		assert.Equal(t, tm.Minute, 0)
		assert.Equal(t, tm.Second, 0)
		assert.True(t, tm.IsZero())
	})
	t.Run("zero time, no seconds", func(t *testing.T) {
		var tm cal.Time
		err := json.Unmarshal([]byte(`"00:00"`), &tm)
		assert.ErrorContains(t, err, `parsing time "00:00" as "15:04:05.999999999": cannot parse "" as ":"`)
	})

	t.Run("valid time", func(t *testing.T) {
		var tm cal.Time
		err := json.Unmarshal([]byte(`"12:34:56"`), &tm)
		require.NoError(t, err)
		assert.Equal(t, tm.Hour, 12)
		assert.Equal(t, tm.Minute, 34)
		assert.Equal(t, tm.Second, 56)
	})

	t.Run("invalid json", func(t *testing.T) {
		var tm cal.Time
		err := tm.UnmarshalJSON([]byte(`"12:34:56`))
		require.ErrorContains(t, err, "unexpected end of JSON input")
	})
}

func TestTimeIsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		var tm cal.Time
		assert.True(t, tm.IsZero())
	})
	t.Run("midnight", func(t *testing.T) {
		tm := cal.MakeTime(0, 0, 0)
		assert.True(t, tm.IsZero())
	})
	t.Run("non-zero", func(t *testing.T) {
		tm := cal.MakeTime(12, 30, 0)
		assert.False(t, tm.IsZero())
	})
	t.Run("only seconds", func(t *testing.T) {
		tm := cal.MakeTime(0, 0, 1)
		assert.False(t, tm.IsZero())
	})
}

func TestTimeOmitZero(t *testing.T) {
	type testStruct struct {
		Name string   `json:"name"`
		Time cal.Time `json:"time,omitzero"`
	}
	t.Run("omits zero time", func(t *testing.T) {
		s := testStruct{Name: "test"}
		data, err := json.Marshal(s)
		require.NoError(t, err)
		assert.JSONEq(t, `{"name":"test"}`, string(data))
	})
	t.Run("includes non-zero time", func(t *testing.T) {
		s := testStruct{
			Name: "test",
			Time: cal.MakeTime(14, 30, 0),
		}
		data, err := json.Marshal(s)
		require.NoError(t, err)
		assert.JSONEq(t, `{"name":"test","time":"14:30:00"}`, string(data))
	})
}

func TestTimeJSONSChema(t *testing.T) {
	// Check the schema for the time type.
	schema := cal.Time{}.JSONSchema()
	out, err := json.Marshal(schema)
	require.NoError(t, err)
	assert.JSONEq(t, `{"description":"Civil time in simplified ISO format, like 13:45:30", "pattern":"^([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$", "title":"Time", "type":"string"}`, string(out))
}
