package bill

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// Lines holds an array of Line objects.
type Lines []*Line

// Line is a single row in an invoice.
type Line struct {
	// Unique identifier for this line
	UUID string `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Line number inside the parent
	Index int `json:"i" jsonschema:"title=Index"`
	// Number of items
	Quantity num.Amount `json:"quantity" jsonschema:"title=Quantity"`
	// Details about what is being sold
	Item *org.Item `json:"item" jsonschema:"title=Item"`
	// Result of quantity multiplied by the item's price
	Sum num.Amount `json:"sum" jsonschema:"title=Sum"`
	// Discounts applied to this line
	Discounts []*LineDiscount `json:"discounts,omitempty" jsonschema:"title=Discounts"`
	// Charges applied to this line
	Charges []*LineCharge `json:"charges,omitempty" jsonschema:"title=Charges"`
	// List of taxes to be applied and used in the invoice totals
	Taxes tax.Rates `json:"taxes,omitempty" jsonschema:"title=Taxes"`
	// Total line amount after applying discounts to the sum
	Total num.Amount `json:"total" jsonschema:"title=Total"`
}

// GetTaxRates responds with the array of tax rates applied to this line.
func (l *Line) GetTaxRates() tax.Rates {
	return l.Taxes
}

// GetTotal provides the final total for this line, excluding any tax calculations.
func (l *Line) GetTotal() num.Amount {
	return l.Total
}

// Validate ensures the line contains everything required.
func (l *Line) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Index, validation.Required),
		validation.Field(&l.Quantity, validation.Required),
		validation.Field(&l.Item, validation.Required),
		validation.Field(&l.Discounts),
		validation.Field(&l.Charges),
		validation.Field(&l.Taxes),
		validation.Field(&l.Sum, validation.Required),
		validation.Field(&l.Total, validation.Required),
	)
}

// calculate figures out the totals according to quantity and discounts.
func (l *Line) calculate() {
	// First we figure out how much the item costs, and get the total
	l.Sum = l.Item.Price.Multiply(l.Quantity)
	l.Total = l.Sum

	for _, d := range l.Discounts {
		if d.Rate != nil && !d.Rate.IsZero() {
			d.Value = d.Rate.Of(l.Sum) // always override
		}
		l.Total = l.Total.Subtract(d.Value)
	}

	for _, c := range l.Charges {
		if c.Rate != nil && !c.Rate.IsZero() {
			c.Value = c.Rate.Of(l.Sum) // always override
		}
		l.Total = l.Total.Add(c.Value)
	}
}
