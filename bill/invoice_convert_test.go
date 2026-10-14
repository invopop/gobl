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
						Rate:     tax.RateGeneral,
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
						Rate:     tax.RateGeneral,
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

	t.Run("conversion with item list price", func(t *testing.T) {
		lines := []*bill.Line{
			{
				Quantity: num.MakeAmount(250, 0),
				Item: &org.Item{
					Name:     "Test Item",
					List:     num.NewAmount(12000, 2),
					Discount: num.NewAmount(1200, 2),
					Per:      num.NewAmount(100, 0),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     tax.RateGeneral,
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
		require.NoError(t, err)
		ip := i2.Lines[0].Item
		assert.Equal(t, "100", ip.Per.String())
		assert.Equal(t, "134.4000", ip.List.String())
		assert.Equal(t, "13.4400", ip.Discount.String())
		assert.Equal(t, "120.9600", i2.Lines[0].Item.Price.String())
		assert.Equal(t, "302.4000", i2.Lines[0].Sum.String())
	})

	t.Run("conversion with item list price and alt prices", func(t *testing.T) {
		lines := []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name: "Test Item",
					AltPrices: []*currency.Amount{
						{Currency: currency.USD, Value: num.MakeAmount(12000, 2)},
					},
					List:     num.NewAmount(12000, 2),
					Discount: num.NewAmount(1200, 2),
					Per:      num.NewAmount(10, 0),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     tax.RateGeneral,
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
		require.NoError(t, err)
		ip := i2.Lines[0].Item
		assert.Equal(t, "10", ip.Per.String())
		assert.Nil(t, ip.List)
		assert.Nil(t, ip.Discount)
		assert.Equal(t, "120.00", i2.Lines[0].Item.Price.String())
		assert.Equal(t, "12.00", i2.Lines[0].Sum.String())
	})

	t.Run("conversion with breakdown and substituted", func(t *testing.T) {
		lines := []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item:     &org.Item{Name: "Group"},
				Breakdown: []*bill.SubLine{
					{
						Quantity: num.MakeAmount(2, 0),
						Item: &org.Item{
							Name:     "Part",
							List:     num.NewAmount(1000, 2),
							Discount: num.NewAmount(100, 2),
						},
						Discounts: []*bill.LineDiscount{
							{Reason: "Promo", Amount: num.MakeAmount(100, 2)},
						},
						Charges: []*bill.LineCharge{
							{Reason: "Handling", Amount: num.MakeAmount(50, 2)},
						},
					},
					{
						Quantity: num.MakeAmount(1, 0),
						Item: &org.Item{
							Name:  "Other",
							Price: num.NewAmount(500, 2),
							AltPrices: []*currency.Amount{
								{Currency: currency.MXN, Value: num.MakeAmount(10000, 2)},
								{Currency: currency.USD, Value: num.MakeAmount(600, 2)},
							},
						},
					},
				},
				Substituted: []*bill.SubLine{
					{
						Quantity: num.MakeAmount(1, 0),
						Item: &org.Item{
							Name:     "Old",
							Currency: currency.EUR,
							Price:    num.NewAmount(300, 2),
						},
					},
					{
						Quantity: num.MakeAmount(1, 0),
						Item:     &org.Item{Name: "Unpriced"},
					},
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     tax.RateGeneral,
					},
				},
			},
		}
		inv := baseInvoice(t, lines...)
		inv.ExchangeRates = append(inv.ExchangeRates, &currency.ExchangeRate{
			From:   currency.EUR,
			To:     currency.USD,
			Amount: num.MakeAmount(2, 0),
		})

		i2, err := inv.ConvertInto(currency.USD)
		require.NoError(t, err)
		l := i2.Lines[0]
		part := l.Breakdown[0]
		assert.Equal(t, "20.0000", part.Item.List.String())
		assert.Equal(t, "2.0000", part.Item.Discount.String())
		assert.Equal(t, "18.0000", part.Item.Price.String())
		assert.Equal(t, "2.0000", part.Discounts[0].Amount.String())
		assert.Equal(t, "1.0000", part.Charges[0].Amount.String())
		assert.Equal(t, "35.0000", part.Total.String())
		other := l.Breakdown[1]
		assert.Equal(t, "6.00", other.Item.Price.String())
		require.Len(t, other.Item.AltPrices, 2)
		assert.Equal(t, "MXN", other.Item.AltPrices[0].Currency.String())
		assert.Equal(t, "EUR", other.Item.AltPrices[1].Currency.String())
		assert.Equal(t, "5.00", other.Item.AltPrices[1].Value.String())
		assert.Equal(t, "41.0000", l.Item.Price.String())
		assert.Equal(t, "6.0000", l.Substituted[0].Item.Price.String())
		assert.Equal(t, "USD", l.Substituted[0].Item.Currency.String())
		assert.Nil(t, l.Substituted[1].Item.Price)

		// The original invoice is left untouched
		assert.Len(t, inv.Lines[0].Breakdown[1].Item.AltPrices, 2)
		assert.Equal(t, "USD", inv.Lines[0].Breakdown[1].Item.AltPrices[1].Currency.String())
		assert.Equal(t, "10.00", inv.Lines[0].Breakdown[0].Item.List.String())
	})

	t.Run("conversion with substituted under unpriced line", func(t *testing.T) {
		lines := []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item:     &org.Item{Name: "Group"},
				Substituted: []*bill.SubLine{
					{
						Quantity: num.MakeAmount(1, 0),
						Item:     &org.Item{Name: "Old", Price: num.NewAmount(1000, 2)},
					},
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     tax.RateGeneral,
					},
				},
			},
		}
		inv := baseInvoice(t, lines...)
		inv.ExchangeRates = append(inv.ExchangeRates, &currency.ExchangeRate{
			From:   currency.EUR,
			To:     currency.USD,
			Amount: num.MakeAmount(2, 0),
		})

		i2, err := inv.ConvertInto(currency.USD)
		require.NoError(t, err)
		assert.Nil(t, i2.Lines[0].Item.Price)
		assert.Equal(t, "20.0000", i2.Lines[0].Substituted[0].Item.Price.String())
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
							Rate:     "general",
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
				Advances: []*pay.Record{
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

	t.Run("Conversion with Item Currency set", func(t *testing.T) {
		lines := []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:     "Test Item",
					Currency: currency.EUR,
					Price:    num.NewAmount(12050, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     tax.RateGeneral,
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
		assert.Equal(t, "134.96", i2.Totals.Payable.String())
		assert.Equal(t, "USD", i2.Currency.String())
		assert.Equal(t, "USD", i2.Lines[0].Item.Currency.String())
		assert.Equal(t, "134.9600", i2.Lines[0].Item.Price.String())
		assert.Len(t, i2.Lines[0].Item.AltPrices, 1)
		assert.Equal(t, "EUR", i2.Lines[0].Item.AltPrices[0].Currency.String())
		assert.Equal(t, "120.50", i2.Lines[0].Item.AltPrices[0].Value.String())

	})

	t.Run("Conversion with Item Currency not set", func(t *testing.T) {
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
						Rate:     tax.RateGeneral,
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
		assert.Equal(t, "134.96", i2.Totals.Payable.String())
		assert.Equal(t, "USD", i2.Currency.String())
		assert.Empty(t, i2.Lines[0].Item.Currency)
		assert.Equal(t, "134.9600", i2.Lines[0].Item.Price.String())
		assert.Len(t, i2.Lines[0].Item.AltPrices, 1)
		assert.Equal(t, "EUR", i2.Lines[0].Item.AltPrices[0].Currency.String())
		assert.Equal(t, "120.50", i2.Lines[0].Item.AltPrices[0].Value.String())

	})

}
