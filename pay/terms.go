package pay

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
)

// Terms defines when we expect the customer to pay, or have paid, for
// the contents of the document.
type Terms struct {
	// Type of terms to be applied.
	Key org.Key `json:"key" jsonschema:"title=Key"`
	// Text detail of the chosen payment terms.
	Detail string `json:"detail,omitempty" jsonschema:"title=Detail"`
	// Set of dates for agreed payments.
	DueDates []*DueDate `json:"due_dates,omitempty" jsonschema:"title=Due Dates"`
	// Description of the conditions for payment.
	Notes string `json:"notes,omitempty" jsonschema:"title=Notes"`
}

// Pre-defined Payment Terms based on UNTDID 4279
const (
	TermKeyNA         org.Key = ""             // None defined
	TermKeyEndOfMonth org.Key = "end-of-month" // End of Month
	TermKeyDueDate    org.Key = "due-date"     // Due on a specific date
	TermKeyDeferred   org.Key = "deferred"     // Deferred until after the due dates
	TermKeyProximo    org.Key = "proximo"      // Month after the present
	TermKeyInstant    org.Key = "instant"      // on receipt of invoice
	TermKeyElective   org.Key = "elective"     // chosen by buyer
	TermKeyPending    org.Key = "pending"      // Seller to advise buyer in separate transaction
	TermKeyAdvance    org.Key = "advance"      // Payment made in advance
	TermKeyDelivery   org.Key = "delivery"     // Payment on Delivery
)

// Source: https://service.unece.org/trade/untdid/d15b/tred/tred4279.htm
var untdid4279Terms = map[org.Key]string{
	TermKeyNA:         "16", // Not Yet Defined
	TermKeyEndOfMonth: "2",  // End of month
	TermKeyDueDate:    "3",  // Fixed date
	TermKeyDeferred:   "4",  // Deferred
	TermKeyProximo:    "9",  // Proximo
	TermKeyInstant:    "10", // Instant
	TermKeyElective:   "11", // Elective
	TermKeyPending:    "13", // Seller to advise buyer
	TermKeyAdvance:    "32", // Advanced payment
	TermKeyDelivery:   "52", // Cash on Delivery (COD)
}

// DueDate contains an amount that should be paid by the given date.
type DueDate struct {
	Date     *cal.Date       `json:"date" jsonschema:"title=Date,description=When the payment is due."`
	Notes    string          `json:"notes,omitempty" jsonschema:"title=Notes,description=Other details to take into account for the due date."`
	Amount   num.Amount      `json:"amount" jsonschema:"title=Amount,description=How much needs to be paid by the date."`
	Percent  *num.Percentage `json:"percent,omitempty" jsonschema:"title=Percent,description=Percentage of the total that should be paid by the date."`
	Currency currency.Code   `json:"currency,omitempty" jsonschema:"title=Currency,description=If different from the parent document's base currency."`
}

// UNTDID4279 returns the UNTDID 4270 code associated with the terms key.
func (t *Terms) UNTDID4279() string {
	return untdid4279Terms[t.Key]
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
	validTermKeys := make([]interface{}, len(untdid4279Terms))
	i := 0
	for v := range untdid4279Terms {
		validTermKeys[i] = v
		i++
	}
	return validation.ValidateStruct(t,
		validation.Field(&t.Key, validation.In(validTermKeys...)),
		validation.Field(&t.DueDates),
	)
}

// Validate checks the DueDate has the required fields.
func (dd *DueDate) Validate() error {
	return validation.ValidateStruct(dd,
		validation.Field(&dd.Date, validation.Required),
		validation.Field(&dd.Amount, validation.Required),
		validation.Field(&dd.Percent),
		validation.Field(&dd.Currency),
	)
}
