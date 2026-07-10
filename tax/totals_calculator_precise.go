package tax

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

	t.Calculate(tc.Currency, RoundingRulePrecise)
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
