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

	t.Run("a rounding total above the currency's precision", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(1000, 2)))
		existing := num.MakeAmount(5, 3)
		inv.Totals = &bill.Totals{Rounding: &existing}
		require.NoError(t, inv.Calculate())
		payable := inv.Totals.Payable

		require.NoError(t, inv.RoundToCurrency())

		require.NotNil(t, inv.Totals.Rounding)
		// The rounding total is the only amount out of range, so nothing else
		// moves and the payable total already absorbed the half cent.
		assert.Equal(t, "0.01", inv.Totals.Rounding.String())
		assert.Equal(t, payable.String(), inv.Totals.Payable.String())
	})

	t.Run("an existing rounding total without decimals", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(10555, 3)))
		existing := num.MakeAmount(-5, 0)
		inv.Totals = &bill.Totals{Rounding: &existing}
		require.NoError(t, inv.Calculate())
		payable := inv.Totals.Payable

		require.NoError(t, inv.RoundToCurrency())

		require.NotNil(t, inv.Totals.Rounding)
		// The cent gained by the line must survive being added to a rounding
		// total written without decimals.
		assert.Equal(t, "-5.01", inv.Totals.Rounding.String())
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

	t.Run("totals that hold only a rounding amount", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(10555, 3)))
		rnd := num.MakeAmount(-5, 2)
		inv.Totals = &bill.Totals{Rounding: &rnd}

		require.NoError(t, inv.RoundToCurrency())

		// The rounding amount is the only total a caller supplies, so these
		// totals still need calculating before anything can be read from them.
		assert.Equal(t, "31.67", inv.Lines[0].Total.String())
		assert.Equal(t, "-0.06", inv.Totals.Rounding.String())
	})

	t.Run("totals that hold only an amount payable", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(10555, 3)))
		inv.Totals = &bill.Totals{Payable: num.MakeAmount(10000, 2)}

		require.NoError(t, inv.RoundToCurrency())

		// Lines without a total mean nothing has been calculated, whatever the
		// totals claim to hold.
		assert.Equal(t, "31.67", inv.Lines[0].Total.String())
		assert.Equal(t, "38.31", inv.Totals.Payable.String())
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

