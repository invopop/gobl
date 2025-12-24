package bill

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// LineCharge represents an amount added to the line, and will be
// applied before taxes.
type LineCharge struct {
	// Key for grouping or identifying charges for tax purposes. A suggested list of
	// keys is provided, but these are for reference only and may be extended by
	// the issuer.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// Reference or ID for this charge defined by the issuer
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`
	// Text description as to why the charge was applied
	Reason string `json:"reason,omitempty" jsonschema:"title=Reason"`
	// Base for percent calculations instead of the line's sum
	Base *num.Amount `json:"base,omitempty" jsonschema:"title=Base"`
	// Percentage of base or parent line's sum
	Percent *num.Percentage `json:"percent,omitempty" jsonschema:"title=Percent"`
	// Quantity of units to apply the charge to when using the rate instead of
	// the line's quantity.
	Quantity *num.Amount `json:"quantity,omitempty" jsonschema:"title=Quantity"`
	// Unit to associate with the quantity when using the rate.
	Unit org.Unit `json:"unit,omitempty" jsonschema:"title=Unit"`
	// Rate defines a price per unit to use instead of the percentage.
	Rate *num.Amount `json:"rate,omitempty" jsonschema:"title=Rate"`
	// Fixed or resulting charge amount to apply (calculated if percent present).
	Amount num.Amount `json:"amount" jsonschema:"title=Amount" jsonschema_extras:"calculated=true"`
	// Extension codes that apply to the charge
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`
}

// Normalize performs normalization on the charge and embedded objects using the
// provided list of normalizers.
func (lc *LineCharge) Normalize(normalizers tax.Normalizers) {
	lc.Code = cbc.NormalizeCode(lc.Code)
	lc.Ext = tax.CleanExtensions(lc.Ext)
	normalizers.Each(lc)
}

// Validate checks the line charge's fields.
func (lc *LineCharge) Validate() error {
	return validation.ValidateStruct(lc,
		validation.Field(&lc.Key),
		validation.Field(&lc.Code),
		validation.Field(&lc.Base),
		validation.Field(&lc.Percent,
			validation.When(
				lc.Base != nil,
				validation.Required,
			),
		),
		validation.Field(&lc.Quantity,
			validation.When(
				lc.Base != nil || lc.Percent != nil,
				validation.Empty.Error("must be blank with base or percent"),
			),
		),
		validation.Field(&lc.Unit,
			validation.When(
				lc.Quantity == nil,
				validation.Empty.Error("must be blank without quantity"),
			),
		),
		validation.Field(&lc.Rate,
			validation.When(
				lc.Base != nil || lc.Percent != nil,
				validation.Empty.Error("must be blank with base or percent"),
			),
			validation.When(
				lc.Quantity != nil,
				validation.Required.Error("cannot be blank with quantity"),
			),
		),
		validation.Field(&lc.Amount),
		validation.Field(&lc.Ext),
	)
}

// IsEmpty returns true if the charge is empty.
func (lc *LineCharge) IsEmpty() bool {
	return lc.Key.IsEmpty() &&
		lc.Code.IsEmpty() &&
		lc.Reason == "" &&
		(lc.Percent == nil || lc.Percent.IsZero()) &&
		lc.Amount.IsZero() &&
		len(lc.Ext) == 0
}

// CleanLineCharges removes any empty charges from the list.
func CleanLineCharges(lines []*LineCharge) []*LineCharge {
	var cleaned []*LineCharge
	for _, l := range lines {
		if l.IsEmpty() {
			continue
		}
		cleaned = append(cleaned, l)
	}
	return cleaned
}

// JSONSchemaExtend adds the charge key definitions to the schema.
func (LineCharge) JSONSchemaExtend(schema *jsonschema.Schema) {
	extendJSONSchemaWithChargeKey(schema)
}
