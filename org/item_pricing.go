package org

import (
	"strconv"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// ItemPricing describes how an item's price was determined.
type ItemPricing struct {
	// Number of the item's units the price applies to, e.g. 100 for a price
	// per 100 kg. Assumed to be 1 when empty.
	Per *num.Amount `json:"per,omitempty" jsonschema:"title=Per"`
	// Price before the discount is applied. When set, the item's price is
	// calculated as the gross price less the discount.
	Gross *num.Amount `json:"gross,omitempty" jsonschema:"title=Gross"`
	// Amount deducted from the gross price to determine the item's price.
	Discount *num.Amount `json:"discount,omitempty" jsonschema:"title=Discount"`
}

func itemPricingRules() *rules.Set {
	return rules.For(new(ItemPricing),
		rules.Field("per",
			rules.AssertIfPresent("01", "item pricing per must be positive", num.Positive),
		),
		rules.Field("gross",
			rules.AssertIfPresent("02", "item pricing gross must be zero or positive", num.ZeroOrPositive),
		),
		rules.Field("discount",
			rules.AssertIfPresent("03", "item pricing discount must be zero or positive", num.ZeroOrPositive),
		),
		rules.When(is.Expr(`Discount != nil`),
			rules.Assert("04", "item pricing discount requires a gross price",
				is.Expr(`Gross != nil`),
			),
			rules.Assert("05", "item pricing discount must not exceed the gross price",
				is.Func("discount within gross", itemPricingDiscountWithinGross),
			),
		),
	)
}

func itemPricingDiscountWithinGross(val any) bool {
	ip, ok := val.(*ItemPricing)
	if !ok || ip.Gross == nil || ip.Discount == nil {
		return true
	}
	return ip.Discount.Compare(*ip.Gross) <= 0
}

// IsEmpty returns true if no pricing details are set.
func (ip *ItemPricing) IsEmpty() bool {
	return ip == nil || (ip.Per == nil && ip.Gross == nil && ip.Discount == nil)
}

// PerUnit divides an amount expressed for the pricing's Per number of units
// so that it applies to a single unit. The amount is upscaled to cover the
// decimal places the division may introduce.
func (ip *ItemPricing) PerUnit(a num.Amount) num.Amount {
	if ip == nil || ip.Per == nil || !ip.Per.IsPositive() {
		return a
	}
	extra := uint32(len(strconv.FormatInt(ip.Per.Rescale(0).Value(), 10)))
	return a.Upscale(extra).Divide(*ip.Per)
}

// Net provides the gross price less the discount, or nil if there is no
// gross price.
func (ip *ItemPricing) Net() *num.Amount {
	if ip == nil || ip.Gross == nil {
		return nil
	}
	p := *ip.Gross
	if ip.Discount != nil {
		p = p.MatchPrecision(*ip.Discount).Subtract(*ip.Discount)
	}
	return &p
}
