package xrechnung

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateTaxCombo(tc *tax.Combo) error {
	if tc == nil {
		return nil
	}
	return validation.ValidateStruct(tc,
		validation.Field(&tc.Category,
			validation.When(tc.Category == tax.CategoryVAT,
				validation.By(validateVATRate),
			),
		),
	)
}

// BR-DE-14
func validateVATRate(value interface{}) error {
	rate, _ := value.(cbc.Key)
	if rate == "" {
		return validation.NewError("required", "VAT category rate is required")
	}
	return nil
}
