package fr

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
			Tags:       []cbc.Key{tax.TagReverseCharge},
			Categories: []cbc.Code{tax.CategoryVAT},
			Note: &tax.Note{
				Category: tax.CategoryVAT,
				Key:      tax.KeyReverseCharge,
				Text:     "Reverse Charge / Autoliquidation de la TVA - Article 283-1 du CGI. Le client est redevable de la TVA.",
			},
		},
	},
}
