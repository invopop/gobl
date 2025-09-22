package au

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.By(validateInvoiceSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.By(validateInvoiceCustomer(inv)),
			validation.Skip,
		),
	)
}

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

func validateInvoiceCustomer(inv *bill.Invoice) validation.RuleFunc {
	return func(value any) error {
		p, ok := value.(*org.Party)
		if !ok || p == nil {
			return nil
		}

		if requiresCustomerIdentity(inv) {
			return validation.ValidateStruct(p,
				validation.Field(&p.Name,
					validation.When(
						p.TaxID == nil || p.TaxID.Code == "",
						validation.Required.Error("customer identity or ABN required for invoices $1,000 AUD or more"),
					),
				),
				validation.Field(&p.TaxID,
					validation.When(
						p.Name == "",
						validation.Required.Error("customer identity or ABN required for invoices $1,000 AUD or more"),
						tax.RequireIdentityCode,
					),
					validation.Skip,
				),
			)
		}

		return nil
	}
}

// For invoices $1,000 AUD or more, customer identity or ABN is required
// Source: https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/gst/tax-invoices
func requiresCustomerIdentity(inv *bill.Invoice) bool {
	threshold := num.MakeAmount(100000, 2)

	if !inv.Totals.TotalWithTax.IsZero() {
		return inv.Totals.TotalWithTax.Compare(threshold) >= 0
	}

	if !inv.Totals.Total.IsZero() {
		return inv.Totals.Total.Compare(threshold) >= 0
	}

	return false
}
