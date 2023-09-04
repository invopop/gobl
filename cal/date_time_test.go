package cal_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
)

func TestDateTimeJSONParsing(t *testing.T) {
	// Handle a zero date
	t.Run("zero datetime", func(t *testing.T) {
		var d cal.DateTime
		data, err := json.Marshal(d)
		assert.NoError(t, err)
		assert.EqualValues(t, string(data), `"0000-00-00T00:00:00"`)

		err = json.Unmarshal([]byte(`"0000-00-00T00:00:00"`), &d)
		assert.NoError(t, err)
	})

	t.Run("valid datetime", func(t *testing.T) {
		d := cal.MakeDateTime(2021, time.May, 26, 10, 30, 20, 0)
		data, err := json.Marshal(d)
		assert.NoError(t, err)
		assert.EqualValues(t, string(data), `"2021-05-26T10:30:20"`)

		err = json.Unmarshal([]byte(`"2021-05-26T10:30:20"`), &d)
		assert.NoError(t, err)
		assert.Equal(t, d.Date.Year, 2021)
		assert.Equal(t, d.Date.Month, time.May)
		assert.Equal(t, d.Date.Day, 26)
		assert.Equal(t, d.Time.Hour, 10)
		assert.Equal(t, d.Time.Minute, 30)
		assert.Equal(t, d.Time.Second, 20)
	})
}

func TestDateTimeValidation(t *testing.T) {
	t.Run("basics", func(t *testing.T) {
		d := cal.MakeDateTime(2021, time.May, 26, 10, 20, 30, 1000)
		err := validation.Validate(d)
		assert.NoError(t, err)

		d = cal.MakeDateTime(2021, time.May, 1, 24, 0, 0, 0)
		err = d.Validate()
		assert.Error(t, err)
		err = validation.Validate(d)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid date")

		// Pointer
		dp := cal.NewDateTime(2021, time.May, 1, 24, 0, 0, 0)
		assert.Error(t, dp.Validate())
		assert.Error(t, validation.Validate(dp))

		dp = nil
		assert.NoError(t, validation.Validate(dp))
	})
}
