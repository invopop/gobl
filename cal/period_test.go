package cal_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPeriodValidation(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		p := cal.Period{
			Start: cal.NewDate(2022, 1, 25),
			End:   cal.NewDate(2022, 2, 28),
		}
		assert.NoError(t, rules.Validate(p))
		assert.NoError(t, rules.Validate(&p))
	})

	t.Run("same day", func(t *testing.T) {
		p := cal.Period{
			Start: cal.NewDate(2022, 1, 25),
			End:   cal.NewDate(2022, 1, 25),
		}
		assert.NoError(t, rules.Validate(p))
		assert.NoError(t, rules.Validate(&p))
	})

	t.Run("end only", func(t *testing.T) {
		p := cal.Period{
			End: cal.NewDate(2022, 2, 28),
		}
		assert.NoError(t, rules.Validate(p))
		assert.NoError(t, rules.Validate(&p))
	})

	t.Run("start only", func(t *testing.T) {
		p := cal.Period{
			Start: cal.NewDate(2022, 1, 25),
		}
		assert.NoError(t, rules.Validate(p))
		assert.NoError(t, rules.Validate(&p))
	})

	t.Run("zero start", func(t *testing.T) {
		p := cal.Period{
			Start: new(cal.Date),
			End:   cal.NewDate(2022, 2, 28),
		}
		faults := rules.Validate(p)
		require.NotNil(t, faults)
		assert.True(t, faults.HasCode("GOBL-CAL-PERIOD-01"))
		assert.Equal(t, "start date cannot be zero", faults.First().Message())
	})

	t.Run("zero end", func(t *testing.T) {
		p := cal.Period{
			Start: cal.NewDate(2022, 1, 25),
			End:   new(cal.Date),
		}
		faults := rules.Validate(p)
		require.NotNil(t, faults)
		assert.True(t, faults.HasCode("GOBL-CAL-PERIOD-02"))
		assert.Equal(t, "end date cannot be zero", faults.First().Message())
	})

	t.Run("end before start", func(t *testing.T) {
		p := cal.Period{
			Start: cal.NewDate(2022, 1, 25),
			End:   cal.NewDate(2022, 1, 20),
		}
		faults := rules.Validate(p)
		require.NotNil(t, faults)
		assert.True(t, faults.HasCode("GOBL-CAL-PERIOD-10"))
		assert.Equal(t, "end date must be on or after start date", faults.First().Message())

		faults = rules.Validate(&p)
		require.NotNil(t, faults)
		assert.True(t, faults.HasCode("GOBL-CAL-PERIOD-10"))
	})

	t.Run("empty", func(t *testing.T) {
		p := cal.Period{}
		faults := rules.Validate(p)
		require.NotNil(t, faults)
		assert.True(t, faults.HasCode("GOBL-CAL-PERIOD-11"))
		assert.True(t, faults.HasPath("$"))
		assert.Equal(t, "either a start or end date is required", faults.First().Message())
		assert.False(t, faults.HasCode("GOBL-CAL-PERIOD-01"))
		assert.False(t, faults.HasCode("GOBL-CAL-PERIOD-02"))
		assert.False(t, faults.HasCode("GOBL-CAL-PERIOD-10"))

		faults = rules.Validate(&p)
		require.NotNil(t, faults)
		assert.True(t, faults.HasCode("GOBL-CAL-PERIOD-11"))
	})

	t.Run("label only", func(t *testing.T) {
		p := cal.Period{Label: "Q3"}
		faults := rules.Validate(p)
		require.NotNil(t, faults)
		assert.True(t, faults.HasCode("GOBL-CAL-PERIOD-11"))
	})
}

func TestPeriodJSON(t *testing.T) {
	t.Run("both bounds", func(t *testing.T) {
		p := cal.Period{
			Start: cal.NewDate(2022, 1, 25),
			End:   cal.NewDate(2022, 2, 28),
		}
		data, err := json.Marshal(p)
		require.NoError(t, err)
		assert.JSONEq(t, `{"start":"2022-01-25","end":"2022-02-28"}`, string(data))
	})

	t.Run("end only omits start", func(t *testing.T) {
		p := cal.Period{
			End: cal.NewDate(2022, 2, 28),
		}
		data, err := json.Marshal(p)
		require.NoError(t, err)
		assert.JSONEq(t, `{"end":"2022-02-28"}`, string(data))

		var out cal.Period
		require.NoError(t, json.Unmarshal(data, &out))
		assert.Equal(t, p, out)
		assert.Nil(t, out.Start)
	})

	t.Run("start only omits end", func(t *testing.T) {
		p := cal.Period{
			Start: cal.NewDate(2022, 1, 25),
		}
		data, err := json.Marshal(p)
		require.NoError(t, err)
		assert.JSONEq(t, `{"start":"2022-01-25"}`, string(data))

		var out cal.Period
		require.NoError(t, json.Unmarshal(data, &out))
		assert.Equal(t, p, out)
		assert.Nil(t, out.End)
	})

	t.Run("zero date input", func(t *testing.T) {
		var p cal.Period
		require.NoError(t, json.Unmarshal([]byte(`{"start":"0000-00-00","end":"2022-02-28"}`), &p))
		require.NotNil(t, p.Start)
		assert.True(t, p.Start.IsZero())

		faults := rules.Validate(p)
		require.NotNil(t, faults)
		assert.True(t, faults.HasCode("GOBL-CAL-PERIOD-01"))
	})
}
