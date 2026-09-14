package bill

import (
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
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
	payable := doc.getTotals().Payable

	// Bases are provided externally and never reduced by the calculator.
	rescaleBases(doc, cd)

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
	// Only reachable for totals assembled by hand.
	return taxTotalExceedsCurrencyPrecision(doc.getTotals().Taxes, over)
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
	for _, sl := range l.Breakdown {
		if sl != nil && (over(sl.Sum) || over(sl.Total)) {
			return true
		}
	}
	for _, d := range l.Discounts {
		if d != nil && (over(&d.Amount) || over(d.Base)) {
			return true
		}
	}
	for _, c := range l.Charges {
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
		for _, rt := range ct.Rates {
			if over(&rt.Base) || over(&rt.Amount) {
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
		for _, d := range l.Discounts {
			if d != nil {
				d.Base = rescaledBase(d.Base, cd)
			}
		}
		for _, c := range l.Charges {
			if c != nil {
				c.Base = rescaledBase(c.Base, cd)
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

func rescaledBase(a *num.Amount, cd *currency.Def) *num.Amount {
	if a == nil {
		return nil
	}
	b := cd.Rescale(*a)
	return &b
}
