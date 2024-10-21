package tax

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/jsonschema"
)

// Regime defines a structure that can be embedded inside another structure to enable
// methods and the `$regime` attribute to be able to determine a tax regime definition
// to associate with the document.
type Regime struct {
	Country l10n.TaxCountryCode `json:"$regime,omitempty" jsonschema:"title=Tax Regime"`
}

// WithRegime prepares a Regime struct with the provided country code.
func WithRegime(country l10n.TaxCountryCode) Regime {
	return Regime{Country: country}
}

// GetRegime returns the regime country code.
func (r Regime) GetRegime() l10n.TaxCountryCode {
	return r.Country
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
	r.Country = rd.Country
}

// RegimeDef provides the associated regime definition.
func (r Regime) RegimeDef() *RegimeDef {
	return Regimes().For(r.Country.Code())
}

// IsEmpty returns true if the regime is empty.
func (r Regime) IsEmpty() bool {
	return r.Country.Empty()
}

// JSONSchemaExtend will add the addon options to the JSON list.
func (r Regime) JSONSchemaExtend(js *jsonschema.Schema) {
	props := js.Properties
	if asl, ok := props.Get("$regime"); ok {
		asl.OneOf = make([]*jsonschema.Schema, len(AllRegimeDefs()))
		for i, rd := range AllRegimeDefs() {
			asl.OneOf[i] = &jsonschema.Schema{
				Const: rd.Code().String(),
				Title: rd.Name.String(),
			}
		}
	}
}
