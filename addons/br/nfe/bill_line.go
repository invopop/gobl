package nfe

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/regimes/br"
	"github.com/invopop/validation"
)

func validateLine(line *bill.Line) error {
	if line == nil {
		return nil
	}
	return validation.Validate(line,
		bill.RequireLineTaxCategory(br.TaxCategoryICMS),
		bill.RequireLineTaxCategory(br.TaxCategoryPIS),
		bill.RequireLineTaxCategory(br.TaxCategoryCOFINS),
		validation.Skip,
	)
}
