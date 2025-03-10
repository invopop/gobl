package bill

import (
	"testing"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func exampleRates(t *testing.T) []*currency.ExchangeRate {
	t.Helper()
	return []*currency.ExchangeRate{
		{
			From:   currency.USD,
			To:     currency.EUR,
			Amount: num.MakeAmount(875967, 6),
		},
		{
			From:   currency.EUR,
			To:     currency.USD,
			Amount: num.MakeAmount(1141860, 6),
		},
		{
			From:   currency.EUR,
			To:     currency.CLP,
			Amount: num.MakeAmount(100629, 2),
		},
	}
}

func TestLineValidation(t *testing.T) {
	t.Run("sublines: basic", func(t *testing.T) {
		lines := []*Line{
			{
				Item: &org.Item{
					Name: "Test Group Item",
				},
				Breakdown: []*SubLine{
					{
						Quantity: num.MakeAmount(1, 0),
						Item: &org.Item{
							Name:  "Test Item",
							Price: num.NewAmount(1000, 2),
						},
					},
				},
			},
		}
		require.NoError(t, calculateLines(lines, currency.EUR, exampleRates(t), tax.RoundingRuleSumThenRound))
		require.NoError(t, validation.Validate(lines))
	})
	t.Run("sublines: with error", func(t *testing.T) {
		lines := []*Line{
			{
				Item: &org.Item{
					Name: "Test Group Item",
				},
				Breakdown: []*SubLine{
					{
						Quantity: num.MakeAmount(1, 0),
						Item: &org.Item{
							Price: num.NewAmount(1000, 2),
						},
					},
				},
			},
		}
		require.NoError(t, calculateLines(lines, currency.EUR, exampleRates(t), tax.RoundingRuleSumThenRound))
		require.ErrorContains(t, validation.Validate(lines), "0: (breakdown: (0: (item: (name: cannot be blank.).).).)")
	})
	t.Run("sublines: missing sum and total", func(t *testing.T) {
		lines := []*Line{
			{
				Item: &org.Item{
					Name: "Test Group Item",
				},
				Breakdown: []*SubLine{
					{
						Quantity: num.MakeAmount(1, 0),
						Item: &org.Item{
							Name:  "Test Item",
							Price: num.NewAmount(1000, 2),
						},
					},
				},
			},
		}
		require.NoError(t, calculateLines(lines, currency.EUR, exampleRates(t), tax.RoundingRuleSumThenRound))
		lines[0].Breakdown[0].Total = nil
		lines[0].Breakdown[0].Sum = nil
		require.ErrorContains(t, validation.Validate(lines), "0: (breakdown: (0: (sum: cannot be blank; total: cannot be blank.).).)")
	})
}

func TestLinePriceNormalization(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		line := &Line{
			Item: &org.Item{
				Name:  "Test Item",
				Price: num.NewAmount(10, 0),
			},
		}
		err := calculateLineItemPrice(line.Item, currency.EUR, exampleRates(t))
		require.NoError(t, err)
		assert.Equal(t, "10.00", line.Item.Price.String())
	})
	t.Run("basic with currency", func(t *testing.T) {
		line := &Line{
			Item: &org.Item{
				Name:     "Test Item",
				Currency: currency.EUR,
				Price:    num.NewAmount(10, 0),
			},
		}
		err := calculateLineItemPrice(line.Item, currency.EUR, exampleRates(t))
		require.NoError(t, err)
		assert.Equal(t, "10.00", line.Item.Price.String())
	})

	t.Run("alt prices", func(t *testing.T) {
		line := &Line{
			Item: &org.Item{
				Name:     "Test Item",
				Currency: currency.USD,
				Price:    num.NewAmount(10, 0),
				AltPrices: []*currency.Amount{
					{
						Currency: currency.EUR,
						Value:    num.MakeAmount(8, 0),
					},
				},
			},
		}
		err := calculateLineItemPrice(line.Item, currency.EUR, exampleRates(t))
		require.NoError(t, err)
		assert.Equal(t, "8.00", line.Item.Price.String())
		require.Len(t, line.Item.AltPrices, 1)
		assert.Equal(t, "10.00", line.Item.AltPrices[0].Value.String())
		assert.Equal(t, "USD", string(line.Item.AltPrices[0].Currency))
	})

	t.Run("use exchange rates", func(t *testing.T) {
		line := &Line{
			Item: &org.Item{
				Name:     "Test Item",
				Currency: currency.USD,
				Price:    num.NewAmount(10, 0),
			},
		}
		err := calculateLineItemPrice(line.Item, currency.EUR, exampleRates(t))
		require.NoError(t, err)
		assert.Equal(t, "8.76", line.Item.Price.String())
		require.Len(t, line.Item.AltPrices, 1)
		assert.Equal(t, "10.00", line.Item.AltPrices[0].Value.String())
		assert.Equal(t, "USD", string(line.Item.AltPrices[0].Currency))
	})

	t.Run("missing exchange rate", func(t *testing.T) {
		line := &Line{
			Item: &org.Item{
				Name:     "Test Item",
				Currency: currency.MXN,
				Price:    num.NewAmount(100, 0),
			},
		}
		err := calculateLineItemPrice(line.Item, currency.EUR, exampleRates(t))
		assert.ErrorContains(t, err, "no exchange rate found from 'MXN' to 'EUR")
	})
}

