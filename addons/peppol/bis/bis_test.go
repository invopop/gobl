package bis_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/addons/peppol/bis"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
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
	var found bool
	for _, id := range def.Identities {
		if id != nil && id.Key == bis.IdentityKeyGreekMARK {
			found = true
			break
		}
	}
	assert.True(t, found, "Greek MARK identity definition should be registered")
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
