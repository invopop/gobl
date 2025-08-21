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
		require.ErrorContains(t, inv.Calculate(), "lines: (0: (item: no exchange rate found from 'USD' to 'EUR'.).)")
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
			Advances: []*pay.Advance{
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
			Advances: []*pay.Advance{
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
			Advances: []*pay.Advance{
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

	t.Run("with informative tax", func(t *testing.T) {
		inv := baseInvoice(t, &bill.Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "test item 1",
				Price: num.NewAmount(10000, 2),
			},
			Taxes: tax.Set{
				{
					Category: br.TaxCategoryISS,
					Percent:  num.NewPercentage(50, 2),
				},
			},
		})
		inv.Supplier.TaxID.Country = "BR"
		require.NoError(t, inv.Calculate())

		assert.Equal(t, "0.00", inv.Totals.Tax.String())
		assert.Equal(t, "100.00", inv.Totals.TotalWithTax.String())
		assert.Equal(t, "100.00", inv.Totals.Payable.String())
		if iss := inv.Totals.Taxes.Category(br.TaxCategoryISS); assert.NotNil(t, iss) {
			assert.Equal(t, "50.00", iss.Amount.String())
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
						Rate:     "standard",
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
						Rate:     "standard",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.RemoveIncludedTaxes())
		assert.Equal(t, "8.26", inv.Totals.Charge.String())
	})
}
