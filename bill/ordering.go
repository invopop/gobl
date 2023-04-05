package bill

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// Ordering provides additional information about the ordering process including references
// to other documents and alternative parties involved in the order-to-delivery process.
type Ordering struct {
	// Identifier assigned by the customer or buyer for internal routing purposes.
	Code string `json:"code,omitempty" jsonschema:"title=Code"`
	// Period of time that the invoice document refers to.
	Period *cal.Period `json:"period,omitempty" jsonschema:"title=Period"`
	// Project this invoice refers to.
	Project *DocumentReference `json:"project,omitempty" jsonschema:"title=Project"`
	// The identification of a contract.
	Contract *DocumentReference `json:"contract,omitempty" jsonschema:"title=Contract"`
	// Purchase order issued by the customer or buyer.
	Purchase *DocumentReference `json:"purchase,omitempty" jsonschema:"title=Purchase Order"`
	// Sales order issued by the supplier or seller.
	Sale *DocumentReference `json:"sale,omitempty" jsonschena:"title=Sales Order"`
	// Receiving Advice.
	Receiving *DocumentReference `json:"receiving,omitempty" jsonschame:"title=Receiving Advice"`
	// Despatch advice.
	Despatch *DocumentReference `json:"despatch,omitempty" jsonschema:"title=Despatch Advice"`
	// Tender advice, the identification of the call for tender or lot the invoice relates to.
	Tender *DocumentReference `json:"tender,omitempty" jsonscheme:"title=Tender Advice"`

	// Party who is responsible for making the purchase, but is not responsible
	// for handling taxes.
	Buyer *org.Party `json:"buyer,omitempty" jsonschema:"title=Buyer"`
	// Party who is selling the goods but is not responsible for taxes like the
	// supplier.
	Seller *org.Party `json:"seller,omitempty" jsonschema:"title=Seller"`
}

// DocumentReference provides a link to a existing document.
type DocumentReference struct {
	// Unique ID copied from the source document.
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Series the reference document belongs to.
	Series string `json:"series,omitempty" jsonschema:"title=Series"`
	// Source document's code or other identifier.
	Code string `json:"code,omitempty" jsonschema:"title=Code"`
	// Link to the source document.
	URL string `json:"url,omitempty" jsonschema:"title=URL,format=uri"`
}

// Validate the ordering details.
func (o *Ordering) Validate() error {
	return validation.ValidateStruct(o,
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

// Validate ensures the Document Reference looks correct.
func (dr *DocumentReference) Validate() error {
	return validation.ValidateStruct(dr,
		validation.Field(&dr.UUID),
		validation.Field(&dr.Code),
		validation.Field(&dr.URL, is.URL),
	)
}
