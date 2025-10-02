package bill

import (
	"context"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Ordering provides additional information about the ordering process including references
// to other documents and alternative parties involved in the order-to-delivery process.
type Ordering struct {
	// Identifier assigned by the customer or buyer for internal routing purposes.
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`
	// Any additional Codes, IDs, SKUs, or other regional or custom
	// identifiers that may be used to identify the order.
	Identities []*org.Identity `json:"identities,omitempty" jsonschema:"title=Identities"`
	// Buyer accounting reference cost code associated with the document.
	Cost cbc.Code `json:"cost,omitempty" jsonschema:"title=Cost,example=1287:65464"`
	// Period of time that the invoice document refers to often used in addition to the details
	// provided in the individual line items.
	Period *cal.Period `json:"period,omitempty" jsonschema:"title=Period"`
	// Party who is responsible for issuing payment, if not the same as the customer.
	Buyer *org.Party `json:"buyer,omitempty" jsonschema:"title=Buyer"`
	// Seller is the party liable to pay taxes on the transaction if not the same as the supplier.
	Seller *org.Party `json:"seller,omitempty" jsonschema:"title=Seller"`
	// Issuer represents a third party responsible for issuing the invoice, but is not
	// responsible for tax. Some tax regimes and formats require this field.
	Issuer *org.Party `json:"issuer,omitempty" jsonschema:"title=Issuer"`
	// Projects this invoice refers to.
	Projects []*org.DocumentRef `json:"projects,omitempty" jsonschema:"title=Projects"`
	// The identification of contracts.
	Contracts []*org.DocumentRef `json:"contracts,omitempty" jsonschema:"title=Contracts"`
	// Purchase orders issued by the customer or buyer.
	Purchases []*org.DocumentRef `json:"purchases,omitempty" jsonschema:"title=Purchase Orders"`
	// Sales orders issued by the supplier or seller.
	Sales []*org.DocumentRef `json:"sales,omitempty" jsonschema:"title=Sales Orders"`
	// Receiving Advice.
	Receiving []*org.DocumentRef `json:"receiving,omitempty" jsonschema:"title=Receiving Advice"`
	// Despatch advice.
	Despatch []*org.DocumentRef `json:"despatch,omitempty" jsonschema:"title=Despatch Advice"`
	// Tender advice, the identification of the call for tender or lot the invoice relates to.
	Tender []*org.DocumentRef `json:"tender,omitempty" jsonschema:"title=Tender Advice"`
}

// Normalize attempts to clean and normalize the Ordering data.
func (o *Ordering) Normalize(normalizers tax.Normalizers) {
	if o == nil {
		return
	}
	o.Code = cbc.NormalizeCode(o.Code)
	o.Cost = cbc.NormalizeCode(o.Cost)
	normalizers.Each(o)
	tax.Normalize(normalizers, o.Identities)
	tax.Normalize(normalizers, o.Projects)
	tax.Normalize(normalizers, o.Contracts)
	tax.Normalize(normalizers, o.Purchases)
	tax.Normalize(normalizers, o.Sales)
	tax.Normalize(normalizers, o.Receiving)
	tax.Normalize(normalizers, o.Despatch)
	tax.Normalize(normalizers, o.Tender)
	tax.Normalize(normalizers, o.Buyer)
	tax.Normalize(normalizers, o.Seller)
	tax.Normalize(normalizers, o.Issuer)
}

// Validate the ordering details.
func (o *Ordering) Validate() error {
	return o.ValidateWithContext(context.Background())
}

// ValidateWithContext the ordering details with context.
func (o *Ordering) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithContext(ctx, o,
		validation.Field(&o.Code),
		validation.Field(&o.Identities),
		validation.Field(&o.Cost),
		validation.Field(&o.Projects),
		validation.Field(&o.Contracts),
		validation.Field(&o.Purchases),
		validation.Field(&o.Sales),
		validation.Field(&o.Receiving),
		validation.Field(&o.Despatch),
		validation.Field(&o.Tender),
		validation.Field(&o.Buyer),
		validation.Field(&o.Seller),
		validation.Field(&o.Issuer),
	)
}
