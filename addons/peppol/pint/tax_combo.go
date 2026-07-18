package pint

import (
	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// normalizeTaxCombo re-maps the UNTDID tax category code from the tax combo's
// key, reusing EN 16931's key-to-code mapping.
func normalizeTaxCombo(tc *tax.Combo) {
	if tc == nil || tc.Key.IsEmpty() {
		return
	}
	if code := en16931.TaxCategoryMap.Get(tc.Key); code != cbc.CodeEmpty {
		tc.Ext = tc.Ext.Set(untdid.ExtKeyTaxCategory, code)
	}
}

func taxComboRules() *rules.Set {
	return rules.For(new(tax.Combo),
		// PINT recognises GST (and other non-VAT schemes) as standard, zero or
		// exempt rated, so the EN 16931 rule that forces every non-VAT category
		// to outside scope must not apply here.
		rules.Ignore("GOBL-EU-EN16931-TAX-COMBO-05"),

		rules.Field("ext",
			rules.Assert("01", "invoice line tax category code must be a recognised PINT UNTDID 5305 code",
				tax.ExtensionsHasCodes(untdid.ExtKeyTaxCategory, en16931.TaxCategoryMap.Values()...),
			),
		),
	)
}
