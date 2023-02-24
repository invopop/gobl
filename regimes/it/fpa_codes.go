package it

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/i18n"
)

// FPACodeDefinition defines properties of an alphanumeric codes used by
// FatturaPA, Italy's e-invoicing system. The codes are used to classify
// various aspects of an invoice, namely the tax system, fund type, payment
// method, document type, nature, and withholding type. An FPACode can include
// uppercase letters, numbers, and a "." separator for the numeric portion.
type FPACode string

type FPACodeDefinition struct {
	// Alphanumeric code as required by FatturaPA.
	Code FPACode
	// Description offering more details about when the key should be used.
	Desc i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`
}

var (
	codePattern          = `^[A-Z]+[0-9]+(\.[0-9]+)?$`
	codeValidationRegexp = regexp.MustCompile(codePattern)
)

const (
	// Tax System (RegimeFiscale) Codes
	FPACodeTaxSystemOrdinary FPACode = "RF01"

	// Payment Method (ModalitaPagamento) Codes
	FPACodePaymentCash                 FPACode = "MP01"
	FPACodePaymentBankTransfer         FPACode = "MP05"
	FPACodePaymentCard                 FPACode = "MP08"
	FPACodePaymentDirectDebit          FPACode = "MP10"
	FPACodePaymentDirectDebitUtilities FPACode = "MP10"
	FPACodePaymentDirectDebitFast      FPACode = "MP11"
	FPACodePaymentDirectDebitSepa      FPACode = "MP19"
	FPACodePaymentDirectDebitSepaCore  FPACode = "MP20"
	FPACodePaymentDirectDebitSepaB2B   FPACode = "MP21"

	// Document Type (TipoDocumento) Codes
	FPACodeDocumentTypeInvoice    FPACode = "TD01"
	FPACodeDocumentTypeCreditNote FPACode = "TD04"

	// Nature (Natura) Codes
	// Reverse Charges
	FPACodeNatureRCScrapMaterials             FPACode = "N6.1"
	FPACodeNatureRCGoldSilver                 FPACode = "N6.2"
	FPACodeNatureRCConstructionSubcontracting FPACode = "N6.3"
	FPACodeNatureRCBuildings                  FPACode = "N6.4"
	FPACodeNatureRCMobile                     FPACode = "N6.5"
	FPACodeNatureRCElectronics                FPACode = "N6.6"
	FPACodeNatureRCConstructionProvisions     FPACode = "N6.7"
	FPACodeNatureRCEnergy                     FPACode = "N6.8"
	FPACodeNatureRCOther                      FPACode = "N6.9"

	// Withholding Tax (TipoRitenuta) Codes
	FPACodeWithholdingNaturalPersons       FPACode = "TR01"
	FPACodeWithholdingLegalPersons         FPACode = "TR02"
	FPACodeWithholdingINPSContribution     FPACode = "TR03"
	FPACodeWithholdingENASARCOContribution FPACode = "TR04"
	FPACodeWithholdingENPAMContribution    FPACode = "TR05"
	FPACodeWithholdingOtherSocialSecurity  FPACode = "TR06"
)

// FPACodeDefs includes all FatturaPA codes currently supported in GOBL
var FPACodeDefs = []*FPACodeDefinition{
	// Tax System Codes
	{
		Code: FPACodeTaxSystemOrdinary,
		Desc: i18n.String{
			i18n.EN: "Ordinary tax system",
			i18n.IT: "Regime ordinario",
		},
	},
	// Payment Method Codes
	{
		Code: FPACodePaymentCash,
		Desc: i18n.String{
			i18n.EN: "Cash",
			i18n.IT: "Contanti",
		},
	},
	{
		Code: FPACodePaymentBankTransfer,
		Desc: i18n.String{
			i18n.EN: "Bank transfer",
			i18n.IT: "Bonifico bancario",
		},
	},
	{
		Code: FPACodePaymentCard,
		Desc: i18n.String{
			i18n.EN: "Card",
			i18n.IT: "Carta di credito",
		},
	},
	{
		Code: FPACodePaymentDirectDebit,
		Desc: i18n.String{
			i18n.EN: "Direct debit",
			i18n.IT: "RID",
		},
	},
	{
		Code: FPACodePaymentDirectDebitUtilities,
		Desc: i18n.String{
			i18n.EN: "Direct debit utilities",
			i18n.IT: "RID utenze",
		},
	},
	{
		Code: FPACodePaymentDirectDebitFast,
		Desc: i18n.String{
			i18n.EN: "Direct debit fast",
			i18n.IT: "RID veloce",
		},
	},
	{
		Code: FPACodePaymentDirectDebitSepa,
		Desc: i18n.String{
			i18n.EN: "SEPA Direct Debit",
			i18n.IT: "SEPA Direct Debit",
		},
	},
	{
		Code: FPACodePaymentDirectDebitSepaCore,
		Desc: i18n.String{
			i18n.EN: "SEPA Direct Debit CORE",
			i18n.IT: "SEPA Direct Debit CORE",
		},
	},
	{
		Code: FPACodePaymentDirectDebitSepaB2B,
		Desc: i18n.String{
			i18n.EN: "SEPA Direct Debit B2B",
			i18n.IT: "SEPA Direct Debit B2B",
		},
	},
	// Document Type Codes
	{
		Code: FPACodeDocumentTypeInvoice,
		Desc: i18n.String{
			i18n.EN: "Invoice",
			i18n.IT: "Fattura",
		},
	},
	{
		Code: FPACodeDocumentTypeCreditNote,
		Desc: i18n.String{
			i18n.EN: "Credit cote",
			i18n.IT: "Nota di credito",
		},
	},
	// Withholding Tax Codes
	{
		Code: FPACodeWithholdingNaturalPersons,
		Desc: i18n.String{
			i18n.EN: "Withholding tax natural persons",
			i18n.IT: "Ritenuta persone fisiche",
		},
	},
	{
		Code: FPACodeWithholdingLegalPersons,
		Desc: i18n.String{
			i18n.EN: "Withholding tax legal persons",
			i18n.IT: "Ritenuta persone giuridiche",
		},
	},
	{
		Code: FPACodeWithholdingINPSContribution,
		Desc: i18n.String{
			i18n.EN: "INPS contribution",
			i18n.IT: "Contributo INPS",
		},
	},
	{
		Code: FPACodeWithholdingENASARCOContribution,
		Desc: i18n.String{
			i18n.EN: "ENASARCO contribution",
			i18n.IT: "Contributo ENASARCO",
		},
	},
	{
		Code: FPACodeWithholdingENPAMContribution,
		Desc: i18n.String{
			i18n.EN: "ENPAM contribution",
			i18n.IT: "Contributo ENPAM",
		},
	},
	{
		Code: FPACodeWithholdingOtherSocialSecurity,
		Desc: i18n.String{
			i18n.EN: "Other social security contribution",
			i18n.IT: "Altro contributo previdenziale",
		},
	},
	// Nature Codes
	{
		Code: FPACodeNatureRCScrapMaterials,
		Desc: i18n.String{
			i18n.EN: "Reverse charge - transfer of scrap and other recyclable materials",
			i18n.IT: "Inversione contabile - cessione di rottami e altri materiali di recupero",
		},
	},
	{
		Code: FPACodeNatureRCGoldSilver,
		Desc: i18n.String{
			i18n.EN: "Reverse charge - trasnfer of gold, pure silver, and jewelery",
			i18n.IT: "Inversione contabile - cessione di oro e argento puro",
		},
	},
	{
		Code: FPACodeNatureRCConstructionSubcontracting,
		Desc: i18n.String{
			i18n.EN: "Reverse charge - subcontracting in the construction sector",
			i18n.IT: "Inversione contabile - subappalto nel settore edile",
		},
	},
	{
		Code: FPACodeNatureRCBuildings,
		Desc: i18n.String{
			i18n.EN: "Reverse charge - transfer of buildings",
			i18n.IT: "Inversione contabile - cessione di fabbricati",
		},
	},
	{
		Code: FPACodeNatureRCMobile,
		Desc: i18n.String{
			i18n.EN: "Reverse charge - transfer of mobile phones",
			i18n.IT: "Inversione contabile - cessione di telefoni cellulari",
		},
	},
	{
		Code: FPACodeNatureRCElectronics,
		Desc: i18n.String{
			i18n.EN: "Reverse charge - transfer of electronic products",
			i18n.IT: "Inversione contabile - cessione di prodotti elettronici",
		},
	},
	{
		Code: FPACodeNatureRCConstructionProvisions,
		Desc: i18n.String{
			i18n.EN: "Reverse charge - provisions in the construction and related sectors",
			i18n.IT: "Inversione contabile - prestazioni comparto edile e settori connessi",
		},
	},
	{
		Code: FPACodeNatureRCEnergy,
		Desc: i18n.String{
			i18n.EN: "Reverse charge - transactions in the energy sector",
			i18n.IT: "Inversione contabile - operazioni settore energetico",
		},
	},
	{
		Code: FPACodeNatureRCOther,
		Desc: i18n.String{
			i18n.EN: "Reverse charge - other cases",
			i18n.IT: "Inversione contabile - altri casi",
		},
	},
}

// Validate ensures that the code complies with the expected rules.
func (c FPACode) Validate() error {
	return validation.Validate(string(c),
		validation.Match(codeValidationRegexp),
	)
}
