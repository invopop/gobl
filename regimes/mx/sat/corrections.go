package sat

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

// CorrectionDefinitions provides the array of correction definitions that apply to SAT
func CorrectionDefinitions() []*tax.CorrectionDefinition {
	return correctionDefinitions
}

var correctionDefinitions = []*tax.CorrectionDefinition{
	{
		Schema: bill.ShortSchemaInvoice,
		Types: []cbc.Key{
			bill.InvoiceTypeCreditNote,
		},
		Stamps: []cbc.Key{
			StampUUID,
		},
	},
}
