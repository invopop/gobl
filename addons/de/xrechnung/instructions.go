package xrechnung

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/validation"
)

// Payment keys for XRechnung SEPA direct debit and credit transfer
const (
	KeyPaymentMeansSEPACreditTransfer cbc.Key = "sepa-credit-transfer"
	KeyPaymentMeansSEPADirectDebit    cbc.Key = "sepa-direct-debit"
)

var validPaymentKeys = []cbc.Key{
	pay.MeansKeyCash,
	pay.MeansKeyCheque,
	pay.MeansKeyCreditTransfer,
	pay.MeansKeyCard,
	pay.MeansKeyDirectDebit,
	pay.MeansKeyOther,
	KeyPaymentMeansSEPACreditTransfer,
	KeyPaymentMeansSEPADirectDebit,
}

// ValidatePaymentInstructions validates the payment instructions according to the XRechnung standard
func validatePaymentInstructions(value interface{}) error {
	instr, ok := value.(*pay.Instructions)
	if !ok || instr == nil {
		return nil
	}
	return validation.ValidateStruct(instr,
		validation.Field(&instr.Key,
			validation.Required,
			validation.By(validatePaymentKey),
			validation.Skip,
		),
		// BR-DE-23
		validation.Field(&instr.CreditTransfer,
			validation.When(instr.Key == KeyPaymentMeansSEPACreditTransfer,
				validation.Required,
				validation.Each(validation.By(validateCreditTransfer)),
			),
			validation.Skip,
		),
		// BR-DE-24
		validation.Field(&instr.Card,
			validation.When(instr.Key == pay.MeansKeyCard,
				validation.Required,
			),
			validation.Skip,
		),
		// BR-DE-25
		validation.Field(&instr.DirectDebit,
			validation.When(instr.Key == KeyPaymentMeansSEPADirectDebit || instr.Key == pay.MeansKeyDirectDebit,
				validation.Required,
				validation.By(validateDirectDebit),
				validation.Skip,
			),
		),
	)
}

func validatePaymentKey(value interface{}) error {
	t, ok := value.(cbc.Key)
	if !ok {
		return validation.NewError("invalid_key", "invalid payment key")
	}
	if !t.In(validPaymentKeys...) {
		return validation.NewError("invalid", "invalid payment key")
	}
	return nil
}

func validateDirectDebit(value interface{}) error {
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
	creditTransfer, _ := value.(*pay.CreditTransfer)
	if creditTransfer == nil {
		return nil
	}
	return validation.ValidateStruct(creditTransfer,
		validation.Field(&creditTransfer.Number,
			validation.When(creditTransfer.IBAN == "",
				validation.Required,
			),
		),
	)
}
