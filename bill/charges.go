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

// Charge keys for identifying the type of charge being applied.
// These are based on a subset of the UN/CEFACT UNTDID 7161 codes,
// and are intentionally kept lean.
const (
	ChargeKeyStampDuty cbc.Key = "stamp-duty"
	ChargeKeyOutlay    cbc.Key = "outlay"
	ChargeKeyTax       cbc.Key = "tax"
	ChargeKeyCustoms   cbc.Key = "customs"
	ChargeKeyDelivery  cbc.Key = "delivery"
	ChargeKeyPacking   cbc.Key = "packing"
	ChargeKeyHandling  cbc.Key = "handling"
	ChargeKeyInsurance cbc.Key = "insurance"
	ChargeKeyStorage   cbc.Key = "storage"
	ChargeKeyAdmin     cbc.Key = "admin" // administration
	ChargeKeyCleaning  cbc.Key = "cleaning"
)

var chargeKeyDefinitions = []*cbc.Definition{
	{
		Key:  ChargeKeyStampDuty,
		Name: i18n.NewString("Stamp Duty"),
	},
	{
		Key:  ChargeKeyOutlay,
		Name: i18n.NewString("Outlay"),
	},
	{
		Key:  ChargeKeyTax,
		Name: i18n.NewString("Tax"),
	},
	{
		Key:  ChargeKeyCustoms,
		Name: i18n.NewString("Customs"),
	},
	{
		Key:  ChargeKeyDelivery,
		Name: i18n.NewString("Delivery"),
	},
	{
		Key:  ChargeKeyPacking,
		Name: i18n.NewString("Packing"),
	},
	{
		Key:  ChargeKeyHandling,
		Name: i18n.NewString("Handling"),
	},
	{
		Key:  ChargeKeyInsurance,
		Name: i18n.NewString("Insurance"),
	},
	{
		Key:  ChargeKeyStorage,
		Name: i18n.NewString("Storage"),
	},
	{
		Key:  ChargeKeyAdmin,
		Name: i18n.NewString("Administration"),
	},
	{
		Key:  ChargeKeyCleaning,
		Name: i18n.NewString("Cleaning"),
	},
}

// Charge represents a surchange applied to the complete document
// independent from the individual lines.
type Charge struct {
	uuid.Identify
	// Line number inside the list of charges (calculated).
	Index int `json:"i" jsonschema:"title=Index" jsonschema_extras:"calculated=true"`
	// Key for grouping or identifying charges for tax purposes. A suggested list of
	// keys is provided, but these may be extended by the issuer.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// Code to used to refer to the this charge by the issuer
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`
	// Text description as to why the charge was applied
	Reason string `json:"reason,omitempty" jsonschema:"title=Reason"`
	// Base represents the value used as a base for percent calculations instead
	// of the invoice's sum of lines.
	Base *num.Amount `json:"base,omitempty" jsonschema:"title=Base"`
	// Percentage to apply to the sum of all lines
	Percent *num.Percentage `json:"percent,omitempty" jsonschema:"title=Percent"`
	// Amount to apply (calculated if percent present)
	Amount num.Amount `json:"amount" jsonschema:"title=Amount" jsonschema_extras:"calculated=true"`
	// List of taxes to apply to the charge
	Taxes tax.Set `json:"taxes,omitempty" jsonschema:"title=Taxes"`
	// Extension codes that apply to the charge
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`
	// Additional semi-structured information.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Normalize performs normalization on the line and embedded objects using the
// provided list of normalizers.
func (m *Charge) Normalize(normalizers tax.Normalizers) {
	m.Code = cbc.NormalizeCode(m.Code)
	m.Taxes = tax.CleanSet(m.Taxes)
	m.Ext = tax.CleanExtensions(m.Ext)
	normalizers.Each(m)
	tax.Normalize(normalizers, m.Taxes)
}

// ValidateWithContext checks the charge's fields.
func (m *Charge) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithContext(ctx, m,
		validation.Field(&m.UUID),
		validation.Field(&m.Key),
		validation.Field(&m.Code),
		validation.Field(&m.Base),
		validation.Field(&m.Percent,
			validation.When(
				m.Base != nil,
				validation.Required,
			),
		),
		validation.Field(&m.Amount, validation.Required),
		validation.Field(&m.Taxes),
		validation.Field(&m.Ext),
		validation.Field(&m.Meta),
	)
}

// GetTaxes responds with the array of tax rates applied to this line.
func (m *Charge) GetTaxes() tax.Set {
	return m.Taxes
}

// GetTotal provides the final total for this line, excluding any tax calculations.
func (m *Charge) GetTotal() num.Amount {
	return m.Amount
}

func (m *Charge) removeIncludedTaxes(cat cbc.Code) *Charge {
	accuracy := defaultTaxRemovalAccuracy
	rate := m.Taxes.Get(cat)
	if rate == nil || rate.Percent == nil {
		return m
	}
	m2 := *m
	m2.Amount = m2.Amount.Upscale(accuracy).Remove(*rate.Percent)
	return &m2
}

// JSONSchemaExtend adds the charge key definitions to the schema.
func (Charge) JSONSchemaExtend(schema *jsonschema.Schema) {
	extendJSONSchemaWithChargeKey(schema)
}

func calculateCharges(lines []*Charge, cur currency.Code, sum num.Amount, rr cbc.Key) {
	// COPIED FROM discount.go
	zero := cur.Def().Zero()
	if len(lines) == 0 {
		return
	}
	for i, l := range lines {
		l.Index = i + 1
		if l.Percent != nil && !l.Percent.IsZero() {
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

func calculateChargeSum(charges []*Charge, cur currency.Code) *num.Amount {
	if len(charges) == 0 {
		return nil
	}
	total := cur.Def().Zero()
	for _, l := range charges {
		total = total.MatchPrecision(l.Amount)
		total = total.Add(l.Amount)
	}
	return &total
}

func (c *Charge) round(cur currency.Code) {
	cd := cur.Def()
	c.Amount = cd.Rescale(c.Amount)
}

func roundCharges(lines []*Charge, cur currency.Code) {
	for _, l := range lines {
		l.round(cur)
	}
}

func extendJSONSchemaWithChargeKey(schema *jsonschema.Schema) {
	prop, ok := schema.Properties.Get("key")
	if !ok {
		return
	}
	prop.AnyOf = make([]*jsonschema.Schema, len(chargeKeyDefinitions))
	for i, v := range chargeKeyDefinitions {
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