func TestRoundToCurrencyPrecisionChecks(t *testing.T) {
	t.Run("a line discount base above the currency's precision", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(1000, 2)))
		base := num.MakeAmount(101234, 4)
		inv.Lines[0].Discounts = []*bill.LineDiscount{
			{
				Base:    &base,
				Percent: num.NewPercentage(10, 2),
				Reason:  "testing",
			},
		}
		require.NoError(t, inv.Calculate())

		require.NoError(t, inv.RoundToCurrency())

		require.NotNil(t, inv.Lines[0].Discounts[0].Base)
		assert.Equal(t, "10.12", inv.Lines[0].Discounts[0].Base.String())
	})

	t.Run("a line charge base above the currency's precision", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(1000, 2)))
		base := num.MakeAmount(101234, 4)
		inv.Lines[0].Charges = []*bill.LineCharge{
			{
				Base:    &base,
				Percent: num.NewPercentage(10, 2),
				Reason:  "testing",
			},
		}
		require.NoError(t, inv.Calculate())

		require.NoError(t, inv.RoundToCurrency())

		require.NotNil(t, inv.Lines[0].Charges[0].Base)
		assert.Equal(t, "10.12", inv.Lines[0].Charges[0].Base.String())
	})

	t.Run("a document charge base above the currency's precision", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(1000, 2)))
		base := num.MakeAmount(101234, 4)
		inv.Charges = []*bill.Charge{
			{
				Base:    &base,
				Percent: num.NewPercentage(10, 2),
				Reason:  "testing",
			},
		}
		require.NoError(t, inv.Calculate())

		require.NoError(t, inv.RoundToCurrency())

		require.NotNil(t, inv.Charges[0].Base)
		assert.Equal(t, "10.12", inv.Charges[0].Base.String())
	})

	t.Run("a breakdown sub-line above the currency's precision", func(t *testing.T) {
		price := num.MakeAmount(10555, 3)
		inv := roundingInvoice(t, &bill.Line{
			Quantity: num.MakeAmount(1, 0),
			Item:     &org.Item{Name: "Grouped item"},
			Breakdown: []*bill.SubLine{
				{
					Quantity: num.MakeAmount(3, 0),
					Item: &org.Item{
						Name:  "Sub item",
						Price: &price,
					},
				},
			},
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Percent:  num.NewPercentage(210, 3),
				},
			},
		})
		require.NoError(t, inv.Calculate())
		require.Equal(t, "31.665", inv.Lines[0].Breakdown[0].Total.String())
		payable := inv.Totals.Payable

		require.NoError(t, inv.RoundToCurrency())

		assert.Equal(t, tax.RoundingRuleCurrency, inv.Tax.Rounding)
		assert.Equal(t, payable.String(), inv.Totals.Payable.String())
	})

	t.Run("a sub-line discount base above the currency's precision", func(t *testing.T) {
		price := num.MakeAmount(1000, 2)
		base := num.MakeAmount(10005, 4)
		inv := roundingInvoice(t, &bill.Line{
			Quantity: num.MakeAmount(1, 0),
			Item:     &org.Item{Name: "Grouped item"},
			Breakdown: []*bill.SubLine{
				{
					Quantity: num.MakeAmount(3, 0),
					Item: &org.Item{
						Name:  "Sub item",
						Price: &price,
					},
					Discounts: []*bill.LineDiscount{
						{
							Base:    &base,
							Percent: num.NewPercentage(10, 2),
							Reason:  "testing",
						},
					},
				},
			},
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Percent:  num.NewPercentage(210, 3),
				},
			},
		})
		require.NoError(t, inv.Calculate())
		require.Equal(t, "0.1001", inv.Lines[0].Breakdown[0].Discounts[0].Amount.String())

		require.NoError(t, inv.RoundToCurrency())

		d := inv.Lines[0].Breakdown[0].Discounts[0]
		require.NotNil(t, d.Base)
		assert.Equal(t, "1.00", d.Base.String())
		assert.Equal(t, "0.10", d.Amount.String())
	})

	t.Run("a fixed line discount above the currency's precision", func(t *testing.T) {
		line := roundingLine(num.MakeAmount(3, 0), num.MakeAmount(10555, 3))
		line.Discounts = []*bill.LineDiscount{
			{
				Amount: num.MakeAmount(5, 3),
				Reason: "testing",
			},
		}
		inv := roundingInvoice(t, line)
		require.NoError(t, inv.Calculate())
		payable := inv.Totals.Payable

		require.NoError(t, inv.RoundToCurrency())

		// Fixed amounts are never reduced by the calculator on their own.
		assert.Equal(t, "0.01", inv.Lines[0].Discounts[0].Amount.String())
		assert.Equal(t, "31.66", inv.Lines[0].Total.String())
		assert.Equal(t, payable.String(), inv.Totals.Payable.String())
	})

	t.Run("an invoice without any lines", func(t *testing.T) {
		inv := roundingInvoice(t)
		require.NoError(t, inv.RoundToCurrency())
		assert.Nil(t, inv.Totals)
	})

	t.Run("an invoice that cannot be calculated", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(10555, 3)))
		inv.Supplier = nil
		inv.Currency = currency.CodeEmpty
		assert.Error(t, inv.RoundToCurrency())
	})

	t.Run("a recalculation that cannot resolve its taxes is reported", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(10555, 3)))
		require.NoError(t, inv.Calculate())
		// A rate key the regime does not define: the recalculation fails, and
		// the failure must surface rather than a half-rounded document.
		inv.Lines[0].Taxes = tax.Set{{Category: tax.CategoryVAT, Rate: "nonsense"}}

		assert.Error(t, inv.RoundToCurrency())
	})
}

func TestRoundToCurrencyWithoutTaxes(t *testing.T) {
	// No taxes at all, so the totals carry no tax breakdown to inspect.
	inv := roundingInvoice(t, &bill.Line{
		Quantity: num.MakeAmount(3, 0),
		Item: &org.Item{
			Name:  "Test Item",
			Price: num.NewAmount(10555, 3),
		},
	})
	require.NoError(t, inv.Calculate())
	require.Nil(t, inv.Totals.Taxes)
	payable := inv.Totals.Payable

	require.NoError(t, inv.RoundToCurrency())

	assert.Equal(t, "31.67", inv.Lines[0].Total.String())
	assert.Equal(t, payable.String(), inv.Totals.Payable.String())
}

