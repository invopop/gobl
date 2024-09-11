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
			Ext: tax.Extensions{
				ExtKeySAFTInvoiceType: "FT",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{tax.TagSimplified},
			Ext: tax.Extensions{
				ExtKeySAFTInvoiceType: "FS",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Filter: func(doc any) bool {
				inv, ok := doc.(*bill.Invoice)
				if !ok {
					return false
				}
				return inv.Tax.ContainsTag(TagInvoiceReceipt) || inv.Totals.Paid()
			},
			Ext: tax.Extensions{
				ExtKeySAFTInvoiceType: "FR",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeDebitNote},
			Ext: tax.Extensions{
				ExtKeySAFTInvoiceType: "ND",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Ext: tax.Extensions{
				ExtKeySAFTInvoiceType: "NC",
			},
		},

		// Extension texts
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M01",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Artigo 16.°, n.° 6 do CIVA",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M02",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Artigo 6.° do Decreto-Lei n.° 198/90, de 19 de junho",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M04",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Exempt / Isento artigo 13.° do CIVA",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M05",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Exempt / Isento artigo 14.° do CIVA",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M06",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Exempt / Isento artigo 15.° do CIVA",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M07",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Exempt / Isento artigo 9.° do CIVA",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M09",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Não confere direito a dedução / Artigo 62.° alínea b) do CIVA",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M10",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Regime de isenção / Artigo 57.° do CIVA",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M11",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Regime particular do tabaco / Decreto-Lei n.° 346/85, de 23 de agosto",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M12",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Regime da margem de lucro - Agências de viagens / Decreto-Lei n.° 221/85, de 3 de julho",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M13",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Regime da margem de lucro - Bens em segunda mão / Decreto-Lei n.° 199/96, de 18 de outubro",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M14",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Regime da margem de lucro - Objetos de arte / Decreto-Lei n.° 199/96, de 18 de outubro",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M15",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Regime da margem de lucro - Objetos de coleção e antiguidades / Decreto-Lei n.° 199/96, de 18 de outubro",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M16",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Exempt / Isento artigo 14.° do RITI",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M19",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Outras isenções - Isenções temporárias determinadas em diploma próprio",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M20",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Regime forfetário / Artigo 59.°-D n.°2 do CIVA",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M21",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Não confere direito à dedução (ou expressão similar) - Artigo 72.° n.° 4 do CIVA",
			},
		},
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M25",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Mercadorias à consignação - Artigo 38.° n.° 1 alínea a) do CIVA",
			},
		},
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
		{
			ExtKey:   ExtKeyExemptionCode,
			ExtValue: "M99",
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  ExtKeyExemptionCode,
				Text: "Não sujeito ou não tributado",
			},
		},
	},
}
