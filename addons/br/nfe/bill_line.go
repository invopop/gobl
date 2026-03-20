package nfe

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/br"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

func billLineRules() *rules.Set {
	return rules.For(new(bill.Line),
		rules.Assert("01", "line taxes must include the ICMS category",
			is.Func("has ICMS", lineHasTaxCategory(br.TaxCategoryICMS)),
		),
		rules.Assert("02", "line taxes must include the PIS category",
			is.Func("has PIS", lineHasTaxCategory(br.TaxCategoryPIS)),
		),
		rules.Assert("03", "line taxes must include the COFINS category",
			is.Func("has COFINS", lineHasTaxCategory(br.TaxCategoryCOFINS)),
		),
	)
}

// lineHasTaxCategory returns a function that checks whether a bill.Line
// contains a tax combo with the given category code.
func lineHasTaxCategory(cat cbc.Code) func(any) bool {
	return func(val any) bool {
		line, ok := val.(*bill.Line)
		if !ok || line == nil {
			return true
		}
		for _, tc := range line.Taxes {
			if tc != nil && tc.Category == cat {
				return true
			}
		}
		return false
	}
}
