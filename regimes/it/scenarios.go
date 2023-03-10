package it

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
		// **** TIPO DOCUMENTO ****
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Meta: cbc.Meta{
				KeyFatturaPATipoDocumento: "TD01",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{common.TagPartial},
			Meta: cbc.Meta{
				KeyFatturaPATipoDocumento: "TD02",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Meta: cbc.Meta{
				KeyFatturaPATipoDocumento: "TD04",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeDebitNote},
			Meta: cbc.Meta{
				KeyFatturaPATipoDocumento: "TD05",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{TagFreelance},
			Meta: cbc.Meta{
				KeyFatturaPATipoDocumento: "TD06",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{common.TagPartial, TagFreelance},
			Meta: cbc.Meta{
				KeyFatturaPATipoDocumento: "TD03",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{common.TagSimplified},
			Meta: cbc.Meta{
				KeyFatturaPATipoDocumento: "TD07",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Tags:  []cbc.Key{common.TagSimplified},
			Meta: cbc.Meta{
				KeyFatturaPATipoDocumento: "TD08",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeDebitNote},
			Tags:  []cbc.Key{common.TagSimplified},
			Meta: cbc.Meta{
				KeyFatturaPATipoDocumento: "TD09",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{common.TagReverseCharge},
			Meta: cbc.Meta{
				KeyFatturaPATipoDocumento: "TD16",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{common.TagSelfBilled, TagCeilingExceeded},
			Meta: cbc.Meta{
				KeyFatturaPATipoDocumento: "TD21",
			},
		},

		// **** REGIME FISCALE ****
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard, bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote},
			Meta: cbc.Meta{
				KeyFatturaPARegimeFiscale: "RF01", // default
			},
		},
		{
			Tags: []cbc.Key{TagMinimumTaxPayers},
			Meta: cbc.Meta{
				KeyFatturaPARegimeFiscale: "RF02",
			},
		},
	},
}
