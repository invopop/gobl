package nfe

import (
"github.com/invopop/gobl/bill"
"github.com/invopop/gobl/regimes/br"
"github.com/invopop/gobl/rules"
"github.com/invopop/gobl/tax"
)

func billLineRules() *rules.Set {
return rules.For(new(bill.Line),
rules.Field("taxes",
rules.Assert("01", "ICMS tax category is required",
tax.SetHasCategory(br.TaxCategoryICMS),
),
rules.Assert("02", "PIS tax category is required",
tax.SetHasCategory(br.TaxCategoryPIS),
),
rules.Assert("03", "COFINS tax category is required",
tax.SetHasCategory(br.TaxCategoryCOFINS),
),
),
)
}
