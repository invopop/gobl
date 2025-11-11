package pay

import (
	"context"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Terms defines when we expect the customer to pay, or have paid, for
// the contents of the document.
type Terms struct {
	// Type of terms to be applied.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// Text detail of the chosen payment terms (Deprecated).
	Detail string `json:"detail,omitempty" jsonschema:"title=Detail"`
	// Set of dates for agreed payments.
	DueDates []*DueDate `json:"due_dates,omitempty" jsonschema:"title=Due Dates"`
	// Description of the conditions for payment.
	Notes string `json:"notes,omitempty" jsonschema:"title=Notes"`
	// Extensions to the terms for local codes.
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`
}

// Pre-defined Payment Terms based on UNTDID 4279
const (
	// End of Month
	TermKeyEndOfMonth cbc.Key = "end-of-month"
	// Due on a specific date
	TermKeyDueDate cbc.Key = "due-date"
	// Deferred until after the due dates
	TermKeyDeferred cbc.Key = "deferred"
	// Month after the present
	TermKeyProximo cbc.Key = "proximo"
	// on receipt of invoice
	TermKeyInstant cbc.Key = "instant"
	// chosen by buyer
	TermKeyElective cbc.Key = "elective"
	// Seller to advise buyer in separate transaction
	TermKeyPending cbc.Key = "pending"
	// Payment made in advance
	TermKeyAdvanced cbc.Key = "advanced"
	// Payment on Delivery
	TermKeyDelivery cbc.Key = "delivery"
	// Not yet defined
	TermKeyUndefined cbc.Key = "undefined"
)

// TermKeyDef holds a definition of a single payment term key
type TermKeyDef struct {
	// The key being defined
	Key cbc.Key `json:"key" jsonschema:"title=Key"`
	// Human readable title for the key
	Title string `json:"title" jsonschema:"title=Title"`
	// Human text for the key
	Description string `json:"description" jsonschema:"title=Description"`
	// The equivalent UNTDID 4279 Code
	UNTDID4279 cbc.Code `json:"untdid4279" jsonschema:"title=UNTDID 4279 Code"`
}

// TermKeyDefinitions includes all the currently accepted
// GOBL Payment Term definitions.
var TermKeyDefinitions = []TermKeyDef{
	{TermKeyEndOfMonth, "End of Month", "End of month", "2"},
	{TermKeyDueDate, "Due Date", "Due on a specific date", "3"},
	{TermKeyDeferred, "Deferred", "Deferred until after the due date", "4"},
	{TermKeyProximo, "Proximo", "Month after the present", "9"},
	{TermKeyInstant, "Instant", "On receipt of invoice", "10"},
	{TermKeyElective, "Elective", "Chosen by the buyer", "11"},
	{TermKeyPending, "Pending", "Seller to advise buyer in separate transaction", "13"},
	{TermKeyAdvanced, "Advanced", "Payment made in advance", "32"},
	{TermKeyDelivery, "Delivery", "Payment on Delivery", "52"}, // Cash on Delivery (COD)
	{TermKeyUndefined, "Undefined", "Not yet defined", "16"},
}

// DueDate contains an amount that should be paid by the given date.
type DueDate struct {
	Date     *cal.Date       `json:"date" jsonschema:"title=Date,description=When the payment is due."`
	Notes    string          `json:"notes,omitempty" jsonschema:"title=Notes,description=Other details to take into account for the due date."`
	Amount   num.Amount      `json:"amount" jsonschema:"title=Amount,description=How much needs to be paid by the date."`
	Percent  *num.Percentage `json:"percent,omitempty" jsonschema:"title=Percent,description=Percentage of the total that should be paid by the date."`
	Currency currency.Code   `json:"currency,omitempty" jsonschema:"title=Currency,description=If different from the parent document's base currency."`
}

// Normalize will try to normalize the payment terms.
func (t *Terms) Normalize() {
	if t == nil {
		return
	}

	if t.Detail != "" && t.Notes == "" {
		t.Notes = t.Detail
		t.Detail = ""
	}

	t.Detail = cbc.NormalizeString(t.Detail)
	t.Notes = cbc.NormalizeString(t.Notes)
	t.Ext = tax.CleanExtensions(t.Ext)
}

// UNTDID4279 returns the UNTDID 4279 code associated with the terms key.
func (t *Terms) UNTDID4279() cbc.Code {
	for _, v := range TermKeyDefinitions {
		if t.Key == v.Key {
			return v.UNTDID4279
		}
	}
	return cbc.CodeEmpty
}

// CalculateDues goes through each DueDate. If it has a percentage
// value set, it'll be used to calculate the amount.
func (t *Terms) CalculateDues(zero num.Amount, sum num.Amount) {
	if t == nil {
		return
	}
	for _, dd := range t.DueDates {
		if dd.Percent != nil && !dd.Percent.IsZero() {
			dd.Amount = dd.Percent.Of(sum)
		}
		dd.Amount = dd.Amount.Rescale(zero.Exp())
	}
}

// Validate ensures that the terms contain everything required.
func (t *Terms) Validate() error {
	return t.ValidateWithContext(context.Background())
}

// ValidateWithContext ensures that the terms contain everything required.
func (t *Terms) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithContext(ctx, t,
		validation.Field(&t.Key, isValidTermKey),
		validation.Field(&t.DueDates),
		validation.Field(&t.Ext),
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
		validation.Field(&dd.Amount, validation.Required, num.NotZero),
		validation.Field(&dd.Percent),
		validation.Field(&dd.Currency),
	)
}

// JSONSchemaExtend adds the payment terms key list to the schema.
func (Terms) JSONSchemaExtend(schema *jsonschema.Schema) {
	prop, ok := schema.Properties.Get("key")
	if ok {
		prop.OneOf = make([]*jsonschema.Schema, len(TermKeyDefinitions))
		for i, v := range TermKeyDefinitions {
			prop.OneOf[i] = &jsonschema.Schema{
				Const:       v.Key,
				Title:       v.Title,
				Description: v.Description,
			}
		}
	}
}
