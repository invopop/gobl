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

func TestDateJSONParsing(t *testing.T) {
	// Handle a zero date
	t.Run("zero date", func(t *testing.T) {
		var d cal.Date
		data, err := json.Marshal(d)
		assert.NoError(t, err)
		assert.EqualValues(t, string(data), `"0000-00-00"`)

		err = json.Unmarshal([]byte(`"0000-00-00"`), &d)
		assert.NoError(t, err)
	})

	t.Run("valid date", func(t *testing.T) {
		d := cal.MakeDate(2021, time.May, 26)
		data, err := json.Marshal(d)
		assert.NoError(t, err)
		assert.EqualValues(t, string(data), `"2021-05-26"`)

		err = json.Unmarshal([]byte(`"2021-05-26"`), &d)
		assert.NoError(t, err)
		assert.Equal(t, d.Year, 2021)
		assert.Equal(t, d.Month, time.May)
		assert.Equal(t, d.Day, 26)
	})
}

func TestDateValidation(t *testing.T) {
	t.Run("basics", func(t *testing.T) {
		d := cal.MakeDate(2021, time.May, 26)
		err := validation.Validate(d)
		assert.NoError(t, err)

		d = cal.MakeDate(2021, 0, 1)
		err = d.Validate()
		assert.Error(t, err)
		err = validation.Validate(d)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid date")

		d = cal.MakeDate(2021, 1, 0)
		err = validation.Validate(d)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid date")

		// Pointer
		dp := cal.NewDate(2021, 1, 0)
		assert.Error(t, dp.Validate())
		assert.Error(t, validation.Validate(dp))

		dp = nil
		assert.NoError(t, validation.Validate(dp))
	})

	t.Run("date not zero", func(t *testing.T) {
		d := cal.MakeDate(2021, time.May, 26)
		err := validation.Validate(d, cal.DateNotZero())
		assert.NoError(t, err)

		d = cal.Date{}
		err = validation.Validate(d, cal.DateNotZero())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "required")

		dp := new(cal.Date)
		err = validation.Validate(dp, cal.DateNotZero())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "required")

		dp = nil
		err = validation.Validate(dp, cal.DateNotZero())
		assert.NoError(t, err)
	})

	t.Run("date after", func(t *testing.T) {
		d := cal.MakeDate(2023, time.March, 25)
		err := validation.Validate(d, cal.DateAfter(cal.MakeDate(2023, time.March, 24)))
		assert.NoError(t, err)

		err = validation.Validate(d, cal.DateAfter(cal.MakeDate(2023, time.March, 25)))
		assert.NoError(t, err)

		err = validation.Validate(d, cal.DateAfter(cal.MakeDate(2023, time.March, 26)))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "too early")
	})

	t.Run("date before", func(t *testing.T) {
		d := cal.MakeDate(2023, time.March, 25)
		err := validation.Validate(d, cal.DateBefore(cal.MakeDate(2023, time.March, 26)))
		assert.NoError(t, err)

		err = validation.Validate(d, cal.DateBefore(cal.MakeDate(2023, time.March, 25)))
		assert.NoError(t, err)

		err = validation.Validate(d, cal.DateBefore(cal.MakeDate(2023, time.March, 24)))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "too late")
	})
}

func TestDateToday(t *testing.T) {
	d := cal.Today()
	assert.Equal(t, d.Year, time.Now().Year())
	assert.Equal(t, d.Month, time.Now().Month())
	assert.Equal(t, d.Day, time.Now().Day())

	l, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)
	tn := time.Now().In(l)
	d = cal.TodayIn(l)
	assert.Equal(t, d.Year, tn.Year())
	assert.Equal(t, d.Month, tn.Month())
	assert.Equal(t, d.Day, tn.Day())
}
