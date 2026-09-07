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
			Start: cal.MakeDate(2022, 1, 25),
			End:   cal.MakeDate(2022, 2, 28),
		}
		assert.NoError(t, rules.Validate(p))
		assert.NoError(t, rules.Validate(&p))
	})

	t.Run("same day", func(t *testing.T) {
		p := cal.Period{
			Start: cal.MakeDate(2022, 1, 25),
			End:   cal.MakeDate(2022, 1, 25),
		}
		assert.NoError(t, rules.Validate(p))
		assert.NoError(t, rules.Validate(&p))
	})

	// EN 16931 BR-CO-19 / BR-CO-20: a period only needs one of its bounds.
	t.Run("end only", func(t *testing.T) {
		p := cal.Period{
			End: cal.MakeDate(2022, 2, 28),
		}
		assert.NoError(t, rules.Validate(p))
		assert.NoError(t, rules.Validate(&p))
	})

	t.Run("start only", func(t *testing.T) {
		p := cal.Period{
			Start: cal.MakeDate(2022, 1, 25),
		}
		assert.NoError(t, rules.Validate(p))
		assert.NoError(t, rules.Validate(&p))
	})

	t.Run("end before start", func(t *testing.T) {
		p := cal.Period{
			Start: cal.MakeDate(2022, 1, 25),
			End:   cal.MakeDate(2022, 1, 20),
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
		assert.True(t, faults.HasCode("GOBL-CAL-PERIOD-03"))
		assert.True(t, faults.HasPath("$"))
		assert.Equal(t, "either a start or end date is required", faults.First().Message())
		// A single fault: the retired field-level codes must not reappear.
		assert.False(t, faults.HasCode("GOBL-CAL-PERIOD-01"))
		assert.False(t, faults.HasCode("GOBL-CAL-PERIOD-02"))
		assert.False(t, faults.HasCode("GOBL-CAL-PERIOD-10"))

		faults = rules.Validate(&p)
		require.NotNil(t, faults)
		assert.True(t, faults.HasCode("GOBL-CAL-PERIOD-03"))
	})

	t.Run("label only", func(t *testing.T) {
		p := cal.Period{Label: "Q3"}
		faults := rules.Validate(p)
		require.NotNil(t, faults)
		assert.True(t, faults.HasCode("GOBL-CAL-PERIOD-03"))
	})
}

func TestPeriodJSON(t *testing.T) {
	t.Run("both bounds", func(t *testing.T) {
		p := cal.Period{
			Start: cal.MakeDate(2022, 1, 25),
			End:   cal.MakeDate(2022, 2, 28),
		}
		data, err := json.Marshal(p)
		require.NoError(t, err)
		assert.JSONEq(t, `{"start":"2022-01-25","end":"2022-02-28"}`, string(data))
	})

	t.Run("end only omits start", func(t *testing.T) {
		p := cal.Period{
			End: cal.MakeDate(2022, 2, 28),
		}
		data, err := json.Marshal(p)
		require.NoError(t, err)
		assert.JSONEq(t, `{"end":"2022-02-28"}`, string(data))
		assert.NotContains(t, string(data), "0000-00-00")

		var out cal.Period
		require.NoError(t, json.Unmarshal(data, &out))
		assert.Equal(t, p, out)
		assert.True(t, out.Start.IsZero())
	})

	t.Run("start only omits end", func(t *testing.T) {
		p := cal.Period{
			Start: cal.MakeDate(2022, 1, 25),
		}
		data, err := json.Marshal(p)
		require.NoError(t, err)
		assert.JSONEq(t, `{"start":"2022-01-25"}`, string(data))
		assert.NotContains(t, string(data), "0000-00-00")

		var out cal.Period
		require.NoError(t, json.Unmarshal(data, &out))
		assert.Equal(t, p, out)
		assert.True(t, out.End.IsZero())
	})

	t.Run("legacy zero date input", func(t *testing.T) {
		// Older documents may carry "0000-00-00" for the missing bound; it
		// must still parse and then be dropped on re-serialization.
		var p cal.Period
		require.NoError(t, json.Unmarshal([]byte(`{"start":"0000-00-00","end":"2022-02-28"}`), &p))
		assert.True(t, p.Start.IsZero())
		assert.NoError(t, rules.Validate(p))
		data, err := json.Marshal(p)
		require.NoError(t, err)
		assert.JSONEq(t, `{"end":"2022-02-28"}`, string(data))
	})
}
