package bill

import (
	"fmt"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// Lines holds an array of Line objects.
type Lines []*Line

// Line is a single row in an invoice. For tax calculations, it
// represents the base.
type Line struct {
	UUID     string        `json:"uuid,omitempty"`
	Index    int           `json:"index" jsonschema:"title=Index,description=Line number inside the invoice."`
	Quantity num.Amount    `json:"quantity"`
	Item     *org.Item     `json:"item"`
	Sum      num.Amount    `json:"sum" jsonschema:"title=Sum,description=Result of quantity multiplied by item price"`
	Discount *org.Discount `json:"discount,omitempty" jsonschema:"title=Discount,description=Discount applied to this line."`
	Taxes    tax.Rates     `json:"taxes,omitempty" jsonschema:"title=Taxes,description=List of taxes to be applied to the line in the invoice totals."`
	Total    num.Amount    `json:"total" jsonschema:"title=Total,description=Total line amount after applying discounts to the sum."`
}

// calculate takes the provided region and date to correctly
// assign taxes and totals for the line. Both of these fields
// should be extractable from the parent invoice.
func (l *Line) calculate(reg *tax.Region, date org.Date) error {
	// First we figure out how much the item costs, and get the total
	l.Sum = l.Item.Price.Multiply(l.Quantity)

	if l.Discount != nil {
		d := l.Discount
		if !d.Rate.IsZero() {
			// always override value with calculated rate
			d.Value = d.Rate.Of(l.Sum)
		}
		l.Total = l.Sum.Subtract(d.Value)
	} else {
		l.Total = l.Sum
	}

	// Figure out the taxes
	for _, r := range l.Taxes {
		cat, ok := reg.Category(r.Category)
		if !ok {
			return fmt.Errorf("failed to find category, invalid code: %v", r.Category)
		}
		def, ok := cat.Def(r.Code)
		if !ok {
			return fmt.Errorf("failed to find rate definition, invalid code: %v", r.Code)
		}
		v, ok := def.On(date)
		if !ok {
			return fmt.Errorf("tax rate cannot be provided for date")
		}
		r.Retained = cat.Retained
		r.Percent = v.Percent
		r.Base = l.Total
		r.Value = r.Percent.Of(r.Base)
	}

	return nil
}
