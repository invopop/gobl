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

func TestItemPricingNormalize(t *testing.T) {
	t.Run("empty pricing removed", func(t *testing.T) {
		i := &org.Item{Name: "Item", Pricing: &org.ItemPricing{}}
		norm.Normalize(i)
		assert.Nil(t, i.Pricing)
	})
	t.Run("gross less discount sets price", func(t *testing.T) {
		i := &org.Item{
			Name: "Item",
			Pricing: &org.ItemPricing{
				Gross:    num.NewAmount(12000, 2),
				Discount: num.NewAmount(1200, 2),
			},
		}
		norm.Normalize(i)
		require.NotNil(t, i.Price)
		assert.Equal(t, "108.00", i.Price.String())
	})
	t.Run("gross without discount", func(t *testing.T) {
		i := &org.Item{
			Name:    "Item",
			Price:   num.NewAmount(100, 0),
			Pricing: &org.ItemPricing{Gross: num.NewAmount(12000, 2)},
		}
		norm.Normalize(i)
		assert.Equal(t, "120.00", i.Price.String())
	})
	t.Run("keeps price precision", func(t *testing.T) {
		i := &org.Item{
			Name:  "Item",
			Price: num.NewAmount(108000, 3),
			Pricing: &org.ItemPricing{
				Gross:    num.NewAmount(120, 0),
				Discount: num.NewAmount(12, 0),
			},
		}
		norm.Normalize(i)
		assert.Equal(t, "108.000", i.Price.String())
	})
	t.Run("per only leaves price", func(t *testing.T) {
		i := &org.Item{
			Name:    "Item",
			Price:   num.NewAmount(12000, 2),
			Pricing: &org.ItemPricing{Per: num.NewAmount(100, 0)},
		}
		norm.Normalize(i)
		assert.Equal(t, "120.00", i.Price.String())
	})
}

func TestItemPricingValidation(t *testing.T) {
	item := func(ip *org.ItemPricing) *org.Item {
		return &org.Item{Name: "Item", Price: num.NewAmount(100, 2), Pricing: ip}
	}
	t.Run("valid", func(t *testing.T) {
		assert.NoError(t, rules.Validate(item(&org.ItemPricing{
			Per:      num.NewAmount(100, 0),
			Gross:    num.NewAmount(12000, 2),
			Discount: num.NewAmount(1200, 2),
		})))
	})
	t.Run("zero per", func(t *testing.T) {
		err := rules.Validate(item(&org.ItemPricing{Per: num.NewAmount(0, 0)}))
		assert.ErrorContains(t, err, "item pricing per must be positive")
	})
	t.Run("negative gross", func(t *testing.T) {
		err := rules.Validate(item(&org.ItemPricing{Gross: num.NewAmount(-100, 2)}))
		assert.ErrorContains(t, err, "item pricing gross must be zero or positive")
	})
	t.Run("negative discount", func(t *testing.T) {
		err := rules.Validate(item(&org.ItemPricing{
			Gross:    num.NewAmount(100, 2),
			Discount: num.NewAmount(-10, 2),
		}))
		assert.ErrorContains(t, err, "item pricing discount must be zero or positive")
	})
	t.Run("discount without gross", func(t *testing.T) {
		err := rules.Validate(item(&org.ItemPricing{Discount: num.NewAmount(10, 2)}))
		assert.ErrorContains(t, err, "item pricing discount requires a gross price")
	})
	t.Run("discount exceeds gross", func(t *testing.T) {
		err := rules.Validate(item(&org.ItemPricing{
			Gross:    num.NewAmount(100, 2),
			Discount: num.NewAmount(101, 2),
		}))
		assert.ErrorContains(t, err, "item pricing discount must not exceed the gross price")
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
	t.Run("no pricing", func(t *testing.T) {
		i := &org.Item{Name: "Item", Price: num.NewAmount(120, 2)}
		assert.Equal(t, "1.20", i.UnitPrice().String())
	})
	t.Run("per 100", func(t *testing.T) {
		i := &org.Item{
			Name:    "Item",
			Price:   num.NewAmount(12000, 2),
			Pricing: &org.ItemPricing{Per: num.NewAmount(100, 0)},
		}
		assert.Equal(t, "1.20000", i.UnitPrice().String())
	})
	t.Run("fractional per", func(t *testing.T) {
		i := &org.Item{
			Name:    "Item",
			Price:   num.NewAmount(500, 2),
			Pricing: &org.ItemPricing{Per: num.NewAmount(25, 1)},
		}
		assert.Equal(t, "2.000", i.UnitPrice().String())
	})
}
