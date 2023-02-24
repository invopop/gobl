package it

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/common"
)

// Fattura PA codes are used by the Italian tax system to classify various
// aspects of the invoice, namely the type of document, the payment method,
// the nature of the transaction, the type of fund (if applicable), and the
// type of withholding taxes (if applicable).
const (
	FPACodeTypeTaxSystem         cbc.Key = "fpa-tax-system"         // RegimeFiscale
	FPACodeTypeFundType          cbc.Key = "fpa-fund-type"          // TipoCassa
	FPACodeTypePaymentMethod     cbc.Key = "fpa-payment-method"     // ModalitaPagamento
	FPACodeTypeDocumentType      cbc.Key = "fpa-document-type"      // TipoDocumento
	FPACodeTypeTransactionNature cbc.Key = "fpa-transaction-nature" // Natura
	FPACodeTypeWithholdingType   cbc.Key = "fpa-withholding-type"   // TipoRitenuta
)

// FPACodeGroup is a group of FPA codes under the same type that are variations
// of a related concept.
type FPACodeGroup struct {
	Type  cbc.Key
	Codes []FPACode
}

// InvoiceTypeMap maps invoice types to FPA codes.
var InvoiceTypeMap = map[bill.InvoiceType]*FPACodeGroup{
	bill.InvoiceTypeNone: {
		Type: FPACodeTypeDocumentType,
		Codes: []FPACode{
			FPACodeDocumentTypeInvoice,
		},
	},
	bill.InvoiceTypeCreditNote: {
		Type: FPACodeTypeDocumentType,
		Codes: []FPACode{
			FPACodeDocumentTypeCreditNote,
		},
	},
}

// PaymentMethodMap maps the invoice's payment instruction keys to FPA codes.
var PaymentMethodMap = map[pay.MethodKey]*FPACodeGroup{
	pay.MethodKeyCash: {
		Type: FPACodeTypePaymentMethod,
		Codes: []FPACode{
			FPACodePaymentMethodCash,
		},
	},
	pay.MethodKeyCard: {
		Type: FPACodeTypePaymentMethod,
		Codes: []FPACode{
			FPACodePaymentMethodCard,
		},
	},
	pay.MethodKeyDirectDebit: {
		Type: FPACodeTypePaymentMethod,
		Codes: []FPACode{
			FPACodePaymentDirectDebit,
		},
	},
}

// SchemeMap maps the invoice's scheme keys to FPA codes. There is a limitation
// here, in that the FPA codes of "Nature" do not necessarily map onto the
// general concept of "Scheme" and vice versa.
var SchemeMap = map[cbc.Key]*FPACodeGroup{
	common.SchemeReverseCharge: {
		Type: FPACodeTypeTransactionNature,
		Codes: []FPACode{
			FPACodeNatureRCScrapMaterials,
			FPACodeNatureRCGoldSilver,
			FPACodeNatureRCConstructionSubcontracting,
			FPACodeNatureRCBuildings,
			FPACodeNatureRCMobile,
			FPACodeNatureRCElectronics,
			FPACodeNatureRCConstructionProvisions,
			FPACodeNatureRCEnergy,
			FPACodeNatureRCOther,
		},
	},
}

// TaxCategoryMap maps the invoice's tax category keys to FPA codes related to
// withholding taxes.
var TaxCategoryMap = map[cbc.Code]*FPACodeGroup{
	TaxCategoryRA: {
		Type: FPACodeTypeWithholdingType,
		Codes: []FPACode{
			FPACodeWithholdingNaturalPersons,
			FPACodeWithholdingLegalPersons,
			FPACodeWithholdingINPSContribution,
			FPACodeWithholdingENASARCOContribution,
			FPACodeWithholdingENPAMContribution,
			FPACodeWithholdingOtherSocialSecurity,
		},
	},
}
