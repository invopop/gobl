package xrechnung

import (
	"github.com/invopop/gobl/pay"
	"github.com/invopop/validation"
)

// ValidatePaymentInstructions validates the payment instructions according to the XRechnung standard
func validatePaymentInstructions(value interface{}) error {
	instr, ok := value.(*pay.Instructions)
	if !ok || instr == nil {
		return nil
	}
	return validation.ValidateStruct(instr,
		// BR-DE-23
		validation.Field(&instr.CreditTransfer,
			validation.When(
				instr.Key.Has(pay.MeansKeyCreditTransfer),
				validation.Required,
				validation.Each(validation.By(validateCreditTransfer)),
			),
			validation.Skip,
		),
		// BR-DE-24
		validation.Field(&instr.Card,
			validation.When(
				instr.Key.Has(pay.MeansKeyCard),
				validation.Required,
			),
			validation.Skip,
		),
		// BR-DE-25
		validation.Field(&instr.DirectDebit,
			validation.When(
				instr.Key.Has(pay.MeansKeyDirectDebit),
				validation.Required,
				validation.By(validateInstructionsDirectDebit),
				validation.Skip,
			),
		),
	)
}

func validateInstructionsDirectDebit(value interface{}) error {
	dd, ok := value.(*pay.DirectDebit)
	if !ok || dd == nil {
		return nil
	}
	return validation.ValidateStruct(dd,
		// BR-DE-29 - Changed to Peppol-EN16931-R061
		validation.Field(&dd.Ref,
			validation.Required,
		),
		// BR-DE-30
		validation.Field(&dd.Creditor,
			validation.Required,
		),
		// BR-DE-31
		validation.Field(&dd.Account,
			validation.Required,
		),
	)
}

// BR-DE-19
func validateCreditTransfer(value interface{}) error {
	ct, ok := value.(*pay.CreditTransfer)
	if ct == nil || !ok {
		return nil
	}
	return validation.ValidateStruct(ct,
		validation.Field(&ct.Number,
			validation.When(
				ct.IBAN == "",
				validation.Required,
			),
		),
	)
}
