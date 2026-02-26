package ro

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateInvoice(inv *bill.Invoice) error {
	simplified := inv.Tags.HasTags(tax.TagSimplified)
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.By(validateInvoiceSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.When(!simplified, validation.Required),
			validation.Skip,
		),
	)
}

// validateInvoiceSupplier ensures the supplier has a valid tax identity.
// The CUI (Codul Unic de Înregistrare) is always required: it is assigned to
// every business upon registration with ONRC, regardless of VAT status.
// Non-VAT-registered businesses (below the RON 395,000 threshold) still have
// a CUI and issue invoices without VAT. The "RO" prefix indicates VAT
// registration and is stripped during normalization.

func validateInvoiceSupplier(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
	)
}
