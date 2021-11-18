package org

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
)

// Item is used to describe a single product or service. Minimal usage
// implies just adding a name and price, more complete usage consists
// of adding descriptions, supplier IDs, SKUs, dimensions, etc.
//
// A set of additional code, ID, or SKU can be included in the `codes` property.
// Each `ItemCode` can be defined with an optional type agreed upon between the
// supplier and customer.
// For general purpose use, the Item's `Ref` property is much
// easier to use.
//
// We recommend setting prices with the item's "net" value, without tax,
// unless the document you're building supports the `price_includes_tax`
// option included in the `bill.Invoice` definition for example.
type Item struct {
	UUID        string       `json:"uuid,omitempty" jsonschema:"title=UUID,description=Unique identify of this item independent of the Supplier IDs"`
	Ref         string       `json:"ref,omitempty" jsonschema:"title=Ref,description=Primary reference code that identifies this item. Additional codes can be provided in the 'codes' field."`
	Name        string       `json:"name"`
	Description string       `json:"desc,omitempty"`
	Currency    string       `json:"currency,omitempty" jsonschema:"title=Currency,description=Only required if this line has a different currency from the rest."`
	Price       num.Amount   `json:"price" jsonschema:"title=Price,description=Price of item being sold."`
	Unit        string       `json:"unit,omitempty" jsonschema:"title=Unit,description=Code for unit of the item being sold"`
	Codes       []*ItemCode  `json:"codes,omitempty" jsonschema:"title=Codes,description=List of additional codes, IDs, or SKUs which can be used to identify the item. The should be agreed upon between supplier and customer."`
	Origin      l10n.Country `json:"origin,omitempty" jsonschema:"title=Country of Origin,description=Country code of where this item was from originally."`
	Meta        Meta         `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// ItemCode contains a value and optional type property that means additional
// codes can be added to an item.
type ItemCode struct {
	Type  string `json:"typ,omitempty" jsonschema:"title=Type"`
	Value string `json:"val" jsonschema:"title=Value"`
}
