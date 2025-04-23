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
	t.Run("valid time", func(t *testing.T) {
		tm := cal.TimeNow()
		assert.NotZero(t, tm.Hour)
		assert.NotZero(t, tm.Minute)
		assert.NotZero(t, tm.Second)
	})
}

func TestTimeNowIn(t *testing.T) {
	t.Run("valid time", func(t *testing.T) {
		loc, err := time.LoadLocation("America/New_York")
		require.NoError(t, err)
		tm := cal.TimeNowIn(loc)
		assert.NotZero(t, tm.Hour)
		assert.NotZero(t, tm.Minute)
		assert.NotZero(t, tm.Second)
	})
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

func TestTimeJSONSChema(t *testing.T) {
	// Check the schema for the time type.
	schema := cal.Time{}.JSONSchema()
	out, err := json.Marshal(schema)
	require.NoError(t, err)
	assert.JSONEq(t, `{"description":"Civil time in simplified ISO format, like 13:45:30", "format":"time", "title":"Time", "type":"string"}`, string(out))
}
