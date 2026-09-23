package org

import (
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
)

const (
	// ItemKeyServices indicates that the item is a service.
	ItemKeyServices cbc.Key = "services"
	// ItemKeyGoods indicates that the item is a physical good.
	ItemKeyGoods cbc.Key = "goods"
)

// Item is used to describe a single product or service. Minimal usage
// implies just adding a name and price, more complete usage consists
// of adding descriptions, supplier IDs, SKUs, dimensions, etc.
//
// A set of additional code, ID, or SKU can be included in the `identities` property.
// Each `Identity` can be defined with an optional type agreed upon between the
// supplier and customer.
// For general purpose use, the Item's `Ref` property is easier to use.
type Item struct {
	uuid.Identify
	// Primary reference code that identifies this item.
	// Additional codes can be provided in the 'identities' property.
	Ref cbc.Code `json:"ref,omitempty" jsonschema:"title=Ref"`
	// Special key used to classify the item sometimes required by some regimes.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// Brief name of the item
	Name string `json:"name" jsonschema:"title=Name"`
	// List of additional codes, IDs, or SKUs which can be used to identify the item. They should be agreed upon between supplier and customer.
	Identities []*Identity `json:"identities,omitempty" jsonschema:"title=Identities"`
	// Attributes describe named features or properties of the item, such as
	// color or size.
	Attributes []*Attribute `json:"attributes,omitempty" jsonschema:"title=Attributes"`
	// Detailed description of the item.
	Description string `json:"description,omitempty" jsonschema:"title=Description"`
	// Images associated with the item.
	Images []*Image `json:"images,omitempty" jsonschema:"title=Images"`
	// Currency used for the item's price.
	Currency currency.Code `json:"currency,omitempty" jsonschema:"title=Currency"`
	// List price before the price discount is applied. When set, the price
	// is calculated as the list price less the discount.
	List *num.Amount `json:"list,omitempty" jsonschema:"title=List Price"`
	// Amount deducted from the list price to determine the price.
	Discount *num.Amount `json:"discount,omitempty" jsonschema:"title=Price Discount"`
	// Net price to be sold at, for a single unit or the number of units defined
	// by per. Must be either zero or positive.
	Price *num.Amount `json:"price,omitempty" jsonschema:"title=Net Price"`
	// AltPrices defines a list of prices with their currencies that may be used
	// as an alternative to the item's base price.
	AltPrices []*currency.Amount `json:"alt_prices,omitempty" jsonschema:"title=Alternative Prices"`
	// Number of units the prices apply to, e.g. 100 for a price per 100 kg.
	// Assumed to be 1 when empty.
	Per *num.Amount `json:"per,omitempty" jsonschema:"title=Price Base Quantity"`
	// Unit of measure using a GOBL key. Standard UN/ECE codes may be preserved
	// in the untdid-unit extension.
	Unit cbc.Key `json:"unit,omitempty" jsonschema:"title=Unit"`
	// Country code of where this item was from originally.
	Origin l10n.ISOCountryCode `json:"origin,omitempty" jsonschema:"title=Country of Origin"`
	// Extension code map for any additional regime specific codes that may be required.
	Ext tax.Extensions `json:"ext,omitzero" jsonschema:"title=Extensions"`
	// Additional meta information that may be useful
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

func itemRules() *rules.Set {
	return rules.For(new(Item),
		rules.Field("name",
			rules.Assert("01", "item name is required", is.Present),
		),
		rules.Field("price",
			rules.AssertIfPresent("02", "item price must be zero or positive", num.ZeroOrPositive),
		),
		rules.Field("attributes",
			rules.Assert("03", "item attributes must not contain duplicate keys",
				AttributesHaveUniqueKeys(),
			),
		),
		rules.Field("unit",
			rules.AssertIfPresent("04", "item unit must be valid", HasValidUnitKey),
		),
		rules.Field("per",
			rules.AssertIfPresent("05", "item per must be positive", num.Positive),
		),
		rules.Field("list",
			rules.AssertIfPresent("06", "item list price must be zero or positive", num.ZeroOrPositive),
		),
		rules.Field("discount",
			rules.AssertIfPresent("07", "item price discount must be zero or positive", num.ZeroOrPositive),
		),
		rules.When(is.Expr(`Discount != nil`),
			rules.Assert("08", "item price discount requires a list price",
				is.Expr(`List != nil`),
			),
			rules.Assert("09", "item price discount must not exceed the list price",
				is.Func("discount within list price", itemDiscountWithinList),
			),
		),
	)
}

func itemDiscountWithinList(val any) bool {
	i, ok := val.(*Item)
	if !ok || i.List == nil || i.Discount == nil {
		return true
	}
	return i.Discount.Compare(*i.List) <= 0
}

// JSONSchemaExtend adds extra details to the schema.
func (Item) JSONSchemaExtend(js *jsonschema.Schema) {
	ExtendUnitKeySchema(js, "unit")
	prop, ok := js.Properties.Get("key")
	if ok {
		prop.AnyOf = []*jsonschema.Schema{
			{
				Const: ItemKeyGoods,
				Title: "Goods",
			},
			{
				Const: ItemKeyServices,
				Title: "Services",
			},
			{
				Title:   "Other",
				Pattern: cbc.KeyPattern,
			},
		}
	}
}

func normalizeItem(i *Item) {
	i.Unit, i.Ext = normalizeUnit(i.Unit, i.Ext)
	i.Name = cbc.NormalizeString(i.Name)
	i.Description = cbc.NormalizeString(i.Description)
	i.Attributes = CleanAttributes(i.Attributes)
	normalizeItemPrice(i)
}

func normalizeItemPrice(i *Item) {
	p := i.PriceFromList()
	if p == nil {
		return
	}
	if i.Price != nil {
		*p = p.MatchPrecision(*i.Price)
	}
	i.Price = p
}

// PriceFromList provides the list price less the discount, or nil if there
// is no list price.
func (i *Item) PriceFromList() *num.Amount {
	if i == nil || i.List == nil {
		return nil
	}
	p := *i.List
	if i.Discount != nil {
		p = p.MatchPrecision(*i.Discount).Subtract(*i.Discount)
	}
	return &p
}

// PerUnit divides an amount that applies to the item's Per number of units
// so that it applies to a single unit. The amount is upscaled to cover the
// decimal places the division may introduce.
func (i *Item) PerUnit(a num.Amount) num.Amount {
	if i == nil || i.Per == nil || !i.Per.IsPositive() {
		return a
	}
	extra := uint32(len(strconv.FormatInt(i.Per.Rescale(0).Value(), 10)))
	return a.Upscale(extra).Divide(*i.Per)
}

// UnitPrice provides the item's price for a single unit, dividing by the
// Per value when set, or nil if the item has no price.
func (i *Item) UnitPrice() *num.Amount {
	if i == nil || i.Price == nil {
		return nil
	}
	p := i.PerUnit(*i.Price)
	return &p
}
