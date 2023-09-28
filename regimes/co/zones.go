package co

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

// The Colombian tax agency (DIAN) requires a Municipality code for both suppliers and customers.
// This file contains the latest available data from:
//
// * https://www.dian.gov.co/atencionciudadano/formulariosinstructivos/Formularios/2007/Codigos_municipios_2007.pdf
// * https://github.com/ALAxHxC/MunicipiosDane
//

// Keys used in meta data
const (
	KeyZoneISO cbc.Key = "iso"
	KeyZoneDep cbc.Key = "dep"
)

var zones = tax.NewZoneStoreEmbedded(Data, "data/zones.json")
