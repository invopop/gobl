package no

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateBillInvoice(inv *bill.Invoice) error {
	if inv == nil {
		return nil
	}

	simplified := inv.Tags.HasTags(tax.TagSimplified)
	reverseCharge := inv.Tags.HasTags(tax.TagReverseCharge)

	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.Required,
			validation.By(validateBillInvoiceParty(!simplified, reverseCharge)),
			validation.By(validateBillInvoiceSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.When(
				simplified,
				validation.Empty,
			).Else(
				validation.Required,
				validation.By(validateBillInvoiceParty(true, reverseCharge)),
				validation.By(validateBillInvoiceCustomer(reverseCharge)),
			),
			validation.Skip,
		),
	)
}

func validateBillInvoiceParty(withAddress bool, reverseCharge bool) validation.RuleFunc {
	return func(value any) error {
		party, ok := value.(*org.Party)
		if !ok || party == nil {
			return nil
		}

		return validation.ValidateStruct(party,
			validation.Field(&party.Name,
				validation.Required,
				validation.Skip,
			),
			validation.Field(&party.Addresses,
				validation.When(withAddress,
					validation.Required,
				),
				validation.Skip,
			),
		)
	}
}

func validateBillInvoiceSupplier(value any) error {
	party, _ := value.(*org.Party)
	return validation.ValidateStruct(party,
		validation.Field(&party.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
	)
}

func validateBillInvoiceCustomer(reverseCharge bool) validation.RuleFunc {
	return func(value any) error {
		if !reverseCharge {
			return nil
		}
		party, ok := value.(*org.Party)
		if !ok || party == nil {
			return nil
		}
		// For reverse charge invoices we require the customer TaxID as well,
		// since the buyer accounts for VAT.
		return validation.ValidateStruct(party,
			validation.Field(&party.TaxID,
				validation.Required,
				tax.RequireIdentityCode,
				validation.Skip,
			),
		)
	}
}
