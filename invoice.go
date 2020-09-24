package gobl

// Invoice represents a payment claim for goods or services supplied under
// conditions agreed between the supplier and the customer. In most cases
// the resulting document describes the actual financial commitment of goods
// or services ordered from the supplier.
type Invoice struct {
	ID       string        `json:"id" jsonschema:"title=ID"`
	Code     string        `json:"code" jsonschema:"title=Code"`
	Currency string        `json:"currency" jsonschema:"title=Currency,description=Currency for all invoice totals."`
	Rates    ExchangeRates `json:"rates,omitempty" jsonschema:"title=Exchange Rates,description=Simple array of values used to convert other currencies into the invoice's main currency."`

	IssueDate     *Date `json:"issue_date" jsonschema:"title=Issue Date"`
	OperationDate *Date `json:"op_date,omitempty" jsonschema:"title=Operation Date"`
	ValueDate     *Date `json:"value_date" jsonschema:"title=Value Date"`

	Supplier *Party `json:"supplier"`
	Customer *Party `json:"customer"`

	Lines InvoiceLines `json:"lines,omitempty"`

	Taxes  *InvoiceTaxes  `json:"taxes" jsonschema:"title=Taxes"`
	Totals *InvoiceTotals `json:"totals" jsonschema:"title=Totals"`

	Payment *InvoicePayment `json:"payment,omitempty" jsonschema:"title=Payment Details"`
}

// ExchangeRates represents an array of currency exchange rates
type ExchangeRates []ExchangeRate

// ExchangeRate contains data on the rate to be used when converting data.
// The rate is always multipled by the
type ExchangeRate struct {
	Currency string `json:"currency"`
	Rate     Amount `json:"rate"`
}

// InvoiceLines holds an array of InvoiceLine objects.
type InvoiceLines []*InvoiceLine

// InvoiceLine represents a single row in an invoice.
type InvoiceLine struct {
	ID       string            `json:"id,omitempty"`
	Quantity Amount            `json:"amount"`
	Item     *Item             `json:"item"`
	Taxes    []*InvoiceLineTax `json:"taxes,omitempty"`
}

// InvoiceLineTax describes a single type of tax applied to the line,
// including the rate that should be used.
type InvoiceLineTax struct {
	Code TaxCode `json:"code"`
	Rate Amount  `json:"rate"`
}

type InvoiceTaxes struct {
}

type InvoiceTotals struct {
}

type InvoicePayment struct {
}

// Type provides the body type used for mapping.
func (i *Invoice) Type() BodyType {
	return BodyTypeInvoice
}
