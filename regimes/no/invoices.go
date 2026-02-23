package no

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateBillInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.By(validateBillInvoiceSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.Required,
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

func validateBillInvoiceSupplier(value any) error {
	party, _ := value.(*org.Party)
	if party == nil {
		return nil
	}
	return validation.ValidateStruct(party,
		validation.Field(&party.Name,
			validation.Required,
			validation.Skip,
		),
		validation.Field(&party.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
		validation.Field(&party.Addresses,
			validation.Required,
			validation.Skip,
		),
	)
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
