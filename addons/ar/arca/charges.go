package arca

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Charge keys for Argentina ARCA v4
const (
	ChargeKeyNationalTaxes               cbc.Key = "national-taxes"
	ChargeKeyProvincialTaxes             cbc.Key = "provincial-taxes"
	ChargeKeyMunicipalTaxes              cbc.Key = "municipal-taxes"
	ChargeKeyInternalTaxes               cbc.Key = "internal-taxes"
	ChargeKeyGrossIncomeTax              cbc.Key = "gross-income-tax"
	ChargeKeyVATPrepayment               cbc.Key = "vat-prepayment"
	ChargeKeyGrossIncomeTaxPrepayment    cbc.Key = "gross-income-tax-prepayment"
	ChargeKeyMunicipalTaxesPrepayment    cbc.Key = "municipal-taxes-prepayment"
	ChargeKeyOtherPrepayments            cbc.Key = "other-prepayments"
	ChargeKeyVATNotCategorizedPrepayment cbc.Key = "vat-not-categorized-prepayment"
	ChargeKeyOther                       cbc.Key = "other"
)

var chargeKeyMap = tax.Extensions{
	ChargeKeyNationalTaxes:               TributeTypeNationalTaxes,
	ChargeKeyProvincialTaxes:             TributeTypeProvincialTaxes,
	ChargeKeyMunicipalTaxes:              TributeTypeMunicipalTaxes,
	ChargeKeyInternalTaxes:               TributeTypeInternalTaxes,
	ChargeKeyGrossIncomeTax:              TributeTypeGrossIncomeTax,
	ChargeKeyVATPrepayment:               TributeTypeVATPrepayment,
	ChargeKeyGrossIncomeTaxPrepayment:    TributeTypeGrossIncomeTaxPrepayment,
	ChargeKeyMunicipalTaxesPrepayment:    TributeTypeMunicipalTaxesPrepayment,
	ChargeKeyOtherPrepayments:            TributeTypeOtherPrepayments,
	ChargeKeyVATNotCategorizedPrepayment: TributeTypeVATNotCategorizedPrepayment,
	ChargeKeyOther:                       TributeTypeOther,
}

func normalizeCharge(charge *bill.Charge) {
	if val, ok := chargeKeyMap[charge.Key]; ok {
		charge.Ext = charge.Ext.Merge(tax.Extensions{
			ExtKeyTributeType: val,
		})
	}
}

func validateCharge(charge *bill.Charge) error {
	return validation.ValidateStruct(charge,
		validation.Field(&charge.Percent,
			validation.When(
				charge.Ext.Has(ExtKeyTributeType),
				validation.Required,
			),
		),
		validation.Field(&charge.Reason,
			validation.When(
				charge.Ext.Get(ExtKeyTributeType) == TributeTypeOther,
				validation.Required.Error("reason is required when tribute type is 'other'"),
			),
		),
	)
}
