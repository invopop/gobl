package xrechnung

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func isValidInvoiceType(t cbc.Key) bool {
	validTypes := []cbc.Key{
		bill.InvoiceTypeStandard,
		bill.InvoiceTypeCreditNote,
		bill.InvoiceTypeCorrective,
		invoiceTypeSelfBilled,
		invoiceTypePartial,
	}
	for _, validType := range validTypes {
		if t == validType {
			return true
		}
	}
	return false
}

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Type,
			validation.In(bill.InvoiceTypeStandard, bill.InvoiceTypeCreditNote, bill.InvoiceTypeCorrective, invoiceTypeSelfBilled, invoiceTypePartial),
		),
		// BR-DE-01
		validation.Field(&inv.Payment.Instructions,
			validation.Required,
		),
		validation.Field(&inv.Supplier,
			validation.By(validateSupplier),
		),
		validation.Field(&inv.Customer,
			validation.By(validateCustomer),
		),
		validation.Field(&inv.Delivery,
			validation.When(inv.Delivery != nil,
				validation.By(validateDeliveryParty),
			),
		),
	)
}

func validateSupplier(value interface{}) error {
	party, _ := value.(*org.Party)
	if party == nil {
		return nil
	}
	return validation.ValidateStruct(party,
		// BR-DE-05
		validation.Field(&party.Name,
			validation.Required,
		),
		// BR-DE-02
		validation.Field(&party.Addresses,
			validation.Required,
			validation.Each(validation.By(validatePartyAddress)),
		),
		// BR-DE-06
		validation.Field(&party.People,
			validation.Required,
		),
		validation.Field(&party.Telephones,
			validation.Required,
		),
		validation.Field(&party.Emails,
			validation.Required,
		),
	)
}

func validateCustomer(value interface{}) error {
	party, _ := value.(*org.Party)
	if party == nil {
		return nil
	}
	return validation.ValidateStruct(party,
		validation.Field(&party.Addresses,
			validation.Required,
			validation.Each(validation.By(validatePartyAddress)),
		),
	)
}

func validatePartyAddress(value interface{}) error {
	addr, _ := value.(*org.Address)
	if addr == nil {
		return nil
	}
	return validation.ValidateStruct(addr,
		validation.Field(&addr.Locality,
			validation.Required,
		),
		validation.Field(&addr.Code,
			validation.Required,
		),
	)
}

func validateDeliveryParty(value interface{}) error {
	party, _ := value.(*org.Party)
	if party == nil {
		return nil
	}
	return validation.ValidateStruct(party,
		validation.Field(&party.Addresses,
			validation.Required,
			validation.Length(1, 1),
			validation.Each(validation.By(validateDeliveryAddress)),
		),
	)
}

func validateDeliveryAddress(value interface{}) error {
	addr, _ := value.(*org.Address)
	if addr == nil {
		return nil
	}
	return validation.ValidateStruct(addr,
		// BR-DE-10
		validation.Field(&addr.Locality,
			validation.Required,
		),
		// BR-DE-11
		validation.Field(&addr.Code,
			validation.Required,
		),
	)
}

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

func validateVATRate(value interface{}) error {
	rate, _ := value.(cbc.Key)
	if rate == "" {
		return validation.NewError("required", "VAT category rate is required")
	}
	return nil
}

func validatePaymentMeans(inv *bill.Invoice) error {
	if inv.Payment == nil || inv.Payment.Instructions == nil {
		return nil
	}

	instr := inv.Payment.Instructions
	return validation.ValidateStruct(instr,
		validation.Field(&instr.Key, validation.Required),
		validation.Field(&instr.CreditTransfer,
			validation.When(instr.Key == pay.MeansKeyCreditTransfer,
				validation.Required.Error("Credit transfer details are required when payment means is credit transfer"),
				validation.Length(1, 1).Error("Exactly one credit transfer detail must be provided"),
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
		validation.Field(&instr.Online,
			validation.When(instr.Key != pay.MeansKeyCreditTransfer && instr.Key != pay.MeansKeyCard && instr.Key != pay.MeansKeyDirectDebit,
				validation.Empty.Error("Online payment details should not be present for this payment means"),
			),
		),
	)
}

func validateCorrectiveInvoice(inv *bill.Invoice) error {
	if inv.Type.In(bill.InvoiceTypeCorrective) {
		if inv.Preceding == nil {
			return validation.NewError("required", "Preceding invoice details are required for corrective invoices")
		}
	}
	return nil
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
