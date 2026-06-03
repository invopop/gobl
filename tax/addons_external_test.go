package tax_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApprovedAddons(t *testing.T) {
	approved := tax.ApprovedAddons()
	require.NotEmpty(t, approved)

	byKey := make(map[cbc.Key]*tax.ExternalAddon, len(approved))
	for _, ea := range approved {
		byKey[ea.Key] = ea
	}

	for _, key := range []cbc.Key{"fr-ctc-v1", "fr-ctc-flow2-v1", "fr-ctc-flow6-v1", "fr-ctc-flow10-v1"} {
		ea, ok := byKey[key]
		require.Truef(t, ok, "expected %s on the approved list", key)
		assert.NotEmpty(t, ea.Name.String(), "%s should carry a name", key)
		assert.Equal(t, "github.com/invopop/gobl.fr.ctc", ea.Module, "%s module", key)
	}

	// None of the approved CTC keys are runtime-registered in core (their
	// implementation lives in the external module), so AddonForKey is nil.
	assert.Nil(t, tax.AddonForKey("fr-ctc-v1"))
}

func TestApprovedAddonInJSONSchema(t *testing.T) {
	js := new(jsonschema.Schema)
	require.NoError(t, json.Unmarshal([]byte(`{
		"type": "array",
		"items": { "$ref": "https://gobl.org/draft-0/cbc/key" }
	}`), js))

	tax.AddonList{}.JSONSchemaExtend(js)

	consts := make(map[string]string)
	for _, o := range js.Items.OneOf {
		if k, ok := o.Const.(string); ok {
			consts[k] = o.Title
		}
	}
	// Approved-but-not-registered keys still appear as valid $addons options.
	assert.Contains(t, consts, "fr-ctc-v1")
	assert.Equal(t, "France CTC Flow 6 (Cycle de Vie)", consts["fr-ctc-flow6-v1"])
}

// TestApprovedAddonStillRequiresRegistration locks in the key contract: being
// on the approved list makes a key schema-valid but is NOT a runtime bypass.
// Validating a document that declares an approved-but-unloaded addon must still
// fail "add-on must be registered".
func TestApprovedAddonStillRequiresRegistration(t *testing.T) {
	type testStruct struct {
		tax.Addons
		Name string `json:"test"`
	}
	ts := &testStruct{Addons: tax.WithAddons("fr-ctc-v1"), Name: "Test"}

	err := rules.Validate(ts)
	assert.ErrorContains(t, err, "add-on must be registered")
}
