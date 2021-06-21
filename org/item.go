package org

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
)

// Item is used to describe a single object or service. Minimal usage
// implies just adding a name and price, more complete usage consists
// of adding descriptions, supplier and client IDs, dimensions, etc.
//
// All prices of items should be set as their "net" price, i.e. without
// tax.
//
// The taxes themselves change according to the recipient, so they should
// be defined in their usage context.
type Item struct {
	UUID        string       `json:"uuid,omitempty" jsonschema:"title=UUID,description=Unique identify of this item independent of the Supplier IDs"`
	Name        string       `json:"name"`
	Description string       `json:"desc,omitempty"`
	Currency    string       `json:"currency,omitempty" jsonschema:"title=Currency,description=Only required if this line has a different currency from the rest."`
	Price       num.Amount   `json:"price" jsonschema:"title=Price,description=Price of item being sold."`
	Unit        string       `json:"unit,omitempty" jsonschema:"title=Unit,description=Code for unit of the item being sold"`
	SupplierIDs []*ItemID    `json:"supplier_ids,omitempty" jsonschema:"title=Supplier IDs"`
	Origin      l10n.Country `json:"origin,omitempty" jsonschema:"title=Country of Origin,description=Country code of where this item was from originally."`
	Meta        Meta         `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// ItemID contains
type ItemID struct {
	Value string `json:"value"`
}
