package tax

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
)

var regimes = newRegimeCollection()

type regimeCollection struct {
	list map[l10n.CountryCode][]*Regime
}

func newRegimeCollection() *regimeCollection {
	c := new(regimeCollection)
	c.list = make(map[l10n.CountryCode][]*Regime)
	return c
}

func (c *regimeCollection) add(r *Regime) {
	if _, ok := c.list[r.Country]; !ok {
		c.list[r.Country] = make([]*Regime, 0)
	}
	c.list[r.Country] = append(c.list[r.Country], r)
}

func (c *regimeCollection) forIdentity(country l10n.CountryCode, locality l10n.Code) *Regime {
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

func (c *regimeCollection) all() []*Regime {
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
	return regimes.forIdentity(country, locality)
}

// AllRegimes provides an array of all the regime codes to definitions.
func AllRegimes() []*Regime {
	return regimes.all()
}

// ValidateTaxIdentity attempts to find a matching regime definition (if available)
// and runs tax identity validation.
func ValidateTaxIdentity(tID *org.TaxIdentity) error {
	if tID == nil {
		return nil
	}
	r := RegimeFor(tID.Country, tID.Zone)
	if r == nil {
		return nil
	}
	if r.ValidateTaxIdentity == nil {
		return nil
	}
	return r.ValidateTaxIdentity(tID)
}

// NormalizeTaxIdentity attempts to find a matching regime definition (if available)
// and runs tax identity normalization.
func NormalizeTaxIdentity(tID *org.TaxIdentity) error {
	if tID == nil {
		return nil
	}
	r := RegimeFor(tID.Country, tID.Zone)
	if r == nil {
		return nil
	}
	if r.NormalizeTaxIdentity == nil {
		return nil
	}
	return r.NormalizeTaxIdentity(tID)
}
