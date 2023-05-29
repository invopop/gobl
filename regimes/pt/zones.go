package pt

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

var zones = []tax.Zone{
	{Code: ZoneAveiro, Region: i18n.String{i18n.PT: "Aveiro"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: "PT"}},
	{Code: ZoneBeja, Region: i18n.String{i18n.PT: "Beja"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: "PT"}},
	{Code: ZoneBraga, Region: i18n.String{i18n.PT: "Braga"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: "PT"}},
	{Code: ZoneBraganca, Region: i18n.String{i18n.PT: "Bragança"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: "PT"}},
	{Code: ZoneCasteloBranco, Region: i18n.String{i18n.PT: "Castelo Branco"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: "PT"}},
	{Code: ZoneCoimbra, Region: i18n.String{i18n.PT: "Coimbra"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: "PT"}},
	{Code: ZoneEvora, Region: i18n.String{i18n.PT: "Évora"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: "PT"}},
	{Code: ZoneFaro, Region: i18n.String{i18n.PT: "Faro"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: "PT"}},
	{Code: ZoneGuarda, Region: i18n.String{i18n.PT: "Guarda"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: "PT"}},
	{Code: ZoneLeiria, Region: i18n.String{i18n.PT: "Leiria"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: "PT"}},
	{Code: ZoneLisboa, Region: i18n.String{i18n.PT: "Lisboa"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: "PT"}},
	{Code: ZonePortalegre, Region: i18n.String{i18n.PT: "Portalegre"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: "PT"}},
	{Code: ZonePorto, Region: i18n.String{i18n.PT: "Porto"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: "PT"}},
	{Code: ZoneSantarem, Region: i18n.String{i18n.PT: "Santarém"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: "PT"}},
	{Code: ZoneSetubal, Region: i18n.String{i18n.PT: "Setúbal"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: "PT"}},
	{Code: ZoneVianaDoCastelo, Region: i18n.String{i18n.PT: "Viana do Castelo"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: "PT"}},
	{Code: ZoneVilaReal, Region: i18n.String{i18n.PT: "Vila Real"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: "PT"}},
	{Code: ZoneViseu, Region: i18n.String{i18n.PT: "Viseu"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: "PT"}},
	{Code: ZoneAzores, Region: i18n.String{i18n.PT: "Região Autónoma dos Açores"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: "PT-AC"}},
	{Code: ZoneMadeira, Region: i18n.String{i18n.PT: "Região Autónoma da Madeira"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: "PT-MA"}},
}
