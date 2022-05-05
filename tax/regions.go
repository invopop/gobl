package tax

import "github.com/invopop/gobl/l10n"

var regions = newRegionCollection()

type regionCollection struct {
	list map[l10n.Code][]*Region
}

func newRegionCollection() *regionCollection {
	c := new(regionCollection)
	c.list = make(map[l10n.Code][]*Region)
	return c
}

func (c *regionCollection) add(r *Region) {
	if _, ok := c.list[r.Country]; !ok {
		c.list[r.Country] = make([]*Region, 0)
	}
	c.list[r.Country] = append(c.list[r.Country], r)
}

func (c *regionCollection) forIdentity(country, locality l10n.Code) *Region {
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
func RegionFor(country, locality l10n.Code) *Region {
	return regions.forIdentity(country, locality)
}

// AllRegions provides an array of all the region codes to definitions.
func AllRegions() []*Region {
	return regions.all()
}
