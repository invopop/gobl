package tax_test

import (
	"testing"

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
