package bill

import (
	"errors"
	"fmt"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/region"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"

	validation "github.com/go-ozzo/ozzo-validation"
)

const invoiceType = "bill.Invoice"

// Invoice represents a payment claim for goods or services supplied under
// conditions agreed between the supplier and the customer. In most cases
// the resulting document describes the actual financial commitment of goods
// or services ordered from the supplier.
type Invoice struct {
	UUID          uuid.UUID              `json:"uuid" jsonschema:"title=UUID"`
	Code          string                 `json:"code" jsonschema:"title=Code,description=Sequential ID used to identify this invoice in tax declarations."`
	Region        region.Code            `json:"region" jsonschema:"title=Region,description=GoBL region code used to determine taxes and validation rules."`
	Currency      currency.Code          `json:"currency" jsonschema:"title=Currency,description=Currency for all invoice totals."`
	ExchangeRates currency.ExchangeRates `json:"rates,omitempty" jsonschema:"title=Exchange Rates,description=Exchange rates to be used when converting the invoices monetary values into other currencies."`

	IssueDate     *org.Date `json:"issue_date" jsonschema:"title=Issue Date"`
	OperationDate *org.Date `json:"op_date,omitempty" jsonschema:"title=Operation Date"`
	ValueDate     *org.Date `json:"value_date" jsonschema:"title=Value Date"`

	Supplier *org.Party `json:"supplier" jsonschema:"title=Supplier,description=The taxable entity supplying the goods or services."`
	Customer *org.Party `json:"customer" jsonschema:"title=Customer,description=Legal entity who receives the goods or services."`

	Lines   Lines   `json:"lines,omitempty" jsonschema:"title=Lines,description=The items sold to the customer."`
	Outlays Outlays `json:"outlays,omitempty" jsonschema:"title=Outlays,description=Expenses paid for by the supplier but invoiced directly to the customer."`

	Totals *Totals `json:"totals" jsonschema:"title=Totals"`

	Ordering *Ordering `json:"ordering,omitempty" jsonschema:"title=Ordering Details"`
	Payment  *Payment  `json:"payment,omitempty" jsonschema:"title=Payment Details"`
	Delivery *Delivery `json:"delivery,omitempty" jsonschema:"title=Delivery Details"`
}

// Totals contains the summaries of all calculations for the invoice.
type Totals struct {
	Sum      num.Amount `json:"sum" jsonschema:"title=Sum,description=Sum of all line item sums"`
	Discount num.Amount `json:"discount,omitempty" jsonschema:"title=Discount,description=Sum of all discounts applied to each line."`
	Total    num.Amount `json:"total,omitempty" jsonschema:"title=Total,description=Sum of all line sums minus the discounts."`
	Taxes    *tax.Total `json:"taxes,omitempty" jsonschema:"title=Tax Totals,description=Summary of all taxes with a final sum to add or deduct from the amount payable."`
	Outlays  num.Amount `json:"outlays,omitempty" jsonschema:"title=Outlay Totals,description=Total paid in outlays that need to be reimbursed."`
	Payable  num.Amount `json:"payable" jsonschema:"title=Payable,description=Total amount to be paid after applying taxes."`
}

// Ordering allows additional order details to be appended
type Ordering struct {
	Seller *org.Party `json:"seller,omitempty" jsonschema:"title=Seller,description=Party who is selling the goods and is not responsible for taxes."`
}

// Payment contains details as to how the invoice should be paid.
// TODO: Add terms here.
type Payment struct {
	Terms   *pay.Terms    `json:"terms,omitempty" jsonschema:"title=Terms,description=Payment terms or conditions."`
	Methods []*pay.Method `json:"methods,omitempty" jsonschema:"title=Methods,description=Array of payment options that can be used to pay for this invoice."`
	Payer   *org.Party    `json:"payer,omitempty" jsonschema:"title=Payer,description=The party responsible for paying for the invoice, if not the customer."`
}

// InvoiceDelivery covers the details of the destination for the products described
// in the invoice body.
type Delivery struct {
	Receiver *org.Party `json:"receiver,omitempty" jsonschema:"title=Receiver,description=The party who will receive delivery of the goods defined in the invoice and is not responsible for taxes."`
}

// Type provides the body type used for mapping.
func (Invoice) Type() string {
	return invoiceType
}

// Validate checks to ensure the invoice is valid and contains all the information we need.
func (inv *Invoice) Validate() error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.UUID, validation.Required, uuid.IsV1),
		validation.Field(&inv.Code, validation.Required),
		validation.Field(&inv.Region, validation.Required),
		validation.Field(&inv.Currency, validation.Required),
		validation.Field(&inv.IssueDate, validation.Required),
		validation.Field(&inv.ValueDate, validation.Required),

		validation.Field(&inv.Supplier, validation.Required),
		validation.Field(&inv.Customer, validation.Required),

		validation.Field(&inv.Lines, validation.Required),
		validation.Field(&inv.Totals, validation.Required),
	)
}

// Calculate performs all the calculations required for the invoice totals and taxes. If the original
// invoice only includes partial calculations, this will figure out what's missing.
func (inv *Invoice) Calculate() error {
	r := region.For(inv.Region)
	if r == nil {
		return errors.New("unknown invoice region code")
	}

	if inv.ValueDate == nil {
		inv.ValueDate = inv.IssueDate
	}
	if inv.ValueDate == nil {
		return errors.New("value or issue date cannot be empty")
	}

	// Prepare the totals we'll need with amounts based on currency
	t := new(Totals)
	zero := r.Currency().BaseAmount()
	t.reset(zero)

	// Ensure all the lines are up to date first
	tr := r.Taxes()
	for i, l := range inv.Lines {
		l.Index = i + 1
		if err := l.calculate(tr, *inv.ValueDate); err != nil {
			return fmt.Errorf("line %d: %w", l.Index, err)
		}

		// Basic sum
		t.Sum = t.Sum.Add(l.Sum)
		if l.Discount != nil {
			t.Discount = t.Discount.Add(l.Discount.Value)
		}

		// Taxes
		for _, r := range l.Taxes {
			if err := t.Taxes.AddRate(r, zero); err != nil {
				return fmt.Errorf("line %d: %w", l.Index, err)
			}
		}
	}
	t.Total = t.Sum.Add(t.Discount)

	// Outlays
	for i, o := range inv.Outlays {
		o.Index = i + 1
		t.Outlays = t.Outlays.Add(o.Paid)
	}

	t.Payable = t.Total.Add(t.Taxes.Sum).Add(t.Outlays)
	inv.Totals = t
	return nil
}

// Reset sets all the totals to the provided zero amount with the correct
// decimal places.
func (t *Totals) reset(zero num.Amount) {
	t.Sum = zero
	t.Discount = zero
	t.Taxes = tax.NewTotal(zero)
	t.Total = zero
	t.Outlays = zero
	t.Payable = zero
}
