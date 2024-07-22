package bill

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

// Ordering provides additional information about the ordering process including references
// to other documents and alternative parties involved in the order-to-delivery process.
type Ordering struct {
	// Identifier assigned by the customer or buyer for internal routing purposes.
	Code string `json:"code,omitempty" jsonschema:"title=Code"`
	// Any additional Codes, IDs, SKUs, or other regional or custom
	// identifiers that may be used to identify the order.
	Identities []*org.Identity `json:"identities,omitempty" jsonschema:"title=Identities"`

	// Period of time that the invoice document refers to often used in addition to the details
	// provided in the individual line items.
	Period *cal.Period `json:"period,omitempty" jsonschema:"title=Period"`
	// Project this invoice refers to.
	Project *org.Document `json:"project,omitempty" jsonschema:"title=Project"`
	// The identification of a contract.
	Contract *org.Document `json:"contract,omitempty" jsonschema:"title=Contract"`
	// Purchase order issued by the customer or buyer.
	Purchase *org.Document `json:"purchase,omitempty" jsonschema:"title=Purchase Order"`
	// Sales order issued by the supplier or seller.
	Sale *org.Document `json:"sale,omitempty" jsonschema:"title=Sales Order"`
	// Receiving Advice.
	Receiving *org.Document `json:"receiving,omitempty" jsonschema:"title=Receiving Advice"`
	// Despatch advice.
	Despatch *org.Document `json:"despatch,omitempty" jsonschema:"title=Despatch Advice"`
	// Tender advice, the identification of the call for tender or lot the invoice relates to.
	Tender *org.Document `json:"tender,omitempty" jsonschema:"title=Tender Advice"`

	// Party who is responsible for making the purchase, but is not responsible
	// for handling taxes.
	Buyer *org.Party `json:"buyer,omitempty" jsonschema:"title=Buyer"`
	// Party who is selling the goods but is not responsible for taxes like the
	// supplier.
	Seller *org.Party `json:"seller,omitempty" jsonschema:"title=Seller"`
}

// Validate the ordering details.
func (o *Ordering) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Identities),
		validation.Field(&o.Project),
		validation.Field(&o.Contract),
		validation.Field(&o.Purchase),
		validation.Field(&o.Sale),
		validation.Field(&o.Receiving),
		validation.Field(&o.Despatch),
		validation.Field(&o.Tender),
		validation.Field(&o.Buyer),
		validation.Field(&o.Seller),
	)
}
