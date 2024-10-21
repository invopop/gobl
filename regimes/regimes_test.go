package regimes_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestRegimes(t *testing.T) {
	for _, ad := range tax.AllRegimeDefs() {
		t.Run(ad.Name.String(), func(t *testing.T) {
			assert.NoError(t, ad.Validate())
		})
	}
}
