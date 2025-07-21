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

		// The following scenarios append a note with the applicable regulations
		// for any tax exemption used in the document. These texts are compliant
		// with art. 2.2.14 of Despacho nº8632/2014 (for printed invoices) and
		// point 4.4.19.7 of Portaria nº302/2016 (for SAF-T files).
		//
		// Codes below with a blank text are cases where the user must provide the
		// applicable regulation manually.
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M01",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Artigo 16.º, n.º 6, alíneas a) a d) do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M02",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Artigo 6.º do Decreto-Lei n.º 198/90, de 19 de junho",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M04",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Artigo 13.º do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M05",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Artigo 14.º do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M06",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Artigo 15.º do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M07",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Artigo 9.º do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M09",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Artigo 62.º alínea b) do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M10",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Artigo 57.º do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M11",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Decreto-Lei n.º 346/85, de 23 de agosto",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M12",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Decreto-Lei n.º 221/85, de 3 de julho",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M13",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Decreto-Lei n.º 199/96, de 18 de outubro",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M14",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Decreto-Lei n.º 199/96, de 18 de outubro",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M15",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Decreto-Lei n.º 199/96, de 18 de outubro",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M16",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Artigo 14.º do RITI",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M19",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Outras isenções", // default text, to be overridden by the user
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M20",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Artigo 59.º-D n.º 2 do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M21",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Artigo 72.º n.º 4 do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M25",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Artigo 38.º n.º 1 alínea a) do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M26",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Lei n.º 17/2023, de 14 de abril",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M30",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Artigo 2.º n.º 1 alínea i) do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M31",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Artigo 2.º n.º 1 alínea j) do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M32",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Artigo 2.º n.º 1 alínea l) do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M33",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Artigo 2.º n.º 1 alínea m) do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M34",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Artigo 2.º n.º 1 alínea n) do CIVA",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M40",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Artigo 6.º n.º 6 alínea a) do CIVA, a contrário",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M41",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Artigo 8.º n.º 3 do RITI",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M42",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Decreto-Lei n.º 21/2007, de 29 de janeiro",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M43",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Decreto-Lei n.º 362/99, de 16 de setembro",
			},
		},
		{
			ExtKey:  ExtKeyExemption,
			ExtCode: "M99",
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  ExtKeyExemption,
				Text: "Não sujeito ou não tributado", // default text, to be overridden by the user
			},
		},
	},
}