func TestRoundToCurrencyWithoutTaxesInPrecision(t *testing.T) {
	// Every amount fits, so the check falls through to the tax totals, which
	// an invoice without taxes does not have.
	inv := roundingInvoice(t, &bill.Line{
		Quantity: num.MakeAmount(3, 0),
		Item: &org.Item{
			Name:  "Test Item",
			Price: num.NewAmount(1000, 2),
		},
	})
	require.NoError(t, inv.Calculate())
	require.Nil(t, inv.Totals.Taxes)

	require.NoError(t, inv.RoundToCurrency())

	assert.Nil(t, inv.Tax, "nothing should be recalculated")
}

func TestRoundToCurrencyWithoutDiscountBase(t *testing.T) {
	inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(10555, 3)))
	inv.Discounts = []*bill.Discount{
		{
			Amount: num.MakeAmount(100, 2),
			Reason: "no base",
		},
	}
	inv.Charges = []*bill.Charge{
		{
			Amount: num.MakeAmount(200, 2),
			Reason: "no base",
		},
	}
	require.NoError(t, inv.Calculate())
	payable := inv.Totals.Payable

	require.NoError(t, inv.RoundToCurrency())

	assert.Nil(t, inv.Discounts[0].Base)
	assert.Nil(t, inv.Charges[0].Base)
	assert.Equal(t, payable.String(), inv.Totals.Payable.String())
}

