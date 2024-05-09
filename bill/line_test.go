package bill

import (
	"testing"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
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

func TestLinePriceNormalization(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		line := &Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "Test Item",
				Price: num.MakeAmount(10, 0),
			},
		}
		err := line.normalizeItemPrice(currency.EUR, exampleRates(t))
		require.NoError(t, err)
		assert.Equal(t, "10.00", line.Item.Price.String())
	})
	t.Run("basic with currency", func(t *testing.T) {
		line := &Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:     "Test Item",
				Currency: currency.EUR,
				Price:    num.MakeAmount(10, 0),
			},
		}
		err := line.normalizeItemPrice(currency.EUR, exampleRates(t))
		require.NoError(t, err)
		assert.Equal(t, "10.00", line.Item.Price.String())
	})

	t.Run("alt prices", func(t *testing.T) {
		line := &Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:     "Test Item",
				Currency: currency.USD,
				Price:    num.MakeAmount(10, 0),
				AltPrices: []*currency.Amount{
					{
						Currency: currency.EUR,
						Value:    num.MakeAmount(8, 0),
					},
				},
			},
		}
		err := line.normalizeItemPrice(currency.EUR, exampleRates(t))
		require.NoError(t, err)
		assert.Equal(t, "8.00", line.Item.Price.String())
		require.Len(t, line.Item.AltPrices, 1)
		assert.Equal(t, "10.00", line.Item.AltPrices[0].Value.String())
		assert.Equal(t, "USD", string(line.Item.AltPrices[0].Currency))
	})

	t.Run("use exchange rates", func(t *testing.T) {
		line := &Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:     "Test Item",
				Currency: currency.USD,
				Price:    num.MakeAmount(10, 0),
			},
		}
		err := line.normalizeItemPrice(currency.EUR, exampleRates(t))
		require.NoError(t, err)
		assert.Equal(t, "8.76", line.Item.Price.String())
		require.Len(t, line.Item.AltPrices, 1)
		assert.Equal(t, "10.00", line.Item.AltPrices[0].Value.String())
		assert.Equal(t, "USD", string(line.Item.AltPrices[0].Currency))
	})

	t.Run("missing exchange rate", func(t *testing.T) {
		line := &Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:     "Test Item",
				Currency: currency.MXN,
				Price:    num.MakeAmount(100, 0),
			},
		}
		err := line.normalizeItemPrice(currency.EUR, exampleRates(t))
		assert.ErrorContains(t, err, "no exchange rate found from 'MXN' to 'EUR")
	})
}
