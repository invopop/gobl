package it

import (
	"github.com/invopop/gobl/data"
	"github.com/invopop/gobl/tax"
)

// source: https://www1.agenziaentrate.gov.it/servizi/codici/ricerca/VisualizzaTabella.php?ArcName=ENTI-T2

var zones = tax.NewZoneStore(data.Content, "regimes/it.json")
