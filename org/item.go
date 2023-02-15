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
// A set of additional code, ID, or SKU can be included in the `identities` property.
// Each `Identity` can be defined with an optional type agreed upon between the
// supplier and customer.
// For general purpose use, the Item's `Ref` property is easier to use.
type Item struct {
	// Unique identity of this item
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Primary reference code that identifies this item.
	// Additional codes can be provided in the 'identities' property.
	Ref string `json:"ref,omitempty" jsonschema:"title=Ref"`
	// Brief name of the item
	Name string `json:"name"`
	// List of additional codes, IDs, or SKUs which can be used to identify the item. They should be agreed upon between supplier and customer.
	Identities []*Identity `json:"identities,omitempty" jsonschema:"title=Identities"`
	// Detailed description
	Description string `json:"desc,omitempty"`
	// Currency used for the item's price.
	Currency string `json:"currency,omitempty" jsonschema:"title=Currency"`
	// Base price of a single unit to be sold.
	Price num.Amount `json:"price" jsonschema:"title=Price"`
	// Unit of measure.
	Unit Unit `json:"unit,omitempty" jsonschema:"title=Unit"`
	// Country code of where this item was from originally.
	Origin l10n.CountryCode `json:"origin,omitempty" jsonschema:"title=Country of Origin"`
	// Additional meta information that may be useful
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate checks that an address looks okay.
func (i *Item) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.UUID),
		validation.Field(&i.Name, validation.Required),
		validation.Field(&i.Identities),
		validation.Field(&i.Price, validation.Required),
		validation.Field(&i.Unit),
		validation.Field(&i.Origin),
		validation.Field(&i.Meta),
	)
}
