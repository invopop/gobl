package in_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/in"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestOrgItemValidation(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		i := &org.Item{
			Name: "Test Item",
			Identities: []*org.Identity{
				{
					Type: in.IdentityTypeHSN,
					Code: "1234",
				},
			},
		}
		err := rules.Validate(i, tax.RegimeContext(in.CountryCode))
		assert.NoError(t, err)
	})

	t.Run("invalid HSN", func(t *testing.T) {
		i := &org.Item{
			Identities: []*org.Identity{
				{
					Type: in.IdentityTypeHSN,
					Code: "X", // This will be detected by org_identites check
				},
			},
		}
		err := rules.Validate(i, tax.RegimeContext(in.CountryCode))
		assert.ErrorContains(t, err, "[GOBL-IN-ORG-IDENTITY-02] ($.identities[0].code) identity code must be a valid HSN format")
	})

	t.Run("invalid", func(t *testing.T) {
		i := &org.Item{
			Identities: []*org.Identity{
				{
					Type: in.IdentityTypePAN,
					Code: "CTUGE1616Y",
				},
			},
		}
		err := rules.Validate(i, tax.RegimeContext(in.CountryCode))
		assert.ErrorContains(t, err, "[GOBL-IN-ORG-ITEM-01] ($.identities) all items must have an HSN identity code")
	})
}
