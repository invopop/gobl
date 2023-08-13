package it

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

// Document tag keys
const (
	// Tags for document type
	TagFreelance         cbc.Key = "freelance"
	TagCeilingExceeded   cbc.Key = "ceiling-exceeded"
	TagSanMarinoPaper    cbc.Key = "san-marino-paper"
	TagImport            cbc.Key = "import"
	TagGoods             cbc.Key = "goods"
	TagGoodsEU           cbc.Key = "goods-eu"
	TagGoodsWithTax      cbc.Key = "goods-with-tax"
	TagGoodsExtracted    cbc.Key = "goods-extracted"
	TagRegularization    cbc.Key = "regularization"
	TagDeferred          cbc.Key = "deferred"
	TagThirdPeriod       cbc.Key = "third-period"
	TagDepreciableAssets cbc.Key = "depreciable-assets"
)

// This is only a partial list of all the potential tags that
// could be available for use in Italy. Given the complexity
// involved, we've focussed here on the most useful.
var invoiceTags = []*tax.KeyDefinition{
	// *** Document Type Tags ***
	{
		Key: TagFreelance,
		Name: i18n.String{
			i18n.EN: "Freelancer",
			i18n.IT: "Parcella",
		},
	},
	{
		Key: common.TagReverseCharge,
		Name: i18n.String{
			i18n.EN: "Reverse Charge",
			i18n.IT: "Inversione del soggetto passivo",
		},
	},
	{
		Key: common.TagSelfBilled,
		Name: i18n.String{
			i18n.EN: "Self-billed",
			i18n.IT: "Autofattura",
		},
	},
	{
		Key: TagCeilingExceeded,
		Name: i18n.String{
			i18n.EN: "Ceiling exceeded",
			i18n.IT: "Splafonamento",
		},
	},
	{
		Key: TagSanMarinoPaper,
		Name: i18n.String{
			i18n.EN: "Purchases from San Marino with VAT and paper invoice",
			i18n.IT: "Acquisti da San Marino con IVA e fattura cartacea",
		},
	},
	{
		Key: TagImport,
		Name: i18n.String{
			i18n.EN: "Import",
			i18n.IT: "Importazione",
		},
	},
	{
		Key: TagGoods,
		Name: i18n.String{
			i18n.EN: "Goods",
			i18n.IT: "Beni",
		},
	},
	{
		Key: TagGoodsEU,
		Name: i18n.String{
			i18n.EN: "Goods from EU",
			i18n.IT: "Beni da UE",
		},
	},
	{
		Key: TagGoodsWithTax,
		Name: i18n.String{
			i18n.EN: "Goods with tax",
			i18n.IT: "Beni con imposta",
		},
	},
	{
		Key: TagGoodsExtracted,
		Name: i18n.String{
			i18n.EN: "Goods extracted",
			i18n.IT: "Beni estratti",
		},
	},
	{
		Key: TagDeferred,
		Name: i18n.String{
			i18n.EN: "Deferred",
			i18n.IT: "Differita",
		},
	},
	{
		Key: TagRegularization,
		Name: i18n.String{
			i18n.EN: "Regularization",
			i18n.IT: "Regolarizzazione",
		},
	},
	{
		Key: TagThirdPeriod,
		Name: i18n.String{
			i18n.EN: "Third period",
			i18n.IT: "Terzo periodo",
		},
	},
	{
		Key: TagDepreciableAssets,
		Name: i18n.String{
			i18n.EN: "Depreciable assets",
			i18n.IT: "Beni ammortizzabili",
		},
	},
}

