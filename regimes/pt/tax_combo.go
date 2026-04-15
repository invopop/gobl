package pt

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func taxComboRules() *rules.Set {
	return rules.For(new(tax.Combo),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				is.Expr(`string(Category) == "VAT"`),
				rules.Field("ext",
					rules.Assert("01", "pt-region extension is required for VAT", tax.ExtensionsRequire(ExtKeyRegion)),
				),
			),
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
