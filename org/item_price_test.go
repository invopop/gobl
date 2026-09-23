package org_test

import (
	"testing"

	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItemPriceNormalize(t *testing.T) {
	t.Run("list less discount sets price", func(t *testing.T) {
		i := &org.Item{
			Name:     "Item",
			List:     num.NewAmount(12000, 2),
			Discount: num.NewAmount(1200, 2),
		}
		norm.Normalize(i)
		require.NotNil(t, i.Price)
		assert.Equal(t, "108.00", i.Price.String())
	})
	t.Run("list without discount", func(t *testing.T) {
		i := &org.Item{
			Name:  "Item",
			Price: num.NewAmount(100, 0),
			List:  num.NewAmount(12000, 2),
		}
		norm.Normalize(i)
		assert.Equal(t, "120.00", i.Price.String())
	})
	t.Run("keeps price precision", func(t *testing.T) {
		i := &org.Item{
			Name:     "Item",
			Price:    num.NewAmount(108000, 3),
			List:     num.NewAmount(120, 0),
			Discount: num.NewAmount(12, 0),
		}
		norm.Normalize(i)
		assert.Equal(t, "108.000", i.Price.String())
	})
	t.Run("per only leaves price", func(t *testing.T) {
		i := &org.Item{
			Name:  "Item",
			Price: num.NewAmount(12000, 2),
			Per:   num.NewAmount(100, 0),
		}
		norm.Normalize(i)
		assert.Equal(t, "120.00", i.Price.String())
	})
}

func TestItemPriceValidation(t *testing.T) {
	item := func() *org.Item {
		return &org.Item{Name: "Item", Price: num.NewAmount(100, 2)}
	}
	t.Run("valid", func(t *testing.T) {
		i := item()
		i.Per = num.NewAmount(100, 0)
		i.List = num.NewAmount(12000, 2)
		i.Discount = num.NewAmount(1200, 2)
		assert.NoError(t, rules.Validate(i))
	})
	t.Run("zero per", func(t *testing.T) {
		i := item()
		i.Per = num.NewAmount(0, 0)
		assert.ErrorContains(t, rules.Validate(i), "item per must be positive")
	})
	t.Run("negative list price", func(t *testing.T) {
		i := item()
		i.List = num.NewAmount(-100, 2)
		assert.ErrorContains(t, rules.Validate(i), "item list price must be zero or positive")
	})
	t.Run("negative discount", func(t *testing.T) {
		i := item()
		i.List = num.NewAmount(100, 2)
		i.Discount = num.NewAmount(-10, 2)
		assert.ErrorContains(t, rules.Validate(i), "item price discount must be zero or positive")
	})
	t.Run("discount without list price", func(t *testing.T) {
		i := item()
		i.Discount = num.NewAmount(10, 2)
		assert.ErrorContains(t, rules.Validate(i), "item price discount requires a list price")
	})
	t.Run("discount exceeds list price", func(t *testing.T) {
		i := item()
		i.List = num.NewAmount(100, 2)
		i.Discount = num.NewAmount(101, 2)
		assert.ErrorContains(t, rules.Validate(i), "item price discount must not exceed the list price")
	})
}

func TestItemPriceFromList(t *testing.T) {
	t.Run("nil item", func(t *testing.T) {
		var i *org.Item
		assert.Nil(t, i.PriceFromList())
	})
	t.Run("no list price", func(t *testing.T) {
		i := &org.Item{Name: "Item", Price: num.NewAmount(120, 2)}
		assert.Nil(t, i.PriceFromList())
	})
	t.Run("list less discount", func(t *testing.T) {
		i := &org.Item{Name: "Item", List: num.NewAmount(120, 0), Discount: num.NewAmount(1250, 2)}
		assert.Equal(t, "107.50", i.PriceFromList().String())
	})
}

func TestItemUnitPrice(t *testing.T) {
	t.Run("nil item", func(t *testing.T) {
		var i *org.Item
		assert.Nil(t, i.UnitPrice())
	})
	t.Run("no price", func(t *testing.T) {
		i := &org.Item{Name: "Item"}
		assert.Nil(t, i.UnitPrice())
	})
	t.Run("no per", func(t *testing.T) {
		i := &org.Item{Name: "Item", Price: num.NewAmount(120, 2)}
		assert.Equal(t, "1.20", i.UnitPrice().String())
	})
	t.Run("per 100", func(t *testing.T) {
		i := &org.Item{
			Name:  "Item",
			Price: num.NewAmount(12000, 2),
			Per:   num.NewAmount(100, 0),
		}
		assert.Equal(t, "1.20000", i.UnitPrice().String())
	})
	t.Run("fractional per", func(t *testing.T) {
		i := &org.Item{
			Name:  "Item",
			Price: num.NewAmount(500, 2),
			Per:   num.NewAmount(25, 1),
		}
		assert.Equal(t, "2.000", i.UnitPrice().String())
	})
}
