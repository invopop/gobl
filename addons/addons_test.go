package addons_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestAddons(t *testing.T) {
	for _, ad := range tax.AllAddons() {
		t.Run(ad.Key.String(), func(t *testing.T) {
			assert.NoError(t, ad.Validate())
		})
	}
}
