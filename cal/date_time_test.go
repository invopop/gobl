package cal_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDateTimeJSONParsing(t *testing.T) {
	// Handle a zero date
	t.Run("zero date", func(t *testing.T) {
		var dt cal.DateTime
		data, err := json.Marshal(dt)
		assert.NoError(t, err)
		assert.EqualValues(t, string(data), `"0000-00-00T00:00:00"`)

		err = json.Unmarshal([]byte(`"0000-00-00T00:00:00"`), &dt)
		assert.NoError(t, err)
	})

	t.Run("valid date", func(t *testing.T) {
		dt := cal.MakeDateTime(2023, time.September, 4, 15, 59, 30)
		data, err := json.Marshal(dt)
		assert.NoError(t, err)
		assert.EqualValues(t, string(data), `"2023-09-04T15:59:30"`)

		err = json.Unmarshal([]byte(`"2023-09-04T15:59:30"`), &dt)
		assert.NoError(t, err)
		assert.Equal(t, dt.Date.Year, 2023)
		assert.Equal(t, dt.Date.Month, time.September)
		assert.Equal(t, dt.Date.Day, 4)
	})
}

func TestDateTimeValidation(t *testing.T) {
	t.Run("basics", func(t *testing.T) {
		d := cal.MakeDateTime(2021, time.May, 26, 15, 59, 30)
		err := validation.Validate(d)
		assert.NoError(t, err)

		d = cal.MakeDateTime(2021, 0, 1, 15, 59, 30)
		err = d.Validate()
		assert.Error(t, err)
		err = validation.Validate(d)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid date time")

		d = cal.MakeDateTime(2021, 1, 0, 15, 59, 30)
		err = validation.Validate(d)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid date time")

		// Pointer
		dp := cal.NewDateTime(2021, 1, 0, 15, 59, 30)
		assert.Error(t, dp.Validate())
		assert.Error(t, validation.Validate(dp))

		dp = nil
		assert.NoError(t, validation.Validate(dp))
	})

	t.Run("date time not zero", func(t *testing.T) {
		d := cal.MakeDateTime(2021, time.May, 26, 10, 10, 10)
		err := validation.Validate(d, cal.DateTimeNotZero())
		assert.NoError(t, err)

		d = cal.DateTime{}
		err = validation.Validate(d, cal.DateTimeNotZero())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "required")

		dp := new(cal.DateTime)
		err = validation.Validate(dp, cal.DateTimeNotZero())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "required")

		dp = nil
		err = validation.Validate(dp, cal.DateTimeNotZero())
		assert.NoError(t, err)
	})

	t.Run("date time after", func(t *testing.T) {
		d := cal.MakeDateTime(2023, time.March, 25, 10, 10, 10)
		d2 := cal.MakeDateTime(2023, time.March, 24, 10, 10, 9)
		err := validation.Validate(d, cal.DateTimeAfter(d2))
		assert.NoError(t, err)

		d2 = cal.MakeDateTime(2023, time.March, 25, 10, 10, 9)
		err = validation.Validate(d, cal.DateTimeAfter(d2))
		assert.NoError(t, err)

		d2 = cal.MakeDateTime(2023, time.March, 26, 10, 10, 10)
		err = validation.Validate(d, cal.DateTimeAfter(d2))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "too early")

		d2 = cal.MakeDateTime(2023, time.March, 25, 10, 10, 11)
		err = validation.Validate(d, cal.DateTimeAfter(d2))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "too early")
	})

	t.Run("date time before", func(t *testing.T) {
		d := cal.MakeDateTime(2023, time.March, 25, 10, 10, 10)

		d2 := cal.MakeDateTime(2023, time.March, 26, 10, 10, 10)
		err := validation.Validate(d, cal.DateTimeBefore(d2))
		assert.NoError(t, err)

		assert.NoError(t, err)
		d2 = cal.MakeDateTime(2023, time.March, 25, 10, 10, 11)
		err = validation.Validate(d, cal.DateTimeBefore(d2))
		assert.NoError(t, err)

		d2 = cal.MakeDateTime(2023, time.March, 25, 10, 10, 10)
		err = validation.Validate(d, cal.DateTimeBefore(d2))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "too late")

		d2 = cal.MakeDateTime(2023, time.March, 24, 10, 10, 10)
		err = validation.Validate(d, cal.DateTimeBefore(d2))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "too late")
	})
}

func TestDateTimeThisSecond(t *testing.T) {
	d := cal.ThisSecond()
	tn := time.Now()
	// note: this test may fail if minute changes between
	// the two "Now()" calls
	assert.Equal(t, d.Date.Year, tn.Year())
	assert.Equal(t, d.Date.Month, tn.Month())
	assert.Equal(t, d.Date.Day, tn.Day())
	assert.Equal(t, d.Time.Hour, tn.Hour())
	assert.Equal(t, d.Time.Minute, tn.Minute())

	l, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)
	tn = time.Now().In(l)
	d = cal.ThisSecondIn(l)
	assert.Equal(t, d.Date.Year, tn.Year())
	assert.Equal(t, d.Date.Month, tn.Month())
	assert.Equal(t, d.Date.Day, tn.Day())
	assert.Equal(t, d.Time.Hour, tn.Hour())
	assert.Equal(t, d.Time.Minute, tn.Minute())
}

func TestDateTimeClone(t *testing.T) {
	d := cal.MakeDateTime(2021, time.May, 26, 10, 10, 10)
	d2 := d.Clone()
	assert.Equal(t, d.String(), d2.String())
	d = cal.MakeDateTime(2021, time.May, 27, 10, 10, 10)
	assert.NotEqual(t, d.String(), d2.String())
}

func TestDateTimeWithTimeZ(t *testing.T) {
	d := cal.MakeDateTime(2023, time.July, 28, 10, 10, 5)
	dt := d.TimeZ()
	assert.Equal(t, "2023-07-28 10:10:05 +0000 UTC", dt.String())

	dp := cal.NewDateTime(2023, time.July, 28, 10, 10, 5)
	dt = dp.TimeZ()
	assert.Equal(t, "2023-07-28 10:10:05 +0000 UTC", dt.String())

	loc, err := time.LoadLocation("Europe/Madrid")
	require.NoError(t, err)
	dt = d.In(loc)
	assert.Equal(t, "2023-07-28 10:10:05 +0200 CEST", dt.String())
}

func TestDateTimeOf(t *testing.T) {
	x := time.Date(2023, time.July, 28, 12, 12, 1, 0, time.UTC)
	d := cal.DateTimeOf(x)
	assert.Equal(t, "2023-07-28T12:12:01", d.String())
}
