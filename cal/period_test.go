package cal_test

import (
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

	t.Run("missing start", func(t *testing.T) {
		p := cal.Period{
			End: cal.MakeDate(2022, 2, 28),
		}
		faults := rules.Validate(p)
		require.NotNil(t, faults)
		assert.True(t, faults.HasCode("GOBL-CAL-PERIOD-01"))
		assert.True(t, faults.HasPath("$.start"))
		assert.Equal(t, "start date cannot be zero", faults.First().Message())
	})

	t.Run("missing end", func(t *testing.T) {
		p := cal.Period{
			Start: cal.MakeDate(2022, 1, 25),
		}
		faults := rules.Validate(p)
		require.NotNil(t, faults)
		assert.True(t, faults.HasCode("GOBL-CAL-PERIOD-02"))
		assert.True(t, faults.HasPath("$.end"))
		assert.Equal(t, "end date cannot be zero", faults.First().Message())
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
		assert.True(t, faults.HasCode("GOBL-CAL-PERIOD-01"))
		assert.True(t, faults.HasCode("GOBL-CAL-PERIOD-02"))
	})
}
