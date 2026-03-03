package nz_test

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/nz"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	tID := &tax.Identity{
		Country: "NZ",
		Code:    "12-345-678",
	}

	nz.New().Normalizer(tID)
	assert.Equal(t, "12345678", tID.Code.String())
}

func TestNewRegimeDef(t *testing.T) {
	reg := nz.New()
	assert.Equal(t, l10n.TaxCountryCode("NZ"), reg.Country)
	assert.NotNil(t, reg.Categories)
	assert.NotNil(t, reg.Validator)
	assert.NotNil(t, reg.Normalizer)
}

func TestValidateNonIdentity(t *testing.T) {
	err := nz.Validate("not an identity")
	assert.NoError(t, err)
}
