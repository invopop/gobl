package arca

import (
	"fmt"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func taxComboRules() *rules.Set {
	return rules.For(new(tax.Combo),
		rules.When(is.Func("category is VAT", taxComboIsVAT),
			rules.Field("ext",
				rules.Assert("01",
					fmt.Sprintf("VAT combo requires '%s' extension", ExtKeyVATRate),
					tax.ExtensionsRequire(ExtKeyVATRate),
				),
			),
		),
	)
}

func taxComboIsVAT(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && tc.Category == tax.CategoryVAT
}
