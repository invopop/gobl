package gobl_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestStampDefs(t *testing.T) {
	for _, sd := range head.AllStampDefs() {
		t.Run(sd.Key.String(), func(t *testing.T) {
			assert.NoError(t, rules.Validate(sd))
			for lang, name := range sd.Name {
				assert.Contains(t, name, ": ", "name in '%s' must be prefixed with the issuer", lang)
			}
			if !assert.NotNil(t, sd.Usage) {
				return
			}
			assert.True(t, len(sd.Usage.Keys) > 0 || len(sd.Usage.Codes) > 0, "usage must have keys or codes")
			for _, k := range sd.Usage.Keys {
				assert.NotNil(t, tax.AddonForKey(k), "addon '%s' not registered", k)
			}
			for _, c := range sd.Usage.Codes {
				assert.NotNil(t, tax.RegimeDefFor(l10n.Code(c)), "regime '%s' not registered", c)
			}
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
