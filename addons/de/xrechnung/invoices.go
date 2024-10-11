package xrechnung

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func normalizeInvoice(inv *bill.Invoice) {
	// Ensure payment instructions are present
	if inv.Payment == nil {
		inv.Payment = &bill.Payment{}
	}
	if inv.Payment.Instructions == nil {
		inv.Payment.Instructions = &pay.Instructions{}
	}

	// Ensure invoice type is valid
	if !isValidInvoiceType(inv.Type) {
		inv.Type = bill.InvoiceTypeStandard
	}
}

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
		validation.Field(&inv.Payment.Instructions,
			validation.Required,
		),
		validation.Field(&inv.Supplier,
			validation.By(validateParty),
		),
		validation.Field(&inv.Customer,
			validation.By(validateParty),
		),
		validation.Field(&inv.Delivery,
			validation.When(inv.Delivery != nil,
				validation.By(validateDeliveryParty),
			),
		),
	)
}

func validateParty(value interface{}) error {
	party, _ := value.(*org.Party)
	if party == nil {
		return nil
	}
	return validation.ValidateStruct(party,
		validation.Field(&party.Name,
			validation.Required,
		),
		validation.Field(&party.Addresses,
			validation.Required,
			validation.Length(1, 1),
			validation.Each(validation.By(validateAddress)),
		),
		validation.Field(&party.People,
			validation.Required,
			validation.Length(1, 1),
		),
		validation.Field(&party.Telephones,
			validation.Required,
			validation.Length(1, 1),
		),
		validation.Field(&party.Emails,
			validation.Required,
			validation.Length(1, 1),
		),
	)
}

func validateAddress(value interface{}) error {
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
			validation.Each(validation.By(validateGermanAddress)),
		),
	)
}

func validateGermanAddress(value interface{}) error {
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
		validation.Field(&addr.Country,
			validation.In("DE"),
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
