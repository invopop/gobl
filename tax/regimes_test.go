package tax_test

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestAllRegimes(t *testing.T) {
	for _, r := range tax.AllRegimeDefs() {
		t.Run(r.Name.String(), func(t *testing.T) {
			assert.NoError(t, r.Validate())
		})
	}
}

func TestRegimesAltCountryCodes(t *testing.T) {
	r := tax.RegimeDefFor("GR")
	assert.Equal(t, "EL", r.Country.String())
}

func TestSetRegime(t *testing.T) {
	tests := []struct {
		reg string
		exp string
	}{
		{"ES", "ES"},
		{"GR", "EL"},
		{"YY", ""},
	}

	for _, tt := range tests {
		t.Run(tt.reg, func(t *testing.T) {
			r := new(tax.Regime)
			r.SetRegime(l10n.TaxCountryCode(tt.reg))
			assert.Equal(t, tt.exp, r.GetRegime().String())
		})
	}
}
