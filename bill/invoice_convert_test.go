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
					Price: num.NewAmount(12050, 2),
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

		ex, err := json.Marshal(i2.ExchangeRates)
		require.NoError(t, err)
		assert.JSONEq(t, `[{"amount":"1.12","from":"EUR","to":"USD"}]`, string(ex))
	})

	t.Run("conversion with alt prices", func(t *testing.T) {
		lines := []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(12050, 2),
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
					Country: "ES",
					Code:    "B98602642",
				},
			},
			Customer: &org.Party{
				TaxID: &tax.Identity{
					Country: "ES",
					Code:    "54387763P",
				},
			},
			IssueDate: cal.MakeDate(2022, 6, 13),
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(10, 0),
					Item: &org.Item{
						Name:  "Test Item",
						Price: num.NewAmount(10000, 2),
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
			Charges: []*bill.Charge{
				{
					Reason: "Testing Charge",
					Amount: num.MakeAmount(5000, 2),
				},
			},
			Discounts: []*bill.Discount{
				{
					Reason: "Testing",
					Amount: num.MakeAmount(100, 2),
				},
			},
			Payment: &bill.PaymentDetails{
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
		assert.Equal(t, "671.16", i2.Payment.Advances[0].Amount.String())
		assert.Equal(t, "1064.00", i2.Totals.Sum.String())
		assert.Equal(t, "1342.32", i2.Totals.Payable.String())
		assert.Equal(t, "671.16", i2.Totals.Due.String())
	})
}
