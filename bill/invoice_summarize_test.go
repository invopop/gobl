package bill_test

import (
	"fmt"
	"testing"

	"github.com/LastPossum/kamino"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// summaryTestLines prepares n lines with awkward prices, alternating between
// the general and reduced VAT rates, with a discount on every third line.
func summaryTestLines(n int) []*bill.Line {
	lines := make([]*bill.Line, n)
	for i := range lines {
		rate := tax.RateGeneral
		if i%2 == 1 {
			rate = tax.RateReduced
		}
		l := &bill.Line{
			Quantity: num.MakeAmount(int64(i%7+1), 0),
			Item: &org.Item{
				Name:  fmt.Sprintf("Delivery %d", i),
				Price: num.NewAmount(int64(1234+i*37), 3),
			},
			Taxes: tax.Set{
				{Category: tax.CategoryVAT, Rate: rate},
			},
		}
		if i%3 == 0 {
			l.Discounts = []*bill.LineDiscount{
				{Percent: num.NewPercentage(75, 3)},
			}
		}
		lines[i] = l
	}
	return lines
}

func assertSummarizedTotals(t *testing.T, orig, inv *bill.Invoice) {
	t.Helper()
	assert.Equal(t, orig.Totals.Sum.String(), inv.Totals.Sum.String())
	assert.Equal(t, orig.Totals.Total.String(), inv.Totals.Total.String())
	assert.Equal(t, orig.Totals.Tax.String(), inv.Totals.Tax.String())
	assert.Equal(t, orig.Totals.TotalWithTax.String(), inv.Totals.TotalWithTax.String())
	assert.Equal(t, orig.Totals.Payable.String(), inv.Totals.Payable.String())
}

func TestInvoiceSummarizeLines(t *testing.T) {
	t.Run("groups lines by tax rate", func(t *testing.T) {
		inv := baseInvoice(t, summaryTestLines(1000)...)
		inv.Tax = nil
		inv.Discounts = []*bill.Discount{
			{Reason: "Volume", Percent: num.NewPercentage(5, 2)},
		}
		require.NoError(t, inv.Calculate())
		orig, err := kamino.Clone(inv)
		require.NoError(t, err)

		require.NoError(t, inv.SummarizeLines())

		require.Len(t, inv.Lines, 2)
		assert.Equal(t, "Summary of 500 lines", inv.Lines[0].Item.Name)
		assert.Equal(t, tax.RateGeneral, inv.Lines[0].Taxes[0].Rate)
		assert.Equal(t, "1", inv.Lines[0].Quantity.String())
		assert.Equal(t, tax.RateReduced, inv.Lines[1].Taxes[0].Rate)
		assertSummarizedTotals(t, orig, inv)
		assert.Equal(t, orig.Totals.Discount.String(), inv.Totals.Discount.String())
		assert.True(t, inv.HasTags(tax.TagSummarized))
		require.Len(t, inv.Notes, 1)
		assert.Equal(t, "Lines summarized from 1000 original lines.", inv.Notes[0].Text)
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("prices including tax", func(t *testing.T) {
		inv := baseInvoice(t, summaryTestLines(999)...)
		require.NoError(t, inv.Calculate())
		orig, err := kamino.Clone(inv)
		require.NoError(t, err)

		require.NoError(t, inv.SummarizeLines())

		require.Len(t, inv.Lines, 2)
		assertSummarizedTotals(t, orig, inv)
		assert.Equal(t, orig.Totals.TaxIncluded.String(), inv.Totals.TaxIncluded.String())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("currency rounding", func(t *testing.T) {
		inv := baseInvoice(t, summaryTestLines(500)...)
		inv.Tax = &bill.Tax{Rounding: tax.RoundingRuleCurrency}
		require.NoError(t, inv.Calculate())
		orig, err := kamino.Clone(inv)
		require.NoError(t, err)

		require.NoError(t, inv.SummarizeLines())

		require.Len(t, inv.Lines, 2)
		assertSummarizedTotals(t, orig, inv)
	})

	t.Run("retained taxes", func(t *testing.T) {
		lines := summaryTestLines(300)
		for i, l := range lines {
			if i%5 == 0 {
				l.Taxes = append(l.Taxes, &tax.Combo{Category: "IRPF", Key: "pro"})
			}
		}
		inv := baseInvoice(t, lines...)
		inv.Tax = nil
		require.NoError(t, inv.Calculate())
		orig, err := kamino.Clone(inv)
		require.NoError(t, err)

		require.NoError(t, inv.SummarizeLines())

		// general and reduced, each with and without IRPF
		require.Len(t, inv.Lines, 4)
		assertSummarizedTotals(t, orig, inv)
		assert.Equal(t, orig.Totals.RetainedTax.String(), inv.Totals.RetainedTax.String())
	})

	t.Run("keeps item keys and extensions apart", func(t *testing.T) {
		lines := summaryTestLines(10)
		for i, l := range lines {
			l.Taxes = tax.Set{{Category: tax.CategoryVAT, Rate: tax.RateGeneral}}
			if i%2 == 0 {
				l.Item.Key = "services"
			}
		}
		inv := baseInvoice(t, lines...)
		inv.Tax = nil
		require.NoError(t, inv.Calculate())

		require.NoError(t, inv.SummarizeLines())

		require.Len(t, inv.Lines, 2)
		assert.Equal(t, cbc.Key("services"), inv.Lines[0].Item.Key)
		assert.Empty(t, inv.Lines[1].Item.Key)
	})

	t.Run("credit lines", func(t *testing.T) {
		lines := summaryTestLines(20)
		lines[3].Quantity = num.MakeAmount(-50, 0)
		inv := baseInvoice(t, lines...)
		inv.Tax = nil
		require.NoError(t, inv.Calculate())
		orig, err := kamino.Clone(inv)
		require.NoError(t, err)

		require.NoError(t, inv.SummarizeLines())

		assertSummarizedTotals(t, orig, inv)
	})

	t.Run("calculates the invoice first", func(t *testing.T) {
		inv := baseInvoice(t, summaryTestLines(10)...)
		require.NoError(t, inv.SummarizeLines())
		require.Len(t, inv.Lines, 2)
		assert.NotNil(t, inv.Totals)
	})

	t.Run("already summarized", func(t *testing.T) {
		inv := baseInvoice(t, summaryTestLines(10)...)
		require.NoError(t, inv.SummarizeLines())
		assert.ErrorContains(t, inv.SummarizeLines(), "invoice lines already summarized")
	})

	t.Run("bypass", func(t *testing.T) {
		inv := baseInvoice(t, summaryTestLines(10)...)
		inv.SetTags(tax.TagBypass)
		assert.ErrorContains(t, inv.SummarizeLines(), "cannot summarize an invoice with tag bypass")
	})
}
