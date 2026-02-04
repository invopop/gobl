package nz

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// NZ GST thresholds (GST-inclusive) per Taxable Supply Information rules.
var (
	thresholdMid  = num.MakeAmount(200, 0)  // $200 NZD
	thresholdHigh = num.MakeAmount(1000, 0) // $1,000 NZD
)

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			// Supplier TaxID required for supplies over $200.
			validation.When(
				invoiceTotalExceeds(inv, thresholdMid),
				validation.By(validateSupplierTaxID),
			),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			// Customer name + identifier required for supplies over $1,000.
			validation.When(
				invoiceTotalExceeds(inv, thresholdHigh),
				validation.Required,
				validation.By(validateCustomerDetails),
			),
			validation.Skip,
		),
	)
}

func invoiceTotalExceeds(inv *bill.Invoice, threshold num.Amount) bool {
	if inv.Totals == nil {
		return false
	}
	return inv.Totals.TotalWithTax.Compare(threshold) > 0
}

func validateSupplierTaxID(value any) error {
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

func validateCustomerDetails(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	// For supplies over $1,000, buyer name is required plus at least
	// one identifier: address, phone, email, or website.
	return validation.ValidateStruct(p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Addresses,
			validation.When(
				!hasCustomerIdentifier(p),
				validation.Required,
			),
			validation.Skip,
		),
	)
}

// hasCustomerIdentifier returns true if the customer has at least one
// identifier beyond their name (address, phone, email, website).
func hasCustomerIdentifier(p *org.Party) bool {
	return len(p.Addresses) > 0 ||
		len(p.Telephones) > 0 ||
		len(p.Emails) > 0 ||
		len(p.Websites) > 0 ||
		len(p.Identities) > 0
}
