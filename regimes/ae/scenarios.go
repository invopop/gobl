// Package ae provides tax scenarios specific to UAE VAT regulations.
package ae

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		// Reverse Charges
		{
			Tags: []cbc.Key{tax.TagReverseCharge},
			Note: &tax.Note{
				Category: tax.CategoryVAT,
				Key:      tax.KeyReverseCharge,
				Text:     "Reverse Charge",
			},
		},
		// Simplified Tax Invoice
		{
			Tags: []cbc.Key{tax.TagSimplified},
			Note: &tax.Note{
				Category: tax.CategoryVAT,
				Key:      tax.TagSimplified,
				Text:     "Simplified Tax Invoice",
			},
		},
	},
}
