package pay

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
)

// Terms defines when we expect the customer to pay, or have paid, for
// the contents of the document.
type Terms struct {
	Code     TermCode   `json:"code" jsonschema:"title=Code,description=Type of terms to be applied."`
	Detail   string     `json:"detail,omitempty" jsonschema:"title=Detail,description=Text detail of the chosen payment terms."`
	DueDates []*DueDate `json:"due_dates,omitempty" jsonschema:"title=Due Dates,description=Set of dates for agreed payments."`
	Notes    string     `json:"notes,omitempty" jsonschema:"title=Notes,description=Description of the conditions for payment."`
}

// TermCode is used to define a code that identifies the payment terms.
type TermCode string

// Pre-defined Payment Terms based on UNTDID 4279
const (
	TermNA         TermCode = "na"           // None defined
	TermEndOfMonth TermCode = "end_of_month" // End of Month
	TermDueDate    TermCode = "due_date"     // Due on a specific date
	TermDeferred   TermCode = "deferred"     // Deferred until after the due dates
	TermProximo    TermCode = "proximo"      // Month after the present
	TermInstant    TermCode = "instant"      // on receipt of invoice
	TermElective   TermCode = "elective"     // chosen by buyer
	TermPending    TermCode = "pending"      // Seller to advise buyer in separate transaction
	TermAdvance    TermCode = "advance"      // Payment made in advance
	TermDelivery   TermCode = "delivery"     // Payment on Delivery
)

// Source: https://service.unece.org/trade/untdid/d15b/tred/tred4279.htm
var untdid4279Terms = map[TermCode]string{
	TermNA:         "16", // Not Yet Defined
	TermEndOfMonth: "2",  // End of month
	TermDueDate:    "3",  // Fixed date
	TermDeferred:   "4",  // Deferred
	TermProximo:    "9",  // Proximo
	TermInstant:    "10", // Instant
	TermElective:   "11", // Elective
	TermPending:    "13", // Seller to advise buyer
	TermAdvance:    "32", // Advanced payment
	TermDelivery:   "52", // Cash on Delivery (COD)
}

// Validate checks to ensure the typecode is part of a known list.
func (c TermCode) Validate() error {
	if string(c) == "" {
		return nil
	}
	for k := range untdid4279Terms {
		if k == c {
			return nil
		}
	}
	return errors.New("invalid term code")
}

// DueDate contains an amount that should be paid by the given date.
type DueDate struct {
	Date     *cal.Date       `json:"date" jsonschema:"title=Date,description=When the payment is due."`
	Notes    string          `json:"notes,omitempty" jsonschema:"title=Notes,description=Other details to take into account for the due date."`
	Amount   num.Amount      `json:"amount" jsonschema:"title=Amount,description=How much needs to be paid by the date."`
	Percent  *num.Percentage `json:"percent,omitempty" jsonschema:"title=Percent,description=Percentage of the total that should be paid by the date."`
	Currency currency.Code   `json:"currency,omitempty" jsonschema:"title=Currency,description=If different from the parent document's base currency."`
}

// CalculateDues goes through each DueDate. If it has a percentage
// value set, it'll be used to calculate the amount.
func (t *Terms) CalculateDues(sum num.Amount) {
	if t == nil {
		return
	}
	for _, dd := range t.DueDates {
		if dd.Percent != nil && !dd.Percent.IsZero() {
			dd.Amount = dd.Percent.Of(sum)
		}
	}
}

// Validate ensures that the terms contain everything required.
func (t *Terms) Validate() error {
	return validation.ValidateStruct(t,
		validation.Field(&t.Code, validation.Required),
	)
}

// Validate checks the DueDate has the required fields.
func (dd *DueDate) Validate() error {
	return validation.ValidateStruct(dd,
		validation.Field(&dd.Date, validation.Required),
		validation.Field(&dd.Amount, validation.Required),
	)
}
