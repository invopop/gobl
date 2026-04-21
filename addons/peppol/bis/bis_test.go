package bis_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/addons/peppol/bis"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAddonRegistered verifies the Peppol addon is reachable through the
// global tax registry and declares its dependency on EN 16931.
func TestAddonRegistered(t *testing.T) {
	def := tax.AddonForKey(bis.V3)
	require.NotNil(t, def, "peppol addon must be registered")
	assert.Equal(t, bis.V3, def.Key)
	assert.Contains(t, def.Requires, en16931.V2017, "peppol must require en16931")
}

// TestAddonKeyConstant pins the public key so external references (e.g. the
// gobl.ubl Peppol contexts) remain stable.
func TestAddonKeyConstant(t *testing.T) {
	assert.Equal(t, cbc.Key("peppol-bis-v3"), bis.V3)
}

// TestIdentities checks that the custom identity keys are registered.
func TestIdentities(t *testing.T) {
	def := tax.AddonForKey(bis.V3)
	require.NotNil(t, def)
	keys := map[cbc.Key]bool{}
	for _, id := range def.Identities {
		if id != nil {
			keys[id.Key] = true
		}
	}
	assert.True(t, keys[bis.IdentityKeyGreekMARK], "Greek MARK identity definition should be registered")
	assert.True(t, keys[bis.IdentityKeyFSkatt], "F-skatt identity definition should be registered")
}

// TestFSkattNormalizer verifies that attaching the F-skatt key to a supplier
// auto-fills the boilerplate code via the addon's normalizer.
func TestFSkattNormalizer(t *testing.T) {
	def := tax.AddonForKey(bis.V3)
	require.NotNil(t, def)
	require.NotNil(t, def.Normalizer)

	t.Run("fills code when key set without one", func(t *testing.T) {
		id := &org.Identity{Key: bis.IdentityKeyFSkatt}
		def.Normalizer(id)
		assert.Equal(t, cbc.Code(bis.FSkattText), id.Code)
	})

	t.Run("leaves existing code untouched", func(t *testing.T) {
		id := &org.Identity{Key: bis.IdentityKeyFSkatt, Code: "custom"}
		def.Normalizer(id)
		assert.Equal(t, cbc.Code("custom"), id.Code)
	})

	t.Run("ignores other identities", func(t *testing.T) {
		id := &org.Identity{Key: bis.IdentityKeyGreekMARK, Code: "12345"}
		def.Normalizer(id)
		assert.Equal(t, cbc.Code("12345"), id.Code)
	})
}

// TestEmptyInvoiceNoPanic verifies that running Peppol validation against an
// empty invoice doesn't panic — a baseline sanity check.
func TestEmptyInvoiceNoPanic(t *testing.T) {
	inv := &bill.Invoice{}
	inv.SetAddons(en16931.V2017, bis.V3)
	assert.NotPanics(t, func() {
		_ = rules.Validate(inv)
	})
}
