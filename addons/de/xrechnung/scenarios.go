package xrechnung

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Document tag keys
const (
	// Tags for invoice types
	TagSelfBilled               cbc.Key = "self-billed"
	TagPartial                  cbc.Key = "partial"
	TagPartialConstruction      cbc.Key = "partial-construction"
	TagPartialFinalConstruction cbc.Key = "partial-final-construction"
	TagFinalConstruction        cbc.Key = "final-construction"
)

// Invoice type constants
const (
	invoiceTypeSelfBilled               = "380"
	invoiceTypePartial                  = "326"
	invoiceTypePartialConstruction      = "80"
	invoiceTypePartialFinalConstruction = "84"
	invoiceTypeFinalConstruction        = "389"
)

var invoiceTags = &tax.TagSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*cbc.KeyDefinition{
		{
			Key: TagSelfBilled,
			Name: i18n.String{
				i18n.EN: "Self-billed Invoice",
				i18n.DE: "Gutschrift",
			},
		},
		{
			Key: TagPartial,
			Name: i18n.String{
				i18n.EN: "Partial Invoice",
				i18n.DE: "Abschlagsrechnung",
			},
		},
		{
			Key: TagPartialConstruction,
			Name: i18n.String{
				i18n.EN: "Partial Construction Invoice",
				i18n.DE: "Abschlagsrechnung (Bauleistung)",
			},
		},
		{
			Key: TagPartialFinalConstruction,
			Name: i18n.String{
				i18n.EN: "Partial Final Construction Invoice",
				i18n.DE: "Schlussrechnung (Bauleistung)",
			},
		},
		{
			Key: TagFinalConstruction,
			Name: i18n.String{
				i18n.EN: "Final Construction Invoice",
				i18n.DE: "Schlussrechnung",
			},
		},
	},
}

var scenarios = []*tax.ScenarioSet{
	{
		Schema: bill.ShortSchemaInvoice,
		List: []*tax.Scenario{
			// ** Invoice Document Types **
			{
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
					bill.InvoiceTypeCorrective,
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
			},
			{
				Tags: []cbc.Key{
					tax.TagSelfBilled,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: tax.ExtValue(invoiceTypeSelfBilled),
				},
			},
			{
				Tags: []cbc.Key{
					tax.TagPartial,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: tax.ExtValue(invoiceTypePartial),
				},
			},
			{
				Tags: []cbc.Key{
					tax.TagPartial,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: tax.ExtValue(invoiceTypePartialConstruction),
				},
			},
			{
				Tags: []cbc.Key{
					tax.TagPartial,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: tax.ExtValue(invoiceTypePartialFinalConstruction),
				},
			},
			{
				Ext: tax.Extensions{
					ExtKeyDocType: tax.ExtValue(invoiceTypeFinalConstruction),
				},
			},
			// ** Tax Rates **
			{
				Tags: []cbc.Key{
					tax.RateStandard,
				},
				Ext: tax.Extensions{
					ExtKeyTaxRate: "S",
				},
			},
			{
				Tags: []cbc.Key{
					tax.RateZero,
				},
				Ext: tax.Extensions{
					ExtKeyTaxRate: "Z",
				},
			},
			{
				Tags: []cbc.Key{
					tax.RateExempt,
				},
				Ext: tax.Extensions{
					ExtKeyTaxRate: "E",
				},
			},
			{
				Tags: []cbc.Key{
					tax.TagReverseCharge,
				},
				Ext: tax.Extensions{
					ExtKeyTaxRate: "AE",
				},
			},
			{
				Ext: tax.Extensions{
					ExtKeyTaxRate: "K",
				},
			},
			{
				Ext: tax.Extensions{
					ExtKeyTaxRate: "G",
				},
			},
			{
				Ext: tax.Extensions{
					ExtKeyTaxRate: "O",
				},
			},
		},
	},
}
