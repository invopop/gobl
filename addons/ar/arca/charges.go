package arca

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateBillCharge(charge *bill.Charge) error {
	return validation.ValidateStruct(charge,
		validation.Field(&charge.Ext,
			validation.When(
				charge.Key.Has(bill.ChargeKeyTax),
				tax.ExtensionsRequire(ExtKeyTaxType),
			),
		),
		validation.Field(&charge.Percent,
			validation.When(
				charge.Ext.Has(ExtKeyTaxType),
				validation.Required,
			),
		),
		validation.Field(&charge.Reason,
			validation.When(
				charge.Ext.Get(ExtKeyTaxType) == "99", // Other
				validation.Required.Error("reason is required when tax type is 'other'"),
			),
		),
	)
}
