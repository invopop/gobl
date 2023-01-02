package org

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/uuid"

	validation "github.com/go-ozzo/ozzo-validation/v4"
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
	// Unique identify of this item independent of the Supplier IDs
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Primary reference code that identifies this item. Additional codes can be provided in the 'codes' field.
	Ref string `json:"ref,omitempty" jsonschema:"title=Ref"`
	// Brief name of the item
	Name string `json:"name"`
	// Detailed description
	Description string `json:"desc,omitempty"`
	// Currency used for the item's price.
	Currency string `json:"currency,omitempty" jsonschema:"title=Currency"`
	// Base price of a single unit to be sold.
	Price num.Amount `json:"price" jsonschema:"title=Price"`
	// Unit of measure.
	Unit Unit `json:"unit,omitempty" jsonschema:"title=Unit"`
	//	List of additional codes, IDs, or SKUs which can be used to identify the item. The should be agreed upon between supplier and customer.
	Codes []*ItemCode `json:"codes,omitempty" jsonschema:"title=Codes"`
	// Country code of where this item was from originally.
	Origin l10n.CountryCode `json:"origin,omitempty" jsonschema:"title=Country of Origin"`
	// Additional meta information that may be useful
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// ItemCode contains a value and optional label property that means additional
// codes can be added to an item.
type ItemCode struct {
	// Local or human reference for the type of code the value
	// represents.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// The item code's value.
	Value string `json:"value" jsonschema:"title=Value"`
}

// Validate checks that an address looks okay.
func (i *Item) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.UUID),
		validation.Field(&i.Name, validation.Required),
		validation.Field(&i.Price, validation.Required),
		validation.Field(&i.Unit),
		validation.Field(&i.Codes),
		validation.Field(&i.Origin),
		validation.Field(&i.Meta),
	)
}
