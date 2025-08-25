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
	Tags     []cbc.Key
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

// Calculate the totals
func (tc *TotalCalculator) Calculate(t *Total) error {
	tc.zero = tc.Currency.Def().Zero()

	// reset
	t.Categories = make([]*CategoryTotal, 0)
	t.Sum = tc.zero

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
	t.Calculate(tc.Currency, tc.Rounding)

	return nil
}

func (tc *TotalCalculator) prepareLines(taxLines []*taxLine) error {
	// First, prepare all tax combos using the country, tags, and date
	for _, tl := range taxLines {
		for _, combo := range tl.taxes {
			if err := combo.calculate(tc.Country, tc.Tags, tc.Date); err != nil {
				return err
			}
			// always add 2 decimal places for all tax calculations
			tl.total = tl.total.RescaleUp(tc.zero.Exp() + 2)
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
			if c.informative {
				return ErrInvalidPricesInclude.WithMessage("cannot include informative category '%s'", tc.Includes.String())
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
			rt := t.rateTotalFor(c, tc.zero)
			rt.Base = matchRoundingPrecision(tc.Rounding, rt.Base, tl.total)
			rt.Base = rt.Base.Add(tl.total)
		}
	}
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
