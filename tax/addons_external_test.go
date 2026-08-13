package tax_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These tests exercise the approved external-addon registry mechanism with a
// synthetic key, independent of the curated list (which lives in the addons
// package). See addons/external_test.go for the fr-ctc entries.

func TestApprovedAddonRegistry(t *testing.T) {
	tax.RegisterApprovedAddon(&tax.ExternalAddon{
		Key:    "test-external-v1",
		Name:   i18n.String{i18n.EN: "Test External Addon"},
		Module: "github.com/example/test",
	})

	var got *tax.ExternalAddon
	for _, ea := range tax.ApprovedAddons() {
		if ea.Key == "test-external-v1" {
			got = ea
		}
	}
	require.NotNil(t, got, "approved addon should be listed")
	assert.Equal(t, "github.com/example/test", got.Module)

	// Being approved does not register a usable addon at runtime.
	assert.Nil(t, tax.AddonForKey("test-external-v1"))
}

func TestApprovedAddonInJSONSchema(t *testing.T) {
	tax.RegisterApprovedAddon(&tax.ExternalAddon{
		Key:  "test-external-schema-v1",
		Name: i18n.String{i18n.EN: "Test External Schema Addon"},
	})

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
	// Approved-but-not-registered keys still appear as valid $addons options...
	assert.Equal(t, "Test External Schema Addon", consts["test-external-schema-v1"])
	// ...alongside the runtime-registered ones.
	require.NotEmpty(t, tax.AllAddonDefs())
	assert.Contains(t, consts, tax.AllAddonDefs()[0].Key.String())
}

// TestApprovedAddonStillRequiresRegistration locks in the key contract: being
// on the approved list makes a key schema-valid but is NOT a runtime bypass.
// Validating a document that declares an approved-but-unloaded addon must still
// fail "add-on must be registered".
func TestApprovedAddonStillRequiresRegistration(t *testing.T) {
	tax.RegisterApprovedAddon(&tax.ExternalAddon{
		Key:  "test-external-runtime-v1",
		Name: i18n.String{i18n.EN: "Test External Runtime Addon"},
	})

	type testStruct struct {
		tax.Addons
		Name string `json:"test"`
	}
	ts := &testStruct{Addons: tax.WithAddons("test-external-runtime-v1"), Name: "Test"}

	err := rules.Validate(ts)
	assert.ErrorContains(t, err, "add-on must be registered")
}
