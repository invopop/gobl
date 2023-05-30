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
	{Code: ZoneAveiro, Region: i18n.String{i18n.PT: "Aveiro"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneBeja, Region: i18n.String{i18n.PT: "Beja"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneBraga, Region: i18n.String{i18n.PT: "Braga"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneBraganca, Region: i18n.String{i18n.PT: "Bragança"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneCasteloBranco, Region: i18n.String{i18n.PT: "Castelo Branco"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneCoimbra, Region: i18n.String{i18n.PT: "Coimbra"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneEvora, Region: i18n.String{i18n.PT: "Évora"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneFaro, Region: i18n.String{i18n.PT: "Faro"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneGuarda, Region: i18n.String{i18n.PT: "Guarda"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneLeiria, Region: i18n.String{i18n.PT: "Leiria"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneLisboa, Region: i18n.String{i18n.PT: "Lisboa"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZonePortalegre, Region: i18n.String{i18n.PT: "Portalegre"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZonePorto, Region: i18n.String{i18n.PT: "Porto"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneSantarem, Region: i18n.String{i18n.PT: "Santarém"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneSetubal, Region: i18n.String{i18n.PT: "Setúbal"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneVianaDoCastelo, Region: i18n.String{i18n.PT: "Viana do Castelo"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneVilaReal, Region: i18n.String{i18n.PT: "Vila Real"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneViseu, Region: i18n.String{i18n.PT: "Viseu"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: TaxCountryRegionPT}},
	{Code: ZoneAzores, Region: i18n.String{i18n.PT: "Região Autónoma dos Açores"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: TaxCountryRegionAC}},
	{Code: ZoneMadeira, Region: i18n.String{i18n.PT: "Região Autónoma da Madeira"}, Codes: cbc.CodeSet{KeyATTaxCountryRegion: TaxCountryRegionMA}},
}
