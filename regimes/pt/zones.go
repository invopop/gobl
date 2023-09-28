package pt

import (
	"github.com/invopop/gobl/tax"
)

// AT Tax Country Regions
const (
	TaxCountryRegionPT = "PT"
	TaxCountryRegionAC = "PT-AC"
	TaxCountryRegionMA = "PT-MA"
)

var zones = tax.NewZoneStoreEmbedded(Data, "data/zones.json")
