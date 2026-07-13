package tax

import (
	"github.com/invopop/gobl/num"
)

// calculatePreciseTotals accumulates the rate bases with at least two extra
// decimal places of precision above the currency's own, providing the most
// accurate results possible when rounding the final sums.
func (tc *TotalCalculator) calculatePreciseTotals(t *Total, taxLines []*taxLine) error {
	// always add 2 decimal places for all tax calculations
	for _, tl := range taxLines {
		tl.total = tl.total.RescaleUp(tc.zero.Exp() + 2)
	}

	if err := tc.removeIncludedTaxes(taxLines); err != nil {
		return err
	}

	// Go through each line and add the total to the base of each tax
	for _, tl := range taxLines {
		for _, c := range tl.taxes {
			rt := t.rateTotalFor(c, tc.zero)
			rt.Base = rt.Base.MatchPrecision(tl.total)
			rt.Base = rt.Base.Add(tl.total)
		}
	}

	t.calculatePreciseFinalSums(tc.zero)
	return nil
}

// removeIncludedTaxes performs a pre-loop to update all the line totals
// with the total minus the included tax.
func (tc *TotalCalculator) removeIncludedTaxes(taxLines []*taxLine) error {
	if tc.Includes.IsEmpty() {
		return nil
	}
	for _, tl := range taxLines {
		if c := tl.taxes.Get(tc.Includes); c != nil {
			if err := tc.checkIncludedCombo(c); err != nil {
				return err
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

// calculatePreciseFinalSums provides the final category and document sums
// from the accumulated bases, maintaining their precision until the totals
// are rounded to the currency later.
func (t *Total) calculatePreciseFinalSums(zero num.Amount) {
	t.Sum = zero
	t.Retained = nil
	for _, ct := range t.Categories {
		ct.Amount = zero
		ct.Surcharge = nil
		for _, rt := range ct.Rates {
			if rt.Percent == nil {
				rt.Amount = zero
				continue // exempt, nothing else to do
			}
			rt.Amount = rt.Percent.Of(rt.Base)
			ct.Amount = ct.Amount.MatchPrecision(rt.Amount)
			ct.Amount = ct.Amount.Add(rt.Amount)
			if rt.Surcharge != nil {
				rt.Surcharge.Amount = rt.Surcharge.Percent.Of(rt.Base)
				s := zero
				if ct.Surcharge != nil {
					s = *ct.Surcharge
				}
				s = s.MatchPrecision(rt.Surcharge.Amount)
				s = s.Add(rt.Surcharge.Amount)
				ct.Surcharge = &s
			}
		}

		if ct.Informative {
			// Informative taxes don't affect Sum or Retained
			continue
		}
		if ct.Retained {
			r := zero
			if t.Retained != nil {
				r = *t.Retained
			}
			r = r.MatchPrecision(ct.Amount)
			r = r.Add(ct.Amount)
			if ct.Surcharge != nil {
				r = r.Add(*ct.Surcharge)
			}
			t.Retained = &r
		} else {
			t.Sum = t.Sum.MatchPrecision(ct.Amount)
			t.Sum = t.Sum.Add(ct.Amount)
			if ct.Surcharge != nil {
				t.Sum = t.Sum.Add(*ct.Surcharge)
			}
		}
	}
}
