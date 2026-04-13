package nfse

import (
	"fmt"

	"github.com/invopop/gobl/regimes/br"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// ISSLiabilityDefault is the default value for the ISS liability extension
	ISSLiabilityDefault = "1" // Liable
)

func taxComboRules() *rules.Set {
	return rules.For(new(tax.Combo),
		rules.When(
			is.Func("ISS category", taxComboIsISS),
			rules.Field("ext",
				rules.Assert("01", fmt.Sprintf("ISS tax combo requires '%s' extension", ExtKeyISSLiability),
					tax.ExtensionsRequire(ExtKeyISSLiability),
				),
			),
		),
	)
}

func taxComboIsISS(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && tc.Category == br.TaxCategoryISS
}

func normalizeTaxCombo(tc *tax.Combo) {
	if tc == nil || tc.Category != br.TaxCategoryISS {
		return
	}

	if !tc.Ext.Has(ExtKeyISSLiability) {
		if tc.Ext == nil {
			tc.Ext = make(tax.Extensions)
		}
		tc.Ext[ExtKeyISSLiability] = ISSLiabilityDefault
	}
}
