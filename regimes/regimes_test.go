package regimes_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestRegimes(t *testing.T) {
	for _, ad := range tax.AllRegimeDefs() {
		t.Run(ad.Name.String(), func(t *testing.T) {
			assert.NoError(t, rules.Validate(ad))
		})
	}
}

func TestRegimeStampUsage(t *testing.T) {
	for _, rd := range tax.AllRegimeDefs() {
		for _, sd := range rd.Stamps {
			t.Run(sd.Key.String(), func(t *testing.T) {
				if assert.NotNil(t, sd.Usage) {
					assert.Contains(t, sd.Usage.Codes, cbc.Code(rd.Country))
				}
			})
		}
	}
}
