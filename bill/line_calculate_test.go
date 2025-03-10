package bill

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLineCalculate(t *testing.T) {
	t.Run("lines with errors", func(t *testing.T) {
		lines := []*Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:     "Test Item",
					Currency: currency.MXN,
					Price:    num.NewAmount(1000, 2),
				},
			},
		}
		err := calculateLines(lines, currency.EUR, exampleRates(t), tax.RoundingRulePrecise)
		assert.ErrorContains(t, err, "0: (item: no exchange rate found from 'MXN' to 'EUR'.).")
	})
	t.Run("missing item", func(t *testing.T) {
		line := &Line{}
		err := calculateLine(line, currency.EUR, exampleRates(t), tax.RoundingRulePrecise)
		require.NoError(t, err)
	})
	t.Run("invalid item price", func(t *testing.T) {
		line := &Line{
			Quantity: num.MakeAmount(10, 0),
			Item: &org.Item{
				Name:     "Test Item",
				Currency: currency.MXN,
				Price:    num.NewAmount(1000, 2),
			},
		}
		err := calculateLine(line, currency.EUR, exampleRates(t), tax.RoundingRulePrecise)
		assert.ErrorContains(t, err, "item: no exchange rate found from 'MXN' to 'EUR'")
	})
	t.Run("substituted: basic", func(t *testing.T) {
		line := &Line{
			Quantity: num.MakeAmount(10, 0),
			Item: &org.Item{
				Name:  "New Test Item, less for more",
				Price: num.NewAmount(1000, 2),
			},
			Substituted: []*SubLine{
				{
					Quantity: num.MakeAmount(15, 0),
					Item: &org.Item{
						Name:  "Old Test Item",
						Price: num.NewAmount(10, 0),
					},
				},
			},
		}
		err := calculateLine(line, currency.EUR, exampleRates(t), tax.RoundingRulePrecise)
		require.NoError(t, err)
		assert.Equal(t, 1, line.Substituted[0].Index)
		assert.Equal(t, "150.0000", line.Substituted[0].Total.String())
		assert.Equal(t, "10.00", line.Item.Price.String())
		line.round(currency.EUR)
		assert.Equal(t, "150.00", line.Substituted[0].Total.String())
	})
	t.Run("substituted: currency error", func(t *testing.T) {
		line := &Line{
			Quantity: num.MakeAmount(10, 0),
			Item: &org.Item{
				Name:  "New Test Item, less for more",
				Price: num.NewAmount(1000, 2),
			},
			Substituted: []*SubLine{
				{
					Quantity: num.MakeAmount(15, 0),
					Item: &org.Item{
						Name:     "Old Test Item",
						Currency: currency.MXN,
						Price:    num.NewAmount(10, 0),
					},
				},
			},
		}
		err := calculateLine(line, currency.EUR, exampleRates(t), tax.RoundingRulePrecise)
		require.ErrorContains(t, err, "substituted: (0: no exchange rate found from 'MXN' to 'EUR'.)")
	})
	t.Run("sublines: basic", func(t *testing.T) {
		line := &Line{
			Item: &org.Item{
				Name: "Test Group Item",
			},
			Breakdown: []*SubLine{
				{
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Test Item",
						Price: num.NewAmount(10, 0),
					},
				},
			},
		}
		err := calculateLine(line, currency.EUR, exampleRates(t), tax.RoundingRulePrecise)
		line.round(currency.EUR)
		require.NoError(t, err)
		assert.Equal(t, "10.00", line.Item.Price.String())
	})
	t.Run("sublines: multiple", func(t *testing.T) {
		line := &Line{
			Item: &org.Item{
				Name: "Test Group Item",
			},
			Breakdown: []*SubLine{
				{
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Test Item 1",
						Price: num.NewAmount(10, 0),
					},
				},
				{
					Item: &org.Item{
						Name: "Dummy line",
					},
				},
				{
					Quantity: num.MakeAmount(2, 0),
					Item: &org.Item{
						Name:  "Test Item 2",
						Price: num.NewAmount(10, 0),
					},
				},
			},
		}
		err := calculateLine(line, currency.EUR, exampleRates(t), tax.RoundingRulePrecise)
		require.NoError(t, err)
		assert.Equal(t, "30.00", line.Item.Price.String())
	})
	t.Run("sublines: maintain precision", func(t *testing.T) {
		line := &Line{
			Quantity: num.MakeAmount(2, 0),
			Item: &org.Item{
				Name: "Test Group Item",
			},
			Breakdown: []*SubLine{
				{
					Quantity: num.MakeAmount(2, 0),
					Item: &org.Item{
						Name:  "Test Item",
						Price: num.NewAmount(124567, 4),
					},
				},
			},
		}
		err := calculateLine(line, currency.EUR, exampleRates(t), tax.RoundingRulePrecise)
		require.NoError(t, err)
		assert.Equal(t, "12.4567", line.Breakdown[0].Item.Price.String())
		assert.Equal(t, "24.9134", line.Breakdown[0].Sum.String())
		assert.Equal(t, "24.9134", line.Breakdown[0].Total.String())
		assert.Equal(t, "24.9134", line.Item.Price.String())
		assert.Equal(t, "49.8268", line.Total.String())
		line.round(currency.EUR)
		assert.Equal(t, "12.4567", line.Breakdown[0].Item.Price.String())
		assert.Equal(t, "24.9134", line.Breakdown[0].Sum.String())
		assert.Equal(t, "24.91", line.Breakdown[0].Total.String())
		assert.Equal(t, "24.9134", line.Item.Price.String())
		assert.Equal(t, "49.83", line.Total.String())
	})
	t.Run("sublines: match precision with all sublines", func(t *testing.T) {
		line := &Line{
			Quantity: num.MakeAmount(2, 0),
			Item: &org.Item{
				Name: "Test Group Item",
			},
			Breakdown: []*SubLine{
				{
					Quantity: num.MakeAmount(2, 0),
					Item: &org.Item{
						Name:  "Test Item",
						Price: num.NewAmount(124567, 4),
					},
				},
				{
					Quantity: num.MakeAmount(2, 0),
					Item: &org.Item{
						Name:  "Test Item",
						Price: num.NewAmount(1245678, 5),
					},
				},
			},
		}
		err := calculateLine(line, currency.EUR, exampleRates(t), tax.RoundingRulePrecise)
		require.NoError(t, err)
		assert.Equal(t, "12.4567", line.Breakdown[0].Item.Price.String())
		assert.Equal(t, "24.9134", line.Breakdown[0].Sum.String())
		assert.Equal(t, "24.9134", line.Breakdown[0].Total.String())
		assert.Equal(t, "24.91356", line.Breakdown[1].Total.String())
		assert.Equal(t, "49.82696", line.Item.Price.String())
		assert.Equal(t, "99.65392", line.Total.String())
		line.round(currency.EUR)
		assert.Equal(t, "12.4567", line.Breakdown[0].Item.Price.String())
		assert.Equal(t, "24.9134", line.Breakdown[0].Sum.String())
		assert.Equal(t, "24.91", line.Breakdown[0].Total.String())
		assert.Equal(t, "24.91", line.Breakdown[1].Total.String())
		assert.Equal(t, "49.82696", line.Item.Price.String())
		assert.Equal(t, "99.65", line.Total.String())
	})
	t.Run("sublines: without price", func(t *testing.T) {
		line := &Line{
			Item: &org.Item{
				Name: "Test Group Item",
			},
			Breakdown: []*SubLine{
				{
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name: "Test Item",
					},
				},
			},
		}
		err := calculateLine(line, currency.EUR, exampleRates(t), tax.RoundingRulePrecise)
		require.NoError(t, err)
		assert.Nil(t, line.Item.Price)
	})
	t.Run("sublines: without item", func(t *testing.T) {
		line := &Line{
			Item: &org.Item{
				Name: "Test Group Item",
			},
			Breakdown: []*SubLine{
				{
					Quantity: num.MakeAmount(1, 0),
				},
			},
		}
		err := calculateLine(line, currency.EUR, exampleRates(t), tax.RoundingRulePrecise)
		require.NoError(t, err)
		assert.Nil(t, line.Item.Price)
	})
	t.Run("sublines: missing exchange rate", func(t *testing.T) {
		line := &Line{
			Item: &org.Item{
				Name: "Test Group Item",
			},
			Breakdown: []*SubLine{
				{
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:     "Test Item",
						Currency: currency.MXN,
						Price:    num.NewAmount(1000, 2),
					},
				},
			},
		}
		err := calculateLine(line, currency.EUR, exampleRates(t), tax.RoundingRulePrecise)
		require.ErrorContains(t, err, "breakdown: (0: no exchange rate found from 'MXN' to 'EUR'.).")
	})

	t.Run("lines with round-then-sum", func(t *testing.T) {
		lines := []*Line{
			{
				Quantity: num.MakeAmount(3, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(125934, 4),
				},
				Discounts: []*LineDiscount{
					{
						Percent: num.NewPercentage(7, 2),
					},
				},
			},
		}
		err := calculateLines(lines, currency.EUR, nil, tax.RoundingRuleCurrency)
		assert.NoError(t, err)
		sum := calculateLineSum(lines, currency.EUR)
		assert.Equal(t, "35.14", sum.String())
		assert.Equal(t, "37.78", lines[0].Sum.String())
		assert.Equal(t, "2.64", lines[0].Discounts[0].Amount.String())
		assert.Equal(t, "35.14", lines[0].Total.String())
		roundLines(lines, currency.EUR)
		assert.Equal(t, "35.14", sum.String())
		assert.Equal(t, "37.78", lines[0].Sum.String())
		assert.Equal(t, "2.64", lines[0].Discounts[0].Amount.String())
		assert.Equal(t, "35.14", lines[0].Total.String())
	})

	t.Run("lines with sum-then-round", func(t *testing.T) {
		lines := []*Line{
			{
				Quantity: num.MakeAmount(3, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(125934, 4),
				},
				Discounts: []*LineDiscount{
					{
						Percent: num.NewPercentage(7, 2),
					},
				},
			},
		}
		err := calculateLines(lines, currency.EUR, nil, tax.RoundingRulePrecise)
		assert.NoError(t, err)
		sum := calculateLineSum(lines, currency.EUR)
		assert.Equal(t, "35.1356", sum.String())
		assert.Equal(t, "37.7802", lines[0].Sum.String())
		assert.Equal(t, "2.6446", lines[0].Discounts[0].Amount.String())
		assert.Equal(t, "35.1356", lines[0].Total.String())
		roundLines(lines, currency.EUR)
		assert.Equal(t, "35.1356", sum.String())
		assert.Equal(t, "37.7802", lines[0].Sum.String())
		assert.Equal(t, "2.64", lines[0].Discounts[0].Amount.String())
		assert.Equal(t, "35.14", lines[0].Total.String())
	})
	t.Run("lines with sum-then-round, regular precision", func(t *testing.T) {
		lines := []*Line{
			{
				Quantity: num.MakeAmount(3, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(1259, 2),
				},
				Discounts: []*LineDiscount{
					{
						Percent: num.NewPercentage(7, 2),
					},
				},
			},
		}
		err := calculateLines(lines, currency.EUR, nil, tax.RoundingRulePrecise)
		assert.NoError(t, err)
		sum := calculateLineSum(lines, currency.EUR)
		assert.Equal(t, "35.1261", sum.String())
		assert.Equal(t, "35.1261", lines[0].Total.String())
	})

	t.Run("lines with discount base", func(t *testing.T) {
		lines := []*Line{
			{
				Quantity: num.MakeAmount(3, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(1259, 2),
				},
				Discounts: []*LineDiscount{
					{
						Base:    num.NewAmount(51256, 3),
						Percent: num.NewPercentage(7, 2),
					},
				},
			},
		}
		err := calculateLines(lines, currency.EUR, nil, tax.RoundingRulePrecise)
		assert.NoError(t, err)
		sum := calculateLineSum(lines, currency.EUR)
		assert.Equal(t, "34.1821", sum.String())
		assert.Equal(t, "37.7700", lines[0].Sum.String())
		assert.Equal(t, "51.256", lines[0].Discounts[0].Base.String())
		assert.Equal(t, "3.5879", lines[0].Discounts[0].Amount.String())
		assert.Equal(t, "34.1821", lines[0].Total.String())
		roundLines(lines, currency.EUR)
		assert.Equal(t, "37.77", lines[0].Sum.String())
		assert.Equal(t, "51.256", lines[0].Discounts[0].Base.String(), "maintain original precision")
		assert.Equal(t, "3.59", lines[0].Discounts[0].Amount.String())
		assert.Equal(t, "34.18", lines[0].Total.String())
	})
	t.Run("lines with discount base, currency rounding", func(t *testing.T) {
		lines := []*Line{
			{
				Quantity: num.MakeAmount(3, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(1259, 2),
				},
				Discounts: []*LineDiscount{
					{
						Base:    num.NewAmount(51256, 3),
						Percent: num.NewPercentage(7, 2),
					},
				},
			},
		}
		err := calculateLines(lines, currency.EUR, nil, tax.RoundingRuleCurrency)
		assert.NoError(t, err)
		sum := calculateLineSum(lines, currency.EUR)
		assert.Equal(t, "34.18", sum.String())
		assert.Equal(t, "51.256", lines[0].Discounts[0].Base.String())
		assert.Equal(t, "3.59", lines[0].Discounts[0].Amount.String())
		assert.Equal(t, "34.18", lines[0].Total.String())
		roundLines(lines, currency.EUR)
		assert.Equal(t, "37.77", lines[0].Sum.String())
		assert.Equal(t, "51.256", lines[0].Discounts[0].Base.String(), "maintain original written precision")
		assert.Equal(t, "3.59", lines[0].Discounts[0].Amount.String())
		assert.Equal(t, "34.18", lines[0].Total.String())
	})

	t.Run("lines with discount percent, precise", func(t *testing.T) {
		lines := []*Line{
			{
				Quantity: num.MakeAmount(20, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(903324, 4),
				},
				Discounts: []*LineDiscount{
					{
						Percent: num.NewPercentage(10, 2),
					},
				},
			},
		}
		err := calculateLines(lines, currency.EUR, nil, tax.RoundingRulePrecise)
		assert.NoError(t, err)
		sum := calculateLineSum(lines, currency.EUR)
		assert.Equal(t, "1625.9832", sum.String())
		assert.Equal(t, "180.6648", lines[0].Discounts[0].Amount.String())
		assert.Equal(t, "1625.9832", lines[0].Total.String())
		roundLines(lines, currency.EUR)
		assert.Equal(t, "1806.6480", lines[0].Sum.String())
		assert.Equal(t, "180.66", lines[0].Discounts[0].Amount.String())
		assert.Equal(t, "1625.98", lines[0].Total.String())
	})
	t.Run("lines with discount percent, currency rounding", func(t *testing.T) {
		lines := []*Line{
			{
				Quantity: num.MakeAmount(20, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(903324, 4),
				},
				Discounts: []*LineDiscount{
					{
						Percent: num.NewPercentage(10, 2),
					},
				},
			},
		}
		err := calculateLines(lines, currency.EUR, nil, tax.RoundingRuleCurrency)
		assert.NoError(t, err)
		sum := calculateLineSum(lines, currency.EUR)
		assert.Equal(t, "1625.98", sum.String())
		assert.Equal(t, "180.67", lines[0].Discounts[0].Amount.String())
		assert.Equal(t, "1625.98", lines[0].Total.String())
		roundLines(lines, currency.EUR)
		assert.Equal(t, "1806.65", lines[0].Sum.String())
		assert.Equal(t, "180.67", lines[0].Discounts[0].Amount.String())
		assert.Equal(t, "1625.98", lines[0].Total.String())
	})

	t.Run("lines with charge base", func(t *testing.T) {
		lines := []*Line{
			{
				Quantity: num.MakeAmount(3, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(1259, 2),
				},
				Charges: []*LineCharge{
					{
						Base:    num.NewAmount(512, 2),
						Percent: num.NewPercentage(7, 2),
					},
				},
			},
		}
		err := calculateLines(lines, currency.EUR, nil, tax.RoundingRulePrecise)
		assert.NoError(t, err)
		data, _ := json.Marshal(lines)
		fmt.Printf("%s", string(data))
		sum := calculateLineSum(lines, currency.EUR)
		assert.Equal(t, "38.1284", sum.String())
		assert.Equal(t, "0.3584", lines[0].Charges[0].Amount.String())
		assert.Equal(t, "38.1284", lines[0].Total.String())
		roundLines(lines, currency.EUR)
		assert.Equal(t, "0.36", lines[0].Charges[0].Amount.String())
		assert.Equal(t, "38.13", lines[0].Total.String())
	})

	t.Run("lines with rate", func(t *testing.T) {
		lines := []*Line{
			{
				Quantity: num.MakeAmount(3, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(1259, 2),
				},
				Charges: []*LineCharge{
					{
						Rate: num.NewAmount(2, 2),
					},
				},
			},
		}
		err := calculateLines(lines, currency.EUR, nil, tax.RoundingRulePrecise)
		assert.NoError(t, err)
		sum := calculateLineSum(lines, currency.EUR)
		assert.Equal(t, "37.8300", sum.String())
		assert.Equal(t, "0.06", lines[0].Charges[0].Amount.String())
		assert.Equal(t, "37.7700", lines[0].Sum.String())
		assert.Equal(t, "37.8300", lines[0].Total.String())
	})
	t.Run("lines with quantity and rate", func(t *testing.T) {
		lines := []*Line{
			{
				Quantity: num.MakeAmount(3, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(1259, 2),
				},
				Charges: []*LineCharge{
					{
						Quantity: num.NewAmount(100, 0), // 100g for example
						Rate:     num.NewAmount(2, 2),
					},
				},
			},
		}
		err := calculateLines(lines, currency.EUR, nil, tax.RoundingRulePrecise)
		assert.NoError(t, err)
		sum := calculateLineSum(lines, currency.EUR)
		assert.Equal(t, "39.7700", sum.String())
		assert.Equal(t, "2.00", lines[0].Charges[0].Amount.String())
		assert.Equal(t, "37.7700", lines[0].Sum.String())
		assert.Equal(t, "39.7700", lines[0].Total.String())
	})
}
