package mt_test

import (
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/mt"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVATRateValues(t *testing.T) {
	cat := mt.New().CategoryDef(tax.CategoryVAT)
	require.NotNil(t, cat)

	tests := []struct {
		name string
		rate cbc.Key
		date cal.Date
		want string // percent, or "" when no value applies at the date
	}{
		{"standard 18% (current)", tax.RateGeneral, cal.MakeDate(2026, 1, 1), "18%"},
		{"standard 15% (pre-2004)", tax.RateGeneral, cal.MakeDate(2000, 1, 1), "15%"},
		{"reduced 12% (since 2024)", tax.RateReduced, cal.MakeDate(2026, 1, 1), "12%"},
		{"reduced 12% not yet applicable in 2023", tax.RateReduced, cal.MakeDate(2023, 1, 1), ""},
		{"super-reduced 7%", tax.RateSuperReduced, cal.MakeDate(2026, 1, 1), "7%"},
		{"special 5%", tax.RateSpecial, cal.MakeDate(2026, 1, 1), "5%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rd := cat.RateDef(tax.KeyStandard, tt.rate)
			require.NotNil(t, rd, "rate %s must be defined", tt.rate)
			v := rd.Value(tt.date, tax.Extensions{})
			if tt.want == "" {
				assert.Nil(t, v, "no rate value should apply before its Since date")
				return
			}
			require.NotNil(t, v, "a rate value should apply at %s", tt.date)
			assert.Equal(t, tt.want, v.Percent.String())
		})
	}
}
