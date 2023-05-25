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
	TagFreelance       cbc.Key = "freelance"
	TagCeilingExceeded cbc.Key = "ceiling-exceeded"
	TagSanMerinoPaper  cbc.Key = "san-merino-paper"

	// Tags for Fiscal Regime
	TagMinimumTaxPayers cbc.Key = "minimum-tax-payers"
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
		Key: TagSanMerinoPaper,
		Name: i18n.String{
			i18n.EN: "Purchases from San Merino with VAT and paper invoice",
			i18n.IT: "Acquisti da San Merino con IVA e fattura cartacea",
		},
	},

	// **** Fiscal Regime Tags ****
	{
		Key: TagMinimumTaxPayers,
		Name: i18n.String{
			i18n.EN: "Minimum Taxpayers",
			i18n.IT: "Contribuenti minimi",
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
			Tags:  []cbc.Key{common.TagReverseCharge},
			Name: i18n.String{
				i18n.EN: "Reverse charge internal invoice integration",
				i18n.IT: "Integrazione fattura reverse charge interno",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPATipoDocumento: "TD16",
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

		// **** REGIME FISCALE ****
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard, bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote},
			Name: i18n.String{
				i18n.EN: "Regular regime",
				i18n.IT: "Regime ordinario",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPARegimeFiscale: "RF01", // default
			},
		},
		{
			Tags: []cbc.Key{TagMinimumTaxPayers},
			Name: i18n.String{
				i18n.EN: "Minimum tax payers (art. 1,c.96-117, L. 244/2007)",
				i18n.IT: "Regime dei contribuenti minimi (art. 1,c.96-117, L. 244/2007)",
			},
			Codes: cbc.CodeSet{
				KeyFatturaPARegimeFiscale: "RF02",
			},
		},
	},
}
