package addons_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestAddons(t *testing.T) {
	for _, ad := range tax.AllAddonDefs() {
		t.Run(ad.Key.String(), func(t *testing.T) {
			assert.NoError(t, rules.Validate(ad))
		})
	}
}

func TestAddonStampUsage(t *testing.T) {
	for _, ad := range tax.AllAddonDefs() {
		for _, sd := range ad.Stamps {
			t.Run(sd.Key.String(), func(t *testing.T) {
				if assert.NotNil(t, sd.Usage) {
					assert.Contains(t, sd.Usage.Keys, ad.Key)
				}
			})
		}
	}
}
