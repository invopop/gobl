package bill

import (
	"errors"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/region"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Invoice represents a payment claim for goods or services supplied under
// conditions agreed between the supplier and the customer. In most cases
// the resulting document describes the actual financial commitment of goods
// or services ordered from the supplier.
type Invoice struct {
	// Unique document ID. Not required, but always recommended in addition to the Code.
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Sequential code used to identify this invoice in tax declarations.
	Code string `json:"code" jsonschema:"title=Code"`
	// Used in addition to the Code in some regions.
	Series string `json:"series,omitempty" jsonschema:"title=Series"`
	// Functional type of the invoice, default is always 'commercial'.
	TypeCode TypeCode `json:"type_code,omitempty" jsonschema:"title=Type Code"`
	// Currency for all invoice totals.
	Currency currency.Code `json:"currency" jsonschema:"title=Currency"`
	// Exchange rates to be used when converting the invoices monetary values into other currencies.
	ExchangeRates currency.ExchangeRates `json:"rates,omitempty" jsonschema:"title=Exchange Rates"`
	// When true, implies that all item prices already include non-retained taxes. This is especially
	// useful for retailers where prices are often displayed including tax.
	PricesIncludeTax bool `json:"prices_include_tax,omitempty" jsonschema:"title=Prices Include Tax"`

	// Key information regarding a previous invoice.
	Preceding *Preceding `json:"preceding,omitempty" jsonschema:"title=Preceding Reference"`

	// When the invoice was created.
	IssueDate *org.Date `json:"issue_date" jsonschema:"title=Issue Date"`
	// Date when the operation defined by the invoice became effective.
	OperationDate *org.Date `json:"op_date,omitempty" jsonschema:"title=Operation Date"`
	// When the taxes of this invoice become accountable, if none set, the issue date is used.
	ValueDate *org.Date `json:"value_date,omitempty" jsonschema:"title=Value Date"`

	// The taxable entity supplying the goods or services.
	Supplier *org.Party `json:"supplier" jsonschema:"title=Supplier"`
	// Legal entity receiving the goods or services, may be empty in certain circumstances such as simplified invoices.
	Customer *org.Party `json:"customer,omitempty" jsonschema:"title=Customer"`

	// List of invoice lines representing each of the items sold to the customer.
	Lines Lines `json:"lines,omitempty" jsonschema:"title=Lines"`
	// Discounts or allowances applied to the complete invoice
	Discounts Discounts `json:"discounts,omitempty" jsonschema:"title=Discounts"`
	// Charges or surcharges applied to the complete invoice
	Charges Charges `json:"charges,omitempty" jsonschema:"title=Charges"`
	// Expenses paid for by the supplier but invoiced directly to the customer.
	Outlays Outlays `json:"outlays,omitempty" jsonschema:"title=Outlays"`

	// Summary of all the invoice totals, including taxes.
	Totals *Totals `json:"totals" jsonschema:"title=Totals"`

	Ordering *Ordering `json:"ordering,omitempty" jsonschema:"title=Ordering Details"`
	Payment  *Payment  `json:"payment,omitempty" jsonschema:"title=Payment Details"`
	Delivery *Delivery `json:"delivery,omitempty" jsonschema:"title=Delivery Details"`

	// The EN 16931-1:2017 standard recognises a need to be able to attach additional
	// documents to an invoice. We don't support this yet, but this is where
	// it could go.
	//Attachments Attachments `json:"attachments,omitempty" jsonschema:"title=Attachments"`

	// Unstructured information that is relevant to the invoice, such as correction details.
	Notes org.Notes `json:"notes,omitempty" jsonschema:"title=Notes"`
	// Additional semi-structured data that doesn't fit into the body of the invoice.
	Meta org.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Totals contains the summaries of all calculations for the invoice.
type Totals struct {
	// Sum of all line item sums
	Sum num.Amount `json:"sum" jsonschema:"title=Sum"`
	// Sum of all document level discounts
	Discount *num.Amount `json:"discount,omitempty" jsonschema:"title=Discount"`
	// Sum of all document level charges
	Charge *num.Amount `json:"charge,omitempty" jsonschema:"title=Charge"`
	// Sum of all line sums minus the discounts
	Total num.Amount `json:"total" jsonschema:"title=Total"`
	// Summary of all the taxes with a final sum to add or deduct from the amount payable
	Taxes *tax.Total `json:"taxes,omitempty" jsonschema:"title=Tax Totals"`
	// Total paid in outlays that need to be reimbursed
	Outlays *num.Amount `json:"outlays,omitempty" jsonschema:"title=Outlay Totals"`
	// Total amount to be paid after applying taxes
	Payable num.Amount `json:"payable" jsonschema:"title=Payable"`
	// Total amount paid in advance
	Advances *num.Amount `json:"advance,omitempty" jsonschema:"title=Advance"`
	// How much actually needs to be paid now
	Due *num.Amount `json:"due,omitempty" jsonschema:"title=Due"`
}

// Ordering allows additional order details to be appended
type Ordering struct {
	// Party who is selling the goods and is not responsible for taxes
	Seller *org.Party `json:"seller,omitempty" jsonschema:"title=Seller"`
}

// Payment contains details as to how the invoice should be paid.
type Payment struct {
	Payer        *org.Party        `json:"payer,omitempty" jsonschema:"title=Payer,description=The party responsible for paying for the invoice, if not the customer."`
	Terms        *pay.Terms        `json:"terms,omitempty" jsonschema:"title=Terms,description=Payment terms or conditions."`
	Advances     []*pay.Advance    `json:"advances,omitempty" jsonschema:"title=Advances,description=Any amounts that have been paid in advance and should be deducted from the amount due."`
	Instructions *pay.Instructions `json:"instructions,omitempty" jsonschema:"title=Instructions,description=Details on how payment should be made."`
}

// Delivery covers the details of the destination for the products described
// in the invoice body.
type Delivery struct {
	// The party who will receive delivery of the goods defined in the invoice and is not responsible for taxes.
	Receiver *org.Party `json:"receiver,omitempty" jsonschema:"title=Receiver"`
	// When the goods should be expected
	Date *org.Date `json:"date,omitempty" jsonschema:"title=Date"`
	// Start of a n invoicing or delivery period
	StartDate *org.Date `json:"start_date,omitempty" jsonschema:"title=Start Date"`
	// End of a n invoicing or delivery period
	EndDate *org.Date `json:"end_date,omitempty" jsonschema:"title=End Date"`
}

// Preceding allows for information to be provided about a previous invoice that this one
// will replace or subtract from. If this is used, the invoice type code will most likely need
// to be set to `corrected` or `credit-note`.
type Preceding struct {
	// Preceding document's UUID if available can be useful for tracing.
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Identity code fo the previous invoice.
	Code string `json:"code" jsonschema:"title=Code"`
	// When the preceding invoice was issued.
	IssueDate *org.Date `json:"issue_date" jsonschema:"title=Issue Date"`
	// Additional semi-structured data that may be useful in specific regions
	Meta org.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Type provides the body type used for mapping.
func (Invoice) Type() string {
	return InvoiceType
}

// Validate checks to ensure the invoice is valid and contains all the information we need.
func (inv *Invoice) Validate(r region.Region) error {
	err := validation.ValidateStruct(inv,
		validation.Field(&inv.UUID),
		validation.Field(&inv.Code, validation.Required),
		validation.Field(&inv.TypeCode), // either empty (Commercial) or one of those supported
		validation.Field(&inv.Currency, validation.Required),
		validation.Field(&inv.IssueDate, validation.Required),

		validation.Field(&inv.Supplier, validation.Required),
		validation.Field(&inv.Customer),

		validation.Field(&inv.Lines, validation.Required),
		validation.Field(&inv.Discounts),
		validation.Field(&inv.Charges),
		validation.Field(&inv.Outlays),
		validation.Field(&inv.Totals, validation.Required),
	)
	if err == nil {
		err = r.Validate(inv)
	}
	return err
}

// Calculate performs all the calculations required for the invoice totals and taxes. If the original
// invoice only includes partial calculations, this will figure out what's missing.
func (inv *Invoice) Calculate(r region.Region) error {
	date := inv.ValueDate
	if date == nil {
		date = inv.IssueDate
	}
	if date == nil {
		return errors.New("issue date cannot be empty")
	}

	// Prepare the totals we'll need with amounts based on currency
	t := new(Totals)
	zero := r.Currency().BaseAmount()
	t.reset(zero)

	tr := r.Taxes()
	tls := make([]tax.TaxableLine, 0)

	// Ensure all the lines are up to date first
	for i, l := range inv.Lines {
		l.Index = i + 1
		l.calculate()

		// Basic sum
		t.Sum = t.Sum.Add(l.Total)
		tls = append(tls, l)
	}
	t.Total = t.Sum

	// Subtract discounts
	discounts := zero
	for i, l := range inv.Discounts {
		l.Index = i + 1
		if l.Rate != nil && !l.Rate.IsZero() {
			l.Amount = l.Rate.Of(t.Sum)
		}
		discounts = discounts.Add(l.Amount)
		tls = append(tls, l)
	}
	if !discounts.IsZero() {
		t.Discount = &discounts
		t.Total = t.Total.Subtract(discounts)
	}

	// Add charges
	charges := zero
	for i, l := range inv.Charges {
		l.Index = i + 1
		if l.Rate != nil && !l.Rate.IsZero() {
			l.Amount = l.Rate.Of(t.Sum)
		}
		charges = charges.Add(l.Amount)
		tls = append(tls, l)
	}
	if !charges.IsZero() {
		t.Charge = &charges
		t.Total = t.Total.Add(charges)
	}

	// Now figure out the tax totals (with some interface conversion)
	if err := t.Taxes.Calculate(tr, tls, inv.PricesIncludeTax, *date, zero); err != nil {
		return err
	}

	t.Payable = t.Total.Add(t.Taxes.Sum)

	// Outlays
	if len(inv.Outlays) > 0 {
		t.Outlays = &zero
		for i, o := range inv.Outlays {
			o.Index = i + 1
			v := t.Outlays.Add(o.Paid)
			t.Outlays = &v
		}
		t.Payable = t.Payable.Add(*t.Outlays)
	}

	if inv.Payment != nil {
		// Deal with advances, if any
		if t.Advances = inv.Payment.totalAdvance(zero); t.Advances != nil {
			v := t.Payable.Subtract(*t.Advances)
			t.Due = &v
		}

		// Calculate any due date amounts
		inv.Payment.Terms.CalculateDues(t.Payable)
	}

	inv.Totals = t
	return nil
}

func (p *Payment) totalAdvance(zero num.Amount) *num.Amount {
	if p == nil || len(p.Advances) == 0 {
		return nil
	}
	sum := zero
	for _, a := range p.Advances {
		sum = sum.Add(a.Amount)
	}
	return &sum
}

// Reset sets all the totals to the provided zero amount with the correct
// decimal places.
func (t *Totals) reset(zero num.Amount) {
	t.Sum = zero
	t.Discount = nil
	t.Charge = nil
	t.Total = zero
	t.Taxes = tax.NewTotal(zero)
	t.Outlays = nil
	t.Payable = zero
	t.Advances = nil
	t.Due = nil
}
