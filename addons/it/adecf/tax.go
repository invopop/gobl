package adecf

import (
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateTaxCombo(val any) error {
	c, ok := val.(*tax.Combo)
	if !ok {
		return nil
	}

	if c.Category == tax.CategoryVAT {
		return validation.ValidateStruct(c,
			validation.Field(&c.Ext,
				validation.When(
					c.Percent == nil,
					tax.ExtensionsRequire(ExtKeyExempt),
				),
				validation.Skip,
			),
			validation.Field(&c.Percent,
				validation.By(
					validatePercentage,
				),
				validation.Skip,
			),
		)
	}
	return nil
}

func validatePercentage(val any) error {
	p, ok := val.(*num.Percentage)
	if !ok || p == nil {
		return nil
	}

	for _, vp := range validPercentages() {
		if p.Compare(vp) == 0 {
			return nil
		}
	}
	return validation.NewError("validation_percentage_percentage", "Invalid percentage")
}

// Percentages required by ade cf
func validPercentages() []num.Percentage {
	return []num.Percentage{
		num.MakePercentage(4, 2),
		num.MakePercentage(5, 2),
		num.MakePercentage(10, 2),
		num.MakePercentage(22, 2),
		num.MakePercentage(2, 2),
		num.MakePercentage(64, 3),
		num.MakePercentage(7, 2),
		num.MakePercentage(73, 3),
		num.MakePercentage(75, 3),
		num.MakePercentage(765, 4),
		num.MakePercentage(795, 4),
		num.MakePercentage(83, 3),
		num.MakePercentage(85, 3),
		num.MakePercentage(88, 3),
		num.MakePercentage(95, 3),
		num.MakePercentage(123, 3),
	}
}
