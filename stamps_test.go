package gobl_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestStampOwners(t *testing.T) {
	for _, sd := range head.AllStampDefs() {
		t.Run(sd.Key.String(), func(t *testing.T) {
			if sd.Addon != "" {
				assert.Empty(t, sd.Regime)
				assert.NotNil(t, tax.AddonForKey(sd.Addon), "unknown addon '%s'", sd.Addon)
				return
			}
			assert.NotNil(t, tax.RegimeDefFor(sd.Regime.Code()), "unknown regime '%s'", sd.Regime)
		})
	}
}
