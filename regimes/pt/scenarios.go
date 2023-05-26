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
			Meta: cbc.Meta{
				KeyATInvoiceType: "FT",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{common.TagSimplified},
			Meta: cbc.Meta{
				KeyATInvoiceType: "FS",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{TagInvoiceReceipt},
			Meta: cbc.Meta{
				KeyATInvoiceType: "FR",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeDebitNote},
			Meta: cbc.Meta{
				KeyATInvoiceType: "ND",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Meta: cbc.Meta{
				KeyATInvoiceType: "NC",
			},
		},
	},
}
