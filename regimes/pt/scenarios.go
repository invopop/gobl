package pt

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

// Invoice type tags
const (
	TagInvoiceReceipt cbc.Key = "invoice-receipt"
	TagAzores         cbc.Key = "azores"
	TagMadeira        cbc.Key = "madeira"
)

var scenarios = []*tax.ScenarioSet{
	invoiceScenarios,
}

var invoiceTags = common.InvoiceTagsWith([]*tax.KeyDefinition{
	{
		Key: TagInvoiceReceipt,
		Name: i18n.String{
			i18n.EN: "Invoice-receipt",
			i18n.PT: "Fatura-recibo",
		},
	},
	{
		Key: TagAzores,
		Name: i18n.String{
			i18n.EN: "Azores",
			i18n.PT: "Açores",
		},
		Desc: i18n.String{
			i18n.EN: "Tag for use when the invoice is issued in the Azores region with special local rates.",
			i18n.PT: "Tag para uso quando a fatura é emitida na região dos Açores com taxas locais especiais.",
		},
	},
	{
		Key: TagMadeira,
		Name: i18n.String{
			i18n.EN: "Madeira",
			i18n.PT: "Madeira",
		},
		Desc: i18n.String{
			i18n.EN: "Tag for use when the invoice is issued in the Madeira region with special local rates.",
			i18n.PT: "Tag para uso quando a fatura é emitida na região da Madeira com taxas locais especiais.",
		},
	},
})

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Codes: cbc.CodeMap{
				KeyATInvoiceType: "FT",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{tax.TagSimplified},
			Codes: cbc.CodeMap{
				KeyATInvoiceType: "FS",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{TagInvoiceReceipt},
			Codes: cbc.CodeMap{
				KeyATInvoiceType: "FR",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeDebitNote},
			Codes: cbc.CodeMap{
				KeyATInvoiceType: "ND",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Codes: cbc.CodeMap{
				KeyATInvoiceType: "NC",
			},
		},

		// Reverse Charges
		{
			Tags: []cbc.Key{tax.TagReverseCharge},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  tax.TagReverseCharge,
				Text: "Reverse charge / Autoliquidação - Artigo 2.º n.º 1 alínea j) do Código do IVA",
			},
		},
	},
}
