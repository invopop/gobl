package bill

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

const (
	// linePrecisionExtra is the number of decimal places added to the
	// currency's own when using the `precise` rounding rule.
	linePrecisionExtra uint32 = 2

	// defaultTaxRemovalAccuracy is the number of decimal places added to a
	// price before dividing out the tax included in it.
	defaultTaxRemovalAccuracy uint32 = 2
)

// RoundToCurrency recalculates the invoice using the `currency` rounding rule
// so that every amount fits the number of decimal places supported by the
// currency, as required by formats such as UBL or CII. Any difference in the
// amount payable is carried in the totals' rounding amount, known as BT-114 in
// EN 16931.
//
// Documents already within the currency's precision are left untouched.
//
// This method will replace the invoice contents in place, or return an error.
func (inv *Invoice) RoundToCurrency() error {
	return roundToCurrency(inv)
}

func roundToCurrency(doc billable) error {
	if doc.HasTags(tax.TagBypass) {
		// Calculations are skipped entirely.
		return nil
	}
	if err := ensureCalculated(doc); err != nil {
		return err
	}
	if doc.getTotals() == nil {
		// Nothing priced, such as an order or delivery without amounts.
		return nil
	}
	cd := doc.GetCurrency().Def()
	if cd == nil || !exceedsCurrencyPrecision(doc, cd) {
		return nil
	}
	t := doc.getTotals()
	payable := t.Payable

	// Bases and the rounding amount are provided externally and never
	// reduced by the calculator.
	rescaleBases(doc, cd)
	if t.Rounding != nil {
		r := cd.Rescale(*t.Rounding)
		t.Rounding = &r
	}

	tx := doc.getTax()
	if tx == nil {
		tx = new(Tax)
		doc.setTax(tx)
	}
	tx.Rounding = tax.RoundingRuleCurrency

	return recalculateKeepingPayable(doc, payable)
}

// exceedsCurrencyPrecision reports whether the document holds an amount with
// more decimal places than the currency supports. Unit prices are exempt, and
// the document totals are always rounded to the currency.
func exceedsCurrencyPrecision(doc billable, cd *currency.Def) bool {
	over := overCurrencyPrecision(cd)
	for _, l := range doc.getLines() {
		if lineExceedsCurrencyPrecision(l, over) {
			return true
		}
	}
	for _, d := range doc.getDiscounts() {
		if d != nil && (over(&d.Amount) || over(d.Base)) {
			return true
		}
	}
	for _, c := range doc.getCharges() {
		if c != nil && (over(&c.Amount) || over(c.Base)) {
			return true
		}
	}
	t := doc.getTotals()
	if over(t.Rounding) {
		return true
	}
	// Only reachable for totals assembled by hand.
	return taxTotalExceedsCurrencyPrecision(t.Taxes, over)
}

// overCurrencyPrecision reports if an amount has more decimal places than the
// currency supports.
func overCurrencyPrecision(cd *currency.Def) func(*num.Amount) bool {
	return func(a *num.Amount) bool {
		return a != nil && a.Exp() > cd.Subunits
	}
}

func lineExceedsCurrencyPrecision(l *Line, over func(*num.Amount) bool) bool {
	if l == nil {
		return false
	}
	if over(l.Sum) || over(l.Total) {
		return true
	}
	if lineAdjustmentsExceedCurrencyPrecision(l.Discounts, l.Charges, over) {
		return true
	}
	for _, sl := range l.Breakdown {
		if subLineExceedsCurrencyPrecision(sl, over) {
			return true
		}
	}
	for _, sl := range l.Substituted {
		if subLineExceedsCurrencyPrecision(sl, over) {
			return true
		}
	}
	return false
}

func subLineExceedsCurrencyPrecision(sl *SubLine, over func(*num.Amount) bool) bool {
	if sl == nil {
		return false
	}
	if over(sl.Sum) || over(sl.Total) {
		return true
	}
	return lineAdjustmentsExceedCurrencyPrecision(sl.Discounts, sl.Charges, over)
}

func lineAdjustmentsExceedCurrencyPrecision(discounts []*LineDiscount, charges []*LineCharge, over func(*num.Amount) bool) bool {
	for _, d := range discounts {
		if d != nil && (over(&d.Amount) || over(d.Base)) {
			return true
		}
	}
	for _, c := range charges {
		if c != nil && (over(&c.Amount) || over(c.Base)) {
			return true
		}
	}
	return false
}

func taxTotalExceedsCurrencyPrecision(t *tax.Total, over func(*num.Amount) bool) bool {
	if t == nil {
		return false
	}
	for _, ct := range t.Categories {
		if over(ct.Surcharge) {
			return true
		}
		for _, rt := range ct.Rates {
			if over(&rt.Base) || over(&rt.Amount) {
				return true
			}
			if rt.Surcharge != nil && over(&rt.Surcharge.Amount) {
				return true
			}
		}
	}
	return false
}

// rescaleBases rounds the base amounts used for percentage calculations to the
// currency's precision.
func rescaleBases(doc billable, cd *currency.Def) {
	for _, l := range doc.getLines() {
		if l == nil {
			continue
		}
		rescaleLineBases(l.Discounts, l.Charges, cd)
		for _, sl := range l.Breakdown {
			if sl != nil {
				rescaleLineBases(sl.Discounts, sl.Charges, cd)
			}
		}
		for _, sl := range l.Substituted {
			if sl != nil {
				rescaleLineBases(sl.Discounts, sl.Charges, cd)
			}
		}
	}
	for _, d := range doc.getDiscounts() {
		if d != nil {
			d.Base = rescaledBase(d.Base, cd)
		}
	}
	for _, c := range doc.getCharges() {
		if c != nil {
			c.Base = rescaledBase(c.Base, cd)
		}
	}
}

func rescaleLineBases(discounts []*LineDiscount, charges []*LineCharge, cd *currency.Def) {
	for _, d := range discounts {
		if d != nil {
			d.Base = rescaledBase(d.Base, cd)
		}
	}
	for _, c := range charges {
		if c != nil {
			c.Base = rescaledBase(c.Base, cd)
		}
	}
}

func rescaledBase(a *num.Amount, cd *currency.Def) *num.Amount {
	if a == nil {
		return nil
	}
	b := cd.Rescale(*a)
	return &b
}

// roundingRule determines the rounding rule to apply to the document, either
// explicitly defined in the tax object, or the one provided by the regime.
func roundingRule(doc billable) cbc.Key {
	if tx := doc.getTax(); tx != nil && tx.Rounding != "" {
		return tx.Rounding
	}
	return doc.RegimeDef().GetRoundingRule()
}

// recalculateKeepingPayable recalculates the document, carrying any change in
// the amount payable into the totals' rounding amount.
func recalculateKeepingPayable(doc billable, payable num.Amount) error {
	if err := calculate(doc); err != nil {
		return err
	}
	t := doc.getTotals()
	if t == nil {
		return nil
	}
	diff := payable.Subtract(t.Payable)
	if diff.IsZero() {
		return nil
	}
	// Add to any rounding amount already present, which is included in both
	// payable amounts.
	rnd := diff
	if t.Rounding != nil {
		rnd = t.Rounding.MatchPrecision(diff).Add(diff)
	}
	t.Rounding = &rnd
	return calculate(doc)
}
