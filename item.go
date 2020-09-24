package gobl

// Item is used to describe a single object or service. Minimal usage
// implies just adding a name and price, more complete usage consists
// of adding descriptions, supplier and client IDs, dimensions, etc.
// Taxes can be complex, so they should be defined in the line item.
type Item struct {
	Name        string    `json:"name"`
	Description string    `json:"desc,omitempty"`
	Currency    string    `json:"currency,omitempty" jsonschema:"title=Currency,description=Only required if this line has a different currency from the rest."`
	Price       *Amount   `json:"price"`
	SupplierIDs []*ItemID `json:"supplier_ids" jsonschema:"title=Supplier IDs"`
	Origin      Country   `json:"origin,omitempty" jsconschema:"title=Country of Origin,description=Country code of where this item was from originally."`
}

// ItemID contains
type ItemID struct {
	Value string `json:"value"`
}
