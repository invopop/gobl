package gobl

import (
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// Invoice represents a payment claim for goods or services supplied under
// conditions agreed between the supplier and the customer. In most cases
// the resulting document describes the actual financial commitment of goods
// or services ordered from the supplier.
type Invoice struct {
	UUID     string        `json:"uuid" jsonschema:"title=UUID"`
	Code     string        `json:"code" jsonschema:"title=Code"`
	RegionID tax.RegionID  `json:"region_id" jsonschema:"title=Region ID,description=Region used for tax purposes."`
	Currency string        `json:"currency" jsonschema:"title=Currency,description=Currency for all invoice totals."`
	Rates    ExchangeRates `json:"rates,omitempty" jsonschema:"title=Exchange Rates,description=Simple array of values used to convert other currencies into the invoice's main currency."`

	IssueDate     *Date `json:"issue_date" jsonschema:"title=Issue Date"`
	OperationDate *Date `json:"op_date,omitempty" jsonschema:"title=Operation Date"`
	ValueDate     *Date `json:"value_date" jsonschema:"title=Value Date"`

	Supplier *Party `json:"supplier"`
	Customer *Party `json:"customer"`

	Lines InvoiceLines `json:"lines,omitempty"`

	Totals *InvoiceTotals `json:"totals" jsonschema:"title=Totals"`

	Payment *InvoicePayment `json:"payment,omitempty" jsonschema:"title=Payment Details"`
}

// InvoiceLines holds an array of InvoiceLine objects.
type InvoiceLines []*InvoiceLine

// InvoiceLine represents a single row in an invoice.
type InvoiceLine struct {
	UUID         string         `json:"uuid,omitempty"`
	Quantity     num.Amount     `json:"quantity"`
	Item         *Item          `json:"item"`
	Sum          num.Amount     `json:"sum" jsonschema:"title=Sum,description=Result of quantity multiplied by item price"`
	DiscountRate num.Percentage `json:"discount_rate,omitempty" jsonschema:"title=Discount Rate,description=Percentage discount applied to sum."`
	Discount     num.Amount     `json:"discount,omitempty" jsonschema:"title=Discount,description=Total discount applied to the line."`
	Taxes        []Tax          `json:"taxes,omitempty" jsonschema:"title=Taxes,description=List of taxes to be applied to the line in the invoice totals."`
	Total        num.Amount     `json:"total" jsonschema:"title=Total,description=Total line amount after applying discounts to the sum."`
}

// InvoiceTotals contains the summaries of all calculations for the invoice.
type InvoiceTotals struct {
	Sum      num.Amount `json:"sum" jsonschema:"title=Sum,description=Sum of all line item sums"`
	Discount num.Amount `json:"discount,omitempty" jsonschema:"title=Discount,description=Sum of all discounts applied to each line."`
	Total    num.Amount `json:"total,omitempty" jsonschema:"title=Total,description=Sum of all line sums minus the discounts."`
	Tax      TaxTotal   `json:"tax,omitempty"`
	Payable  num.Amount `json:"payable" jsonschema:"title=Payable,description=Total amount to be paid after applying taxes."`
}

type InvoicePayment struct {
	Payer *Party `json:"payer,omitempty" jsconschema:"title=Payer,description=The party responsible for paying for the invoice."`
}

// Type provides the body type used for mapping.
func (i *Invoice) Type() BodyType {
	return BodyTypeInvoice
}
