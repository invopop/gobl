package region

import (
	"github.com/invopop/gobl/region/es"
	"github.com/invopop/gobl/region/pack"
	"github.com/invopop/gobl/tax"
)

var packs = map[tax.RegionID]pack.Pack{
	"es": es.New(), // Spain
}

// TaxRegionIDs provides a list of regions that we know about.
func TaxRegionIDs() []tax.RegionID {
	codes := make([]tax.RegionID, len(packs))
	i := 0
	for c := range packs {
		codes[i] = c
		i++
	}
	return codes
}

// PackFor returns the regional pack for the document or nil if the
// region code is invalid.
func PackFor(id tax.RegionID) pack.Pack {
	return packs[id]
}
