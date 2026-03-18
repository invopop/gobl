package tax

import (
	"strings"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/rules"
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

// regimeGetter is an interface for types that expose a GetRegime method,
// including types that embed tax.Regime.
type regimeGetter interface {
	GetRegime() l10n.TaxCountryCode
}

// RulesContext implements rules.ContextAdder so that any struct embedding
// Regime automatically injects the regime code into the validation context.
// This allows guards like rules.HasContext(tax.RegimeIn("ES")) to work on
// nested objects without needing access to the root document.
func (r Regime) RulesContext() rules.WithContext {
	return func(rc *rules.RunCtx) {
		rc.Add(r)
	}
}

// RegimeContext returns a rules.WithContext option that injects the given
// regime code(s) into the validation context. Useful for testing rules against
// specific regimes without a fully calculated document.
func RegimeContext(codes ...l10n.TaxCountryCode) rules.WithContext {
	return func(rc *rules.RunCtx) {
		for _, code := range codes {
			rc.Add(Regime{Country: RegimeCode(code)})
		}
	}
}

// RegimeIn checks if the regime's country code is in the provided list of codes.
func RegimeIn(codes ...l10n.TaxCountryCode) rules.Test {
	str := make([]string, len(codes))
	for i, c := range codes {
		str[i] = c.String()
	}
	return rules.By("regime in ["+strings.Join(str, ",")+"]",
		func(value any) bool {
			rg, ok := value.(regimeGetter)
			if !ok {
				return false
			}
			return rg.GetRegime().In(codes...)
		},
	)
}

// JSONSchema provides a representation of the type for usage in Schema.
func (RegimeCode) JSONSchema() *jsonschema.Schema {
	defs := AllRegimeDefs()
	s := &jsonschema.Schema{
		Title:       "Tax Regime Code",
		Type:        "string",
		OneOf:       make([]*jsonschema.Schema, len(defs)),
		Description: `Identifies a GOBL tax regime`,
	}
	for i, rd := range defs {
		s.OneOf[i] = &jsonschema.Schema{
			Const: rd.Code().String(),
			Title: rd.Name.String(),
		}
	}
	return s
}
