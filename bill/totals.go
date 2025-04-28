package bill

import (
	"context"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Totals contains the summaries of all calculations for the invoice.
type Totals struct {
	// Total of all line item amounts.
	Sum num.Amount `json:"sum" jsonschema:"title=Sum"`
	// Total of all discounts applied at the document level.
	Discount *num.Amount `json:"discount,omitempty" jsonschema:"title=Discount"`
	// Total of all charges applied at the document level.
	Charge *num.Amount `json:"charge,omitempty" jsonschema:"title=Charge"`
	// Total tax amount included in the prices, if prices are tax-inclusive.
	TaxIncluded *num.Amount `json:"tax_included,omitempty" jsonschema:"title=Tax Included"`
	// Net total amount after subtracting discounts and adding charges, excluding tax.
	Total num.Amount `json:"total" jsonschema:"title=Total"`
	// Detailed breakdown of all taxes applied to the invoice.
	Taxes *tax.Total `json:"taxes,omitempty" jsonschema:"title=Tax Totals"`
	// Total indirect tax amount to be applied to the invoice.
	Tax num.Amount `json:"tax,omitempty" jsonschema:"title=Tax"`
	// Final total amount after applying indirect taxes.
	TotalWithTax num.Amount `json:"total_with_tax" jsonschema:"title=Total with Tax"`
	// Total tax amount retained or withheld by the customer to be paid to the tax authority.
	RetainedTax *num.Amount `json:"retained_tax,omitempty" jsonschema:"title=Retained Tax"`
	// Adjustment amount applied to the invoice totals to meet rounding rules or expectations.
	Rounding *num.Amount `json:"rounding,omitempty" jsonschema:"title=Rounding"`
	// Final amount to be paid after retained taxes and rounding adjustments.
	Payable num.Amount `json:"payable" jsonschema:"title=Payable"`
	// Total amount already paid in advance by the customer.
	Advances *num.Amount `json:"advance,omitempty" jsonschema:"title=Advance"`
	// Remaining amount that needs to be paid.
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
		validation.Field(&t.Rounding),
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
	t.RetainedTax = nil
	// t.Rounding = nil // may have been provided externally
	t.Payable = zero
	t.Advances = nil
	t.Due = nil
}

// Paid is a convenience method to quickly determine if the invoice has been
// paid or not based on the total amount due.
func (t *Totals) Paid() bool {
	return t != nil && t.Due != nil && t.Due.IsZero()
}

// round goes through each value that is set and rescales to match
// the zero's exponent
func (t *Totals) round(zero num.Amount) {
	e := zero.Exp()
	t.Sum = t.Sum.Rescale(e)
	if t.Discount != nil {
		*t.Discount = t.Discount.Rescale(e)
	}
	if t.Charge != nil {
		*t.Charge = t.Charge.Rescale(e)
	}
	if t.TaxIncluded != nil {
		*t.TaxIncluded = t.TaxIncluded.Rescale(e)
	}
	t.Total = t.Total.Rescale(e)
	t.Tax = t.Tax.Rescale(e)
	t.TotalWithTax = t.TotalWithTax.Rescale(e)
	if t.RetainedTax != nil {
		*t.RetainedTax = t.RetainedTax.Rescale(e)
	}
	t.Payable = t.Payable.Rescale(e)
	if t.Advances != nil {
		*t.Advances = t.Advances.Rescale(e)
	}
	if t.Due != nil {
		*t.Due = t.Due.Rescale(e)
	}
}
