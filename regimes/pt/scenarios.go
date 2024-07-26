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
)

var scenarios = []*tax.ScenarioSet{
	invoiceScenarios,
}

var invoiceTags = common.InvoiceTagsWith([]*cbc.KeyDefinition{
	{
		Key: TagInvoiceReceipt,
		Name: i18n.String{
			i18n.EN: "Invoice-receipt",
			i18n.PT: "Fatura-recibo",
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

		// Extension texts

		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M30",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Reverse charge / autoliquidação - Artigo 2.° n.° 1 alínea i) do Código do IVA",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M31",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Reverse charge / autoliquidação - Artigo 2.° n.° 1 alínea j) do Código do IVA",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M32",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Reverse charge / autoliquidação - Artigo 2.° n.° 1 alínea I) do Código do IVA",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M33",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Reverse charge / autoliquidação - Artigo 2.° n.° 1 alínea m) do Código do IVA",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M40",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Reverse charge / Autoliquidação - Artigo 6.º n.º 6 alínea a) do Código do IVA, a contrário",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M41",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Reverse charge / Autoliquidação - Artigo 8.º n.º 3 do RITI",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M42",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Reverse charge / Autoliquidação - Decreto-Lei n.º 21/2007, de 29 de janeiro",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M43",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Reverse charge / Autoliquidação - Decreto-Lei n.° 362/99, de 16 de setembro",
			},
		},
	},
}
