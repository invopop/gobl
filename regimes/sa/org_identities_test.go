package sa_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	_ "github.com/invopop/gobl/regimes/sa"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestOrgIdentityCodeFormat(t *testing.T) {
	validate := func(code cbc.Code) error {
		id := &org.Identity{Type: "CRN", Code: code}
		return rules.Validate(id, tax.RegimeContext("SA"))
	}

	t.Run("alphanumeric code is valid", func(t *testing.T) {
		assert.NoError(t, validate("ABC123"))
	})

	t.Run("digits only is valid", func(t *testing.T) {
		assert.NoError(t, validate("1234567890"))
	})

	t.Run("letters only is valid", func(t *testing.T) {
		assert.NoError(t, validate("ABCDEF"))
	})

	t.Run("code with hyphen is invalid", func(t *testing.T) {
		err := validate("ABC-123")
		assert.ErrorContains(t, err, "invoice identity type must be alphanumerical")
	})

	t.Run("code with special characters is invalid", func(t *testing.T) {
		err := validate("ABC@123")
		assert.ErrorContains(t, err, "invoice identity type must be alphanumerical")
	})

	t.Run("code with spaces is invalid", func(t *testing.T) {
		err := validate("ABC 123")
		assert.ErrorContains(t, err, "invoice identity type must be alphanumerical")
	})
}
