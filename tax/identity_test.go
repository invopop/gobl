package tax_test

import (
	"testing"

	_ "github.com/invopop/gobl" // load all mods
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentity(t *testing.T) {
	tID := &tax.Identity{
		Country: l10n.ES,
		Code:    "X3157928M",
	}
	err := tID.Validate()
	assert.NoError(t, err)
	assert.Equal(t, tID.String(), "ESX3157928M")

	// Invalid tax id that should be validated against regional
	// checks.
	tID = &tax.Identity{
		Country: l10n.ES,
		Code:    "X3157928MMM",
	}
	err = tID.Validate()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "code: unknown type")
	}

	tID = &tax.Identity{
		Country: l10n.ES,
		Code:    "  x315-7928 m",
	}
	err = tID.Calculate()
	assert.NoError(t, err)
	assert.Equal(t, tID.Code.String(), "X3157928M")
}

func TestValidationRules(t *testing.T) {
	tID := &tax.Identity{
		Country: l10n.ES,
		Code:    "X3157928M",
	}
	err := validation.Validate(tID, tax.RequireIdentityType)
	assert.Error(t, err)

}
