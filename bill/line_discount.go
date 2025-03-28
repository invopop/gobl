package bill

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// LineDiscount represents an amount deducted from the line, and will be
// applied before taxes.
type LineDiscount struct {
	// Key for identifying the type of discount being applied.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// Code or reference for this discount defined by the issuer
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`
	// Text description as to why the discount was applied
	Reason string `json:"reason,omitempty" jsonschema:"title=Reason"`
	// Base for percent calculations instead of the line's sum.
	Base *num.Amount `json:"base,omitempty" jsonschema:"title=Base"`
	// Percentage to apply to the base or line sum to calculate the discount amount
	Percent *num.Percentage `json:"percent,omitempty" jsonschema:"title=Percent"`
	// Fixed discount amount to apply (calculated if percent present)
	Amount num.Amount `json:"amount" jsonschema:"title=Amount" jsonschema_extras:"calculated=true"`
	// Extension codes that apply to the discount
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`
}

// Normalize performs normalization on the discount and embedded objects using the
// provided list of normalizers.
func (ld *LineDiscount) Normalize(normalizers tax.Normalizers) {
	ld.Code = cbc.NormalizeCode(ld.Code)
	ld.Ext = tax.CleanExtensions(ld.Ext)
	normalizers.Each(ld)
}

// Validate checks the line discount's fields.
func (ld *LineDiscount) Validate() error {
	return validation.ValidateStruct(ld,
		validation.Field(&ld.Key),
		validation.Field(&ld.Code),
		validation.Field(&ld.Base),
		validation.Field(&ld.Percent,
			validation.When(
				ld.Base != nil,
				validation.Required,
			),
		),
		validation.Field(&ld.Amount, validation.Required, num.NotZero),
		validation.Field(&ld.Ext),
	)
}

// IsEmpty returns true if the discount is empty.
func (ld *LineDiscount) IsEmpty() bool {
	return ld.Key.IsEmpty() &&
		ld.Code.IsEmpty() &&
		ld.Reason == "" &&
		(ld.Percent == nil || ld.Percent.IsZero()) &&
		ld.Amount.IsZero() &&
		len(ld.Ext) == 0
}

// CleanLineDiscounts removes any empty discounts from the list.
func CleanLineDiscounts(lines []*LineDiscount) []*LineDiscount {
	var cleaned []*LineDiscount
	for _, d := range lines {
		if d.IsEmpty() {
			continue
		}
		cleaned = append(cleaned, d)
	}
	return cleaned
}

// JSONSchemaExtend adds the discount key definitions to the schema.
func (LineDiscount) JSONSchemaExtend(schema *jsonschema.Schema) {
	extendJSONSchemaWithDiscountKey(schema)
}
