package saft

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/pt"
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
			Ext: tax.Extensions{
				ExtKeyInvoiceType: InvoiceTypeStandard,
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{tax.TagSimplified},
			Ext: tax.Extensions{
				ExtKeyInvoiceType: InvoiceTypeSimplified,
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Filter: func(doc any) bool {
				inv, ok := doc.(*bill.Invoice)
				if !ok {
					return false
				}
				return inv.HasTags(pt.TagInvoiceReceipt) || inv.Totals.Paid()
			},
			Ext: tax.Extensions{
				ExtKeyInvoiceType: InvoiceTypeInvoiceReceipt,
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeDebitNote},
			Ext: tax.Extensions{
				ExtKeyInvoiceType: InvoiceTypeDebitNote,
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Ext: tax.Extensions{
				ExtKeyInvoiceType: InvoiceTypeCreditNote,
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeProforma},
			Ext: tax.Extensions{
				ExtKeyWorkType: WorkTypeProforma,
			},
		},

		// Extension texts
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M01",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Artigo 16.°, n.° 6 do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M02",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Artigo 6.° do Decreto-Lei n.° 198/90, de 19 de junho",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M04",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Exempt / Isento artigo 13.° do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M05",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Exempt / Isento artigo 14.° do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M06",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Exempt / Isento artigo 15.° do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M07",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Exempt / Isento artigo 9.° do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M09",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Não confere direito a dedução / Artigo 62.° alínea b) do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M10",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Regime de isenção / Artigo 57.° do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M11",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Regime particular do tabaco / Decreto-Lei n.° 346/85, de 23 de agosto",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M12",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Regime da margem de lucro - Agências de viagens / Decreto-Lei n.° 221/85, de 3 de julho",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M13",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Regime da margem de lucro - Bens em segunda mão / Decreto-Lei n.° 199/96, de 18 de outubro",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M14",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Regime da margem de lucro - Objetos de arte / Decreto-Lei n.° 199/96, de 18 de outubro",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M15",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Regime da margem de lucro - Objetos de coleção e antiguidades / Decreto-Lei n.° 199/96, de 18 de outubro",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M16",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Exempt / Isento artigo 14.° do RITI",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M19",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Outras isenções - Isenções temporárias determinadas em diploma próprio",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M20",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Regime forfetário / Artigo 59.°-D n.°2 do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M21",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Não confere direito à dedução (ou expressão similar) - Artigo 72.° n.° 4 do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M25",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Mercadorias à consignação - Artigo 38.° n.° 1 alínea a) do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M30",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Reverse charge / autoliquidação - Artigo 2.° n.° 1 alínea i) do Código do IVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M31",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Reverse charge / autoliquidação - Artigo 2.° n.° 1 alínea j) do Código do IVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M32",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Reverse charge / autoliquidação - Artigo 2.° n.° 1 alínea I) do Código do IVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M33",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Reverse charge / autoliquidação - Artigo 2.° n.° 1 alínea m) do Código do IVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M40",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Reverse charge / Autoliquidação - Artigo 6.º n.º 6 alínea a) do Código do IVA, a contrário",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M41",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Reverse charge / Autoliquidação - Artigo 8.º n.º 3 do RITI",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M42",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Reverse charge / Autoliquidação - Decreto-Lei n.º 21/2007, de 29 de janeiro",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M43",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Reverse charge / Autoliquidação - Decreto-Lei n.° 362/99, de 16 de setembro",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M99",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Não sujeito ou não tributado",
			},
		},
	},
}
