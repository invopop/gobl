package org

import (
	"errors"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
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
	// Detailed description of the item.
	Description string `json:"description,omitempty" jsonschema:"title=Description"`
	// Images associated with the item.
	Images []*Image `json:"images,omitempty" jsonschema:"title=Images"`
	// Currency used for the item's price.
	Currency currency.Code `json:"currency,omitempty" jsonschema:"title=Currency"`
	// Base price of a single unit to be sold. Must be either zero or positive.
	Price *num.Amount `json:"price,omitempty" jsonschema:"title=Price"`
	// AltPrices defines a list of prices with their currencies that may be used
	// as an alternative to the item's base price.
	AltPrices []*currency.Amount `json:"alt_prices,omitempty" jsonschema:"title=Alternative Prices"`
	// Unit of measure.
	Unit Unit `json:"unit,omitempty" jsonschema:"title=Unit"`
	// Country code of where this item was from originally.
	Origin l10n.ISOCountryCode `json:"origin,omitempty" jsonschema:"title=Country of Origin"`
	// Extension code map for any additional regime specific codes that may be required.
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`
	// Additional meta information that may be useful
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Normalize performs any required normalizations on the Item.
func (i *Item) Normalize(normalizers tax.Normalizers) {
	if i == nil {
		return
	}
	i.Name = cbc.NormalizeString(i.Name)
	i.Description = cbc.NormalizeString(i.Description)
	i.Ref = cbc.NormalizeCode(i.Ref)
	i.Ext = tax.CleanExtensions(i.Ext)

	tax.Normalize(normalizers, i.Identities)
	tax.Normalize(normalizers, i.Images)
	normalizers.Each(i)
}

func itemRules() *rules.Set {
	return rules.For(new(Item),
		rules.Field("name",
			rules.Assert("01", "item name is required", is.Present),
		),
		rules.Field("price",
			rules.AssertIfPresent("02", "item price must be zero or positive", num.ZeroOrPositive),
		),
	)
}

type itemPriceValidator struct{}

// ItemPriceRequired ensures that the item has a price.
func ItemPriceRequired() validation.Rule {
	return &itemPriceValidator{}
}

// Validate ensures that the item has a price.
func (v *itemPriceValidator) Validate(value any) error {
	i, ok := value.(*Item)
	if i == nil || !ok {
		return nil
	}
	if i.Price == nil {
		return validation.Errors{
			"price": errors.New("cannot be blank"),
		}
	}
	return nil
}

// JSONSchemaExtend adds extra details to the schema.
func (Item) JSONSchemaExtend(js *jsonschema.Schema) {
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
