package it

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/common"
)

const (
	FPACodeTypeTaxSystem       cbc.Key = "fpa-tax-system"
	FPACodeTypeFundType        cbc.Key = "fpa-fund-type"
	FPACodeTypePaymentMethod   cbc.Key = "fpa-payment-method"
	FPACodeTypeDocumentType    cbc.Key = "fpa-document-type"
	FPACodeTypeNature          cbc.Key = "fpa-nature"
	FPACodeTypeWithholdingType cbc.Key = "fpa-withholding-type"
)

type FPACodeOptions struct {
	FPACodeType cbc.Key
	FPACodes    []string
}

var InvoiceTypeMap = map[bill.InvoiceType]*FPACodeOptions{
	bill.InvoiceTypeNone: {
		FPACodeType: FPACodeTypeDocumentType,
		FPACodes:    []string{"TD01"},
	},
	bill.InvoiceTypeCreditNote: {
		FPACodeType: FPACodeTypeDocumentType,
		FPACodes:    []string{"TD04"},
	},
}

var PaymentMethodMap = map[pay.MethodKey]*FPACodeOptions{
	pay.MethodKeyCash: {
		FPACodeType: FPACodeTypePaymentMethod,
		FPACodes:    []string{"MP01"},
	},
	pay.MethodKeyCard: {
		FPACodeType: FPACodeTypePaymentMethod,
		FPACodes:    []string{"MP08"},
	},
	pay.MethodKeyDirectDebit: {
		FPACodeType: FPACodeTypePaymentMethod,
		FPACodes:    []string{"MP10"},
	},
}

var SchemeMap = map[cbc.Key]*FPACodeOptions{
	common.SchemeReverseCharge: {
		FPACodeType: FPACodeTypeNature,
		FPACodes:    []string{"N6.1", "N6.2", "N6.3", "N6.4", "N6.5", "N6.6", "N6.7", "N6.8", "N6.9"},
	},
}

var TaxCategoryMap = map[cbc.Code]*FPACodeOptions{
	TaxCategoryRA: {
		FPACodeType: FPACodeTypeWithholdingType,
		FPACodes:    []string{"RT01", "RT02", "RT03", "RT04", "RT05", "RT06"},
	},
}
