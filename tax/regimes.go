package tax

import (
	"github.com/invopop/gobl/l10n"
)

var regimes = newRegimeCollection()

// RegimeCollection defines how to access details about all the regimes
// currently stored.
type RegimeCollection struct {
	list map[l10n.CountryCode][]*Regime
}

// Regimes provides the current global regime collection object.
func Regimes() *RegimeCollection {
	return regimes
}

func newRegimeCollection() *RegimeCollection {
	c := new(RegimeCollection)
	c.list = make(map[l10n.CountryCode][]*Regime)
	return c
}

func (c *RegimeCollection) add(r *Regime) {
	if _, ok := c.list[r.Country]; !ok {
		c.list[r.Country] = make([]*Regime, 0)
	}
	c.list[r.Country] = append(c.list[r.Country], r)
}

// For provides a single matching regime from the collection, or nil if
// no match is found.
func (c *RegimeCollection) For(country l10n.CountryCode, locality l10n.Code) *Regime {
	set, ok := c.list[country]
	if !ok {
		return nil
	}
	// First sweep
	for _, r := range set {
		if r.Zone == locality {
			return r
		}
	}
	// Second sweep in case there are multiple regimes
	// in a given country, with a main regime that does
	// not define a locality.
	for _, r := range set {
		if r.Zone == "" {
			return r
		}
	}
	return nil
}

// All provides a list of all the registered Regimes.
func (c *RegimeCollection) All() []*Regime {
	all := make([]*Regime, 0)
	for _, set := range c.list {
		all = append(all, set...)
	}
	return all
}

// RegisterRegime adds a new regime to the shared global list of tax regimes.
func RegisterRegime(regime *Regime) {
	regimes.add(regime)
}

// RegimeFor returns the regime definition for country and locality combination
// or nil if no match was found.
func RegimeFor(country l10n.CountryCode, locality l10n.Code) *Regime {
	return regimes.For(country, locality)
}

// AllRegimes provides an array of all the regime codes to definitions.
func AllRegimes() []*Regime {
	return regimes.All()
}
