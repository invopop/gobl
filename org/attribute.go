package org

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/jsonschema"
)

// Standard attribute keys used to optionally identify common types of
// attributes. These are non-binding suggestions that help map attributes
// to other formats; any other valid key agreed upon between supplier and
// customer may also be used.
const (
	// Physical properties
	AttributeKeyColour   cbc.Key = "colour"
	AttributeKeySize     cbc.Key = "size"
	AttributeKeyMaterial cbc.Key = "material"
	AttributeKeyLength   cbc.Key = "length"
	AttributeKeyWidth    cbc.Key = "width"
	AttributeKeyHeight   cbc.Key = "height"
	AttributeKeyDiameter cbc.Key = "diameter"
	AttributeKeyWeight   cbc.Key = "weight"
	AttributeKeyVolume   cbc.Key = "volume"

	// Dates
	AttributeKeyProduction cbc.Key = "production"
	AttributeKeyExpiry     cbc.Key = "expiry"

	// Nutritional declarations, with amounts expressed per 100g or 100ml
	// of the item as per EU Regulation 1169/2011 or similar.
	AttributeKeyNutritionEnergy        cbc.Key = "nutrition+energy"
	AttributeKeyNutritionFat           cbc.Key = "nutrition+fat"
	AttributeKeyNutritionSaturatedFat  cbc.Key = "nutrition+saturated-fat"
	AttributeKeyNutritionCarbohydrates cbc.Key = "nutrition+carbohydrates"
	AttributeKeyNutritionSugars        cbc.Key = "nutrition+sugars"
	AttributeKeyNutritionProtein       cbc.Key = "nutrition+protein"
	AttributeKeyNutritionSalt          cbc.Key = "nutrition+salt"
	AttributeKeyNutritionFibre         cbc.Key = "nutrition+fibre"

	// Environmental impact (experimental)
	AttributeKeyEmissionsCO2e cbc.Key = "emissions+co2e"
)

// AttributeKeyDefinitions describes each of the standard attribute keys.
var AttributeKeyDefinitions = []*cbc.Definition{
	{
		Key:  AttributeKeyColour,
		Name: i18n.NewString("Colour"),
		Desc: i18n.NewString("Colour of the item, either as a name or a code such as RAL or Pantone."),
	},
	{
		Key:  AttributeKeySize,
		Name: i18n.NewString("Size"),
		Desc: i18n.NewString("Size of the item, such as a clothing size."),
	},
	{
		Key:  AttributeKeyMaterial,
		Name: i18n.NewString("Material"),
		Desc: i18n.NewString("Main material the item is made of, such as cotton or steel."),
	},
	{
		Key:  AttributeKeyLength,
		Name: i18n.NewString("Length"),
		Desc: i18n.NewString("Length of the item, usually as an amount with a unit of measure."),
	},
	{
		Key:  AttributeKeyWidth,
		Name: i18n.NewString("Width"),
		Desc: i18n.NewString("Width of the item, usually as an amount with a unit of measure."),
	},
	{
		Key:  AttributeKeyHeight,
		Name: i18n.NewString("Height"),
		Desc: i18n.NewString("Height of the item, usually as an amount with a unit of measure."),
	},
	{
		Key:  AttributeKeyDiameter,
		Name: i18n.NewString("Diameter"),
		Desc: i18n.NewString("Diameter of the item, usually as an amount with a unit of measure."),
	},
	{
		Key:  AttributeKeyWeight,
		Name: i18n.NewString("Weight"),
		Desc: i18n.NewString("Weight of a single unit of the item, usually as an amount with a unit of measure."),
	},
	{
		Key:  AttributeKeyVolume,
		Name: i18n.NewString("Volume"),
		Desc: i18n.NewString("Volume of a single unit of the item, usually as an amount with a unit of measure."),
	},
	{
		Key:  AttributeKeyProduction,
		Name: i18n.NewString("Production Date"),
		Desc: i18n.NewString("Date the item was produced or manufactured."),
	},
	{
		Key:  AttributeKeyExpiry,
		Name: i18n.NewString("Expiry Date"),
		Desc: i18n.NewString("Date the item expires or is best consumed before."),
	},
	{
		Key:  AttributeKeyNutritionEnergy,
		Name: i18n.NewString("Energy"),
		Desc: i18n.NewString("Energy content per 100g or 100ml, usually in kilojoules (kj) or kilocalories (kcal)."),
	},
	{
		Key:  AttributeKeyNutritionFat,
		Name: i18n.NewString("Fat"),
		Desc: i18n.NewString("Total fat content per 100g or 100ml, usually in grams."),
	},
	{
		Key:  AttributeKeyNutritionSaturatedFat,
		Name: i18n.NewString("Saturated Fat"),
		Desc: i18n.NewString("Saturated fat content per 100g or 100ml, usually in grams."),
	},
	{
		Key:  AttributeKeyNutritionCarbohydrates,
		Name: i18n.NewString("Carbohydrates"),
		Desc: i18n.NewString("Carbohydrate content per 100g or 100ml, usually in grams."),
	},
	{
		Key:  AttributeKeyNutritionSugars,
		Name: i18n.NewString("Sugars"),
		Desc: i18n.NewString("Sugar content per 100g or 100ml, usually in grams."),
	},
	{
		Key:  AttributeKeyNutritionProtein,
		Name: i18n.NewString("Protein"),
		Desc: i18n.NewString("Protein content per 100g or 100ml, usually in grams."),
	},
	{
		Key:  AttributeKeyNutritionSalt,
		Name: i18n.NewString("Salt"),
		Desc: i18n.NewString("Salt content per 100g or 100ml, usually in grams."),
	},
	{
		Key:  AttributeKeyNutritionFibre,
		Name: i18n.NewString("Fibre"),
		Desc: i18n.NewString("Fibre content per 100g or 100ml, usually in grams."),
	},
	{
		Key:  AttributeKeyEmissionsCO2e,
		Name: i18n.NewString("CO2e Emissions"),
		Desc: i18n.NewString("Greenhouse gas emissions of a single unit of the item expressed as a carbon dioxide equivalent, usually in kilograms."),
	},
}

