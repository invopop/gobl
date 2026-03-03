package no

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateBillInvoice(inv *bill.Invoice) error {
	simplified := inv.Tags.HasTags(tax.TagSimplified)
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.By(validateBillInvoiceSupplier(!simplified)),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.When(
				!simplified,
				validation.Required,
			),
			validation.By(validateBillInvoiceCustomer),
			validation.Skip,
		),
		validation.Field(&inv.Preceding,
			validation.When(
				inv.Type.In(
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				),
				validation.Required,
			),
			validation.Skip,
		),
	)
}

// validateBillInvoiceSupplier checks supplier fields. TaxID is intentionally
// not required: Norwegian businesses below NOK 50,000 turnover are not
// VAT-registered and cannot charge MVA, but may still issue invoices.
func validateBillInvoiceSupplier(withAddress bool) validation.RuleFunc {
	return func(value any) error {
		party, _ := value.(*org.Party)
		if party == nil {
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

// validateBillInvoiceCustomer checks that the customer has a name. Address is
// not required for the customer per Norwegian VAT regulations (Regnskapsloven);
// only the supplier address is mandatory on domestic invoices.
func validateBillInvoiceCustomer(value any) error {
	party, _ := value.(*org.Party)
	if party == nil {
		return nil
	}
	return validation.ValidateStruct(party,
		validation.Field(&party.Name,
			validation.Required,
			validation.Skip,
		),
	)
}
