package xrechnung

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

const (
	invoiceTypeSelfBilled               cbc.Key = "389"
	invoiceTypePartial                  cbc.Key = "326"
	invoiceTypePartialConstruction      cbc.Key = "875"
	invoiceTypePartialFinalConstruction cbc.Key = "876"
	invoiceTypeFinalConstruction        cbc.Key = "877"
)

var validTypes = []cbc.Key{
	bill.InvoiceTypeStandard,
	bill.InvoiceTypeCreditNote,
	bill.InvoiceTypeCorrective,
	invoiceTypeSelfBilled,
	invoiceTypePartial,
	invoiceTypePartialConstruction,
	invoiceTypePartialFinalConstruction,
	invoiceTypeFinalConstruction,
}

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		// BR-DE-17
		validation.Field(&inv.Type,
			validation.By(validateInvoiceType),
		),
		// BR-DE-01
		validation.Field(&inv.Payment.Instructions,
			validation.Required,
			validation.By(validatePaymentInstructions),
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
		// BR-DE-26
		validation.Field(&inv,
			validation.By(validateCorrectiveInvoice),
		),
	)
}

func validateInvoiceType(value interface{}) error {
	t, ok := value.(cbc.Key)
	if !ok {
		return validation.NewError("type", "Invalid invoice type")
	}
	if t.In(validTypes...) {
		return nil
	}
	return validation.NewError("invalid", "Invalid invoice type")
}

func validateSupplier(value interface{}) error {
	party, _ := value.(*org.Party)
	if party == nil {
		return nil
	}
	return validation.ValidateStruct(party,
		// BR-DE-02
		validation.Field(&party.Name,
			validation.Required,
		),
		// BR-DE-03, BR-DE-04
		validation.Field(&party.Addresses,
			validation.Required,
			validation.Each(validation.By(validatePartyAddress)),
		),
		// BR-DE-06
		validation.Field(&party.People,
			validation.Required,
		),
		// BR-DE-05
		validation.Field(&party.Telephones,
			validation.Required,
		),
		// BR-DE-07
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
		// BR-DE-08, BR-DE-09
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

// BR-DE-14
func validateVATRate(value interface{}) error {
	rate, _ := value.(cbc.Key)
	if rate == "" {
		return validation.NewError("required", "VAT category rate is required")
	}
	return nil
}

func validateCorrectiveInvoice(value interface{}) error {
	inv, ok := value.(*bill.Invoice)
	if !ok || inv == nil {
		return nil
	}
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
