package tax_test

import (
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRateDefValue(t *testing.T) {
	rate := &tax.RateDef{
		Values: []*tax.RateValueDef{
			{
				Since:   cal.NewDate(2021, 1, 1),
				Percent: num.MakePercentage(19, 2),
			},
			{
				Since:   cal.NewDate(2020, 7, 1),
				Percent: num.MakePercentage(16, 2),
			},
			{
				Since:   cal.NewDate(2007, 1, 1),
				Percent: num.MakePercentage(19, 2),
			},
			{
				Since:   cal.NewDate(1993, 1, 1),
				Percent: num.MakePercentage(15, 2),
			},
		},
	}

	tests := []struct {
		name    string
		date    cal.Date
		percent string
	}{
		{
			name:    "on Since date",
			date:    cal.MakeDate(2021, 1, 1),
			percent: "19%",
		},
		{
			name:    "day before Since",
			date:    cal.MakeDate(2020, 12, 31),
			percent: "16%",
		},
		{
			name:    "day after Since",
			date:    cal.MakeDate(2021, 1, 2),
			percent: "19%",
		},
		{
			name:    "mid-period",
			date:    cal.MakeDate(2020, 10, 15),
			percent: "16%",
		},
		{
			name:    "on earliest Since date",
			date:    cal.MakeDate(1993, 1, 1),
			percent: "15%",
		},
		{
			name:    "far future",
			date:    cal.MakeDate(2099, 1, 1),
			percent: "19%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rv := rate.Value(tt.date, nil)
			require.NotNil(t, rv)
			assert.Equal(t, tt.percent, rv.Percent.String())
		})
	}

	t.Run("before any rate returns nil", func(t *testing.T) {
		rv := rate.Value(cal.MakeDate(1992, 12, 31), nil)
		assert.Nil(t, rv)
	})

	t.Run("nil Since always matches", func(t *testing.T) {
		r := &tax.RateDef{
			Values: []*tax.RateValueDef{
				{Percent: num.MakePercentage(10, 2)},
			},
		}
		rv := r.Value(cal.MakeDate(1900, 1, 1), nil)
		require.NotNil(t, rv)
		assert.Equal(t, "10%", rv.Percent.String())
	})
}
