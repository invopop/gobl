package au

import (
	"errors"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

/*
 * Invoice Validation Rules for Australia
 *
 * Source: https://www.ato.gov.au/business/gst/issuing-tax-invoices
 *
 * Key requirement:
 * For tax invoices with a taxable value of A$1,000 or more (GST inclusive),
 * the invoice must include either:
 * - The buyer's name, OR
 * - The buyer's ABN (Australian Business Number)
 *
 * This validator calculates the taxable amount (sum of line totals for lines
 * with GST category) and enforces the buyer identity requirement.
 */

// Threshold for requiring buyer identification (A$1,000.00)
// Stored as 100000 cents with 2 decimal places
var buyerIdentityThreshold = num.MakeAmount(100000, 2)

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Customer,
			validation.By(func(value interface{}) error {
				return validateCustomer(inv, value)
			}),
			validation.Skip,
		),
	)
}

func validateCustomer(inv *bill.Invoice, value interface{}) error {
	customer, ok := value.(*org.Party)
	if !ok || customer == nil {
		return nil
	}

	// Calculate taxable amount (sum of line totals with GST category)
	taxableAmount := calculateTaxableAmount(inv)

	// If taxable amount >= A$1,000, require buyer name OR ABN
	if taxableAmount.Compare(buyerIdentityThreshold) >= 0 {
		hasName := customer.Name != ""
		hasABN := customer.TaxID != nil && customer.TaxID.Code != ""

		if !hasName && !hasABN {
			return errors.New("buyer name or ABN required for invoices with taxable amount >= A$1,000")
		}
	}

	return nil
}

// calculateTaxableAmount sums the totals of all lines with GST category.
// This includes both standard rate (10%) and zero-rated (0%) GST supplies,
// as both are considered "taxable supplies" under Australian GST law.
func calculateTaxableAmount(inv *bill.Invoice) num.Amount {
	if inv == nil || inv.Lines == nil {
		return num.AmountZero
	}

	sum := num.AmountZero
	for _, line := range inv.Lines {
		if line == nil || line.Total == nil {
			continue
		}

		// Check if line has GST tax category
		if hasGSTCategory(line.Taxes) {
			sum = sum.Add(*line.Total)
		}
	}

	return sum
}

// hasGSTCategory checks if the tax set contains a GST category.
// This returns true for both standard rate (10%) and zero-rated (0%) GST,
// as both are taxable supplies under Australian tax law.
func hasGSTCategory(taxes tax.Set) bool {
	if taxes == nil {
		return false
	}

	for _, combo := range taxes {
		if combo != nil && combo.Category == tax.CategoryGST {
			return true
		}
	}
	return false
}
