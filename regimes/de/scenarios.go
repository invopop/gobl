package de

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

var invoiceTags = []*tax.KeyDefinition{
	// Reverse Charge Mechanism
	{
		Key: common.TagReverseCharge,
		Name: i18n.String{
			i18n.EN: "Reverse Charge",
			i18n.DE: "Umkehr der Steuerschuld",
		},
	},
}

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		// ** Special Messages **
		// Reverse Charges
		{
			Tags: []cbc.Key{common.TagReverseCharge},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  common.TagReverseCharge,
				Text: "Reverse Charge / Umkehr der Steuerschuld.",
			},
		},
	},
}
