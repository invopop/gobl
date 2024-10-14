package xrechnung

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/validation"
)

func validatePaymentInstructions(value interface{}) error {
	inv, ok := value.(*bill.Invoice)
	if !ok || inv == nil {
		return nil
	}
	if inv.Payment == nil || inv.Payment.Instructions == nil {
		return nil
	}

	instr := inv.Payment.Instructions
	return validation.ValidateStruct(instr,
		validation.Field(&instr.Key, validation.Required),
		validation.Field(&instr.CreditTransfer,
			validation.When(instr.Key == pay.MeansKeyCreditTransfer,
				validation.Required.Error("Credit transfer details are required when payment means is credit transfer"),
				validation.Each(validation.By(func(ct interface{}) error {
					creditTransfer, _ := ct.(*pay.CreditTransfer)
					if creditTransfer.IBAN == "" && creditTransfer.Number == "" {
						return validation.NewError("required", "Either IBAN or account number must be provided for credit transfer")
					}
					return nil
				})),
			).Else(validation.Empty),
		),
		validation.Field(&instr.Card,
			validation.When(instr.Key == pay.MeansKeyCard,
				validation.Required.Error("Card details are required when payment means is card"),
				validation.By(func(card interface{}) error {
					c, _ := card.(*pay.Card)
					if c == nil || (c.Last4 == "" && c.Holder == "") {
						return validation.NewError("required", "Card details must include either last 4 digits or holder name")
					}
					return nil
				}),
			).Else(validation.Nil),
		),
		validation.Field(&instr.DirectDebit,
			validation.When(instr.Key == pay.MeansKeyDirectDebit,
				validation.Required.Error("Direct debit details are required when payment means is direct debit"),
				validation.By(validateDirectDebit),
			).Else(validation.Nil),
		),
	)
}
