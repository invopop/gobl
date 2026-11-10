package bs_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	_ "github.com/invopop/gobl/regimes/bs"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityRules(t *testing.T) {
	validate := func(code cbc.Code) error {
		return rules.Validate(&tax.Identity{Country: "BS", Code: code})
	}

	t.Run("valid 9-digit TIN", func(t *testing.T) {
		assert.NoError(t, validate("490918850"))
	})

	t.Run("empty code is allowed", func(t *testing.T) {
		assert.NoError(t, validate(""))
	})

	t.Run("too short", func(t *testing.T) {
		err := validate("4909185")
		assert.ErrorContains(t, err, "BS-TAX-IDENTITY-01")
	})

	t.Run("too long", func(t *testing.T) {
		err := validate("136410121320")
		assert.ErrorContains(t, err, "BS-TAX-IDENTITY-01")
	})

	t.Run("contains letters", func(t *testing.T) {
		err := validate("49091185A")
		assert.ErrorContains(t, err, "BS-TAX-IDENTITY-01")
	})

	t.Run("special characters", func(t *testing.T) {
		err := validate("491-091-850")
		assert.ErrorContains(t, err, "BS-TAX-IDENTITY-01")
	})
}

func TestTaxIdentityNormalization(t *testing.T) {
	t.Run("nil identity is safe", func(t *testing.T) {
		var tID *tax.Identity
		assert.NotPanics(t, func() { norm.Normalize(tID) })
	})

	t.Run("normalizes identity", func(t *testing.T) {
		tID := &tax.Identity{Country: "BS", Code: "490-918-501"}
		norm.Normalize(tID)
		assert.Equal(t, cbc.Code("490918501"), tID.Code)
	})

	t.Run("empty code is left untouched", func(t *testing.T) {
		tID := &tax.Identity{Country: "BS", Code: ""}
		norm.Normalize(tID)
		assert.Equal(t, cbc.Code(""), tID.Code)
	})

}
