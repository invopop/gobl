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

func TestApprovedAddons(t *testing.T) {
	byKey := make(map[cbc.Key]*tax.ExternalAddon)
	for _, ea := range tax.ApprovedAddons() {
		byKey[ea.Key] = ea
	}

	groups := []struct {
		name    string
		modules map[cbc.Key]string
	}{
		{
			name: "FR CTC",
			modules: map[cbc.Key]string{
				"fr-ctc-flow2-v1":  "github.com/invopop/gobl.fr.ctc",
				"fr-ctc-flow6-v1":  "github.com/invopop/gobl.fr.ctc",
				"fr-ctc-flow10-v1": "github.com/invopop/gobl.fr.ctc",
			},
		},
		{
			name: "PT SAF-T",
			modules: map[cbc.Key]string{
				"pt-saft-v1": "github.com/invopop/gobl.pt.saft",
			},
		},
		{
			name: "BR",
			modules: map[cbc.Key]string{
				"br-nfe-v4":  "github.com/invopop/gobl.br.nfe",
				"br-nfse-v1": "github.com/invopop/gobl.br.nfse",
			},
		},
		{
			name: "MX CFDI",
			modules: map[cbc.Key]string{
				"mx-cfdi-v4": "github.com/invopop/gobl.mx.cfdi",
			},
		},
	}

	for _, group := range groups {
		t.Run(group.name, func(t *testing.T) {
			for key, module := range group.modules {
				ea, ok := byKey[key]
				require.Truef(t, ok, "expected %s on the approved list", key)
				assert.Equal(t, module, ea.Module, "%s module", key)
				assert.NotEmpty(t, ea.Name.String(), "%s should carry a name", key)

				// The implementation is external, so the key is not runtime-registered
				// in core — recognition only, not a runtime bypass.
				assert.Nil(t, tax.AddonForKey(key), "%s should not be registered in core", key)
			}
		})
	}
}

func TestApprovedDKOIOUBLAddon(t *testing.T) {
	byKey := make(map[cbc.Key]*tax.ExternalAddon)
	for _, ea := range tax.ApprovedAddons() {
		byKey[ea.Key] = ea
	}

	ea, ok := byKey["dk-oioubl-v2"]
	require.True(t, ok, "expected dk-oioubl-v2 on the approved list")
	assert.Equal(t, "github.com/invopop/gobl.dk.oioubl", ea.Module)
	assert.NotEmpty(t, ea.Name.String())
	assert.Nil(t, tax.AddonForKey("dk-oioubl-v2"), "dk-oioubl-v2 should not be registered in core")
}
