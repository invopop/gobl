package addons_test

import (
	"testing"

	// Importing the addons package (the package under test) runs its init,
	// registering the curated approved external-addon list.
	_ "github.com/invopop/gobl/addons"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApprovedFRCTCAddons(t *testing.T) {
	byKey := make(map[cbc.Key]*tax.ExternalAddon)
	for _, ea := range tax.ApprovedAddons() {
		byKey[ea.Key] = ea
	}

	for _, key := range []cbc.Key{"fr-ctc-v1", "fr-ctc-flow2-v1", "fr-ctc-flow6-v1", "fr-ctc-flow10-v1"} {
		ea, ok := byKey[key]
		require.Truef(t, ok, "expected %s on the approved list", key)
		assert.Equal(t, "github.com/invopop/gobl.fr.ctc", ea.Module, "%s module", key)
		assert.NotEmpty(t, ea.Name.String(), "%s should carry a name", key)

		// The implementation is external, so the key is not runtime-registered
		// in core — recognition only, not a runtime bypass.
		assert.Nil(t, tax.AddonForKey(key), "%s should not be registered in core", key)
	}
}

func TestApprovedBRAddons(t *testing.T) {
	byKey := make(map[cbc.Key]*tax.ExternalAddon)
	for _, ea := range tax.ApprovedAddons() {
		byKey[ea.Key] = ea
	}

	// Each Brazil addon ships in its own module.
	modules := map[cbc.Key]string{
		"br-nfe-v4":  "github.com/invopop/gobl.br.nfe",
		"br-nfse-v1": "github.com/invopop/gobl.br.nfse",
	}
	for key, module := range modules {
		ea, ok := byKey[key]
		require.Truef(t, ok, "expected %s on the approved list", key)
		assert.Equal(t, module, ea.Module, "%s module", key)
		assert.NotEmpty(t, ea.Name.String(), "%s should carry a name", key)

		// The implementation is external, so the key is not runtime-registered
		// in core — recognition only, not a runtime bypass.
		assert.Nil(t, tax.AddonForKey(key), "%s should not be registered in core", key)
	}
}
