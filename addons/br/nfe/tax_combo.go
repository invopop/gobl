package nfe

import (
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/br"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// normalizeTaxCombo sets default situation codes on ICMS, PIS and COFINS combos.
func normalizeTaxCombo(tc *tax.Combo) {
	if tc == nil {
		return
	}
	switch tc.Category {
	case br.TaxCategoryICMS:
		if !tc.Ext.Has(ExtKeyICMSCSOSN) {
			tc.Ext = tc.Ext.SetIfEmpty(ExtKeyICMSCST, "00") // Taxed in full
		}
		tc.Ext = tc.Ext.SetIfEmpty(ExtKeyICMSOrigin, "0") // National
	case br.TaxCategoryPIS:
		tc.Ext = tc.Ext.SetIfEmpty(ExtKeyPISCST, "01") // Standard-rate taxable operation
	case br.TaxCategoryCOFINS:
		tc.Ext = tc.Ext.SetIfEmpty(ExtKeyCOFINSCST, "01") // Standard-rate taxable operation
	}
}

func taxComboRules() *rules.Set {
	return rules.For(new(tax.Combo),
		rules.When(
			is.Func("ICMS category", taxCategoryIs(br.TaxCategoryICMS)),
			rules.Field("ext",
				rules.Assert("01", fmt.Sprintf("ICMS tax combo requires '%s' or '%s' extension", ExtKeyICMSCST, ExtKeyICMSCSOSN),
					is.AnyOf(
						tax.ExtensionsRequire(ExtKeyICMSCST),
						tax.ExtensionsRequire(ExtKeyICMSCSOSN),
					),
				),
				rules.Assert("02", fmt.Sprintf("ICMS tax combo requires '%s' extension", ExtKeyICMSOrigin),
					tax.ExtensionsRequire(ExtKeyICMSOrigin),
				),
			),
		),
		rules.When(
			is.Func("PIS category", taxCategoryIs(br.TaxCategoryPIS)),
			rules.Field("ext",
				rules.Assert("03", fmt.Sprintf("PIS tax combo requires '%s' extension", ExtKeyPISCST),
					tax.ExtensionsRequire(ExtKeyPISCST),
				),
			),
		),
		rules.When(
			is.Func("COFINS category", taxCategoryIs(br.TaxCategoryCOFINS)),
			rules.Field("ext",
				rules.Assert("04", fmt.Sprintf("COFINS tax combo requires '%s' extension", ExtKeyCOFINSCST),
					tax.ExtensionsRequire(ExtKeyCOFINSCST),
				),
			),
		),
	)
}

// taxCategoryIs returns a tester that matches a tax combo of the given category.
func taxCategoryIs(cat cbc.Code) func(any) bool {
	return func(val any) bool {
		tc, ok := val.(*tax.Combo)
		return ok && tc != nil && tc.Category == cat
	}
}
