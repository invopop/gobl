package bill

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
)

// Ordering contains information about the commercial process and documents used to prepare
// the invoice, including orders, contracts, projects, preceding documents, and parties
// involved in issuing or addressing the invoice.
type Ordering struct {
	/*** IDENTIFIERS ***/

	// Identifier assigned by the customer or buyer for internal routing purposes.
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`
	// Any additional Codes, IDs, SKUs, or other regional or custom
	// identifiers that may be used to identify the order.
	Identities []*org.Identity `json:"identities,omitempty" jsonschema:"title=Identities"`
	// Buyer accounting reference cost code associated with the document.
	Cost cbc.Code `json:"cost,omitempty" jsonschema:"title=Cost,example=1287:65464"`
	// Period of time that the document refers to often used in addition to the details
	// provided in the individual line items.
	Period *cal.Period `json:"period,omitempty" jsonschema:"title=Period"`

	/*** PARTIES ***/

	// Party who is responsible for issuing payment, if not the same as the customer.
	Buyer *org.Party `json:"buyer,omitempty" jsonschema:"title=Buyer"`
	// Party liable to pay taxes on the transaction, if not the same as the supplier.
	Seller *org.Party `json:"seller,omitempty" jsonschema:"title=Seller"`
	// Third party responsible for issuing the document, but not responsible for tax.
	// Required by some tax regimes and formats.
	Issuer *org.Party `json:"issuer,omitempty" jsonschema:"title=Issuer"`
	// Third party to whom the document is addressed, but who is not the customer.
	Addressee *org.Party `json:"addressee,omitempty" jsonschema:"title=Addressee"`

	/*** DOCUMENTS ***/

	// Projects this document refers to.
	Projects []*org.DocumentRef `json:"projects,omitempty" jsonschema:"title=Projects"`
	// Contracts related to the document.
	Contracts []*org.DocumentRef `json:"contracts,omitempty" jsonschema:"title=Contracts"`
	// Purchase orders issued by the customer or buyer.
	Purchases []*org.DocumentRef `json:"purchases,omitempty" jsonschema:"title=Purchase Orders"`
	// Sales orders issued by the supplier or seller.
	Sales []*org.DocumentRef `json:"sales,omitempty" jsonschema:"title=Sales Orders"`
	// Receiving advice documents.
	Receiving []*org.DocumentRef `json:"receiving,omitempty" jsonschema:"title=Receiving Advice"`
	// Despatch advice documents.
	Despatch []*org.DocumentRef `json:"despatch,omitempty" jsonschema:"title=Despatch Advice"`
	// Call for tender or lot the document relates to.
	Tender []*org.DocumentRef `json:"tender,omitempty" jsonschema:"title=Tender Advice"`
}
