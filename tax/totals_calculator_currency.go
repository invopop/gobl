package tax

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
)

// calculateCurrencyTotals performs all calculations using the currency's
// precision, rounding each line's amounts before adding them to the sums,
// so that the results can always be recalculated from the data presented.
func (tc *TotalCalculator) calculateCurrencyTotals(t *Total, taxLines []*taxLine) error {
	// round all the line totals to the currency's precision
	for _, tl := range taxLines {
		tl.total = tl.total.Rescale(tc.zero.Exp())
	}

	if err := tc.extractIncludedTaxes(t, taxLines); err != nil {
		return err
	}

	// Go through each line and add the total to the base of each tax,
	// rounded to the currency's precision.
	for _, tl := range taxLines {
		for _, c := range tl.taxes {
			if c.Category == tc.Includes && c.Percent != nil {
				// base already accumulated during extraction
				continue
			}
			rt := t.rateTotalFor(c, tc.zero)
			rt.Base = rt.Base.Add(tl.total)
		}
	}

	t.calculateCurrencyFinalSums(tc.zero, tc.Includes)
	return nil
}

// extractIncludedTaxes removes the included tax from each line's total.
// The tax is calculated over the accumulated sum of the rate's tax-inclusive
// line totals, with each line assuming the increment, so that the resulting
// bases and tax amounts will always add up to the sum of the original totals.
func (tc *TotalCalculator) extractIncludedTaxes(t *Total, taxLines []*taxLine) error {
	if tc.Includes.IsEmpty() {
		return nil
	}
	for _, tl := range taxLines {
		c := tl.taxes.Get(tc.Includes)
		if c == nil {
			continue
		}
		if err := tc.checkIncludedCombo(c); err != nil {
			return err
		}
		if c.Percent == nil {
			// no taxes, skip
			continue
		}
		// use the rate total to accumulate the running sums of the
		// tax-inclusive totals and their taxes, with four extra decimal
		// places of precision to avoid double-rounding errors close to
		// the currency's rounding boundaries
		rt := t.rateTotalFor(c, tc.zero)
		rt.Base = rt.Base.Add(tl.total)
		tax := c.Percent.From(rt.Base.Upscale(4)).Rescale(tc.zero.Exp())
		tl.total = tl.total.Subtract(tax.Subtract(rt.Amount))
		rt.Amount = tax
	}

	// swap the accumulated tax-inclusive totals for the final bases
	if ct := t.Category(tc.Includes); ct != nil {
		for _, rt := range ct.Rates {
			rt.Base = rt.Base.Subtract(rt.Amount)
		}
	}
	return nil
}

// calculateCurrencyFinalSums provides the final category and document sums
// from the accumulated bases at the currency's precision, keeping any tax
// amounts already extracted from the tax-inclusive totals of the included
// category.
func (t *Total) calculateCurrencyFinalSums(zero num.Amount, includes cbc.Code) {
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
			if ct.Code != includes {
				rt.Amount = rt.Percent.Of(rt.Base)
			}
			ct.Amount = ct.Amount.Add(rt.Amount)
			if rt.Surcharge != nil {
				rt.Surcharge.Amount = rt.Surcharge.Percent.Of(rt.Base)
				s := zero
				if ct.Surcharge != nil {
					s = *ct.Surcharge
				}
				s = s.Add(rt.Surcharge.Amount)
				ct.Surcharge = &s
			}
		}

		if ct.Informative {
			// Informative taxes don't affect Sum or Retained
			continue
		}
		a := ct.Amount
		if ct.Surcharge != nil {
			a = a.Add(*ct.Surcharge)
		}
		if ct.Retained {
			r := zero
			if t.Retained != nil {
				r = *t.Retained
			}
			r = r.Add(a)
			t.Retained = &r
		} else {
			t.Sum = t.Sum.Add(a)
		}
	}
}