func TestLineNormalize(t *testing.T) {
	t.Run("basic", func(t *testing.T) {

		line := &Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:     "Test Item",
				Currency: currency.USD,
				Price:    num.NewAmount(10, 0),
			},
			Breakdown: []*SubLine{
				{
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Test Item",
						Price: num.NewAmount(10, 0),
					},
					Discounts: []*LineDiscount{
						{
							Amount: num.MakeAmount(0, 0),
						},
					},
				},
			},
			Discounts: []*LineDiscount{
				{
					Amount: num.MakeAmount(0, 0),
				},
			},
			Charges: []*LineCharge{
				{
					Amount: num.MakeAmount(0, 0),
				},
			},
		}
		line.Normalize(nil)
		assert.Len(t, line.Discounts, 0)
		assert.Len(t, line.Charges, 0)
		assert.Len(t, line.Breakdown[0].Discounts, 0)
	})
}

func TestLineRemoveIncludedTaxes(t *testing.T) {
	t.Run("basic with VAT", func(t *testing.T) {
		line := &Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "Test Item",
				Price: num.NewAmount(1000, 2),
			},
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Percent:  num.NewPercentage(210, 3),
				},
			},
		}
		require.NoError(t, calculateLine(line, currency.EUR, nil, tax.RoundingRuleSumThenRound))
		line = removeLineIncludedTaxes(line, tax.CategoryVAT)
		require.NoError(t, calculateLine(line, currency.EUR, nil, tax.RoundingRuleSumThenRound))
		assert.Equal(t, "8.2645", line.Item.Price.String())
		assert.Equal(t, "8.2645", line.Total.String())
	})
	t.Run("basic with VAT and discounts", func(t *testing.T) {
		line := &Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "Test Item",
				Price: num.NewAmount(1000, 2),
			},
			Discounts: []*LineDiscount{
				{
					Percent: num.NewPercentage(10, 3),
					Reason:  "test",
				},
			},
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Percent:  num.NewPercentage(210, 3),
				},
			},
		}
		require.NoError(t, calculateLine(line, currency.EUR, nil, tax.RoundingRuleSumThenRound))
		line = removeLineIncludedTaxes(line, tax.CategoryVAT)
		require.NoError(t, calculateLine(line, currency.EUR, nil, tax.RoundingRuleSumThenRound))
		assert.Equal(t, "8.2645", line.Item.Price.String())
		assert.Equal(t, "8.1819", line.Total.String())
	})
	t.Run("basic with VAT and charges", func(t *testing.T) {
		line := &Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "Test Item",
				Price: num.NewAmount(1000, 2),
			},
			Charges: []*LineCharge{
				{
					Percent: num.NewPercentage(10, 3),
					Reason:  "test",
				},
			},
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Percent:  num.NewPercentage(210, 3),
				},
			},
		}
		require.NoError(t, calculateLine(line, currency.EUR, nil, tax.RoundingRuleSumThenRound))
		line = removeLineIncludedTaxes(line, tax.CategoryVAT)
		require.NoError(t, calculateLine(line, currency.EUR, nil, tax.RoundingRuleSumThenRound))
		assert.Equal(t, "8.2645", line.Item.Price.String())
		assert.Equal(t, "8.3471", line.Total.String())
	})

	t.Run("sublines: basic with VAT", func(t *testing.T) {
		line := &Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name: "Test Item",
			},
			Breakdown: []*SubLine{
				{
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Test Item",
						Price: num.NewAmount(1000, 2),
					},
				},
			},
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Percent:  num.NewPercentage(210, 3),
				},
			},
		}
		require.NoError(t, calculateLine(line, currency.EUR, nil, tax.RoundingRuleSumThenRound))
		line = removeLineIncludedTaxes(line, tax.CategoryVAT)
		require.NoError(t, calculateLine(line, currency.EUR, nil, tax.RoundingRuleSumThenRound))
		assert.Equal(t, "8.2645", line.Item.Price.String())
		assert.Equal(t, "8.2645", line.Breakdown[0].Item.Price.String())
	})
}