// Attribute describes a named feature or property of the parent object,
// such as the colour or size of an item. Attributes are identified by
// either a key or a type, and hold exactly one of the text, code, amount,
// or date value fields.
type Attribute struct {
	// Label for the attribute for internal use, not included in output documents.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`

	// Identifier fields; either key or type must be provided, but not both.

	// Key that identifies the attribute, either from the list pre-defined by
	// GOBL or an alternative agreed upon between the supplier and customer.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// Type defines a code used to identify the attribute when the "key"
	// field is empty, from a code list agreed upon between the supplier
	// and customer.
	Type cbc.Code `json:"type,omitempty" jsonschema:"title=Type"`

	// Value fields; exactly one must be provided.

	// Text value of the attribute.
	Text string `json:"text,omitempty" jsonschema:"title=Text"`
	// Code value of the attribute, used as an alternative to text when the
	// value comes from a code list agreed upon between the supplier and customer.
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`
	// Amount used when the attribute represents a numeric or measurable value.
	Amount *num.Amount `json:"amount,omitempty" jsonschema:"title=Amount"`
	// Unit of measure that accompanies the amount.
	Unit Unit `json:"unit,omitempty" jsonschema:"title=Unit"`
	// Date value of the attribute.
	Date *cal.Date `json:"date,omitempty" jsonschema:"title=Date"`
}

func attributeRules() *rules.Set {
	return rules.For(new(Attribute),
		rules.Assert("01", "attribute must have either a key or a type, but not both",
			is.Expr(`(string(Key) == "") != (string(Type) == "")`),
		),
		rules.Assert("02", "attribute must have exactly one of the text, code, amount, or date values",
			is.Expr(`(Text == "" ? 0 : 1) + (string(Code) == "" ? 0 : 1) + (Amount == nil ? 0 : 1) + (Date == nil ? 0 : 1) == 1`),
		),
		rules.When(is.Expr(`string(Unit) != ""`),
			rules.Assert("03", "attribute unit may only be used alongside an amount",
				is.Expr(`Amount != nil`),
			),
		),
	)
}

func normalizeAttribute(a *Attribute) {
	a.Label = cbc.NormalizeString(a.Label)
	a.Text = cbc.NormalizeString(a.Text)
}

// IsEmpty returns true if the attribute has no meaningful content.
func (a *Attribute) IsEmpty() bool {
	return a == nil || (a.Label == "" &&
		a.Key == "" &&
		a.Type == "" &&
		a.Text == "" &&
		a.Code == "" &&
		a.Amount == nil &&
		a.Unit == UnitEmpty &&
		a.Date == nil)
}

// AttributesHaveUniqueKeys provides a test that ensures no two attributes
// in a list share the same key. Attributes without a key are ignored.
func AttributesHaveUniqueKeys() rules.Test {
	return is.Func("attribute keys must be unique", func(value any) bool {
		attrs, ok := value.([]*Attribute)
		if !ok {
			return false
		}
		seen := make(map[cbc.Key]bool, len(attrs))
		for _, a := range attrs {
			if a == nil || a.Key == "" {
				continue
			}
			if seen[a.Key] {
				return false
			}
			seen[a.Key] = true
		}
		return true
	})
}

// CleanAttributes removes any nil or empty attributes from the list.
func CleanAttributes(attrs []*Attribute) []*Attribute {
	var cleaned []*Attribute
	for _, a := range attrs {
		if a.IsEmpty() {
			continue
		}
		cleaned = append(cleaned, a)
	}
	return cleaned
}

// JSONSchemaExtend adds extra details to the schema.
func (Attribute) JSONSchemaExtend(js *jsonschema.Schema) {
	prop, ok := js.Properties.Get("key")
	if !ok {
		return
	}
	anyOf := make([]*jsonschema.Schema, 0, len(AttributeKeyDefinitions)+1)
	for _, def := range AttributeKeyDefinitions {
		anyOf = append(anyOf, &jsonschema.Schema{
			Const:       def.Key,
			Title:       def.Name.String(),
			Description: def.Desc.String(),
		})
	}
	anyOf = append(anyOf, &jsonschema.Schema{
		Title:   "Other",
		Pattern: cbc.KeyPattern,
	})
	prop.AnyOf = anyOf
}
