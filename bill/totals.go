package bill

import (
	"context"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Totals contains the summaries of all calculations for the invoice.
type Totals struct {
	// Sum of all line item sums
	Sum num.Amount `json:"sum" jsonschema:"title=Sum"`
	// Sum of all document level discounts
	Discount *num.Amount `json:"discount,omitempty" jsonschema:"title=Discount"`
	// Sum of all document level charges
	Charge *num.Amount `json:"charge,omitempty" jsonschema:"title=Charge"`
	// If prices include tax, this is the total tax included in the price.
	TaxIncluded *num.Amount `json:"tax_included,omitempty" jsonschema:"title=Tax Included"`
	// Sum of all line sums minus the discounts, plus the charges, without tax.
	Total num.Amount `json:"total" jsonschema:"title=Total"`
	// Summary of all the taxes included in the invoice.
	Taxes *tax.Total `json:"taxes,omitempty" jsonschema:"title=Tax Totals"`
	// Total amount of tax to apply to the invoice.
	Tax num.Amount `json:"tax,omitempty" jsonschema:"title=Tax"`
	// Grand total after all taxes have been applied.
	TotalWithTax num.Amount `json:"total_with_tax" jsonschema:"title=Total with Tax"`
	// Total paid in outlays that need to be reimbursed
	Outlays *num.Amount `json:"outlays,omitempty" jsonschema:"title=Outlay Totals"`
	// Total amount to be paid after applying taxes and outlays.
	Payable num.Amount `json:"payable" jsonschema:"title=Payable"`
	// Total amount already paid in advance.
	Advances *num.Amount `json:"advance,omitempty" jsonschema:"title=Advance"`
	// How much actually needs to be paid now.
	Due *num.Amount `json:"due,omitempty" jsonschema:"title=Due"`
}

// ValidateWithContext checks the totals calculated for the invoice.
func (t *Totals) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, t,
		validation.Field(&t.Sum, validation.Required),
		validation.Field(&t.Discount),
		validation.Field(&t.Charge),
		validation.Field(&t.TaxIncluded),
		validation.Field(&t.Total, validation.Required),
		validation.Field(&t.Taxes),
		validation.Field(&t.Tax),
		validation.Field(&t.TotalWithTax),
		validation.Field(&t.Outlays),
		validation.Field(&t.Payable),
		validation.Field(&t.Advances),
		validation.Field(&t.Due),
	)
}

// Reset sets all the totals to the provided zero amount with the correct
// decimal places.
func (t *Totals) reset(zero num.Amount) {
	t.Sum = zero
	t.Discount = nil
	t.Charge = nil
	t.TaxIncluded = nil
	t.Total = zero
	t.Taxes = nil
	t.Tax = zero
	t.TotalWithTax = zero
	t.Outlays = nil
	t.Payable = zero
	t.Advances = nil
	t.Due = nil
}
