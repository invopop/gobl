package bill

import (
	"context"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
)

// Predefined list of the receipt types supported.
const (
	ReceiptTypePayment    = "payment"
	ReceiptTypeRemittance = "remittance"
)

// ReceiptTypes defines the list of potential payment types.
var ReceiptTypes = []*cbc.KeyDefinition{
	{
		Key: ReceiptTypePayment,
		Name: i18n.String{
			i18n.EN: "Payment",
		},
		Desc: i18n.String{
			i18n.EN: "A payment receipt sent from the supplier to a customer reflecting that the referenced documents have been paid.",
		},
	},
	{
		Key: ReceiptTypeRemittance,
		Name: i18n.String{
			i18n.EN: "Remittance",
		},
		Desc: i18n.String{
			i18n.EN: "A remittance advice sent from the customer to the supplier reflecting that the referenced documents have been paid.",
		},
	},
}

var isValidReceiptType = validation.In(validReceiptTypes()...)

func validReceiptTypes() []interface{} {
	list := make([]interface{}, len(ReceiptTypes))
	for i, d := range ReceiptTypes {
		list[i] = d.Key
	}
	return list
}

// A Receipt is used to link an invoice or invoices with a payment transaction.
type Receipt struct {
	uuid.Identify

	// Type of receipt document being issued.
	Type cbc.Key `json:"type" jsonschema:"title=Type" jsonschema_extras:"calculated=true"`
	// Used as a prefix to group codes.
	Series string `json:"series,omitempty" jsonschema:"title=Series"`
	// Sequential code used to identify this payment in tax declarations.
	Code string `json:"code" jsonschema:"title=Code"`
	// When the payment was issued.
	IssueDate cal.Date `json:"issue_date" jsonschema:"title=Issue Date" jsonschema_extras:"calculated=true"`
	// When the taxes of this payment become accountable, if none set, the issue date is assumed.
	ValueDate *cal.Date `json:"value_date,omitempty" jsonschema:"title=Value Date"`
	// Currency for all payment totals.
	Currency currency.Code `json:"currency" jsonschema:"title=Currency" jsonschema_extras:"calculated=true"`
	// Exchange rates to be used when converting the payment's monetary values into other currencies.
	ExchangeRates []*currency.ExchangeRate `json:"exchange_rates,omitempty" jsonschema:"title=Exchange Rates"`

	// The taxable entity who is responsible for supplying goods or services.
	Supplier *org.Party `json:"supplier" jsonschema:"title=Supplier"`
	// Legal entity that receives the goods or services.
	Customer *org.Party `json:"customer,omitempty" jsonschema:"title=Customer"`

	// List of lines that detail what has been paid for.
	Lines []*Line `json:"lines" jsonschema:"title=Lines"`

	// Summary of all the payment totals, including tax calculations (calculated)
	Totals *Totals `json:"totals" jsonschema:"title=Totals" jsonschema_extras:"calculated=true"`

	// Unstructured information that is relevant to the receipt, such as correction or additional
	// legal details.
	Notes []*cbc.Note `json:"notes,omitempty" jsonschema:"title=Notes"`

	// Additional complementary objects that add relevant information to the receipt.
	Complements []*schema.Object `json:"complements,omitempty" jsonschema:"title=Complements"`

	// Additional semi-structured data that doesn't fit into the body of the invoice.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// ValidateWithContext ensures that the fields contained in the Receipt look correct.
func (d *Receipt) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, d,
		validation.Field(&d.Type, validation.Required, isValidReceiptType),
	)
}
