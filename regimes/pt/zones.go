package pt

import (
	"github.com/invopop/gobl/data"
	"github.com/invopop/gobl/tax"
)

// AT Tax Country Regions
const (
	TaxCountryRegionPT = "PT"
	TaxCountryRegionAC = "PT-AC"
	TaxCountryRegionMA = "PT-MA"
)

var zones = tax.NewZoneStoreEmbedded(data.Content, "regimes/pt/zones.json")
