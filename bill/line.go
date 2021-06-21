package bill

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// Lines holds an array of Line objects.
type Lines []*Line

// Line is a single row in an invoice.
type Line struct {
	UUID     string        `json:"uuid,omitempty"`
	Index    int           `json:"i" jsonschema:"title=Index,description=Line number inside the invoice."`
	Quantity num.Amount    `json:"quantity"`
	Item     *org.Item     `json:"item"`
	Sum      num.Amount    `json:"sum" jsonschema:"title=Sum,description=Result of quantity multiplied by item price"`
	Discount *org.Discount `json:"discount,omitempty" jsonschema:"title=Discount,description=Discount applied to this line."`
	Taxes    tax.Rates     `json:"taxes,omitempty" jsonschema:"title=Taxes,description=List of taxes to be applied to the line in the invoice totals."`
	Total    num.Amount    `json:"total" jsonschema:"title=Total,description=Total line amount after applying discounts to the sum."`
}

// Validate ensures the line contains everything required.
func (l *Line) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Index, validation.Required),
		validation.Field(&l.Quantity, validation.Required),
		validation.Field(&l.Item, validation.Required),
		validation.Field(&l.Sum, validation.Required),
		validation.Field(&l.Total, validation.Required),
	)
}

// calculate takes the provided region and date to correctly
// assign taxes and totals for the line. Both of these fields
// should be extractable from the parent invoice. When taxIncluded
// is true, we assume the item prices include non-retained taxes
// and attempt to calculate the correct base and taxable amount.
// Most electronic invoicing formats do not support including taxes,
// so this is a bit of an experiment.
func (l *Line) calculate(reg *tax.Region, date org.Date, taxIncluded bool) error {
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

	// group taxes
	regular := make([]*lineRate, 0)
	retained := make([]*lineRate, 0)
	for _, r := range l.Taxes {
		lr := new(lineRate)
		if err := lr.with(reg, r, date); err != nil {
			return err
		}
		if lr.rate.Retained {
			retained = append(retained, lr)
		} else {
			regular = append(regular, lr)
		}
	}

	// figure out the base for regular taxes
	var base num.Amount
	for _, lr := range regular {
		r := lr.rate
		if taxIncluded {
			r.Value = r.Percent.From(l.Total)
			r.Base = l.Total.Subtract(r.Value)
		} else {
			r.Base = l.Total
			r.Value = r.Percent.Of(r.Base)
		}
		base = r.Base
	}

	// Deal with retained taxes, using the base
	if base.IsZero() {
		base = l.Total
	}
	for _, lr := range retained {
		r := lr.rate
		r.Base = base
		r.Value = r.Percent.Of(base)
	}

	return nil
}

// lineRate is used to make calculations easier
type lineRate struct {
	rate *tax.Rate
	cat  tax.Category
	def  tax.Def
	val  tax.Value
}

func (l *lineRate) with(reg *tax.Region, r *tax.Rate, date org.Date) error {
	l.rate = r
	var ok bool
	if l.cat, ok = reg.Category(r.Category); !ok {
		return fmt.Errorf("failed to find category, invalid code: %v", r.Category)
	}
	if l.def, ok = l.cat.Def(r.Code); !ok {
		return fmt.Errorf("failed to find rate definition, invalid code: %v", r.Code)
	}
	if l.val, ok = l.def.On(date); !ok {
		return fmt.Errorf("tax rate cannot be provided for date")
	}
	r.Retained = l.cat.Retained
	r.Percent = l.val.Percent
	return nil
}
