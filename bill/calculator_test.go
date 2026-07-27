package bill_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/br"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: many calculation tests are distributed throughout this package.

func TestCalculate(t *testing.T) {
	t.Run("with round-then-sum rounding rule", func(t *testing.T) {
		inv := baseInvoice(t, &bill.Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "test item 1",
				Price: num.NewAmount(942, 2),
			},
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Percent:  num.NewPercentage(24, 2),
				},
			},
		}, &bill.Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "test item 2",
				Price: num.NewAmount(942, 2),
			},
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Percent:  num.NewPercentage(13, 2),
				},
			},
		})
		inv.Tax.PricesInclude = ""
		inv.Tax.Rounding = tax.RoundingRuleCurrency
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "3.48", inv.Totals.Tax.String())

		inv.Tax.Rounding = tax.RoundingRulePrecise
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "3.49", inv.Totals.Tax.String())
	})
	t.Run("with line errors", func(t *testing.T) {
		inv := baseInvoice(t, &bill.Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:     "test item 1",
				Currency: "USD",
				Price:    num.NewAmount(942, 2),
			},
		})
		require.ErrorContains(t, inv.Calculate(), "lines: 0: item: no exchange rate found from 'USD' to 'EUR'")
	})
	t.Run("with preceding docs and taxes", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Preceding = []*org.DocumentRef{
			{
				Code:      "ABC",
				IssueDate: cal.NewDate(2022, 11, 6),
				Currency:  currency.EUR,
				Tax: &tax.Total{
					Categories: []*tax.CategoryTotal{
						{
							Code: tax.CategoryVAT,
							Rates: []*tax.RateTotal{
								{
									Base:    num.MakeAmount(1000, 2),
									Percent: num.NewPercentage(21, 2),
								},
							},
						},
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "2.10", inv.Preceding[0].Tax.Sum.String())
	})
	t.Run("with nil preceding", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Preceding = []*org.DocumentRef{nil}
		require.NoError(t, inv.Calculate())
	})

	t.Run("update issue date and time", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.IssueDate = cal.MakeDate(2022, 11, 6)
		inv.IssueTime = cal.NewTime(0, 0, 0)
		require.NoError(t, inv.Calculate())
		tn := cal.ThisSecondIn(inv.RegimeDef().TimeLocation())
		assert.Equal(t, tn.Date().String(), inv.IssueDate.String())
		assert.Equal(t, tn.Time().String(), inv.IssueTime.String())
	})

	t.Run("with retained taxes", func(t *testing.T) {
		inv := baseInvoice(t, &bill.Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "test item 1",
				Price: num.NewAmount(942, 2),
			},
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Percent:  num.NewPercentage(21, 2),
				},
				{
					Category: es.TaxCategoryIRPF,
					Percent:  num.NewPercentage(15, 2),
				},
			},
		})
		inv.Tax.PricesInclude = ""
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "1.98", inv.Totals.Tax.String())
		assert.Equal(t, "1.41", inv.Totals.RetainedTax.String())
		assert.Equal(t, "9.99", inv.Totals.Payable.String())
	})

	t.Run("with advances and rounding", func(t *testing.T) {
		inv := baseInvoice(t, &bill.Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "test item 1",
				Price: num.NewAmount(90005, 3),
			},
		})
		inv.Tax.PricesInclude = ""
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Record{
				{
					Amount: num.MakeAmount(9001, 2),
				},
			},
		}
		require.NoError(t, inv.Calculate())
		data, _ := json.MarshalIndent(inv.Totals, "", "  ")
		t.Logf("TOTALS: %s", string(data))
		assert.Equal(t, "90.01", inv.Totals.Payable.String())
		assert.Equal(t, "0.00", inv.Totals.Due.String())
	})

	t.Run("with precision advances, calculated twice", func(t *testing.T) {
		inv := baseInvoice(t, &bill.Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "test item 1",
				Price: num.NewAmount(90005, 3),
			},
		})
		inv.Tax.PricesInclude = ""
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Record{
				{
					Amount: num.MakeAmount(900050, 4),
				},
			},
		}
		require.NoError(t, inv.Calculate())
		data, _ := json.MarshalIndent(inv.Totals, "", "  ")
		t.Logf("TOTALS: %s", string(data))
		assert.Equal(t, "90.01", inv.Totals.Advances.String())
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "90.01", inv.Totals.Advances.String())
		assert.Equal(t, "90.01", inv.Totals.Payable.String())
		assert.Equal(t, "0.00", inv.Totals.Due.String())
	})

	t.Run("with retained taxes and advances", func(t *testing.T) {
		inv := baseInvoice(t, &bill.Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "test item 1",
				Price: num.NewAmount(10000, 2),
			},
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Percent:  num.NewPercentage(21, 2),
				},
				{
					Category: es.TaxCategoryIRPF,
					Percent:  num.NewPercentage(15, 2),
				},
			},
		})
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Record{
				{
					Description: "Half paid",
					Percent:     num.NewPercentage(50, 2),
				},
			},
		}
		inv.Tax.PricesInclude = ""
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "21.00", inv.Totals.Tax.String())
		assert.Equal(t, "15.00", inv.Totals.RetainedTax.String())
		assert.Equal(t, "106.00", inv.Totals.Payable.String())
		assert.Equal(t, "53.00", inv.Totals.Due.String())
	})

	t.Run("with multiple informative taxes", func(t *testing.T) {
		inv := baseInvoice(t,
			&bill.Line{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "test item 1",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: br.TaxCategoryISS,
						Percent:  num.NewPercentage(50, 3),
					},
				},
			},
			&bill.Line{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "test item 2",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: br.TaxCategoryISS,
						Percent:  num.NewPercentage(30, 3),
					},
				},
			},
		)
		inv.Supplier.TaxID.Country = "BR"
		inv.Tax.PricesInclude = ""
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "0.00", inv.Totals.Tax.String())
		assert.Equal(t, "200.00", inv.Totals.TotalWithTax.String())
		assert.Equal(t, "200.00", inv.Totals.Payable.String())
		if iss := inv.Totals.Taxes.Category(br.TaxCategoryISS); assert.NotNil(t, iss) {
			assert.Equal(t, "8.00", iss.Amount.String())
			assert.Equal(t, "5.00", iss.Rates[0].Amount.String())
			assert.Equal(t, "100.00", iss.Rates[0].Base.String())
			assert.Equal(t, "3.00", iss.Rates[1].Amount.String())
			assert.Equal(t, "100.00", iss.Rates[1].Base.String())
			assert.True(t, iss.Informative)
		}
	})
}

