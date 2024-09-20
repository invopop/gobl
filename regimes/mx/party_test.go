package mx_test

import (
	"testing"

	"github.com/invopop/gobl/addons/mx/cfdi"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestMigratePartyIdentities(t *testing.T) {
	customer := &org.Party{
		Name: "Test Customer",
		TaxID: &tax.Identity{
			Country: "MX",
			Code:    "ZZZ010101ZZZ",
			Zone:    "65000",
		},
	}

	mx := tax.RegimeDefFor("MX")
	mx.Normalizer(customer)

	assert.Equal(t, "65000", customer.Ext[cfdi.ExtKeyPostCode].String())
	assert.Empty(t, customer.TaxID.Zone) //nolint:staticcheck
}
