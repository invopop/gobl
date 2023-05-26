package pt

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

var scenarios = []*tax.ScenarioSet{
	invoiceScenarios,
}

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Codes: cbc.CodeSet{
				KeyATInvoiceType: "FT",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{common.TagSimplified},
			Codes: cbc.CodeSet{
				KeyATInvoiceType: "FS",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{TagInvoiceReceipt},
			Codes: cbc.CodeSet{
				KeyATInvoiceType: "FR",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeDebitNote},
			Codes: cbc.CodeSet{
				KeyATInvoiceType: "ND",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Codes: cbc.CodeSet{
				KeyATInvoiceType: "NC",
			},
		},
	},
}