var scenarios = []*tax.ScenarioSet{
	invoiceScenarios,
}

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		// **** TIPO DOCUMENTO ****
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Name: i18n.String{
				i18n.EN: "Regular Invoice",
				i18n.IT: "Fattura",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD01",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{common.TagPartial},
			Name: i18n.String{
				i18n.EN: "Advance or down payment on invoice",
				i18n.IT: "Acconto / anticipo su fattura",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD02",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Name: i18n.String{
				i18n.EN: "Credit Note",
				i18n.IT: "Nota di credito",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD04",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeDebitNote},
			Name: i18n.String{
				i18n.EN: "Debit Note",
				i18n.IT: "Nota di debito",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD05",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{TagFreelance},
			Name: i18n.String{
				i18n.EN: "Freelancer invoice with retained taxes",
				i18n.IT: "Parcella",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD06",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{common.TagPartial, TagFreelance},
			Name: i18n.String{
				i18n.EN: "Advance or down payment on freelance invoice",
				i18n.IT: "Acconto / anticipo su parcella",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD03",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{common.TagSimplified},
			Name: i18n.String{
				i18n.EN: "Simplified Invoice",
				i18n.IT: "Fattura Semplificata",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD07",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Tags:  []cbc.Key{common.TagSimplified},
			Name: i18n.String{
				i18n.EN: "Simplified Credit Note",
				i18n.IT: "Nota di credito semplificata",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD08",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeDebitNote},
			Tags:  []cbc.Key{common.TagSimplified},
			Name: i18n.String{
				i18n.EN: "Simplified Debit Note",
				i18n.IT: "Nota di debito semplificata",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD09",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{common.TagSelfBilled},
			Name: i18n.String{
				i18n.EN: "Self-billed for self consumption or for free transfer without recourse",
				i18n.IT: "Fattura per autoconsumo o per cessioni gratuite senza rivalsa",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD27", // order is important
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{common.TagSelfBilled, common.TagReverseCharge},
			Name: i18n.String{
				i18n.EN: "Reverse charge",
				i18n.IT: "Integrazione fattura reverse charge interno",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD16",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{common.TagSelfBilled, TagImport},
			Name: i18n.String{
				i18n.EN: "Self-billed Import",
				i18n.IT: "Integrazione/autofattura per acquisto servizi da estero",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD17",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{common.TagSelfBilled, TagImport, TagGoodsEU},
			Name: i18n.String{
				i18n.EN: "Self-billed EU Goods Import",
				i18n.IT: "Integrazione per acquisto beni intracomunitari",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD18",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{common.TagSelfBilled, TagImport, TagGoods},
			Name: i18n.String{
				i18n.EN: "Self-billed Goods Import",
				i18n.IT: "Integrazione/autofattura per acquisto beni ex art.17 c.2 DPR 633/72",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD19",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{common.TagSelfBilled, TagRegularization},
			Name: i18n.String{
				i18n.EN: "Self-billed Regularization",
				i18n.IT: "Autofattura per regolarizzazione e integrazione delle fatture - art.6 c.8 d.lgs.471/97 o art.46 c.5 D.L.331/93",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD20",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{common.TagSelfBilled, TagCeilingExceeded},
			Name: i18n.String{
				i18n.EN: "Self-billed invoice when ceiling exceeded",
				i18n.IT: "Autofattura per splafonamento",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD21",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{common.TagSelfBilled, TagGoodsExtracted},
			Name: i18n.String{
				i18n.EN: "Self-billed for goods extracted from VAT warehouse",
				i18n.IT: "Estrazione beni da Deposito IVA",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD22",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{common.TagSelfBilled, TagGoodsWithTax},
			Name: i18n.String{
				i18n.EN: "Self-billed for goods extracted from VAT warehouse with VAT payment",
				i18n.IT: "Estrazione beni da Deposito IVA con versamento IVA",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD23",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{TagDeferred},
			Name: i18n.String{
				i18n.EN: "Deferred invoice ex art.21, c.4, lett. a) DPR 633/72",
				i18n.IT: "Fattura differita - art.21 c.4 lett. a",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD24",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{TagDeferred, TagThirdPeriod},
			Name: i18n.String{
				i18n.EN: "Deferred invoice ex art.21, c.4, third period lett. b) DPR 633/72",
				i18n.IT: "Fattura differita - art.21 c.4 terzo periodo lett. b",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD25",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{TagDepreciableAssets},
			Name: i18n.String{
				i18n.EN: "Sale of depreciable assets and for internal transfers (ex art.36 DPR 633/72",
				i18n.IT: "Cessione di beni ammortizzabili e per passaggi interni - art.36 DPR 633/72",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD26",
			},
		},

		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{common.TagSelfBilled, TagSanMarinoPaper},
			Name: i18n.String{
				i18n.EN: "Purchases from San Marino with VAT (paper invoice)",
				i18n.IT: "Acquisti da San Marino con IVA (fattura cartacea)",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD28",
			},
		},

		// **** MESSAGES ****
		{
			Tags: []cbc.Key{common.TagReverseCharge},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  common.TagReverseCharge,
				Text: "Reverse Charge / Inversione del soggetto passivo",
			},
		},
	},
}
