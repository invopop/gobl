package bill_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// roundingInvoice provides an invoice without taxes included in the prices,
// so that the rounding rule is the only thing under test.
func roundingInvoice(t *testing.T, lines ...*bill.Line) *bill.Invoice {
	t.Helper()
	inv := &bill.Invoice{
		Series:    "TEST",
		Code:      "00123",
		Currency:  currency.EUR,
		IssueDate: cal.MakeDate(2022, 6, 13),
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
		},
		Lines: lines,
	}
	return inv
}

func roundingLine(qty, price num.Amount) *bill.Line {
	return &bill.Line{
		Quantity: qty,
		Item: &org.Item{
			Name:  "Test Item",
			Price: &price,
		},
		Taxes: tax.Set{
			{
				Category: tax.CategoryVAT,
				Percent:  num.NewPercentage(210, 3),
			},
		},
	}
}

func TestRoundToCurrency(t *testing.T) {
	t.Run("line total above the currency's precision", func(t *testing.T) {
		// 3 × 10.555 = 31.665, one decimal more than EUR supports.
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(10555, 3)))
		require.NoError(t, inv.Calculate())
		require.Equal(t, "31.665", inv.Lines[0].Total.String())
		payable := inv.Totals.Payable

		require.NoError(t, inv.RoundToCurrency())

		require.NotNil(t, inv.Tax)
		assert.Equal(t, tax.RoundingRuleCurrency, inv.Tax.Rounding)
		assert.Equal(t, "31.67", inv.Lines[0].Total.String())
		assert.Equal(t, "31.67", inv.Totals.Sum.String())
		assert.Equal(t, payable.String(), inv.Totals.Payable.String(),
			"the amount owed must not move")
		// The unit price keeps its extra decimals, which EN 16931 allows.
		assert.Equal(t, "10.555", inv.Lines[0].Item.Price.String())
	})

	t.Run("an invoice that already fits is left untouched", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(1000, 2)))
		require.NoError(t, inv.Calculate())
		before := inv.Totals.Payable

		require.NoError(t, inv.RoundToCurrency())

		assert.Nil(t, inv.Tax, "no tax object should be invented")
		assert.Nil(t, inv.Totals.Rounding, "no rounding should be invented")
		assert.Equal(t, before.String(), inv.Totals.Payable.String())
	})

	t.Run("an existing rounding total is added to, not replaced", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(10555, 3)))
		existing := num.MakeAmount(-5, 2)
		inv.Totals = &bill.Totals{Rounding: &existing}
		require.NoError(t, inv.Calculate())
		payable := inv.Totals.Payable

		require.NoError(t, inv.RoundToCurrency())

		require.NotNil(t, inv.Totals.Rounding)
		// Rounding the line to 31.67 gains a cent, so the -0.05 the document
		// already carried becomes -0.06. Replacing it would read -0.01 and
		// move the payable total by the 0.05 it dropped.
		assert.Equal(t, "-0.06", inv.Totals.Rounding.String())
		assert.Equal(t, payable.String(), inv.Totals.Payable.String())
	})

	t.Run("is idempotent", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(10555, 3)))
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.RoundToCurrency())
		first, err := json.Marshal(inv)
		require.NoError(t, err)

		require.NoError(t, inv.RoundToCurrency())
		require.NoError(t, inv.RoundToCurrency())
		second, err := json.Marshal(inv)
		require.NoError(t, err)

		assert.JSONEq(t, string(first), string(second))
	})

	t.Run("a document that has not been calculated", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(10555, 3)))
		require.NoError(t, inv.RoundToCurrency())
		assert.Equal(t, "31.67", inv.Lines[0].Total.String())
	})

	t.Run("a discount base above the currency's precision", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(1000, 2)))
		base := num.MakeAmount(101234, 4)
		inv.Discounts = []*bill.Discount{
			{
				Base:    &base,
				Percent: num.NewPercentage(10, 2),
				Reason:  "testing",
			},
		}
		require.NoError(t, inv.Calculate())

		require.NoError(t, inv.RoundToCurrency())

		require.NotNil(t, inv.Discounts[0].Base)
		assert.Equal(t, "10.12", inv.Discounts[0].Base.String())
		assert.Equal(t, "1.01", inv.Discounts[0].Amount.String())
	})

	t.Run("with the tax bypass tag", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(10555, 3)))
		inv.SetTags(tax.TagBypass)
		require.NoError(t, inv.Calculate())

		require.NoError(t, inv.RoundToCurrency())
		assert.Nil(t, inv.Tax, "nothing should be recalculated")
	})

	t.Run("with a currency without subunits", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(10555, 3)))
		inv.Currency = currency.JPY
		require.NoError(t, inv.Calculate())
		payable := inv.Totals.Payable

		require.NoError(t, inv.RoundToCurrency())

		assert.Equal(t, "32", inv.Lines[0].Total.String())
		assert.Equal(t, payable.String(), inv.Totals.Payable.String())
	})
}

func TestRoundToCurrencyWithIncludedTaxes(t *testing.T) {
	// The pairing of RemoveIncludedTaxes and RoundToCurrency must reproduce
	// exactly what removing the taxes under the currency rule used to give.
	lines := make([]*bill.Line, 12)
	for i := range lines {
		lines[i] = &bill.Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "Room rate",
				Price: num.NewAmount(12500, 2),
			},
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Percent:  num.NewPercentage(6, 2),
				},
			},
		}
	}
	inv := baseInvoice(t, lines...)
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.RemoveIncludedTaxes())

	// Removal alone keeps the precise amounts and the original tax total.
	assert.Equal(t, "1415.09", inv.Totals.Sum.String())
	assert.Equal(t, "84.91", inv.Totals.Tax.String())

	require.NoError(t, inv.RoundToCurrency())

	assert.Equal(t, tax.RoundingRuleCurrency, inv.Tax.Rounding)
	assert.Equal(t, "117.9245", inv.Lines[0].Item.Price.String())
	assert.Equal(t, "117.92", inv.Lines[0].Total.String())
	assert.Equal(t, "1415.04", inv.Totals.Sum.String())
	assert.Equal(t, "84.90", inv.Totals.Tax.String())
	assert.Equal(t, "1499.94", inv.Totals.TotalWithTax.String())
	require.NotNil(t, inv.Totals.Rounding)
	assert.Equal(t, "0.06", inv.Totals.Rounding.String())
	assert.Equal(t, "1500.00", inv.Totals.Payable.String())
}
