package gobl_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestStampsRegistered(t *testing.T) {
	for _, r := range tax.AllRegimeDefs() {
		for _, k := range r.Stamps {
			assert.NotNil(t, head.StampDefForKey(k), "regime '%s' stamp '%s' not registered", r.Code(), k)
		}
	}
	for _, a := range tax.AllAddonDefs() {
		for _, k := range a.Stamps {
			assert.NotNil(t, head.StampDefForKey(k), "addon '%s' stamp '%s' not registered", a.Key, k)
		}
	}
}
