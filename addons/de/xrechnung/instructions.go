package xrechnung

import (
	"github.com/invopop/gobl/pay"
	"github.com/invopop/validation"
)

func validatePayInstructions(instructions *pay.Instructions) error {
	return validation.ValidateStruct(instructions,
		validation.Field(&instructions.CreditTransfer,
			validation.When(instructions.Key == pay.MeansKeyCreditTransfer,
				validation.By(validateCreditTransfer),
			),
		),
	)
}

func validateCreditTransfer(value interface{}) error {
	credit, _ := value.(*pay.CreditTransfer)
	if credit == nil {
		return nil
	}
	return nil
	// return validation.ValidateStruct(credit,
	// 	validation.Field(&credit.IBAN,
	// 		validation.When(credit.Key == pay.MeansKeyCreditTransfer,
	// 			validation.Required,
	// 		),
	// 	),
	// )
}
