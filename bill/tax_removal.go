package bill

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

const (
	// maxIncludedTaxAlignPasses limits how many times a document is
	// recalculated while lining its rows up with their net totals.
	maxIncludedTaxAlignPasses = 4
	// maxIncludedTaxPriceAccuracy is the most extra decimal places that may be
	// added to a net item price in order to reproduce a line's sum.
	maxIncludedTaxPriceAccuracy uint32 = 6
)

// includedTaxNetTotals holds the total that each of a document's rows should
// end up with once the taxes included in their prices have been removed. The
// lists are aligned with the document's own, with nil entries for rows whose
// total is of no interest.
type includedTaxNetTotals struct {
	cur       currency.Code
	cat       cbc.Code
	lines     []*num.Amount
	discounts []*num.Amount
	charges   []*num.Amount
}

// netTotalSlot points at the entry reserved for a row inside one of the
// net total lists.
type netTotalSlot struct {
	list  []*num.Amount
	index int
}

// prepareIncludedTaxNetTotals works out the total that each of the document's
// rows should have once the taxes included in their prices have been removed.
//
// Only the currency rounding rule needs this: it extracts the tax from the sum
// of each rate's tax-inclusive totals and shares it back over the rows, so
// removing the tax from each price on its own rounds every row down and leaves
// the document's tax bases short by a subunit or two. With precise rounding the
// extra accuracy maintained in the prices already keeps the sums aligned.
func prepareIncludedTaxNetTotals(doc billable, cat cbc.Code) (*includedTaxNetTotals, error) {
	if roundingRule(doc) != tax.RoundingRuleCurrency {
		return nil, nil
	}

	lines := doc.getLines()
	discounts := doc.getDiscounts()
	charges := doc.getCharges()

	nt := &includedTaxNetTotals{
		cur:       doc.GetCurrency(),
		cat:       cat,
		lines:     make([]*num.Amount, len(lines)),
		discounts: make([]*num.Amount, len(discounts)),
		charges:   make([]*num.Amount, len(charges)),
	}

	// Build the same list of taxable rows the calculator works with, keeping
	// track of where each of the results belongs.
	tls := make([]tax.TaxableLine, 0, len(lines)+len(discounts)+len(charges))
	slots := make([]netTotalSlot, 0, cap(tls))
	for i, l := range lines {
		if l == nil || l.Total == nil {
			continue
		}
		tls = append(tls, l)
		slots = append(slots, netTotalSlot{nt.lines, i})
	}
	for i, d := range discounts {
		if d == nil {
			continue
		}
		tls = append(tls, d)
		slots = append(slots, netTotalSlot{nt.discounts, i})
	}
	for i, c := range charges {
		if c == nil {
			continue
		}
		tls = append(tls, c)
		slots = append(slots, netTotalSlot{nt.charges, i})
	}
	if len(tls) == 0 {
		return nil, nil
	}

	date := doc.getValueDate()
	if date == nil {
		id := doc.getIssueDate()
		date = &id
	}
	tc := &tax.TotalCalculator{
		Currency: nt.cur,
		Rounding: tax.RoundingRuleCurrency,
		Country:  doc.RegimeDef().GetCountry(),
		Date:     *date,
		Lines:    tls,
		Includes: cat,
	}
	totals, err := tc.ExtractIncludedTaxes()
	if err != nil {
		return nil, err
	}
	for i, total := range totals {
		slots[i].list[slots[i].index] = &total
	}

	return nt, nil
}

// align nudges the prices and amounts of the document's rows until their totals
// match the net totals, recalculating in between as line discounts and charges
// may shift the results. Line prices are adjusted instead of the line totals
// directly, so that every amount presented can still be recalculated from the
// data alongside it.
func (nt *includedTaxNetTotals) align(doc billable) error {
	if nt == nil {
		return nil
	}
	for range maxIncludedTaxAlignPasses {
		if !nt.nudge(doc) {
			return nil
		}
		if err := calculate(doc); err != nil {
			return err
		}
	}
	return nil
}

// nudge moves each of the document's rows towards its net total, reporting
// whether anything was changed.
func (nt *includedTaxNetTotals) nudge(doc billable) bool {
	changed := false
	for i, l := range doc.getLines() {
		if i < len(nt.lines) && nt.nudgeLine(l, nt.lines[i]) {
			changed = true
		}
	}
	for i, d := range doc.getDiscounts() {
		if i >= len(nt.discounts) || nt.discounts[i] == nil {
			continue
		}
		// Percentages are always recalculated from the document's new sum.
		if d == nil || d.Percent != nil || !hasTaxPercent(d.Taxes, nt.cat) {
			continue
		}
		if a := nt.discounts[i].Negate(); !a.Equals(d.Amount) {
			d.Amount = a
			changed = true
		}
	}
	for i, c := range doc.getCharges() {
		if i >= len(nt.charges) || nt.charges[i] == nil {
			continue
		}
		if c == nil || c.Percent != nil || !hasTaxPercent(c.Taxes, nt.cat) {
			continue
		}
		if a := *nt.charges[i]; !a.Equals(c.Amount) {
			c.Amount = a
			changed = true
		}
	}
	return changed
}

// nudgeLine adjusts the line's net item price so that the line's total matches
// the net total provided, reporting whether the price was changed.
func (nt *includedTaxNetTotals) nudgeLine(l *Line, net *num.Amount) bool {
	if net == nil || l == nil || l.Sum == nil || l.Total == nil ||
		l.Item == nil || l.Item.Price == nil || !hasTaxPercent(l.Taxes, nt.cat) {
		return false
	}
	diff := net.Subtract(*l.Total)
	if diff.IsZero() {
		return false
	}
	// The net total accounts for the line's discounts and charges, but it is
	// the sum that the item's price determines, so shift it by the same
	// difference.
	price := netItemPrice(l.Sum.Add(diff), l.Quantity, nt.cur)
	if price == nil {
		return false
	}
	l.Item.Price = price
	return true
}

// netItemPrice provides the item price that will reproduce the given line sum
// for the quantity, using the least accuracy needed so that the price
// presented can be multiplied back into the sum.
func netItemPrice(sum, quantity num.Amount, cur currency.Code) *num.Amount {
	if quantity.IsZero() {
		return nil
	}
	for extra := uint32(0); extra <= maxIncludedTaxPriceAccuracy; extra++ {
		price := sum.Upscale(extra).Divide(quantity)
		got := tax.ApplyRoundingRule(tax.RoundingRuleCurrency, cur, price.Multiply(quantity))
		if got.Equals(sum) {
			return &price
		}
	}
	return nil
}

// hasTaxPercent determines if the set includes a percentage for the category,
// which is what makes a row's price subject to tax removal.
func hasTaxPercent(ts tax.Set, cat cbc.Code) bool {
	c := ts.Get(cat)
	return c != nil && c.Percent != nil
}