// TestRoundToCurrencyModifiedAfterCalculation covers the guards for documents
// whose contents were altered after they were calculated, which the regular
// calculation would otherwise have normalized away.
func TestRoundToCurrencyModifiedAfterCalculation(t *testing.T) {
	t.Run("with a nil line", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(10555, 3)))
		require.NoError(t, inv.Calculate())
		// First, so the loop reaches it before a line that exceeds.
		inv.Lines = append([]*bill.Line{nil}, inv.Lines...)
		payable := inv.Totals.Payable

		require.NoError(t, inv.RoundToCurrency())

		assert.Equal(t, payable.String(), inv.Totals.Payable.String())
	})

	t.Run("with a breakdown sub-line above the currency's precision", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(1000, 2)))
		require.NoError(t, inv.Calculate())
		require.Equal(t, "30.00", inv.Lines[0].Total.String())

		total := num.MakeAmount(300001, 4)
		inv.Lines[0].Breakdown = []*bill.SubLine{
			{
				Quantity: num.MakeAmount(1, 0),
				Item:     &org.Item{Name: "Sub item", Price: num.NewAmount(300001, 4)},
				Total:    &total,
			},
		}

		require.NoError(t, inv.RoundToCurrency())

		assert.Equal(t, tax.RoundingRuleCurrency, inv.Tax.Rounding)
	})

	t.Run("with a substituted sub-line above the currency's precision", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(1000, 2)))
		require.NoError(t, inv.Calculate())

		total := num.MakeAmount(300001, 4)
		inv.Lines[0].Substituted = []*bill.SubLine{
			{
				Quantity: num.MakeAmount(1, 0),
				Item:     &org.Item{Name: "Substituted item", Price: num.NewAmount(300001, 4)},
				Total:    &total,
			},
		}

		require.NoError(t, inv.RoundToCurrency())

		assert.Equal(t, tax.RoundingRuleCurrency, inv.Tax.Rounding)
		assert.Equal(t, "30.00", inv.Lines[0].Substituted[0].Total.String())
	})

	t.Run("with a rate surcharge above the currency's precision", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(1000, 2)))
		require.NoError(t, inv.Calculate())
		require.NotNil(t, inv.Totals.Taxes)

		rt := inv.Totals.Taxes.Categories[0].Rates[0]
		rt.Surcharge = &tax.RateTotalSurcharge{
			Percent: num.MakePercentage(52, 3),
			Amount:  num.MakeAmount(15600, 4),
		}
		payable := inv.Totals.Payable

		require.NoError(t, inv.RoundToCurrency())

		assert.Equal(t, tax.RoundingRuleCurrency, inv.Tax.Rounding)
		assert.Equal(t, payable.String(), inv.Totals.Payable.String())
	})

	t.Run("with a category surcharge above the currency's precision", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(1000, 2)))
		require.NoError(t, inv.Calculate())
		require.NotNil(t, inv.Totals.Taxes)

		surcharge := num.MakeAmount(15600, 4)
		inv.Totals.Taxes.Categories[0].Surcharge = &surcharge
		payable := inv.Totals.Payable

		require.NoError(t, inv.RoundToCurrency())

		assert.Equal(t, tax.RoundingRuleCurrency, inv.Tax.Rounding)
		assert.Equal(t, payable.String(), inv.Totals.Payable.String())
	})

	t.Run("with an aggregate tax sum above the currency's precision", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(1000, 2)))
		require.NoError(t, inv.Calculate())
		require.NotNil(t, inv.Totals.Taxes)

		inv.Totals.Taxes.Sum = num.MakeAmount(63001, 4)
		payable := inv.Totals.Payable

		require.NoError(t, inv.RoundToCurrency())

		assert.Equal(t, tax.RoundingRuleCurrency, inv.Tax.Rounding)
		assert.Equal(t, "6.30", inv.Totals.Taxes.Sum.String())
		assert.Equal(t, payable.String(), inv.Totals.Payable.String())
	})

	t.Run("with a retained tax total above the currency's precision", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(1000, 2)))
		require.NoError(t, inv.Calculate())
		require.NotNil(t, inv.Totals.Taxes)

		retained := num.MakeAmount(15001, 4)
		inv.Totals.Taxes.Retained = &retained
		payable := inv.Totals.Payable

		require.NoError(t, inv.RoundToCurrency())

		assert.Equal(t, tax.RoundingRuleCurrency, inv.Tax.Rounding)
		assert.Equal(t, payable.String(), inv.Totals.Payable.String())
	})

	t.Run("with a category amount above the currency's precision", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(1000, 2)))
		require.NoError(t, inv.Calculate())
		require.NotNil(t, inv.Totals.Taxes)

		inv.Totals.Taxes.Categories[0].Amount = num.MakeAmount(63001, 4)
		payable := inv.Totals.Payable

		require.NoError(t, inv.RoundToCurrency())

		assert.Equal(t, tax.RoundingRuleCurrency, inv.Tax.Rounding)
		assert.Equal(t, "6.30", inv.Totals.Taxes.Categories[0].Amount.String())
		assert.Equal(t, payable.String(), inv.Totals.Payable.String())
	})

	t.Run("with nil nested entries", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(1000, 2)))
		require.NoError(t, inv.Calculate())
		// Normalization prunes nil entries, so only a document modified
		// afterwards can hold them. Every amount here fits the currency, so
		// the scan walks them without triggering a recalculation.
		line := inv.Lines[0]
		line.Discounts = []*bill.LineDiscount{nil}
		line.Charges = []*bill.LineCharge{nil}
		line.Breakdown = []*bill.SubLine{nil}
		line.Substituted = []*bill.SubLine{nil}
		inv.Lines = append([]*bill.Line{nil}, inv.Lines...)
		inv.Discounts = []*bill.Discount{nil}
		inv.Charges = []*bill.Charge{nil}
		payable := inv.Totals.Payable

		require.NoError(t, inv.RoundToCurrency())

		assert.Nil(t, inv.Tax)
		assert.Equal(t, payable.String(), inv.Totals.Payable.String())
	})

	t.Run("with a tax rate base above the currency's precision", func(t *testing.T) {
		inv := roundingInvoice(t, roundingLine(num.MakeAmount(3, 0), num.MakeAmount(1000, 2)))
		require.NoError(t, inv.Calculate())
		require.NotNil(t, inv.Totals.Taxes)

		rt := inv.Totals.Taxes.Categories[0].Rates[0]
		rt.Base = num.MakeAmount(300001, 4)
		payable := inv.Totals.Payable

		require.NoError(t, inv.RoundToCurrency())

		assert.Equal(t, tax.RoundingRuleCurrency, inv.Tax.Rounding)
		assert.Equal(t, payable.String(), inv.Totals.Payable.String())
	})
}
