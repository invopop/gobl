package nfse

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/regimes/br"
	"github.com/invopop/validation"
)

func validateLine(value any) error {
	line, _ := value.(*bill.Line)
	if line == nil {
		return nil
	}

	return validation.Validate(line,
		bill.RequireLineTaxCategory(br.TaxCategoryISS),
	)
}
