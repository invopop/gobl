package pt

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateTaxCombo(tc *tax.Combo) error {
	if tc == nil {
		return nil
	}

	return validation.ValidateStruct(tc,
		validation.Field(&tc.Ext,
			validation.When(
				tc.Category == tax.CategoryVAT,
				tax.ExtensionsRequire(ExtKeyRegion),
			),
			validation.Skip,
		),
	)
}

func normalizeTaxCombo(tc *tax.Combo) {
	if tc == nil || tc.Category != tax.CategoryVAT {
		return
	}

	if tc.Ext == nil {
		tc.Ext = make(tax.Extensions)
	}

	// Set default region
	if _, ok := tc.Ext[ExtKeyRegion]; !ok {
		tc.Ext[ExtKeyRegion] = RegionMainland
	}

	// Override region with foreign country if present
	if tc.Country != "" && tc.Country != l10n.PT.Tax() {
		tc.Ext[ExtKeyRegion] = cbc.Code(isoCountry(tc.Country))
	}
}

func isoCountry(c l10n.TaxCountryCode) l10n.ISOCountryCode {
	isoCode := c.Code().ISO()
	if isoCode.Validate() == nil {
		// Code is already an ISO country code.
		return isoCode
	}

	// Code is not an ISO country code. Try with the alternative code.
	def := l10n.Countries().Code(c.Code())
	if def != nil && def.AltCode != l10n.CodeEmpty {
		isoCode := def.AltCode.ISO()
		if isoCode.Validate() == nil {
			return isoCode
		}
	}

	return l10n.CodeEmpty.ISO()
}
