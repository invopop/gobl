package bill

import (
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

const invoiceBodyType = "bill.Invoice"

// Invoice represents a payment claim for goods or services supplied under
// conditions agreed between the supplier and the customer. In most cases
// the resulting document describes the actual financial commitment of goods
// or services ordered from the supplier.
type Invoice struct {
	UUID       string                 `json:"uuid" jsonschema:"title=UUID"`
	Code       string                 `json:"code" jsonschema:"title=Code,description=Sequential ID used to identify this invoice in tax declarations."`
	RegionCode tax.Code               `json:"region_code" jsonschema:"title=Region Code,description=Region used for tax purposes."`
	Currency   currency.Code          `json:"currency" jsonschema:"title=Currency,description=Currency for all invoice totals."`
	Rates      currency.ExchangeRates `json:"rates,omitempty" jsonschema:"title=Exchange Rates,description=Simple array of values used to convert other currencies into the invoice's main currency."`

	IssueDate     *org.Date `json:"issue_date" jsonschema:"title=Issue Date"`
	OperationDate *org.Date `json:"op_date,omitempty" jsonschema:"title=Operation Date"`
	ValueDate     *org.Date `json:"value_date" jsonschema:"title=Value Date"`

	Supplier *org.Party `json:"supplier" jsonschema:"title=Supplier,description=The taxable entity supplying the goods or services."`
	Customer *org.Party `json:"customer" jsonschema:"title=Customer,description=Legal entity who receives the goods or services."`

	Lines Lines `json:"lines,omitempty"`

	Totals *Totals `json:"totals" jsonschema:"title=Totals"`

	Ordering *Ordering `json:"ordering,omitempty" jsonschema:"title=Ordering Details"`
	Payment  *Payment  `json:"payment,omitempty" jsonschema:"title=Payment Details"`
	Delivery *Delivery `json:"delivery,omitempty" jsonschema:"title=Delivery Details"`
}

// InvoiceLines holds an array of InvoiceLine objects.
type Lines []*Line

// Line represents a single row in an invoice.
type Line struct {
	UUID         string         `json:"uuid,omitempty"`
	Quantity     num.Amount     `json:"quantity"`
	Item         *org.Item      `json:"item"`
	Sum          num.Amount     `json:"sum" jsonschema:"title=Sum,description=Result of quantity multiplied by item price"`
	DiscountRate num.Percentage `json:"discount_rate,omitempty" jsonschema:"title=Discount Rate,description=Percentage discount applied to sum."`
	Discount     num.Amount     `json:"discount,omitempty" jsonschema:"title=Discount,description=Total discount applied to the line."`
	Taxes        []tax.Rates    `json:"taxes,omitempty" jsonschema:"title=Taxes,description=List of taxes to be applied to the line in the invoice totals."`
	Total        num.Amount     `json:"total" jsonschema:"title=Total,description=Total line amount after applying discounts to the sum."`
}

// Totals contains the summaries of all calculations for the invoice.
type Totals struct {
	Sum      num.Amount `json:"sum" jsonschema:"title=Sum,description=Sum of all line item sums"`
	Discount num.Amount `json:"discount,omitempty" jsonschema:"title=Discount,description=Sum of all discounts applied to each line."`
	Total    num.Amount `json:"total,omitempty" jsonschema:"title=Total,description=Sum of all line sums minus the discounts."`
	Taxes    tax.Totals `json:"taxes,omitempty"`
	Payable  num.Amount `json:"payable" jsonschema:"title=Payable,description=Total amount to be paid after applying taxes."`
}

// Ordering allows additional order details to be appended
type Ordering struct {
	Seller *org.Party `json:"seller,omitempty" jsonschema:"title=Seller,description=Party who is selling the goods and is not responsible for taxes."`
}

// Payment contains details as to how the invoice should be paid.
// TODO: Add terms here.
type Payment struct {
	Payer *org.Party `json:"payer,omitempty" jsonschema:"title=Payer,description=The party responsible for paying for the invoice."`
}

// InvoiceDelivery covers the details of the destination for the products described
// in the invoice body.
type Delivery struct {
	Receiver *org.Party `json:"receiver,omitempty" jsonschema:"title=Receiver,description=The party who will receive delivery of the goods defined in the invoice and is not responsible for taxes."`
}

// BodyType provides the body type used for mapping.
func (i *Invoice) BodyType() string {
	return invoiceBodyType
}
