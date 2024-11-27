package verifactu

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
)

const (
	// TagSubstitution is used to identify the case for:
	// Factura emitida en sustitución de facturas simplificadas facturadas y declaradas.
	// To be used when a simplified invoice has been issued and declared.
	TagSubstitution cbc.Key = "substitution"
)

var invoiceTags = &tax.TagSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*cbc.KeyDefinition{
		{
			Key: TagSubstitution,
			Name: i18n.String{
				i18n.EN: "Invoice issued in substitution of simplified invoices issued and declared",
				i18n.ES: "Factura emitida en sustitución de facturas simplificadas facturadas y declaradas",
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
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "F1",
				},
			},
			{
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
				},
				Tags: []cbc.Key{
					tax.TagSimplified,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "F2",
				},
			},
			{
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
				},
				Tags: []cbc.Key{
					TagSubstitution,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "F3",
				},
			},
			{
				Types: es.InvoiceCorrectionTypes,
				Ext: tax.Extensions{
					ExtKeyDocType: "R1",
				},
			},
			{
				Types: es.InvoiceCorrectionTypes,
				Tags: []cbc.Key{
					tax.TagSimplified,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "R5",
				},
			},
		},
	},
}
