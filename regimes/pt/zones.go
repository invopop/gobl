package pt

import (
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

var zones = []tax.Zone{
	{Code: ZoneAveiro, Region: i18n.String{i18n.PT: "Aveiro"}},
	{Code: ZoneBeja, Region: i18n.String{i18n.PT: "Beja"}},
	{Code: ZoneBraga, Region: i18n.String{i18n.PT: "Braga"}},
	{Code: ZoneBraganca, Region: i18n.String{i18n.PT: "Bragança"}},
	{Code: ZoneCasteloBranco, Region: i18n.String{i18n.PT: "Castelo Branco"}},
	{Code: ZoneCoimbra, Region: i18n.String{i18n.PT: "Coimbra"}},
	{Code: ZoneEvora, Region: i18n.String{i18n.PT: "Évora"}},
	{Code: ZoneFaro, Region: i18n.String{i18n.PT: "Faro"}},
	{Code: ZoneGuarda, Region: i18n.String{i18n.PT: "Guarda"}},
	{Code: ZoneLeiria, Region: i18n.String{i18n.PT: "Leiria"}},
	{Code: ZoneLisboa, Region: i18n.String{i18n.PT: "Lisboa"}},
	{Code: ZonePortalegre, Region: i18n.String{i18n.PT: "Portalegre"}},
	{Code: ZonePorto, Region: i18n.String{i18n.PT: "Porto"}},
	{Code: ZoneSantarem, Region: i18n.String{i18n.PT: "Santarém"}},
	{Code: ZoneSetubal, Region: i18n.String{i18n.PT: "Setúbal"}},
	{Code: ZoneVianaDoCastelo, Region: i18n.String{i18n.PT: "Viana do Castelo"}},
	{Code: ZoneVilaReal, Region: i18n.String{i18n.PT: "Vila Real"}},
	{Code: ZoneViseu, Region: i18n.String{i18n.PT: "Viseu"}},
	{Code: ZoneAzores, Region: i18n.String{i18n.PT: "Região Autónoma dos Açores"}},
	{Code: ZoneMadeira, Region: i18n.String{i18n.PT: "Região Autónoma da Madeira"}},
}