func TestRemoveIncludedTaxes(t *testing.T) {
	t.Run("no included tax", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Tax = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.RemoveIncludedTaxes())
	})

	t.Run("with currency rounding", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Tax = &bill.Tax{
			PricesInclude: tax.CategoryVAT,
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.RemoveIncludedTaxes())
		assert.Equal(t, "826.45", inv.Totals.Sum.String())
		assert.Equal(t, "173.55", inv.Totals.Tax.String())
		assert.Equal(t, "1000.00", inv.Totals.Payable.String())
	})

	t.Run("from discounts", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Discounts = []*bill.Discount{
			{
				Amount: num.MakeAmount(1000, 2),
				Reason: "testing",
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     "general",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.RemoveIncludedTaxes())
		assert.Equal(t, "8.26", inv.Totals.Discount.String())
	})

	t.Run("from charges", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Charges = []*bill.Charge{
			{
				Amount: num.MakeAmount(1000, 2),
				Reason: "testing",
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     "general",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.RemoveIncludedTaxes())
		assert.Equal(t, "8.26", inv.Totals.Charge.String())
	})
}

func TestRemoveIncludedTaxesCurrencyRounding(t *testing.T) {
	// roomRates provides the tax-inclusive lines of the hotel invoice that
	// motivated sharing the removed tax out between the lines: 12 nights at
	// 125.00 with 6% VAT included, whose base is 1415.09 even though removing
	// the tax from each price on its own only adds up to 1415.04.
	roomRates := func(n int) []*bill.Line {
		lines := make([]*bill.Line, n)
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
		return lines
	}
	// lineSum adds up the totals of every line in the invoice.
	lineSum := func(inv *bill.Invoice) num.Amount {
		sum := num.MakeAmount(0, 2)
		for _, l := range inv.Lines {
			sum = sum.Add(*l.Total)
		}
		return sum
	}

	t.Run("multiple lines", func(t *testing.T) {
		inv := baseInvoice(t, roomRates(12)...)
		inv.Tax.Rounding = tax.RoundingRuleCurrency
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.RemoveIncludedTaxes())

		assert.Equal(t, "1415.09", inv.Totals.Sum.String())
		assert.Equal(t, "1415.09", inv.Totals.Total.String())
		assert.Equal(t, "84.91", inv.Totals.Tax.String())
		assert.Equal(t, "1500.00", inv.Totals.TotalWithTax.String())
		assert.Equal(t, "1500.00", inv.Totals.Payable.String())
		assert.Nil(t, inv.Totals.Rounding, "no compensation needed")

		rt := inv.Totals.Taxes.Categories[0].Rates[0]
		assert.Equal(t, "1415.09", rt.Base.String())
		assert.Equal(t, "84.91", rt.Amount.String())

		// The lines carry the remainder, so they still add up to the base and
		// each of their sums can be recalculated from the price presented.
		assert.Equal(t, "1415.09", lineSum(inv).String())
		assert.Equal(t, "117.9245", inv.Lines[0].Item.Price.String())
		assert.Equal(t, "117.92", inv.Lines[0].Total.String())
		assert.Equal(t, "117.93", inv.Lines[1].Item.Price.String())
		assert.Equal(t, "117.93", inv.Lines[1].Total.String())
		for _, l := range inv.Lines {
			recalc := l.Item.Price.Multiply(l.Quantity).Rescale(2)
			assert.Equal(t, l.Sum.String(), recalc.String(), "line %d", l.Index)
		}
	})

	t.Run("large quantity", func(t *testing.T) {
		// A single line needs more accuracy in its price to be able to absorb
		// the remainder of a large quantity.
		lines := roomRates(1)
		lines[0].Quantity = num.MakeAmount(1000, 0)
		inv := baseInvoice(t, lines...)
		inv.Tax.Rounding = tax.RoundingRuleCurrency
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.RemoveIncludedTaxes())

		assert.Equal(t, "117.92453", inv.Lines[0].Item.Price.String())
		assert.Equal(t, "117924.53", inv.Totals.Total.String())
		assert.Equal(t, "125000.00", inv.Totals.TotalWithTax.String())
		assert.Nil(t, inv.Totals.Rounding)
	})

	t.Run("with line discounts", func(t *testing.T) {
		lines := roomRates(10)
		lines[0].Quantity = num.MakeAmount(3, 0)
		lines[0].Discounts = []*bill.LineDiscount{
			{
				Percent: num.NewPercentage(10, 2),
				Reason:  "promo",
			},
		}
		inv := baseInvoice(t, lines...)
		inv.Tax.Rounding = tax.RoundingRuleCurrency
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.RemoveIncludedTaxes())

		assert.Equal(t, "1379.72", inv.Totals.Total.String())
		assert.Equal(t, "1462.50", inv.Totals.TotalWithTax.String())
		assert.Nil(t, inv.Totals.Rounding)
		assert.Equal(t, "1379.72", lineSum(inv).String())
		assert.Equal(t, "117.927", inv.Lines[0].Item.Price.String())
		assert.Equal(t, "353.78", inv.Lines[0].Sum.String())
	})

	t.Run("with a document discount", func(t *testing.T) {
		inv := baseInvoice(t, roomRates(12)...)
		inv.Tax.Rounding = tax.RoundingRuleCurrency
		inv.Discounts = []*bill.Discount{
			{
				Amount: num.MakeAmount(3333, 2),
				Reason: "goodwill",
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Percent:  num.NewPercentage(6, 2),
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.RemoveIncludedTaxes())

		assert.Equal(t, "31.44", inv.Totals.Discount.String())
		assert.Equal(t, "1383.65", inv.Totals.Total.String())
		assert.Equal(t, "1466.67", inv.Totals.TotalWithTax.String())
		assert.Nil(t, inv.Totals.Rounding)
	})

	t.Run("with a document charge", func(t *testing.T) {
		inv := baseInvoice(t, roomRates(12)...)
		inv.Tax.Rounding = tax.RoundingRuleCurrency
		inv.Charges = []*bill.Charge{
			{
				Amount: num.MakeAmount(3333, 2),
				Reason: "late checkout",
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Percent:  num.NewPercentage(6, 2),
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.RemoveIncludedTaxes())

		assert.Equal(t, "31.45", inv.Totals.Charge.String())
		assert.Equal(t, "1446.54", inv.Totals.Total.String())
		assert.Equal(t, "1533.33", inv.Totals.TotalWithTax.String())
		assert.Nil(t, inv.Totals.Rounding)
	})

	t.Run("multiple rates", func(t *testing.T) {
		extra := func(name string, price *num.Amount, percent *num.Percentage) *bill.Line {
			return &bill.Line{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  name,
					Price: price,
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Percent:  percent,
					},
				},
			}
		}
		lines := append(roomRates(2),
			extra("Breakfast", num.NewAmount(1250, 2), num.NewPercentage(13, 2)),
			extra("Breakfast", num.NewAmount(1250, 2), num.NewPercentage(13, 2)),
			extra("Dinner menu", num.NewAmount(3580, 2), num.NewPercentage(13, 2)),
			extra("Minibar drinks", num.NewAmount(795, 2), num.NewPercentage(23, 2)),
			extra("Spa access", num.NewAmount(4500, 2), num.NewPercentage(23, 2)),
		)
		inv := baseInvoice(t, lines...)
		inv.Tax.Rounding = tax.RoundingRuleCurrency
		require.NoError(t, inv.Calculate())
		bases := make([]string, 0, 3)
		for _, rt := range inv.Totals.Taxes.Categories[0].Rates {
			bases = append(bases, rt.Base.String())
		}
		require.NoError(t, inv.RemoveIncludedTaxes())

		// Every rate keeps the base the invoice was issued with, sharing the
		// remainder out between its own lines.
		for i, rt := range inv.Totals.Taxes.Categories[0].Rates {
			assert.Equal(t, bases[i], rt.Base.String(), "rate %s", rt.Percent)
		}
		assert.Equal(t, []string{"235.85", "53.81", "43.05"}, bases)
		assert.Equal(t, "332.71", inv.Totals.Total.String())
		assert.Equal(t, "332.71", lineSum(inv).String())

		// The 13% tax of 6.99 that the tax-inclusive totals produced cannot be
		// recalculated from a base of 53.81, which gives 7.00, so the resulting
		// cent still needs to be compensated for.
		assert.Equal(t, "-0.01", inv.Totals.Rounding.String())
		assert.Equal(t, "363.75", inv.Totals.Payable.String())
	})
}

