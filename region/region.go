package region

import (
	"github.com/invopop/gobl"
	"github.com/invopop/gobl/region/es"
	"github.com/invopop/gobl/tax"
)

var regions = map[tax.RegionID]Region{
	"es": es.New(), // Spain
}

// Region represents the methods we expect to be available from a region.
type Region interface {
	TaxDefs() *tax.Defs
	Validate(gobl.Document) error
}

// IDs provides a list of region IDs that we know about.
func IDs() []tax.RegionID {
	codes := make([]tax.RegionID, len(regions))
	i := 0
	for c := range regions {
		codes[i] = c
		i++
	}
	return codes
}

// For returns the region definition for the document or nil if the
// region code is invalid.
func For(id tax.RegionID) Region {
	return regions[id]
}
