package bill

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

type lineValidation struct {
	taxKey cbc.Code
}

// RequireLineTaxCategory is a validation rule that ensures the line has the required
// tax category.
func RequireLineTaxCategory(key cbc.Code) validation.Rule {
	return &lineValidation{taxKey: key}
}

func (v *lineValidation) Validate(value any) error {
	line, ok := value.(*Line)
	if !ok {
		return nil
	}
	if v.taxKey == cbc.CodeEmpty {
		return nil
	}
	return validation.ValidateStruct(line,
		validation.Field(&line.Taxes,
			tax.SetHasCategory(v.taxKey),
			validation.Required,
			validation.Skip,
		),
	)
}