func TestCalculatePricesIncludeCurrencyRounding(t *testing.T) {
	t.Run("multiple lines", func(t *testing.T) {
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
		inv.Tax.Rounding = tax.RoundingRuleCurrency
		require.NoError(t, inv.Calculate())

		assert.Equal(t, "1500.00", inv.Totals.Sum.String())
		assert.Equal(t, "84.91", inv.Totals.TaxIncluded.String())
		assert.Equal(t, "1415.09", inv.Totals.Total.String())
		assert.Equal(t, "84.91", inv.Totals.Tax.String())
		assert.Equal(t, "1500.00", inv.Totals.TotalWithTax.String())
		assert.Equal(t, "1500.00", inv.Totals.Payable.String())
		rt := inv.Totals.Taxes.Categories[0].Rates[0]
		assert.Equal(t, "1415.09", rt.Base.String())
		assert.Equal(t, "84.91", rt.Amount.String())
	})

	t.Run("with retained taxes", func(t *testing.T) {
		lines := make([]*bill.Line, 10)
		for i := range lines {
			lines[i] = &bill.Line{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Service",
					Price: num.NewAmount(193, 2),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     tax.RateGeneral,
					},
					{
						Category: es.TaxCategoryIRPF,
						Rate:     "pro",
					},
				},
			}
		}
		inv := baseInvoice(t, lines...)
		inv.Tax.Rounding = tax.RoundingRuleCurrency
		require.NoError(t, inv.Calculate())

		assert.Equal(t, "19.30", inv.Totals.Sum.String())
		assert.Equal(t, "3.35", inv.Totals.TaxIncluded.String())
		assert.Equal(t, "15.95", inv.Totals.Total.String())
		assert.Equal(t, "19.30", inv.Totals.TotalWithTax.String())
		assert.Equal(t, "2.39", inv.Totals.RetainedTax.String())
		assert.Equal(t, "16.91", inv.Totals.Payable.String())
		vat := inv.Totals.Taxes.Categories[0].Rates[0]
		irpf := inv.Totals.Taxes.Categories[1].Rates[0]
		assert.Equal(t, "15.95", vat.Base.String())
		assert.Equal(t, "3.35", vat.Amount.String())
		assert.Equal(t, "15.95", irpf.Base.String())
	})

	t.Run("multiple lines with retained taxes", func(t *testing.T) {
		lines := make([]*bill.Line, 12)
		for i := range lines {
			lines[i] = &bill.Line{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Professional services",
					Price: num.NewAmount(12500, 2),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     tax.RateGeneral,
					},
					{
						Category: es.TaxCategoryIRPF,
						Rate:     "pro",
					},
				},
			}
		}
		inv := baseInvoice(t, lines...)
		inv.Tax.Rounding = tax.RoundingRuleCurrency
		require.NoError(t, inv.Calculate())

		assert.Equal(t, "1500.00", inv.Totals.Sum.String())
		assert.Equal(t, "260.33", inv.Totals.TaxIncluded.String())
		assert.Equal(t, "1239.67", inv.Totals.Total.String())
		assert.Equal(t, "1500.00", inv.Totals.TotalWithTax.String())
		assert.Equal(t, "185.95", inv.Totals.RetainedTax.String())
		assert.Equal(t, "1314.05", inv.Totals.Payable.String())
		vat := inv.Totals.Taxes.Categories[0].Rates[0]
		irpf := inv.Totals.Taxes.Categories[1].Rates[0]
		assert.Equal(t, "1239.67", vat.Base.String())
		assert.Equal(t, "260.33", vat.Amount.String())
		assert.Equal(t, "1239.67", irpf.Base.String())
	})
}
