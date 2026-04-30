package pay

import (
	"encoding/json"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
)

// Record represents a single payment that has been or will be made via a
// specific means (cash, card, credit-transfer, etc.), with its own amount and
// optional currency. It is used both for advance payments declared on an
// invoice (via bill.PaymentDetails.Advances) and for the methods recorded on
// a bill.Payment document (via bill.Payment.Methods).
type Record struct {
	uuid.Identify
	// When the payment was made.
	Date *cal.Date `json:"date,omitempty" jsonschema:"title=Date"`
	// The payment means used.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// ID or reference for the payment.
	Ref string `json:"ref,omitempty" jsonschema:"title=Reference"`
	// Description about the payment.
	Description string `json:"description,omitempty" jsonschema:"title=Description"`
	// Percentage of the total amount payable that was paid. Note that
	// multiple records with percentages may lead to rounding errors,
	// especially when the total sums to 100%. We recommend only including one
	// record with a percent value per document.
	Percent *num.Percentage `json:"percent,omitempty" jsonschema:"title=Percent"`
	// How much was paid.
	Amount num.Amount `json:"amount" jsonschema:"title=Amount"`
	// If different from the parent document's base currency.
	Currency currency.Code `json:"currency,omitempty" jsonschema:"title=Currency"`
	// Details of the payment that was made via a credit or debit card.
	Card *Card `json:"card,omitempty" jsonschema:"title=Card"`
	// Details about how the payment was made by credit (bank) transfer.
	CreditTransfer *CreditTransfer `json:"credit_transfer,omitempty" jsonschema:"title=Credit Transfer"`
	// Details of the payment that was made via direct debit.
	DirectDebit *DirectDebit `json:"direct_debit,omitempty" jsonschema:"title=Direct Debit"`
	// Details of the payment that was made via an online provider.
	Online []*Online `json:"online,omitempty" jsonschema:"title=Online"`
	// Tax extensions required by tax regimes or addons.
	Ext tax.Extensions `json:"ext,omitzero" jsonschema:"title=Extensions"`
	// Additional details useful for the parties involved.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

func recordRules() *rules.Set {
	return rules.For(new(Record),
		rules.Field("key",
			rules.AssertIfPresent("01", "key must be valid", HasValidMeansKey),
		),
	)
}

// Normalize will try to normalize the record's data, threading the supplied
// addon/regime normalizers through to the means-specific nested structs (Card,
// CreditTransfer, DirectDebit, Online) so addons can extend or transform them.
func (r *Record) Normalize(normalizers tax.Normalizers) {
	if r == nil {
		return
	}
	uuid.Normalize(&r.UUID)
	r.Ref = cbc.NormalizeString(r.Ref)
	r.Description = cbc.NormalizeString(r.Description)
	r.Ext = r.Ext.Clean()

	tax.Normalize(normalizers, r.Card)
	tax.Normalize(normalizers, r.CreditTransfer)
	tax.Normalize(normalizers, r.DirectDebit)
	tax.Normalize(normalizers, r.Online)

	normalizers.Each(r)
}

// CalculateFrom will update the amount using the rate of the provided
// total, if defined.
func (r *Record) CalculateFrom(payable num.Amount) {
	if r != nil && r.Percent != nil {
		r.Amount = r.Percent.Of(payable)
	}
}

// UnmarshalJSON helps migrate the desc field to description.
func (r *Record) UnmarshalJSON(data []byte) error {
	type Alias Record
	aux := struct {
		Desc string `json:"desc,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Desc != "" {
		r.Description = aux.Desc
	}
	return nil
}

// JSONSchemaExtend extends the JSONSchema for the Record type.
func (Record) JSONSchemaExtend(schema *jsonschema.Schema) {
	extendJSONSchemaWithMeansKey(schema, "key")
}
