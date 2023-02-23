package it

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// KeyDefinition defines properties of a key that is specific for a regime.
type KeyDefinition struct {
	// Actual key value.
	Key cbc.Key `json:"key" jsonschema:"title=Key"`
	// There is usually a mapping between a key and some local code.
	Code string `json:"code,omitempty" jsonschema:"title=Code"`
	// Short name for the key, if relevant.
	Name i18n.String `json:"name,omitempty" jsonschema:"title=Name"`
	// Description offering more details about when the key should be used.
	Desc i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`
	// Any additional data that might be relevant in some regimes?
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

const (
	// Tax System (RegimeFiscale) Codes
	FPACodeTaxSystemOrdinary cbc.Key = "tax-system-ordinary" // RF01

	// Payment Method (ModalitaPagamento) Codes
	FPACodePaymentMethodCash         cbc.Key = "tax-system-cash"          // MP01
	FPACodePaymentMethodBankTransfer cbc.Key = "tax-system-bank-transfer" // MP05
	FPACodePaymentMethodCard         cbc.Key = "tax-system-card"          // MP08
	FPACodePaymentDirectDebit        cbc.Key = "tax-system-direct-debit"  // MP10

	// Document Type (TipoDocumento) Codes
	FPACodeDocumentTypeInvoice    cbc.Key = "document-type-invoice"     // TD01
	FPACodeDocumentTypeCreditNote cbc.Key = "document-type-credit-note" // TD04

	// Withholding Tax (TipoRitenuta) Codes
	FPAWithholdingTaxNaturalPersons    cbc.Key = "withholding-tax-natural-persons"       // TR01
	FPAWithholdingTaxLegalPersons      cbc.Key = "withholding-tax-legal-persons"         // TR02
	FPAWithholdingINPSContribution     cbc.Key = "withholding-tax-inps-contribution"     // TR03
	FPAWithholdingENASARCOContribution cbc.Key = "withholding-tax-enasarco-contribution" // TR04
	FPAWithholdingENPAMContribution    cbc.Key = "withholding-tax-enpam-contribution"    // TR05
	FPAWithholdingOtherSocialSecurity  cbc.Key = "withholding-tax-other-social-security" // TR06
)

// FPACode defines alphanumeric codes used by FatturaPA, Italy's e-invoicing
// system
var FPACodes = []*KeyDefinition{
	// Tax System Codes
	{
		Key:  FPACodeTaxSystemOrdinary,
		Code: "RF01",
		Desc: i18n.String{
			i18n.EN: "Ordinary tax system",
			i18n.IT: "Regime ordinario",
		},
	},
	// Payment Method Codes
	{
		Key:  FPACodePaymentMethodCash,
		Code: "MP01",
		Desc: i18n.String{
			i18n.EN: "Cash",
			i18n.IT: "Contanti",
		},
	},
	{
		Key:  FPACodePaymentMethodBankTransfer,
		Code: "MP05",
		Desc: i18n.String{
			i18n.EN: "Bank transfer",
			i18n.IT: "Bonifico bancario",
		},
	},
	{
		Key:  FPACodePaymentMethodCard,
		Code: "MP08",
		Desc: i18n.String{
			i18n.EN: "Card",
			i18n.IT: "Carta di credito",
		},
	},
	{
		Key:  FPACodePaymentDirectDebit,
		Code: "MP10",
		Desc: i18n.String{
			i18n.EN: "Direct debit",
			i18n.IT: "Rid",
		},
	},
	// Document Type Codes
	{
		Key:  FPACodeDocumentTypeInvoice,
		Code: "TD01",
		Desc: i18n.String{
			i18n.EN: "Invoice",
			i18n.IT: "Fattura",
		},
	},
	{
		Key:  FPACodeDocumentTypeCreditNote,
		Code: "TD04",
		Desc: i18n.String{
			i18n.EN: "Credit cote",
			i18n.IT: "Nota di credito",
		},
	},
	// Withholding Tax Codes
	{
		Key:  FPAWithholdingTaxNaturalPersons,
		Code: "TR01",
		Desc: i18n.String{
			i18n.EN: "Withholding tax natural persons",
			i18n.IT: "Ritenuta persone fisiche",
		},
	},
	{
		Key:  FPAWithholdingTaxLegalPersons,
		Code: "TR02",
		Desc: i18n.String{
			i18n.EN: "Withholding tax legal persons",
			i18n.IT: "Ritenuta persone giuridiche",
		},
	},
	{
		Key:  FPAWithholdingINPSContribution,
		Code: "TR03",
		Desc: i18n.String{
			i18n.EN: "INPS contribution",
			i18n.IT: "Contributo INPS",
		},
	},
	{
		Key:  FPAWithholdingENASARCOContribution,
		Code: "TR04",
		Desc: i18n.String{
			i18n.EN: "ENASARCO contribution",
			i18n.IT: "Contributo ENASARCO",
		},
	},
	{
		Key:  FPAWithholdingENPAMContribution,
		Code: "TR05",
		Desc: i18n.String{
			i18n.EN: "ENPAM contribution",
			i18n.IT: "Contributo ENPAM",
		},
	},
	{
		Key:  FPAWithholdingOtherSocialSecurity,
		Code: "TR06",
		Desc: i18n.String{
			i18n.EN: "Other social security contribution",
			i18n.IT: "Altro contributo previdenziale",
		},
	},
}
