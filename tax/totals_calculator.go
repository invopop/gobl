package tax

import (
	"errors"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
)

// Tax total calculator options
const (
	TotalCalculatorTotal cbc.Key = "total" // default
	TotalCalculatorLine  cbc.Key = "line"
)

// TotalCalculatorDefs to use in the schema
var TotalCalculatorDefs = []*KeyDefinition{
	{
		Key: TotalCalculatorTotal,
		Name: i18n.String{
			i18n.EN: "Total",
		},
		Desc: i18n.String{
			i18n.EN: "Calculate the taxes based on the sum of all the line items (default).",
		},
	},
	{
		Key: TotalCalculatorLine,
		Name: i18n.String{
			i18n.EN: "Line",
		},
		Desc: i18n.String{
			i18n.EN: "Calculate the taxes based on each line item.",
		},
	},
}

// TotalCalculator defines the base structure with the available
// data for calculating tax totals.
type TotalCalculator struct {
	Regime     *Regime
	Zone       l10n.Code
	Zero       num.Amount
	Date       cal.Date
	Lines      []TaxableLine
	Includes   cbc.Code // Tax included in price
	Calculator cbc.Key  // Calculation model to use
}

// TaxableLine defines what we expect from a line in order to subsequently calculate
// the taxes that need to be added or retained.
type TaxableLine interface {
	GetTaxes() Set
	GetTotal() num.Amount
}

// Calculate the totals
func (tc *TotalCalculator) Calculate(t *Total) error {
	if tc.Regime == nil {
		return ErrMissingRegime
	}
	if tc.Calculator == cbc.KeyEmpty {
		tc.Calculator = TotalCalculatorTotal
	}

	// reset
	t.Categories = make([]*CategoryTotal, 0)
	t.Sum = tc.Zero

	// get simplified list of lines
	taxLines := mapTaxLines(tc.Lines)

	if err := tc.prepareLines(taxLines); err != nil {
		return err
	}

	// Pre-Process each line for tax calculations
	var err error
	switch tc.Calculator {
	case TotalCalculatorLine:
		err = tc.calculateLineTaxes(taxLines)
	case TotalCalculatorTotal:
		err = tc.removeIncludedTaxes(taxLines)
	default:
		err = errors.New("unknown tax calculator type")
	}
	if err != nil {
		return err
	}

	// Go through each line to calculate rate totals
	switch tc.Calculator {
	case TotalCalculatorLine:
		tc.calculateLineRateTotals(taxLines, t)
	case TotalCalculatorTotal:
		tc.calculateBaseRateTotals(taxLines, t)
	}

	tc.calculateFinalSum(t)
	tc.round(t)

	return nil
}

func (tc *TotalCalculator) prepareLines(taxLines []*taxLine) error {
	// First, prepare all tax combos using the regime, zone, and date
	for _, tl := range taxLines {
		for _, combo := range tl.taxes {
			if err := combo.prepare(tc.Regime, tc.Zone, tc.Date); err != nil {
				return err
			}
			// always add 2 decimal places for all tax calculations
			tl.total = tl.total.RescaleUp(tc.Zero.Exp() + 2)
		}
	}
	return nil
}

func (tc *TotalCalculator) calculateLineTaxes(taxLines []*taxLine) error {
	// Go through each line, and figure out the totals for each tax combo
	for _, tl := range taxLines {
		// prepare included taxes first so we can update the total
		if tc.Includes != cbc.CodeEmpty {
			if combo := tl.taxes.Get(tc.Includes); combo != nil && combo.Percent != nil {
				if combo.category.Retained {
					return ErrInvalidPricesInclude.WithMessage("cannot include retained category '%s'", tc.Includes.String())
				}
				tl.total = tl.total.Remove(*combo.Percent)
			}
		}

		// Make calculations
		for _, c := range tl.taxes {
			c.base = tl.total
			if c.Percent != nil {
				c.amount = c.Percent.Of(c.base)
			}
			if c.Surcharge != nil {
				sc := c.Surcharge.Of(tl.total)
				c.surcharge = &sc
			}
		}
	}
	return nil
}

func (tc *TotalCalculator) removeIncludedTaxes(taxLines []*taxLine) error {
	// If prices include a tax, perform a pre-loop to update all the line prices with
	// the price minus the defined tax.
	if tc.Includes.IsEmpty() {
		return nil
	}
	for _, tl := range taxLines {
		if c := tl.taxes.Get(tc.Includes); c != nil {
			if c.category.Retained {
				return ErrInvalidPricesInclude.WithMessage("cannot include retained category '%s'", tc.Includes.String())
			}
			if c.Percent == nil {
				// can't work without a percent value, just skip
				continue
			}
			tl.total = tl.total.Remove(*c.Percent)
		}
	}
	return nil
}

