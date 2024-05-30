package bill

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
)

// Predefined list of the payment types supported.
const (
	PaymentTypeReceipt    = "receipt"
	PaymentTypeRemittance = "remittance"
	PaymentTypeRefund     = "refund"
)

// PaymentTypes defines the list of potential payment types.
var PaymentTypes = []*cbc.KeyDefinition{
	{
		Key: PaymentTypeReceipt,
		Name: i18n.String{
			i18n.EN: "Receipt",
		},
		Desc: i18n.String{
			i18n.EN: "A payment receipt sent from the supplier to a customer reflecting that the referenced documents have been paid.",
		},
	},
	{
		Key: PaymentTypeRemittance,
		Name: i18n.String{
			i18n.EN: "Remittance",
		},
		Desc: i18n.String{
			i18n.EN: "A remittance advice sent from the customer to the supplier reflecting that the referenced documents have been paid.",
		},
	},
	{
		Key: PaymentTypeRefund,
		Name: i18n.String{
			i18n.EN: "Refund",
		},
		Desc: i18n.String{
			i18n.EN: "A refund payment sent from the supplier to a customer reflecting that the referenced documents have been refunded.",
		},
	},
}

var isValidPaymentType = validation.In(validPaymentTypes()...)

func validPaymentTypes() []interface{} {
	list := make([]interface{}, len(InvoiceTypes))
	for i, d := range InvoiceTypes {
		list[i] = d.Key
	}
	return list
}

// A Payment is used to link an invoice or invoices with a payment transaction.
type Payment struct {
	uuid.Identify

	// Type of payment document being issued.
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

	Documents []*DocumentReference `json:"lines" jsonschema:"title=Lines"`

	// Summary of all the payment totals, including tax calculations (calculated)
	Totals *Totals `json:"totals" jsonschema:"title=Totals" jsonschema_extras:"calculated=true"`

	// Unstructured information that is relevant to the invoice, such as correction or additional
	// legal details.
	Notes []*cbc.Note `json:"notes,omitempty" jsonschema:"title=Notes"`

	// Additional complementary objects that add relevant information to the invoice.
	Complements []*schema.Object `json:"complements,omitempty" jsonschema:"title=Complements"`

	// Additional semi-structured data that doesn't fit into the body of the invoice.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}
