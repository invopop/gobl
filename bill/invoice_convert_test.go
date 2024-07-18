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
			Tax: &bill.Tax{
				PricesInclude: tax.CategoryVAT,
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
						Amount:      num.MakeAmount(25000, 2),
					},
				},
			},
		}
		i2, err := i.ConvertInto(currency.USD)
		assert.NoError(t, err)
		require.NotNil(t, i2)
		assert.Equal(t, "11.20", i2.Outlays[0].Amount.String())
		assert.Equal(t, "280.00", i2.Payment.Advances[0].Amount.String())
		assert.Equal(t, i2.Totals.Sum.String(), "1064.00")
		assert.Equal(t, i2.Totals.Due.String(), "795.20")
	})
}
