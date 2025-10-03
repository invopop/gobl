package verifactu_test

import (
	"testing"

	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatePartyThroughAddon(t *testing.T) {
	// Get the addon
	addon := tax.AddonForKey(verifactu.V1)
	require.NotNil(t, addon, "Verifactu addon should be registered")

	t.Run("valid party name", func(t *testing.T) {
		party := &org.Party{
			Name: "Valid Company Name S.L.",
		}

		err := addon.Validator(party)
		assert.NoError(t, err)
	})

	t.Run("empty party name", func(t *testing.T) {
		party := &org.Party{
			Name: "",
		}

		err := addon.Validator(party)
		assert.NoError(t, err) // Empty name should be valid (validation.Skip is used)
	})

	t.Run("party name with forbidden characters", func(t *testing.T) {
		testCases := []struct {
			name         string
			partyName    string
			expectedChar rune
		}{
			{"less than", "Company < Name", '<'},
			{"greater than", "Company > Name", '>'},
			{"double quote", "Company \" Name", '"'},
			{"single quote", "Company ' Name", '\''},
			{"equals", "Company = Name", '='},
			{"multiple forbidden chars", "Company <> \"Name\"", '<'}, // Should catch the first one
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				party := &org.Party{
					Name: tc.partyName,
				}

				err := addon.Validator(party)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "contains forbidden character")
				assert.Contains(t, err.Error(), string(tc.expectedChar))
			})
		}
	})
}
