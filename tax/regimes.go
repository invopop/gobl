package tax

import (
	"github.com/invopop/gobl/l10n"
)

var regimes = newRegimeCollection()

// RegimeDefCollection defines how to access details about all the regimes
// currently stored. Currently only a single tax regime per country is
// supported as we've not yet come across situations where multiple
// regimes exist within a single country.
type RegimeDefCollection struct {
	list map[l10n.Code]*RegimeDef
}

// Regimes provides the current global regime collection object.
func Regimes() *RegimeDefCollection {
	return regimes
}

func newRegimeCollection() *RegimeDefCollection {
	c := new(RegimeDefCollection)
	c.list = make(map[l10n.Code]*RegimeDef)
	return c
}

func (c *RegimeDefCollection) add(r *RegimeDef) {
	c.list[r.Country.Code()] = r
	for _, cc := range r.AltCountryCodes {
		c.list[cc] = r
	}
}

// For provides a single matching regime from the collection, or nil if
// no match is found.
func (c *RegimeDefCollection) For(country l10n.Code) *RegimeDef {
	r, ok := c.list[country]
	if !ok {
		return nil
	}
	return r
}

// All provides a list of all the registered Regimes.
func (c *RegimeDefCollection) All() []*RegimeDef {
	all := make([]*RegimeDef, len(c.list))
	i := 0
	for _, r := range c.list {
		all[i] = r
		i++
	}
	return all
}

// RegisterRegimeDef adds a new regime to the shared global list of tax regimes.
func RegisterRegimeDef(regime *RegimeDef) {
	for _, ext := range regime.Extensions {
		RegisterExtension(ext)
	}
	regimes.add(regime)
}

// RegimeDefFor returns the regime definition for country and locality combination
// or nil if no match was found.
func RegimeDefFor(country l10n.Code) *RegimeDef {
	return regimes.For(country)
}

// AllRegimeDefs provides an array of all the regime codes to definitions.
func AllRegimeDefs() []*RegimeDef {
	return regimes.All()
}
