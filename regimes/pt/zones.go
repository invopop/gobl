package pt

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// AT Tax Country Regions
const (
	TaxCountryRegionPT = "PT"
	TaxCountryRegionAC = "PT-AC"
	TaxCountryRegionMA = "PT-MA"
)

var zones = []tax.Zone{
	{Code: ZoneAveiro, Region: i18n.String{i18n.PT: "Aveiro"}, Codes: cbc.CodeMap{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneBeja, Region: i18n.String{i18n.PT: "Beja"}, Codes: cbc.CodeMap{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneBraga, Region: i18n.String{i18n.PT: "Braga"}, Codes: cbc.CodeMap{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneBraganca, Region: i18n.String{i18n.PT: "Bragança"}, Codes: cbc.CodeMap{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneCasteloBranco, Region: i18n.String{i18n.PT: "Castelo Branco"}, Codes: cbc.CodeMap{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneCoimbra, Region: i18n.String{i18n.PT: "Coimbra"}, Codes: cbc.CodeMap{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneEvora, Region: i18n.String{i18n.PT: "Évora"}, Codes: cbc.CodeMap{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneFaro, Region: i18n.String{i18n.PT: "Faro"}, Codes: cbc.CodeMap{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneGuarda, Region: i18n.String{i18n.PT: "Guarda"}, Codes: cbc.CodeMap{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneLeiria, Region: i18n.String{i18n.PT: "Leiria"}, Codes: cbc.CodeMap{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneLisboa, Region: i18n.String{i18n.PT: "Lisboa"}, Codes: cbc.CodeMap{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZonePortalegre, Region: i18n.String{i18n.PT: "Portalegre"}, Codes: cbc.CodeMap{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZonePorto, Region: i18n.String{i18n.PT: "Porto"}, Codes: cbc.CodeMap{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneSantarem, Region: i18n.String{i18n.PT: "Santarém"}, Codes: cbc.CodeMap{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneSetubal, Region: i18n.String{i18n.PT: "Setúbal"}, Codes: cbc.CodeMap{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneVianaDoCastelo, Region: i18n.String{i18n.PT: "Viana do Castelo"}, Codes: cbc.CodeMap{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneVilaReal, Region: i18n.String{i18n.PT: "Vila Real"}, Codes: cbc.CodeMap{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneViseu, Region: i18n.String{i18n.PT: "Viseu"}, Codes: cbc.CodeMap{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneAzores, Region: i18n.String{i18n.PT: "Região Autónoma dos Açores"}, Codes: cbc.CodeMap{KeyATTaxCountryRegion: TaxCountryRegionAC}},
	{Code: ZoneMadeira, Region: i18n.String{i18n.PT: "Região Autónoma da Madeira"}, Codes: cbc.CodeMap{KeyATTaxCountryRegion: TaxCountryRegionMA}},
}
