package pay

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/jsonschema"
)

// TermKey defines the type of terms being handled
type TermKey org.Key

// Terms defines when we expect the customer to pay, or have paid, for
// the contents of the document.
type Terms struct {
	// Type of terms to be applied.
	Key TermKey `json:"key" jsonschema:"title=Key"`
	// Text detail of the chosen payment terms.
	Detail string `json:"detail,omitempty" jsonschema:"title=Detail"`
	// Set of dates for agreed payments.
	DueDates []*DueDate `json:"due_dates,omitempty" jsonschema:"title=Due Dates"`
	// Description of the conditions for payment.
	Notes string `json:"notes,omitempty" jsonschema:"title=Notes"`
}

// Pre-defined Payment Terms based on UNTDID 4279
const (
	// None defined
	TermKeyNA TermKey = ""
	// End of Month
	TermKeyEndOfMonth TermKey = "end-of-month"
	// Due on a specific date
	TermKeyDueDate TermKey = "due-date"
	// Deferred until after the due dates
	TermKeyDeferred TermKey = "deferred"
	// Month after the present
	TermKeyProximo TermKey = "proximo"
	// on receipt of invoice
	TermKeyInstant TermKey = "instant"
	// chosen by buyer
	TermKeyElective TermKey = "elective"
	// Seller to advise buyer in separate transaction
	TermKeyPending TermKey = "pending"
	// Payment made in advance
	TermKeyAdvance TermKey = "advance"
	// Payment on Delivery
	TermKeyDelivery TermKey = "delivery"
)

// TermKeyDef holds a definition of a single payment term key
type TermKeyDef struct {
	// The key being defined
	Key TermKey `json:"key" jsonschema:"Key"`
	// Human text for the key
	Description string `json:"description" jsonschema:"Description"`
	// The equivalent UNTDID 4279 Code
	UNTDID4279 org.Code `json:"untdid4279" jsonschema:"UNTDID 4279 Code"`
}

// TermKeyDefinitions includes all the currently accepted
// GOBL Payment Term definitions.
var TermKeyDefinitions = []TermKeyDef{
	{TermKeyNA, "Not yet defined", "16"},
	{TermKeyEndOfMonth, "End of month", "2"},
	{TermKeyDueDate, "Due on a specific date", "3"},
	{TermKeyDeferred, "Deferred until after the due date", "4"},
	{TermKeyProximo, "Month after the present", "9"},
	{TermKeyInstant, "On receipt of invoice", "10"},
	{TermKeyElective, "Chosen by the buyer", "11"},
	{TermKeyPending, "Seller to advise buyer in separate transaction", "13"},
	{TermKeyAdvance, "Payment made in advance", "32"},
	{TermKeyDelivery, "Payment on Delivery", "52"}, // Cash on Delivery (COD)
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
func (t *Terms) UNTDID4279() org.Code {
	for _, v := range TermKeyDefinitions {
		if t.Key == v.Key {
			return v.UNTDID4279
		}
	}
	return org.CodeEmpty
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
		validation.Field(&t.Key, isValidTermKey),
		validation.Field(&t.DueDates),
	)
}

var isValidTermKey = validation.In(validTermKeys()...)

func validTermKeys() []interface{} {
	list := make([]interface{}, len(TermKeyDefinitions))
	for i, v := range TermKeyDefinitions {
		list[i] = v.Key
	}
	return list
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

// JSONSchema provides a representation of the struct for usage in Schema.
func (TermKey) JSONSchema() *jsonschema.Schema {
	s := &jsonschema.Schema{
		Title:       "Term Key",
		OneOf:       make([]*jsonschema.Schema, len(TermKeyDefinitions)),
		Description: "Payment terms key",
	}
	for i, v := range TermKeyDefinitions {
		s.OneOf[i] = &jsonschema.Schema{
			Const:       v.Key,
			Description: v.Description,
		}
	}
	return s
}
