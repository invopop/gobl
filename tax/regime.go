package tax

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/jsonschema"
)

// RegimeCode defines the tax country code available for use as a regime
// identifier.
type RegimeCode l10n.TaxCountryCode

// Regime defines a structure that can be embedded inside another structure to enable
// methods and the `$regime` attribute to be able to determine a tax regime definition
// to associate with the document.
type Regime struct {
	// Country code that identifies the tax regime applicable to the document.
	// It determines which country-specific tax rules, normalizations, and validations are applied.
	// It may be determined automatically via normalization of a supplier or issuer tax identity
	// country code.
	Country RegimeCode `json:"$regime,omitempty" jsonschema:"title=Tax Regime"`
}

// WithRegime prepares a Regime struct with the provided country code.
func WithRegime(country l10n.TaxCountryCode) Regime {
	return Regime{Country: RegimeCode(country)}
}

// GetRegime returns the regime country code.
func (r Regime) GetRegime() l10n.TaxCountryCode {
	return l10n.TaxCountryCode(r.Country)
}

// SetRegime updates the current regime country code, after first checking to ensure
// that the regime is actually defined. Missing regimes will silently replace
// the current regime with an empty value.
func (r *Regime) SetRegime(country l10n.TaxCountryCode) {
	rd := Regimes().For(country.Code())
	if rd == nil {
		r.Country = ""
		return
	}
	r.Country = RegimeCode(rd.Country)
}

// RegimeDef provides the associated regime definition.
func (r Regime) RegimeDef() *RegimeDef {
	return Regimes().For(r.Country.Code())
}

// IsEmpty returns true if the regime is empty.
func (r Regime) IsEmpty() bool {
	return r.Country.Code().Empty()
}

// Code provides the regime code as an l10n.Code.
func (rc RegimeCode) Code() l10n.Code {
	return l10n.Code(rc)
}

// String provides the string representation of the regime code.
func (rc RegimeCode) String() string {
	return string(rc)
}

/*
// JSONSchemaExtend will add the addon options to the JSON list.
func (r Regime) JSONSchemaExtend(js *jsonschema.Schema) {
	props := js.Properties
	if asl, ok := props.Get("$regime"); ok {
		asl.Ref =
		asl.OneOf = make([]*jsonschema.Schema, len(AllRegimeDefs()))
		for i, rd := range AllRegimeDefs() {
			asl.OneOf[i] = &jsonschema.Schema{
				Const: rd.Code().String(),
				Title: rd.Name.String(),
			}
		}
	}
}
*/

// JSONSchemaExtend will add the addon options to the JSON list.
func (RegimeCode) JSONSchemaExtend(js *jsonschema.Schema) {
	js.OneOf = make([]*jsonschema.Schema, len(AllRegimeDefs()))
	for i, rd := range AllRegimeDefs() {
		js.OneOf[i] = &jsonschema.Schema{
			Const: rd.Code().String(),
			Title: rd.Name.String(),
		}
	}
}
