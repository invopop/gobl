package tax

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
)

// TotalCalculator defines the base structure with the available
// data for calculating tax totals from TaxableLines.
type TotalCalculator struct {
	Country  l10n.TaxCountryCode
	Rounding cbc.Key
	Currency currency.Code
	Date     cal.Date
	Lines    []TaxableLine
	Includes cbc.Code // Tax included in price

	zero num.Amount
}

// TaxableLine defines what we expect from a line in order to subsequently calculate
// the taxes that need to be added or retained.
type TaxableLine interface {
	GetTaxes() Set
	GetTotal() num.Amount
}

// Calculate the totals according to the rounding rule, either using the
// currency's precision for every step, or maintaining precision until
// the final sums.
func (tc *TotalCalculator) Calculate(t *Total) error {
	tc.zero = tc.Currency.Def().Zero()

	// reset
	t.Categories = make([]*CategoryTotal, 0)
	t.Sum = tc.zero

	// get simplified list of lines
	taxLines := mapTaxLines(tc.Lines)
	if err := tc.prepareCombos(taxLines); err != nil {
		return err
	}

	if tc.Rounding == RoundingRuleCurrency {
		return tc.calculateCurrencyTotals(t, taxLines)
	}
	return tc.calculatePreciseTotals(t, taxLines)
}

// prepareCombos ensures the tax combos of every line have been calculated
// using the country and date.
func (tc *TotalCalculator) prepareCombos(taxLines []*taxLine) error {
	for _, tl := range taxLines {
		for _, combo := range tl.taxes {
			if err := combo.calculate(tc.Country, tc.Date); err != nil {
				return err
			}
		}
	}
	return nil
}

// checkIncludedCombo ensures the combo's category can be used for
// taxes included in prices.
func (tc *TotalCalculator) checkIncludedCombo(c *Combo) error {
	if c.retained {
		return ErrInvalidPricesInclude.WithMessage("cannot include retained category '%s'", tc.Includes.String())
	}
	if c.informative {
		return ErrInvalidPricesInclude.WithMessage("cannot include informative category '%s'", tc.Includes.String())
	}
	return nil
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
