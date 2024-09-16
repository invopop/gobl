package tax_test

import (
	"context"
	"testing"

	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmbeddingAddons(t *testing.T) {
	type testStruct struct {
		tax.Addons
		Name string `json:"test"`
	}

	ts := &testStruct{
		Addons: tax.WithAddons("mx-cfdi-v4"),
		Name:   "Test",
	}
	assert.NotNil(t, ts.Addons)
	assert.Equal(t, "Test", ts.Name)

	defs := ts.GetAddons()
	assert.Len(t, defs, 1)
	assert.Equal(t, "mx-cfdi-v4", defs[0].Key.String())

	ts.Addons = tax.WithAddons("mx-cfdi-v4", "invalid-addon")

	err := ts.Addons.Validate()
	assert.ErrorContains(t, err, "1: addon 'invalid-addon' not registered")
}

func TestAddonForKey(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		a := tax.AddonForKey("unknown")
		assert.Nil(t, a)
	})

	t.Run("found", func(t *testing.T) {
		a := tax.AddonForKey("mx-cfdi-v4")
		require.NotNil(t, a)
		assert.NoError(t, a.Validate())
	})
}

func TestAllAddons(t *testing.T) {
	as := tax.AllAddons()
	assert.NotEmpty(t, as)
}

func TestAddonWithContext(t *testing.T) {
	ad := tax.AddonForKey("mx-cfdi-v4")
	ctx := ad.WithContext(context.Background())

	vs := tax.Validators(ctx)
	assert.Len(t, vs, 1)
	// no reliable way to check the function is actually the same :-(
}
