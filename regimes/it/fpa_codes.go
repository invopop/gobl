package it

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// FPACodeDefinition defines properties of an alphanumeric codes used by
// FatturaPA, Italy's e-invoicing system. The codes are used to classify
// various aspects of an invoice, namely the tax system, fund type, payment
// method, document type, nature, and withholding type.
type FPACodeDefinition struct {
	// Actual key value.
	Key cbc.Key `json:"key" jsonschema:"title=Key"`
	//
	Code string `json:"code,omitempty" jsonschema:"title=Code"`
	// Description offering more details about when the key should be used.
	Desc i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`
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

	// Nature (Natura) Codes
	// Reverse Charges
	FPACodeNatureRCScrapMaterials             cbc.Key = "nature-rc-scrap-materials"             // N6.1
	FPACodeNatureRCGoldSilver                 cbc.Key = "nature-rc-gold-silver"                 // N6.2
	FPACodeNatureRCConstructionSubcontracting cbc.Key = "nature-rc-construction-subcontracting" // N6.3
	FPACodeNatureRCBuildings                  cbc.Key = "nature-rc-buildings"                   // N6.4
	FPACodeNatureRCMobile                     cbc.Key = "nature-rc-mobile"                      // N6.5
	FPACodeNatureRCElectronics                cbc.Key = "nature-rc-electronics"                 // N6.6
	FPACodeNatureRCConstructionProvisions     cbc.Key = "nature-rc-construction-provisions"     // N6.7
	FPACodeNatureRCEnergy                     cbc.Key = "nature-rc-energy"                      // N6.8
	FPACodeNatureRCOther                      cbc.Key = "nature-rc-other"                       // N6.9

	// Withholding Tax (TipoRitenuta) Codes
	FPACodeWithholdingNaturalPersons       cbc.Key = "withholding-tax-natural-persons"       // TR01
	FPACodeWithholdingLegalPersons         cbc.Key = "withholding-tax-legal-persons"         // TR02
	FPACodeWithholdingINPSContribution     cbc.Key = "withholding-tax-inps-contribution"     // TR03
	FPACodeWithholdingENASARCOContribution cbc.Key = "withholding-tax-enasarco-contribution" // TR04
	FPACodeWithholdingENPAMContribution    cbc.Key = "withholding-tax-enpam-contribution"    // TR05
	FPACodeWithholdingOtherSocialSecurity  cbc.Key = "withholding-tax-other-social-security" // TR06
)

// FPACodeDefs includes all FatturaPA codes currently supported in GOBL
var FPACodeDefs = []*FPACodeDefinition{
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
		Key:  FPACodeWithholdingNaturalPersons,
		Code: "TR01",
		Desc: i18n.String{
			i18n.EN: "Withholding tax natural persons",
			i18n.IT: "Ritenuta persone fisiche",
		},
	},
	{
		Key:  FPACodeWithholdingLegalPersons,
		Code: "TR02",
		Desc: i18n.String{
			i18n.EN: "Withholding tax legal persons",
			i18n.IT: "Ritenuta persone giuridiche",
		},
	},
	{
		Key:  FPACodeWithholdingINPSContribution,
		Code: "TR03",
		Desc: i18n.String{
			i18n.EN: "INPS contribution",
			i18n.IT: "Contributo INPS",
		},
	},
	{
		Key:  FPACodeWithholdingENASARCOContribution,
		Code: "TR04",
		Desc: i18n.String{
			i18n.EN: "ENASARCO contribution",
			i18n.IT: "Contributo ENASARCO",
		},
	},
	{
		Key:  FPACodeWithholdingENPAMContribution,
		Code: "TR05",
		Desc: i18n.String{
			i18n.EN: "ENPAM contribution",
			i18n.IT: "Contributo ENPAM",
		},
	},
	{
		Key:  FPACodeWithholdingOtherSocialSecurity,
		Code: "TR06",
		Desc: i18n.String{
			i18n.EN: "Other social security contribution",
			i18n.IT: "Altro contributo previdenziale",
		},
	},
	// Nature Codes
	{
		Key:  FPACodeNatureRCScrapMaterials,
		Code: "N6.1",
		Desc: i18n.String{
			i18n.EN: "Reverse charge - transfer of scrap and other recyclable materials",
			i18n.IT: "Inversione contabile - cessione di rottami e altri materiali di recupero",
		},
	},
	{
		Key:  FPACodeNatureRCGoldSilver,
		Code: "N6.2",
		Desc: i18n.String{
			i18n.EN: "Reverse charge - trasnfer of gold, pure silver, and jewelery",
			i18n.IT: "Inversione contabile - cessione di oro e argento puro",
		},
	},
	{
		Key:  FPACodeNatureRCConstructionSubcontracting,
		Code: "N6.3",
		Desc: i18n.String{
			i18n.EN: "Reverse charge - subcontracting in the construction sector",
			i18n.IT: "Inversione contabile - subappalto nel settore edile",
		},
	},
	{
		Key:  FPACodeNatureRCBuildings,
		Code: "N6.4",
		Desc: i18n.String{
			i18n.EN: "Reverse charge - transfer of buildings",
			i18n.IT: "Inversione contabile - cessione di fabbricati",
		},
	},
	{
		Key:  FPACodeNatureRCMobile,
		Code: "N6.5",
		Desc: i18n.String{
			i18n.EN: "Reverse charge - transfer of mobile phones",
			i18n.IT: "Inversione contabile - cessione di telefoni cellulari",
		},
	},
	{
		Key:  FPACodeNatureRCElectronics,
		Code: "N6.6",
		Desc: i18n.String{
			i18n.EN: "Reverse charge - transfer of electronic products",
			i18n.IT: "Inversione contabile - cessione di prodotti elettronici",
		},
	},
	{
		Key:  FPACodeNatureRCConstructionProvisions,
		Code: "N6.7",
		Desc: i18n.String{
			i18n.EN: "Reverse charge - provisions in the construction and related sectors",
			i18n.IT: "Inversione contabile - prestazioni comparto edile e settori connessi",
		},
	},
	{
		Key:  FPACodeNatureRCEnergy,
		Code: "N6.8",
		Desc: i18n.String{
			i18n.EN: "Reverse charge - transactions in the energy sector",
			i18n.IT: "Inversione contabile - operazioni settore energetico",
		},
	},
	{
		Key:  FPACodeNatureRCOther,
		Code: "N6.9",
		Desc: i18n.String{
			i18n.EN: "Reverse charge - other cases",
			i18n.IT: "Inversione contabile - altri casi",
		},
	},
}
