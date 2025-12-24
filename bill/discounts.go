package bill

import (
	"context"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Discount keys for identifying the type of discount being applied.
// These are based on the UN/CEFACT UNTDID 5189 code list subset defined
// in the EN16931 code lists and are mean as suggestions.
const (
	DiscountKeyEarlyCompletion  cbc.Key = "early-completion"
	DiscountKeyMilitary         cbc.Key = "military"
	DiscountKeyWorkAccident     cbc.Key = "work-accident"
	DiscountKeySpecialAgreement cbc.Key = "special-agreement"
	DiscountKeyProductionError  cbc.Key = "production-error"
	DiscountKeyNewOutlet        cbc.Key = "new-outlet"
	DiscountKeySample           cbc.Key = "sample"
	DiscountKeyEndOfRange       cbc.Key = "end-of-range"
	DiscountKeyIncoterm         cbc.Key = "incoterm"
	DiscountKeyPOSThreshold     cbc.Key = "pos-threshold"
	DiscountKeySpecialRebate    cbc.Key = "special-rebate"
	DiscountKeyTemporary        cbc.Key = "temporary"
	DiscountKeyStandard         cbc.Key = "standard"
	DiscountKeyYarlyTurnover    cbc.Key = "yearly-turnover"
)

var discountKeyDefinitions = []*cbc.Definition{
	{
		Key:  DiscountKeyEarlyCompletion,
		Name: i18n.NewString("Bonus for works ahead of schedule"),
	},
	{
		Key:  DiscountKeyMilitary,
		Name: i18n.NewString("Military Discount"),
	},
	{
		Key:  DiscountKeyWorkAccident,
		Name: i18n.NewString("Work Accident Discount"),
	},
	{
		Key:  DiscountKeySpecialAgreement,
		Name: i18n.NewString("Special Agreement Discount"),
	},
	{
		Key:  DiscountKeyProductionError,
		Name: i18n.NewString("Production Error Discount"),
	},
	{
		Key:  DiscountKeyNewOutlet,
		Name: i18n.NewString("New Outlet Discount"),
	},
	{
		Key:  DiscountKeySample,
		Name: i18n.NewString("Sample Discount"),
	},
	{
		Key:  DiscountKeyEndOfRange,
		Name: i18n.NewString("End of Range Discount"),
	},
	{
		Key:  DiscountKeyIncoterm,
		Name: i18n.NewString("Incoterm Discount"),
	},
	{
		Key:  DiscountKeyPOSThreshold,
		Name: i18n.NewString("Point of Sale Threshold Discount"),
	},
	{
		Key:  DiscountKeySpecialRebate,
		Name: i18n.NewString("Special Rebate"),
	},
	{
		Key:  DiscountKeyTemporary,
		Name: i18n.NewString("Temporary"),
	},
	{
		Key:  DiscountKeyStandard,
		Name: i18n.NewString("Standard"),
	},
	{
		Key:  DiscountKeyYarlyTurnover,
		Name: i18n.NewString("Yearly Turnover"),
	},
}

// Discount represents an allowance applied to the complete document
// independent from the individual lines. These have more in common with
// Invoice Lines than anything else, as each discount must have the
// correct taxes defined.
type Discount struct {
	uuid.Identify
	// Line number inside the list of discounts (calculated)
	Index int `json:"i" jsonschema:"title=Index" jsonschema_extras:"calculated=true"`
	// Key for identifying the type of discount being applied.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// Code to used to refer to the this discount by the issuer
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`
	// Text description as to why the discount was applied
	Reason string `json:"reason,omitempty" jsonschema:"title=Reason"`
	// Base represents the value used as a base for percent calculations instead
	// of the invoice's sum of lines.
	Base *num.Amount `json:"base,omitempty" jsonschema:"title=Base"`
	// Percentage to apply to the base or invoice's sum.
	Percent *num.Percentage `json:"percent,omitempty" jsonschema:"title=Percent"`
	// Amount to apply (calculated if percent present).
	Amount num.Amount `json:"amount" jsonschema:"title=Amount" jsonschema_extras:"calculated=true"`
	// List of taxes to apply to the discount
	Taxes tax.Set `json:"taxes,omitempty" jsonschema:"title=Taxes"`
	// Extension codes that apply to the discount
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`
	// Additional semi-structured information.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Normalize performs normalization on the line and embedded objects using the
// provided list of normalizers.
func (m *Discount) Normalize(normalizers tax.Normalizers) {
	if m == nil {
		return
	}
	m.Code = cbc.NormalizeCode(m.Code)
	m.Taxes = tax.CleanSet(m.Taxes)
	m.Ext = tax.CleanExtensions(m.Ext)
	normalizers.Each(m)
	tax.Normalize(normalizers, m.Taxes)
}

// ValidateWithContext checks the discount's fields.
func (m *Discount) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithContext(ctx, m,
		validation.Field(&m.UUID),
		validation.Field(&m.Code),
		validation.Field(&m.Base),
		validation.Field(&m.Percent,
			validation.When(
				m.Base != nil,
				validation.Required,
			),
		),
		validation.Field(&m.Amount),
		validation.Field(&m.Taxes),
		validation.Field(&m.Ext),
		validation.Field(&m.Meta),
	)
}

// GetTaxes responds with the array of tax rates applied to this line.
func (m *Discount) GetTaxes() tax.Set {
	return m.Taxes
}

// GetTotal provides the final total for this line, excluding any tax calculations.
// We return a negative value so that discounts will be applied correctly.
func (m *Discount) GetTotal() num.Amount {
	return m.Amount.Invert()
}

func (m *Discount) removeIncludedTaxes(cat cbc.Code) *Discount {
	accuracy := defaultTaxRemovalAccuracy
	rate := m.Taxes.Get(cat)
	if rate == nil || rate.Percent == nil {
		return m
	}
	m2 := *m
	m2.Amount = m2.Amount.Upscale(accuracy).Remove(*rate.Percent)
	return &m2
}

// JSONSchemaExtend adds the discount key definitions to the schema.
func (Discount) JSONSchemaExtend(schema *jsonschema.Schema) {
	extendJSONSchemaWithDiscountKey(schema)
}

func calculateDiscounts(lines []*Discount, cur currency.Code, sum num.Amount, rr cbc.Key) {
	zero := cur.Def().Zero()
	if len(lines) == 0 {
		return
	}
	for i, l := range lines {
		if l == nil {
			continue
		}
		l.Index = i + 1
		if l.Percent != nil {
			base := sum
			if l.Base != nil {
				base = l.Base.RescaleUp(zero.Exp() + linePrecisionExtra)
				base = tax.ApplyRoundingRule(rr, cur, base)
			}
			l.Amount = l.Percent.Of(base)
		}
		l.Amount = tax.ApplyRoundingRule(rr, cur, l.Amount)
	}
}

func calculateDiscountSum(discounts []*Discount, cur currency.Code) *num.Amount {
	if len(discounts) == 0 {
		return nil
	}
	total := cur.Def().Zero()
	for _, l := range discounts {
		if l == nil {
			continue
		}
		total = total.MatchPrecision(l.Amount)
		total = total.Add(l.Amount)
	}
	return &total
}

func (m *Discount) round(cur currency.Code) {
	// Default round to currency, or use base if present
	e := cur.Def().Subunits
	if m.Base != nil {
		e = m.Base.Exp()
	}
	m.Amount = m.Amount.RescaleDown(e)
}

func roundDiscounts(lines []*Discount, cur currency.Code) {
	for _, l := range lines {
		if l != nil {
			l.round(cur)
		}
	}
}

func extendJSONSchemaWithDiscountKey(schema *jsonschema.Schema) {
	prop, ok := schema.Properties.Get("key")
	if !ok {
		return
	}
	prop.AnyOf = make([]*jsonschema.Schema, len(discountKeyDefinitions))
	for i, v := range discountKeyDefinitions {
		prop.AnyOf[i] = &jsonschema.Schema{
			Const: v.Key,
			Title: v.Name.String(),
		}
	}
	prop.AnyOf = append(prop.AnyOf, &jsonschema.Schema{
		Title:   "Other",
		Pattern: cbc.KeyPattern,
	})
}
