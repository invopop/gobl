package gr

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateTaxCombo(tc *tax.Combo) error {
	return validation.ValidateStruct(tc,
		validation.Field(&tc.Ext,
			tax.ExtensionsRequires(ExtKeyIAPRVATCat),
			validation.When(
				tc.Percent == nil,
				tax.ExtensionsRequires(ExtKeyIAPRExemption),
			),
			validation.Skip,
		),
	)
}

func normalizeTaxCombo(tc *tax.Combo) error {
	if tc == nil || tc.Rate == cbc.KeyEmpty {
		return nil
	}

	reg := tax.RegimeFor(l10n.EL)
	rate := reg.Rate(tax.CategoryVAT, tc.Rate)
	if rate == nil {
		return nil
	}

	if tc.Ext == nil {
		tc.Ext = make(tax.Extensions)
	}

	if tc.Ext.Has(ExtKeyIAPRVATCat) {
		return nil
	}

	tc.Ext[ExtKeyIAPRVATCat] = tax.ExtValue(rate.Map[ExtKeyIAPRVATCat])

	return nil
}
