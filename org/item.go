package org

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
)

// Item is used to describe a single object or service. Minimal usage
// implies just adding a name and price, more complete usage consists
// of adding descriptions, supplier and client IDs, dimensions, etc.
//
// If taxes should be considered to be included in the price, set the
// `TaxIncluded` property to true. This is especially important for consumer
// goods that form part of a simplified invoice. When possible, we recommend
// avoiding using this option as it can get confusing, especially when
// dealing with exported goods and services.
//
// The taxes themselves change according to the recipient, so they should
// be defined in their usage context.
type Item struct {
	UUID        string       `json:"uuid,omitempty" jsonschema:"title=UUID,description=Unique identify of this item independent of the Supplier IDs"`
	Name        string       `json:"name"`
	Description string       `json:"desc,omitempty"`
	Currency    string       `json:"currency,omitempty" jsonschema:"title=Currency,description=Only required if this line has a different currency from the rest."`
	Price       num.Amount   `json:"price"`
	TaxIncluded bool         `json:"tax_included,omitempty" jsonschema:"title=Tax Included,description=When true, the price should be considered to include taxes."`
	SupplierIDs []*ItemID    `json:"supplier_ids" jsonschema:"title=Supplier IDs"`
	Origin      l10n.Country `json:"origin,omitempty" jsonschema:"title=Country of Origin,description=Country code of where this item was from originally."`
	Meta        Meta         `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// ItemID contains
type ItemID struct {
	Value string `json:"value"`
}
