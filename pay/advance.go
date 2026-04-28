package pay

import (
	"encoding/json"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
)

// Advance represents a single payment that has been made already, such
// as a deposit on an intent to purchase, or as credit from a previous
// invoice which was later corrected or cancelled.
type Advance struct {
	uuid.Identify
	// When the advance was made.
	Date *cal.Date `json:"date,omitempty" jsonschema:"title=Date"`
	// The payment means used to make the advance.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// ID or reference for the advance.
	Ref string `json:"ref,omitempty" jsonschema:"title=Reference"`
	// If this "advance" payment has come from a public grant or subsidy, set this to true.
	Grant bool `json:"grant,omitempty" jsonschema:"title=Grant"`
	// Details about the advance.
	Description string `json:"description" jsonschema:"title=Description"`
	// Percentage of the total amount payable that was paid. Note that
	// multiple advances with percentages may lead to rounding errors,
	// especially when the total advances sums to 100%. We recommend only
	// including one advance with a percent value per document.
	Percent *num.Percentage `json:"percent,omitempty" jsonschema:"title=Percent"`
	// How much was paid.
	Amount num.Amount `json:"amount" jsonschema:"title=Amount"`
	// If different from the parent document's base currency.
	Currency currency.Code `json:"currency,omitempty" jsonschema:"title=Currency"`
	// Details of the payment that was made via a credit or debit card.
	Card *Card `json:"card,omitempty" jsonschema:"title=Card"`
	// Details about how the payment was made by credit (bank) transfer.
	CreditTransfer *CreditTransfer `json:"credit_transfer,omitempty" jsonschema:"title=Credit Transfer"`
	// Tax extensions required by tax regimes or addons.
	Ext tax.Extensions `json:"ext,omitzero" jsonschema:"title=Extensions"`
	// Additional details useful for the parties involved.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

func advanceRules() *rules.Set {
	return rules.For(new(Advance),
		rules.Field("description",
			rules.Assert("01", "description is required", is.Present),
		),
		rules.Field("key",
			rules.AssertIfPresent("02", "key must be valid", HasValidMeansKey),
		),
	)
}

// Normalize will try to normalize the advance's data.
func (a *Advance) Normalize() {
	if a == nil {
		return
	}
	uuid.Normalize(&a.UUID)
	a.Ref = cbc.NormalizeString(a.Ref)
	a.Description = cbc.NormalizeString(a.Description)
	a.Ext = a.Ext.Clean()
}

// CalculateFrom will update the amount using the rate of the provided
// total, if defined.
func (a *Advance) CalculateFrom(payable num.Amount) {
	if a != nil && a.Percent != nil {
		a.Amount = a.Percent.Of(payable)
	}
}

// UnmarshalJSON helps migrate the desc field to description.
func (a *Advance) UnmarshalJSON(data []byte) error {
	type Alias Advance
	aux := struct {
		Desc string `json:"desc,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Desc != "" {
		a.Description = aux.Desc
	}
	return nil
}

// JSONSchemaExtend extends the JSONSchema for the Instructions type.
func (Advance) JSONSchemaExtend(schema *jsonschema.Schema) {
	extendJSONSchemaWithMeansKey(schema, "key")
}