// calculateLineRateTotals goes through each line to sum the rate totals.
// This is when the rounding method starts to become important. If we're doing
// post rounding, then then accuracy will be maintained, otherwise each step
// will perform rounding.
func (tc *TotalCalculator) calculateLineRateTotals(taxLines []*taxLine, t *Total) {
	for _, tl := range taxLines {
		for _, combo := range tl.taxes {
			rt := t.rateTotalFor(combo, tc.Zero)
			rt.Base = rt.Base.MatchPrecision(combo.base)
			rt.Base = rt.Base.Add(combo.base)
			if combo.Percent == nil && combo.Rate.IsEmpty() {
				continue // not much to do here!
			}

			rt.Amount = rt.Amount.MatchPrecision(combo.amount)
			rt.Amount = rt.Amount.Add(combo.amount)
			if combo.surcharge != nil {
				rt.Surcharge.Amount = rt.Surcharge.Amount.MatchPrecision(*combo.surcharge)
				rt.Surcharge.Amount = rt.Surcharge.Amount.Add(*combo.surcharge)
			}
		}
	}
}

func (tc *TotalCalculator) calculateBaseRateTotals(taxLines []*taxLine, t *Total) {
	// Go through each line and add the total to the base of each tax
	for _, tl := range taxLines {
		for _, c := range tl.taxes {
			if c.Percent == nil && c.Rate.IsEmpty() {
				continue // not much to do here!
			}
			rt := t.rateTotalFor(c, tc.Zero)
			rt.Base = rt.Base.MatchPrecision(tl.total)
			rt.Base = rt.Base.Add(tl.total)
		}
	}
}

func (tc *TotalCalculator) calculateFinalSum(t *Total) {
	// Now go through each category to apply the percentage and calculate the final sums
	t.Sum = tc.Zero
	for _, ct := range t.Categories {
		if tc.Calculator == TotalCalculatorLine {
			tc.calculateLineCategoryTotal(ct)
		} else {
			tc.calculateBaseCategoryTotal(ct)
		}

		t.Sum = t.Sum.MatchPrecision(ct.Amount)
		if ct.Retained {
			t.Sum = t.Sum.Subtract(ct.Amount)
			if ct.Surcharge != nil {
				t.Sum = t.Sum.Subtract(*ct.Surcharge)
			}
		} else {
			t.Sum = t.Sum.Add(ct.Amount)
			if ct.Surcharge != nil {
				t.Sum = t.Sum.Add(*ct.Surcharge)
			}
		}
	}
}

func (tc *TotalCalculator) calculateLineCategoryTotal(ct *CategoryTotal) {
	zero := tc.Zero
	ct.Amount = zero
	for _, rt := range ct.Rates {
		if rt.Percent == nil {
			rt.Amount = zero
			continue // exempt, nothing else to do
		}
		ct.Amount = ct.Amount.MatchPrecision(rt.Amount)
		ct.Amount = ct.Amount.Add(rt.Amount)
		if rt.Surcharge != nil {
			if ct.Surcharge == nil {
				ct.Surcharge = &zero
			}
			a := rt.Surcharge.Amount
			x := *ct.Surcharge
			x = x.MatchPrecision(a)
			x = x.Add(a)
			ct.Surcharge = &x
		}
	}
}

func (tc *TotalCalculator) calculateBaseCategoryTotal(ct *CategoryTotal) {
	zero := tc.Zero
	ct.Amount = zero
	for _, rt := range ct.Rates {
		if rt.Percent == nil {
			rt.Amount = zero
			continue // exempt, nothing else to do
		}
		base := rt.Base
		rt.Amount = rt.Percent.Of(rt.Base)
		ct.Amount = ct.Amount.MatchPrecision(rt.Amount)
		ct.Amount = ct.Amount.Add(rt.Amount)
		if rt.Surcharge != nil {
			rt.Surcharge.Amount = rt.Surcharge.Percent.Of(base)
			if ct.Surcharge == nil {
				ct.Surcharge = &zero
			}
			a := rt.Surcharge.Amount
			x := *ct.Surcharge
			x = x.MatchPrecision(a)
			x = x.Add(a)
			ct.Surcharge = &x
		}
	}
}

// round will go through all the values generated and round them to the currency's
// preferred precision. The final precise sum will be available in the t.sum variable
// still.
func (tc *TotalCalculator) round(t *Total) {
	zero := tc.Zero
	for _, ct := range t.Categories {
		for _, rt := range ct.Rates {
			rt.Amount = rt.Amount.Rescale(zero.Exp())
			rt.Base = rt.Base.Rescale(zero.Exp())
			if rt.Surcharge != nil {
				rt.Surcharge.Amount = rt.Surcharge.Amount.Rescale(zero.Exp())
			}
		}
		ct.Amount = ct.Amount.Rescale(zero.Exp())
		if ct.Surcharge != nil {
			*ct.Surcharge = ct.Surcharge.Rescale(zero.Exp())
		}
	}
	t.sum = t.Sum
	t.Sum = t.Sum.Rescale(zero.Exp())
}

// taxLine is used to replace
type taxLine struct {
	total num.Amount
	taxes Set
}

func mapTaxLines(lines []TaxableLine) []*taxLine {
	tls := make([]*taxLine, len(lines))
	for i, v := range lines {
		tls[i] = &taxLine{
			total: v.GetTotal(),
			taxes: v.GetTaxes(),
		}
	}
	return tls
}
