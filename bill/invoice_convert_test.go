package bill_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceConvertInto(t *testing.T) {
	t.Run("simple conversion", func(t *testing.T) {
		lines := []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(12050, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     tax.RateStandard,
					},
				},
			},
		}
		inv := baseInvoice(t, lines...)
		_, err := inv.ConvertInto(currency.USD)
		assert.ErrorContains(t, err, "no exchange rate defined for 'EUR' to 'USD'")

		inv.ExchangeRates = append(inv.ExchangeRates, &currency.ExchangeRate{
			From:   currency.EUR,
			To:     currency.USD,
			Amount: num.MakeAmount(112, 2),
		})

		i2, err := inv.ConvertInto(currency.USD)
		assert.NoError(t, err)
		require.NotNil(t, i2)
		assert.Equal(t, "134.96", i2.Totals.Payable.String())
		assert.Equal(t, "USD", i2.Currency.String())
		assert.Equal(t, "134.9600", i2.Lines[0].Item.Price.String())
		assert.Len(t, i2.Lines[0].Item.AltPrices, 1)
		assert.Equal(t, "EUR", i2.Lines[0].Item.AltPrices[0].Currency.String())
		assert.Equal(t, "120.50", i2.Lines[0].Item.AltPrices[0].Value.String())
	})

	t.Run("conversion with alt prices", func(t *testing.T) {
		lines := []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(12050, 2),
					AltPrices: []*currency.Amount{
						{
							Currency: currency.USD,
							Value:    num.MakeAmount(13000, 2),
						},
					},
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     tax.RateStandard,
					},
				},
			},
		}
		inv := baseInvoice(t, lines...)

		inv.ExchangeRates = append(inv.ExchangeRates, &currency.ExchangeRate{
			From:   currency.EUR,
			To:     currency.USD,
			Amount: num.MakeAmount(112, 2),
		})

		i2, err := inv.ConvertInto(currency.USD)
		assert.NoError(t, err)
		require.NotNil(t, i2)
		assert.Equal(t, "130.00", i2.Totals.Payable.String())
		assert.Equal(t, "USD", i2.Currency.String())
		assert.Equal(t, "130.00", i2.Lines[0].Item.Price.String())
		assert.Len(t, i2.Lines[0].Item.AltPrices, 1)
		assert.Equal(t, "EUR", i2.Lines[0].Item.AltPrices[0].Currency.String())
		assert.Equal(t, "120.50", i2.Lines[0].Item.AltPrices[0].Value.String())
	})

	t.Run("complex example", func(t *testing.T) {
		i := &bill.Invoice{
			Code: "123TEST",
			ExchangeRates: []*currency.ExchangeRate{
				{
					From:   currency.EUR,
					To:     currency.USD,
					Amount: num.MakeAmount(112, 2),
				},
			},
			Supplier: &org.Party{
				TaxID: &tax.Identity{
					Country: l10n.ES,
					Code:    "B98602642",
				},
			},
			Customer: &org.Party{
				TaxID: &tax.Identity{
					Country: l10n.ES,
					Code:    "54387763P",
				},
			},
			IssueDate: cal.MakeDate(2022, 6, 13),
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(10, 0),
					Item: &org.Item{
						Name:  "Test Item",
						Price: num.MakeAmount(10000, 2),
					},
					Taxes: tax.Set{
						{
							Category: "VAT",
							Rate:     "standard",
						},
					},
					Discounts: []*bill.LineDiscount{
						{
							Reason: "Testing",
							Amount: num.MakeAmount(10000, 2),
						},
					},
					Charges: []*bill.LineCharge{
						{
							Reason: "Testing Charge",
							Amount: num.MakeAmount(5000, 2),
						},
					},
				},
			},
			Outlays: []*bill.Outlay{
				{
					Description: "Something paid in advance",
					Amount:      num.MakeAmount(1000, 2),
				},
			},
			Payment: &bill.Payment{
				Advances: []*pay.Advance{
					{
						Description: "Test Advance",
						Percent:     num.NewPercentage(50, 2),
					},
				},
			},
		}
		i2, err := i.ConvertInto(currency.USD)
		assert.NoError(t, err)
		require.NotNil(t, i2)
		assert.Equal(t, "11.20", i2.Outlays[0].Amount.String())
		assert.Equal(t, "643.72", i2.Payment.Advances[0].Amount.String())
		assert.Equal(t, "1064.00", i2.Totals.Sum.String())
		assert.Equal(t, "1298.64", i2.Totals.Payable.String())
		assert.Equal(t, "654.92", i2.Totals.Due.String())
	})
}
