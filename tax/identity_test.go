package tax_test

import (
	"testing"

	_ "github.com/invopop/gobl" // load all mods
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentity(t *testing.T) {
	tID := &tax.Identity{
		Country: "ES",
		Code:    "X3157928M",
	}
	err := tID.Validate()
	assert.NoError(t, err)
	assert.Equal(t, tID.String(), "ESX3157928M")

	// Invalid tax id that should be validated against regional
	// checks.
	tID = &tax.Identity{
		Country: "ES",
		Code:    "X3157928MMM",
	}
	err = tID.Validate()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "code: unknown type")
	}

	tID = &tax.Identity{
		Country: "ES",
		Code:    "X3157928M",
		Zone:    "XX",
	}
	err = tID.Validate()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "zone: must be blank.")
	}

	tID = &tax.Identity{
		Country: "ES",
		Code:    "  x315-7928 m",
	}
	err = tID.Calculate()
	assert.NoError(t, err)
	assert.Equal(t, tID.Code.String(), "X3157928M")

	tID = nil
	assert.NoError(t, tID.Normalize(), "should handle nil identities")
}

func TestParseIdentity(t *testing.T) {
	tID, err := tax.ParseIdentity("ESX3157928M")
	assert.NoError(t, err)
	assert.Equal(t, tID.String(), "ESX3157928M")

	tID, err = tax.ParseIdentity("ES-X 315 79. 28M")
	assert.NoError(t, err)
	assert.Equal(t, tID.String(), "ESX3157928M")

	_, err = tax.ParseIdentity("ESX3157928MMM")
	assert.ErrorContains(t, err, "code: unknown type")

	_, err = tax.ParseIdentity("E")
	assert.ErrorContains(t, err, "invalid tax identity code")
}

func TestValidationRules(t *testing.T) {
	tID := &tax.Identity{
		Country: "ES",
	}
	err := validation.Validate(tID, tax.RequireIdentityCode)
	assert.ErrorContains(t, err, "code: cannot be blank")
}
