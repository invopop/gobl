package saft

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// TaxRateExtensions returns the mapping of tax rates defined in PT
// to their extension values used by SAF-T.
//
// Use this to lookup a tax rate key for a SAF-T tax rate code:
//
//	saft.TaxRateExtensions().Lookup("RED") // returns tax.RateReduced
func TaxRateExtensions() tax.Extensions {
	return taxRateMap
}

var taxRateMap = tax.Extensions{
	tax.RateReduced:      TaxRateReduced,
	tax.RateIntermediate: TaxRateIntermediate,
	tax.RateStandard:     TaxRateNormal,
	tax.RateExempt:       TaxRateExempt,
	tax.RateOther:        TaxRateOther,
}

func normalizeTaxCombo(combo *tax.Combo) {
	if combo == nil {
		return
	}

	// copy the SAF-T tax rate code to the line
	switch combo.Category {
	case tax.CategoryVAT:
		if combo.Ext == nil {
			combo.Ext = make(tax.Extensions)
		}
		if combo.Country != "" && combo.Country != l10n.PT.Tax() {
			combo.Ext[ExtKeyTaxRate] = TaxRateOther
			return
		}
		if combo.Rate.IsEmpty() {
			return
		}
		k, ok := taxRateMap[combo.Rate]
		if !ok {
			return
		}
		combo.Ext[ExtKeyTaxRate] = k
	}
}

func validateTaxCombo(val any) error {
	c, ok := val.(*tax.Combo)
	if !ok {
		return nil
	}
	switch c.Category {
	case tax.CategoryVAT:
		return validation.ValidateStruct(c, validateVATExt(&c.Ext))
	}
	return nil
}

func validateVATExt(ext *tax.Extensions) *validation.FieldRules {
	return validation.Field(ext,
		tax.ExtensionsRequire(pt.ExtKeyRegion, ExtKeyTaxRate),
		validation.When(
			(*ext)[ExtKeyTaxRate] == TaxRateExempt,
			tax.ExtensionsRequire(ExtKeyExemption),
		),
		validation.Skip,
	)
}
