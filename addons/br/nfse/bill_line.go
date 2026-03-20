package nfse

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/br"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

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

func billLineRules() *rules.Set {
	return rules.For(new(bill.Line),
		rules.Assert("01", "line taxes must include the ISS category",
			is.Func("has ISS", lineHasTaxCategory(br.TaxCategoryISS)),
		),
	)
}
