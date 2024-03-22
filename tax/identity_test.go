package tax_test

import (
	"encoding/json"
	"testing"

	_ "github.com/invopop/gobl" // load all mods
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		Code:    "X3157928M",
		Zone:    "XX",
	}
	err = tID.Validate()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "zone: must be blank.")
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

func TestIdentityUnmarshalJSON(t *testing.T) {
	// Unmarshal should store the zone, but not return it
	data := []byte(`{"country":"CO","code":"9014514805","zone":"11001"}`)
	tID := &tax.Identity{}
	err := tID.UnmarshalJSON(data)
	require.NoError(t, err)
	assert.Equal(t, tID.Country, l10n.CO)
	assert.Equal(t, tID.Code.String(), "9014514805")
	assert.Equal(t, tID.Zone.String(), "11001")

	out, err := json.Marshal(tID)
	require.NoError(t, err)
	assert.Equal(t, string(out), `{"country":"CO","code":"9014514805"}`)
}
