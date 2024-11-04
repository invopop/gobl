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
	tID.Normalize()
	assert.Equal(t, tID.Code.String(), "X3157928M")

	tID = nil
	assert.NotPanics(t, func() {
		tID.Normalize()
	})

	t.Run("mexican case with custom validation", func(t *testing.T) {
		tID := &tax.Identity{
			Country: "MX",
			Code:    "K&A010301I16",
		}
		assert.NoError(t, tID.Validate())
	})

	t.Run("invalid non-exception case", func(t *testing.T) {
		tID := &tax.Identity{
			Country: "ZW", // update when ZW regime is added
			Code:    "AB&DE",
		}
		assert.ErrorContains(t, tID.Validate(), "code: must be in a valid format")
	})
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

func TestNormalizeIdentity(t *testing.T) {
	t.Run("regular case", func(t *testing.T) {
		tID := &tax.Identity{
			Country: "AU",
			Code:    "  x315-7928 m  ",
		}
		tax.NormalizeIdentity(tID)
		assert.Equal(t, tID.Code.String(), "X3157928M")
	})
	t.Run("with alt country", func(t *testing.T) {
		tID := &tax.Identity{
			Country: "EL",
			Code:    "GR925667500",
		}
		tax.NormalizeIdentity(tID, "GR")
		assert.Equal(t, tID.Code.String(), "925667500")
	})
}

func TestIdentityNormalize(t *testing.T) {
	t.Run("for unkown regime", func(t *testing.T) {
		tID := &tax.Identity{
			Country: "XX",
			Code:    "  x315-7928 m  ",
		}
		tID.Normalize()
		assert.Equal(t, tID.Code.String(), "X3157928M")
	})
	t.Run("for known regime", func(t *testing.T) {
		tID := &tax.Identity{
			Country: "FR",
			Code:    " 356000000 ",
		}
		tID.Normalize()
		assert.Equal(t, tID.Code.String(), "39356000000") // adds 2 0s on end
	})
	t.Run("with calculate method", func(t *testing.T) {
		tID := &tax.Identity{
			Country: "FR",
			Code:    " 356000000 ",
		}
		err := tID.Calculate()
		assert.NoError(t, err)
		assert.Equal(t, tID.Code.String(), "39356000000")
	})
}
