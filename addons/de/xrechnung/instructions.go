package xrechnung

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/validation"
)

const (
	keyPaymentMeansSEPACreditTransfer cbc.Key = "sepa-credit-transfer"
	keyPaymentMeansSEPADirectDebit    cbc.Key = "sepa-direct-debit"
)

func validatePaymentInstructions(value interface{}) error {
	inv, ok := value.(*bill.Invoice)
	if !ok || inv == nil || inv.Payment == nil || inv.Payment.Instructions == nil {
		return nil
	}

	instr := inv.Payment.Instructions
	return validation.ValidateStruct(instr,
		validation.Field(&instr.Key, validation.Required),
		// BR-DE-23
		validation.Field(&instr.CreditTransfer,
			validation.When(instr.Key == keyPaymentMeansSEPACreditTransfer,
				validation.Required,
				validation.By(validateCreditTransfer),
			).Else(validation.Nil),
		),
		// BR-DE-24
		validation.Field(&instr.Card,
			validation.When(instr.Key == pay.MeansKeyCard,
				validation.Required,
			).Else(validation.Nil),
		),
		// BR-DE-25
		validation.Field(&instr.DirectDebit,
			validation.When(instr.Key == keyPaymentMeansSEPADirectDebit || instr.Key == pay.MeansKeyDirectDebit,
				validation.Required,
				validation.By(validateDirectDebit),
			).Else(validation.Nil),
		),
	)
}

func validateDirectDebit(value interface{}) error {
	inv, ok := value.(*bill.Invoice)
	if !ok || inv == nil {
		return nil
	}
	if inv.Payment == nil || inv.Payment.Instructions == nil || inv.Payment.Instructions.Key != pay.MeansKeyDirectDebit {
		return nil
	}

	dd := inv.Payment.Instructions.DirectDebit
	return validation.ValidateStruct(dd,
		// BR-DE-29 - Changed to Peppol-EN16931-R061
		validation.Field(&dd.Ref,
			validation.Required.Error("Mandate reference is mandatory for direct debit"),
		),
		// BR-DE-30
		validation.Field(&dd.Creditor,
			validation.Required.Error("Creditor identifier is mandatory for direct debit"),
		),
		// BR-DE-31
		validation.Field(&dd.Account,
			validation.Required.Error("Debited account identifier is mandatory for direct debit"),
		),
	)
}

// BR-DE-19
func validateCreditTransfer(value interface{}) error {
	creditTransfer, _ := value.(*pay.CreditTransfer)
	if creditTransfer == nil {
		return nil
	}
	return validation.ValidateStruct(creditTransfer,
		validation.Field(&creditTransfer.Number,
			validation.When(creditTransfer.IBAN == "",
				validation.Required.Error("IBAN must be provided for SEPA credit transfer"),
			),
		),
	)
}
