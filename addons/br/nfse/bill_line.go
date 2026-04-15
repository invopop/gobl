package nfse

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/regimes/br"
	"github.com/invopop/gobl/rules"
)

func billLineRules() *rules.Set {
	return rules.For(new(bill.Line),
		rules.Assert("01", "line taxes must include the ISS category",
			bill.RequireLineTaxCategory(br.TaxCategoryISS),
		),
	)
}
