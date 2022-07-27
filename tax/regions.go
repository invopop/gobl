package tax

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
)

var regions = newRegionCollection()

type regionCollection struct {
	list map[l10n.CountryCode][]*Region
}

func newRegionCollection() *regionCollection {
	c := new(regionCollection)
	c.list = make(map[l10n.CountryCode][]*Region)
	return c
}

func (c *regionCollection) add(r *Region) {
	if _, ok := c.list[r.Country]; !ok {
		c.list[r.Country] = make([]*Region, 0)
	}
	c.list[r.Country] = append(c.list[r.Country], r)
}

func (c *regionCollection) forIdentity(country l10n.CountryCode, locality l10n.Code) *Region {
	set, ok := c.list[country]
	if !ok {
		return nil
	}
	// First sweep
	for _, r := range set {
		if r.Locality == locality {
			return r
		}
	}
	// Second sweep in case there are multiple regions
	// in a given country, with a main region that does
	// not define a locality.
	for _, r := range set {
		if r.Locality == "" {
			return r
		}
	}
	return nil
}

func (c *regionCollection) all() []*Region {
	all := make([]*Region, 0)
	for _, set := range c.list {
		all = append(all, set...)
	}
	return all
}

// RegisterRegion adds a new region to the shared global list of tax regions.
func RegisterRegion(region *Region) {
	regions.add(region)
}

// RegionFor returns the region definition for country and locality combination
// or nil if no match was found.
func RegionFor(country l10n.CountryCode, locality l10n.Code) *Region {
	return regions.forIdentity(country, locality)
}

// AllRegions provides an array of all the region codes to definitions.
func AllRegions() []*Region {
	return regions.all()
}

// ValidateTaxIdentity attempts to find a matching region definition (if available)
// and runs tax identity validation.
func ValidateTaxIdentity(tID *org.TaxIdentity) error {
	if tID == nil {
		return nil
	}
	r := RegionFor(tID.Country, tID.Locality)
	if r == nil {
		return nil
	}
	if r.ValidateTaxIdentity == nil {
		return nil
	}
	return r.ValidateTaxIdentity(tID)
}
