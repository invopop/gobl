package tax_test

import (
	"testing"

	_ "github.com/invopop/gobl" // load all mods
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentity(t *testing.T) {
	tID := &tax.Identity{
		Country: "ES",
		Code:    "X3157928M",
	}
	err := rules.Validate(tID)
	assert.NoError(t, err)
	assert.Equal(t, tID.String(), "ESX3157928M")

	// Invalid tax id that should be validated against regional
	// checks.
	tID = &tax.Identity{
		Country: "ES",
		Code:    "X3157928MMM",
	}
	err = rules.Validate(tID)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "[GOBL-ES-TAX-IDENTITY-01]")
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

	t.Run("mexican case with strange characters", func(t *testing.T) {
		tID := &tax.Identity{
			Country: "MX",
			Code:    "K&Ñ010301I16",
		}
		assert.NoError(t, rules.Validate(tID))
	})

	t.Run("with scheme", func(t *testing.T) {
		tID := &tax.Identity{
			Country: "ES",
			Code:    "X3157928M",
			Scheme:  tax.CategoryVAT,
		}
		assert.NoError(t, rules.Validate(tID))
	})
	t.Run("with invalid scheme", func(t *testing.T) {
		tID := &tax.Identity{
			Country: "ES",
			Code:    "X3157928M",
			Scheme:  "Fo--o",
		}
		assert.ErrorContains(t, rules.Validate(tID), "GOBL-CBC-CODE-02")
	})

	t.Run("in EU", func(t *testing.T) {
		tID := &tax.Identity{
			Country: "ES",
			Code:    "X3157928M",
		}
		assert.True(t, tID.InEU(cal.MakeDate(2025, 7, 15)))
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
	assert.ErrorContains(t, err, "GOBL-ES-TAX-IDENTITY-01")

	_, err = tax.ParseIdentity("E")
	assert.ErrorContains(t, err, "invalid tax identity code")
}

func TestIdentityGetScheme(t *testing.T) {
	t.Run("use override", func(t *testing.T) {
		tID := &tax.Identity{
			Country: "ES",
			Code:    "X3157928M",
			Scheme:  "IPSI",
		}
		assert.Equal(t, cbc.Code("IPSI"), tID.GetScheme())
	})
	t.Run("use regime default", func(t *testing.T) {
		tID := &tax.Identity{
			Country: "ES",
			Code:    "X3157928M",
		}
		assert.Equal(t, tax.CategoryVAT, tID.GetScheme())
	})
	t.Run("use empty for regime without default", func(t *testing.T) {
		tID := &tax.Identity{
			Country: "US",
			Code:    "1234567",
		}
		assert.Equal(t, cbc.CodeEmpty, tID.GetScheme())
	})
	t.Run("use empty for no regime", func(t *testing.T) {
		tID := &tax.Identity{
			Country: "ZW", // Will need fixing when ZW supported :-)
			Code:    "1234567",
		}
		assert.Equal(t, cbc.CodeEmpty, tID.GetScheme())
	})

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
