// Package ae provides the tax region definition for the United Arab Emirates.
package ae

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// validateInvoice validates the invoice, ensuring the TRN is always required.
func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.By(validateInvoiceSupplier),
		),
	)
}

// validateInvoiceSupplier validates that the supplier (Party) has a Tax Registration Number (TRN).
func validateInvoiceSupplier(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
		),
	)
}
