package is_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	_ "github.com/invopop/gobl/regimes/is"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityRules(t *testing.T) {
	validate := func(code cbc.Code) error {
		return rules.Validate(&tax.Identity{Country: "IS", Code: code})
	}

	// Sample valid Kennitalas constructed against the documented algorithm:
	//   1201743399 → sum=79, rem=2, check=9, century=9 (born 12 Jan 1974)
	//   0301010300 → sum=22, rem=0, check=0, century=0 (born 3 Jan 2001)

	t.Run("valid Kennitala", func(t *testing.T) {
		assert.NoError(t, validate("1201743399"))
	})

	t.Run("valid Kennitala with remainder 0 → check digit 0", func(t *testing.T) {
		assert.NoError(t, validate("0301010300"))
	})

	t.Run("empty code is allowed", func(t *testing.T) {
		assert.NoError(t, validate(""))
	})

	t.Run("wrong check digit", func(t *testing.T) {
		err := validate("1201743299")
		assert.ErrorContains(t, err, "IDENTITY-02")
	})

	// 0101010050: sum = 0*3 + 1*2 + 0*7 + 1*6 + 0*5 + 1*4 + 0*3 + 0*2 = 12,
	// 12 mod 11 = 1, which corresponds to the check-digit-would-be-10 case
	// that the national registry never issues.
	t.Run("remainder 1 → never issued", func(t *testing.T) {
		err := validate("0101010050")
		assert.ErrorContains(t, err, "IDENTITY-02")
	})

	t.Run("too short", func(t *testing.T) {
		err := validate("120174339")
		assert.ErrorContains(t, err, "IDENTITY-01")
	})

	t.Run("too long", func(t *testing.T) {
		err := validate("12017433999")
		assert.ErrorContains(t, err, "IDENTITY-01")
	})

	t.Run("contains letters", func(t *testing.T) {
		err := validate("120174339A")
		assert.ErrorContains(t, err, "IDENTITY-01")
	})

	t.Run("special characters", func(t *testing.T) {
		err := validate("120174-3399")
		assert.ErrorContains(t, err, "IDENTITY-01")
	})
}

func TestTaxIdentityNormalization(t *testing.T) {
	t.Run("nil identity is safe", func(t *testing.T) {
		var tID *tax.Identity
		assert.NotPanics(t, func() { norm.Normalize(tID) })
	})

	t.Run("normalizes identity", func(t *testing.T) {
		tID := &tax.Identity{Country: "IS", Code: "1201743399"}
		norm.Normalize(tID)
		assert.Equal(t, cbc.Code("1201743399"), tID.Code)
	})

	t.Run("empty code is left untouched", func(t *testing.T) {
		tID := &tax.Identity{Country: "IS", Code: ""}
		norm.Normalize(tID)
		assert.Equal(t, cbc.Code(""), tID.Code)
	})

	t.Run("strips whitespace and country prefix", func(t *testing.T) {
		tID := &tax.Identity{Country: "IS", Code: " IS 120174-3399 "}
		norm.Normalize(tID)
		assert.Equal(t, cbc.Code("1201743399"), tID.Code)
	})
}
