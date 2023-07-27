package mx

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

var precedingDefinitions = &tax.PrecedingDefinitions{
	Types: []cbc.Key{
		bill.InvoiceTypeCreditNote,
	},
	Stamps: []cbc.Key{
		StampProviderSATUUID,
	},
}
