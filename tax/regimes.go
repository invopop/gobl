package tax

import (
	"github.com/invopop/gobl/l10n"
)

var regimes = newRegimeCollection()

// RegimeCollection defines how to access details about all the regimes
// currently stored. Currently only a single tax regime per country is
// supported as we've not yet come across situations where multiple
// regimes exist within a single country.
type RegimeCollection struct {
	list map[l10n.Code]*Regime
}

// Regimes provides the current global regime collection object.
func Regimes() *RegimeCollection {
	return regimes
}

func newRegimeCollection() *RegimeCollection {
	c := new(RegimeCollection)
	c.list = make(map[l10n.Code]*Regime)
	return c
}

func (c *RegimeCollection) add(r *Regime) {
	c.list[r.Country.Code()] = r
	for _, cc := range r.AltCountryCodes {
		c.list[cc] = r
	}
}

// For provides a single matching regime from the collection, or nil if
// no match is found.
func (c *RegimeCollection) For(country l10n.Code) *Regime {
	r, ok := c.list[country]
	if !ok {
		return nil
	}
	return r
}

// All provides a list of all the registered Regimes.
func (c *RegimeCollection) All() []*Regime {
	all := make([]*Regime, len(c.list))
	i := 0
	for _, r := range c.list {
		all[i] = r
		i++
	}
	return all
}

// RegisterRegime adds a new regime to the shared global list of tax regimes.
func RegisterRegime(regime *Regime) {
	regimes.add(regime)
}

// RegimeFor returns the regime definition for country and locality combination
// or nil if no match was found.
func RegimeFor(country l10n.Code) *Regime {
	return regimes.For(country)
}

// AllRegimes provides an array of all the regime codes to definitions.
func AllRegimes() []*Regime {
	return regimes.All()
}
