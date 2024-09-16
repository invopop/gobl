package tax

import "github.com/invopop/gobl/l10n"

// Regime defines a structure that can be embedded inside another structure to enable
// methods and the `$regime` attribute to be able to determine a tax regime definition
// to associate with the document.
type Regime struct {
	Country l10n.TaxCountryCode `json:"$regime,omitempty" jsonschema:"title=Regime"`
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
	if Regimes().For(country.Code()) == nil {
		r.Country = ""
		return
	}
	r.Country = country
}

// RegimeDef provides the associated regime definition.
func (r Regime) RegimeDef() *RegimeDef {
	return Regimes().For(r.Country.Code())
}

// IsEmpty returns true if the regime is empty.
func (r Regime) IsEmpty() bool {
	return r.Country.Empty()
}
