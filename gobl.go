package gobl

import (
	"github.com/invopop/gobl/region"
	"github.com/invopop/gobl/regions/es"
)

var regions = region.NewCollection(
	es.New(), // Spain
)

// Regions provides a region collection containing all the known region definitions.
func Regions() *region.Collection {
	return regions
}
