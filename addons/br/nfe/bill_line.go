package nfe

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/regimes/br"
	"github.com/invopop/gobl/rules"
)

func billLineRules() *rules.Set {
	return rules.For(new(bill.Line),
		rules.Assert("01", "line taxes must include the ICMS category",
			bill.RequireLineTaxCategory(br.TaxCategoryICMS),
		),
		rules.Assert("02", "line taxes must include the PIS category",
			bill.RequireLineTaxCategory(br.TaxCategoryPIS),
		),
		rules.Assert("03", "line taxes must include the COFINS category",
			bill.RequireLineTaxCategory(br.TaxCategoryCOFINS),
		),
	)
}
