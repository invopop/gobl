package tax

// calculateCurrencyTotals performs all calculations using the currency's
// precision, rounding each line's amounts before adding them to the sums,
// so that the results can always be recalculated from the data presented.
func (tc *TotalCalculator) calculateCurrencyTotals(t *Total, taxLines []*taxLine) error {
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

	tc.calculateFinalSums(t)
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
		// tax-inclusive totals and their taxes
		rt := t.rateTotalFor(c, tc.zero)
		rt.Base = rt.Base.MatchPrecision(tl.total).Add(tl.total)
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

// calculateFinalSums provides the final category and document sums, keeping
// any tax amounts already extracted from tax-inclusive totals.
func (tc *TotalCalculator) calculateFinalSums(t *Total) {
	t.Sum = tc.zero
	for _, ct := range t.Categories {
		ct.Amount = tc.zero
		for _, rt := range ct.Rates {
			if rt.Percent == nil {
				rt.Amount = tc.zero
				continue // exempt, nothing else to do
			}
			if ct.Code != tc.Includes {
				rt.Amount = rt.Percent.Of(rt.Base)
			}
			ct.Amount = ct.Amount.Add(rt.Amount)
			if rt.Surcharge != nil {
				rt.Surcharge.Amount = rt.Surcharge.Percent.Of(rt.Base)
				s := tc.zero
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
			r := tc.zero
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
