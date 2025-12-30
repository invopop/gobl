package arca

import (
	"github.com/invopop/gobl/regimes/ar"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func normalizeTaxCombo(tc *tax.Combo) {
	if tc == nil {
		return
	}

	if tc.Key == tax.KeyZero {
		tc.Ext = tc.Ext.Set(ExtKeyVATRate, "3")
		return
	}

	switch tc.Rate {
	case tax.RateReduced:
		tc.Ext = tc.Ext.Set(ExtKeyVATRate, "4")
	case tax.RateGeneral:
		tc.Ext = tc.Ext.Set(ExtKeyVATRate, "5")
	case ar.RateIncreased:
		tc.Ext = tc.Ext.Set(ExtKeyVATRate, "6")
	}
}

func validateTaxCombo(tc *tax.Combo) error {
	return validation.ValidateStruct(tc,
		validation.Field(&tc.Ext,
			validation.When(
				tc.Category == tax.CategoryVAT,
				tax.ExtensionsRequire(ExtKeyVATRate),
			),
		),
	)
}
