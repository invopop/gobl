package it

import (
	"github.com/invopop/gobl/pay"
	"github.com/invopop/validation"
)

func validatePayAdvance(a *pay.Advance) error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Code,
			validation.When(
				a.Key == pay.MeansKeyOther,
				validation.Required,
			),
		),
	)
}

func validatePayInstructions(i *pay.Instructions) error {
	return validation.ValidateStruct(i,
		validation.Field(&i.Code,
			validation.When(
				i.Key == pay.MeansKeyOther,
				validation.Required,
			),
		),
	)
}
