package tax

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
)

// CalculatorRoundingRule defines the available methods for calculating the
// totals in the tax calculator.
type CalculatorRoundingRule string

const (
	// CalculatorSumThenRound is the default method of calculating the totals
	// in GOBL, and provides the best results for most cases as the precision
	// is maintained to the maximum amount possible. The tradeoff however is
	// that sometimes the totals may not sum exactly based on what is visible.
	CalculatorSumThenRound CalculatorRoundingRule = "sum-then-round"
	// CalculatorRoundThenSum is the alternative method of calculating the totals
	// that will first round all the amounts to the currency's precision before
	// making the sums. Totals using this approach can always be recalculated using
	// the amounts presented, but can lead to rounding errors in the case of
	// pre-payments and when line item prices include tax.
	CalculatorRoundThenSum CalculatorRoundingRule = "round-then-sum"
)

// TotalCalculator defines the base structure with the available
// data for calculating tax totals.
type TotalCalculator struct {
	Country  l10n.TaxCountryCode
	Tags     []cbc.Key
	Zero     num.Amount
	Date     cal.Date
	Lines    []TaxableLine
	Includes cbc.Code // Tax included in price
}

// TaxableLine defines what we expect from a line in order to subsequently calculate
// the taxes that need to be added or retained.
type TaxableLine interface {
	GetTaxes() Set
	GetTotal() num.Amount
}

// Calculate the totals
func (tc *TotalCalculator) Calculate(t *Total) error {
	// reset
	t.Categories = make([]*CategoryTotal, 0)
	t.Sum = tc.Zero

	// get simplified list of lines
	taxLines := mapTaxLines(tc.Lines)
	if err := tc.prepareLines(taxLines); err != nil {
		return err
	}

	// Remove included taxes
	if err := tc.removeIncludedTaxes(taxLines); err != nil {
		return err
	}

	tc.calculateBaseRateTotals(taxLines, t)
	tc.calculateFinalSum(t)
	tc.round(t)

	return nil
}

func (tc *TotalCalculator) prepareLines(taxLines []*taxLine) error {
	// First, prepare all tax combos using the regime, zone, and date
	for _, tl := range taxLines {
		for _, combo := range tl.taxes {
			if err := combo.calculate(tc.Country, tc.Tags, tc.Date); err != nil {
				return err
			}
			// always add 2 decimal places for all tax calculations
			tl.total = tl.total.RescaleUp(tc.Zero.Exp() + 2)
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
			if c.retained {
				return ErrInvalidPricesInclude.WithMessage("cannot include retained category '%s'", tc.Includes.String())
			}
			if c.Percent == nil {
				// no taxes, skip
				continue
			}
			tl.total = tl.total.Remove(*c.Percent)
		}
	}
	return nil
}

func (tc *TotalCalculator) calculateBaseRateTotals(taxLines []*taxLine, t *Total) {
	// Go through each line and add the total to the base of each tax
	for _, tl := range taxLines {
		for _, c := range tl.taxes {
			if c.Percent == nil && c.Rate.IsEmpty() {
				continue // not much to do here!
			}
			rt := t.rateTotalFor(c, tc.Zero)
			rt.Base = tc.matchPrecision(rt.Base, tl.total)
			rt.Base = rt.Base.Add(tl.total)
		}
	}
}

func (tc *TotalCalculator) calculateFinalSum(t *Total) {
	// Now go through each category to apply the percentage and calculate the final sums
	t.Sum = tc.Zero
	for _, ct := range t.Categories {
		tc.calculateBaseCategoryTotal(ct)

		t.Sum = tc.matchPrecision(t.Sum, ct.Amount)
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
		ct.Amount = tc.matchPrecision(ct.Amount, rt.Amount)
		ct.Amount = ct.Amount.Add(rt.Amount)
		if rt.Surcharge != nil {
			rt.Surcharge.Amount = rt.Surcharge.Percent.Of(base)
			if ct.Surcharge == nil {
				ct.Surcharge = &zero
			}
			a := rt.Surcharge.Amount
			x := *ct.Surcharge
			x = tc.matchPrecision(x, a)
			x = x.Add(a)
			ct.Surcharge = &x
		}
	}
}

// matchPrecision is used to match the precision of two amounts according to the
// current rounding rule.
func (tc *TotalCalculator) matchPrecision(a, b num.Amount) num.Amount {
	r := RegimeFor(tc.Country.Code())
	if r != nil {
		switch r.CalculatorRoundingRule {
		case CalculatorRoundThenSum:
			return a.Rescale(tc.Zero.Exp())
		}
	}
	return a.MatchPrecision(b)
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
		ct.amount = ct.Amount
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
