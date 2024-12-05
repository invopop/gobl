package in_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/in"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestOrgItemValidation(t *testing.T) {
	tr := tax.RegimeDefFor("IN")
	t.Run("valid", func(t *testing.T) {
		i := &org.Item{
			Identities: []*org.Identity{
				{
					Type: in.IdentityTypeHSN,
					Code: "1234",
				},
			},
		}
		err := tr.ValidateObject(i)
		assert.NoError(t, err)
	})

	t.Run("invalid", func(t *testing.T) {
		i := &org.Item{
			Identities: []*org.Identity{
				{
					Type: in.IdentityTypePAN,
					Code: "1234",
				},
			},
		}
		err := tr.ValidateObject(i)
		assert.ErrorContains(t, err, "identities: missing type HSN.")
	})
}
