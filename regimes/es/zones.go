package es

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
)

// Zone code definitions for Spain
const (
	ZoneVI l10n.Code = "VI" // (01) Álava
	ZoneAB l10n.Code = "AB" // (02) Albacete
	ZoneA  l10n.Code = "A"  // (03) Alicante
	ZoneAL l10n.Code = "AL" // (04) Almería
	ZoneAV l10n.Code = "AV" // (05) Ávila
	ZoneBA l10n.Code = "BA" // (06) Badajoz
	ZonePM l10n.Code = "PM" // (07) Baleares
	ZoneIB l10n.Code = "IB" // (07) Baleares
	ZoneB  l10n.Code = "B"  // (08) Barcelona
	ZoneBU l10n.Code = "BU" // (09) Burgos
	ZoneCC l10n.Code = "CC" // (10) Cáceres
	ZoneCA l10n.Code = "CA" // (11) Cádiz
	ZoneCS l10n.Code = "CS" // (12) Castellon
	ZoneCR l10n.Code = "CR" // (13) Ciudad Real
	ZoneCO l10n.Code = "CO" // (14) Cordoba
	ZoneC  l10n.Code = "C"  // (15) La Coruña
	ZoneCU l10n.Code = "CU" // (16) Cuenca
	ZoneGE l10n.Code = "GE" // (17) Gerona
	ZoneGI l10n.Code = "GI" // (17) Girona
	ZoneGR l10n.Code = "GR" // (18) Granada
	ZoneGU l10n.Code = "GU" // (19) Guadalajara
	ZoneSS l10n.Code = "SS" // (20) Guipúzcoa
	ZoneH  l10n.Code = "H"  // (21) Huelva
	ZoneHU l10n.Code = "HU" // (22) Huesca
	ZoneJ  l10n.Code = "J"  // (23) Jaén
	ZoneLE l10n.Code = "LE" // (24) León
	ZoneL  l10n.Code = "L"  // (25) Lérida / Lleida
	ZoneLO l10n.Code = "LO" // (26) La Rioja
	ZoneLU l10n.Code = "LU" // (27) Lugo
	ZoneM  l10n.Code = "M"  // (28) Madrid
	ZoneMA l10n.Code = "MA" // (29) Málaga
	ZoneMU l10n.Code = "MU" // (30) Murcia
	ZoneNA l10n.Code = "NA" // (31) Navarra
	ZoneOR l10n.Code = "OR" // (32) Orense
	ZoneOU l10n.Code = "OU" // (32) Orense
	ZoneO  l10n.Code = "O"  // (33) Asturias
	ZoneP  l10n.Code = "P"  // (34) Palencia
	ZoneGC l10n.Code = "GC" // (35) Las Palmas
	ZonePO l10n.Code = "PO" // (36) Pontevedra
	ZoneSA l10n.Code = "SA" // (37) Salamanca
	ZoneTF l10n.Code = "TF" // (38) Santa Cruz de Tenerife
	ZoneS  l10n.Code = "S"  // (39) Cantabria
	ZoneSG l10n.Code = "SG" // (40) Segovia
	ZoneSE l10n.Code = "SE" // (41) Sevilla
	ZoneSO l10n.Code = "SO" // (42) Soria
	ZoneT  l10n.Code = "T"  // (43) Tarragona
	ZoneTE l10n.Code = "TE" // (44) Teruel
	ZoneTO l10n.Code = "TO" // (45) Toledo
	ZoneV  l10n.Code = "V"  // (46) Valencia
	ZoneVA l10n.Code = "VA" // (47) Valladolid
	ZoneBI l10n.Code = "BI" // (48) Vizcaya
	ZoneZA l10n.Code = "ZA" // (49) Zamora
	ZoneZ  l10n.Code = "Z"  // (50) Zaragoza
	ZoneCE l10n.Code = "CE" // (51) Ceuta
	ZoneML l10n.Code = "ML" // (52) Melilla
)

var zones = []tax.Zone{
	{
		Code:   ZoneVI,
		Region: i18n.String{i18n.ES: "Ávila"},
		Meta:   cbc.Meta{KeyAddressCode: "01"},
	},
	{
		Code:   ZoneAB,
		Region: i18n.String{i18n.ES: "Albacete"},
		Meta:   cbc.Meta{KeyAddressCode: "02"},
	},
	{
		Code:   ZoneA,
		Region: i18n.String{i18n.ES: "Alicante"},
		Meta:   cbc.Meta{KeyAddressCode: "03"},
	},
	{
		Code:   ZoneAL,
		Region: i18n.String{i18n.ES: "Almería"},
		Meta:   cbc.Meta{KeyAddressCode: "04"},
	},
	{
		Code:   ZoneAV,
		Region: i18n.String{i18n.ES: "Ávila"},
		Meta:   cbc.Meta{KeyAddressCode: "05"},
	},
	{
		Code:   ZoneBA,
		Region: i18n.String{i18n.ES: "Badajoz"},
		Meta:   cbc.Meta{KeyAddressCode: "06"},
	},
	{
		Code:   ZonePM,
		Region: i18n.String{i18n.ES: "Baleares"},
		Meta:   cbc.Meta{KeyAddressCode: "07"},
	},
	{
		Code:   ZoneIB,
		Region: i18n.String{i18n.ES: "Baleares"},
		Meta:   cbc.Meta{KeyAddressCode: "07"},
	},
	{
		Code:   ZoneB,
		Region: i18n.String{i18n.ES: "Barcelona"},
		Meta:   cbc.Meta{KeyAddressCode: "08"},
	},
	{
		Code:   ZoneBU,
		Region: i18n.String{i18n.ES: "Burgos"},
		Meta:   cbc.Meta{KeyAddressCode: "09"},
	},
	{
		Code:   ZoneCC,
		Region: i18n.String{i18n.ES: "Cáceres"},
		Meta:   cbc.Meta{KeyAddressCode: "10"},
	},
	{
		Code:   ZoneCA,
		Region: i18n.String{i18n.ES: "Cadiz"},
		Meta:   cbc.Meta{KeyAddressCode: "11"},
	},
	{
		Code:   ZoneCS,
		Region: i18n.String{i18n.ES: "Castellón"},
		Meta:   cbc.Meta{KeyAddressCode: "12"},
	},
	{
		Code:   ZoneCR,
		Region: i18n.String{i18n.ES: "Ciudad Real"},
		Meta:   cbc.Meta{KeyAddressCode: "13"},
	},
	{
		Code:   ZoneCO,
		Region: i18n.String{i18n.ES: "Cordoba"},
		Meta:   cbc.Meta{KeyAddressCode: "14"},
	},
	{
		Code:   ZoneC,
		Region: i18n.String{i18n.ES: "La Coruña"},
		Meta:   cbc.Meta{KeyAddressCode: "15"},
	},
	{
		Code:   ZoneCU,
		Region: i18n.String{i18n.ES: "Cuenca"},
		Meta:   cbc.Meta{KeyAddressCode: "16"},
	},
	{
		Code:   ZoneGE,
		Region: i18n.String{i18n.ES: "Gerona"},
		Meta:   cbc.Meta{KeyAddressCode: "17"},
	},
	{
		Code:   ZoneGI,
		Region: i18n.String{i18n.ES: "Girona"},
		Meta:   cbc.Meta{KeyAddressCode: "17"},
	},
	{
		Code:   ZoneGR,
		Region: i18n.String{i18n.ES: "Granada"},
		Meta:   cbc.Meta{KeyAddressCode: "18"},
	},
	{
		Code:   ZoneGU,
		Region: i18n.String{i18n.ES: "Guadalajara"},
		Meta:   cbc.Meta{KeyAddressCode: "19"},
	},
	{
		Code:   ZoneSS,
		Region: i18n.String{i18n.ES: "Guipúzcoa"},
		Meta:   cbc.Meta{KeyAddressCode: "20"},
	},
	{
		Code:   ZoneH,
		Region: i18n.String{i18n.ES: "Huelva"},
		Meta:   cbc.Meta{KeyAddressCode: "21"},
	},
	{
		Code:   ZoneHU,
		Region: i18n.String{i18n.ES: "Huesca"},
		Meta:   cbc.Meta{KeyAddressCode: "22"},
	},
	{
		Code:   ZoneJ,
		Region: i18n.String{i18n.ES: "Jaén"},
		Meta:   cbc.Meta{KeyAddressCode: "23"},
	},
	{
		Code:   ZoneLE,
		Region: i18n.String{i18n.ES: "León"},
		Meta:   cbc.Meta{KeyAddressCode: "24"},
	},
	{
		Code:   ZoneL,
		Region: i18n.String{i18n.ES: "Lérida / Lleida"},
		Meta:   cbc.Meta{KeyAddressCode: "25"},
	},
	{
		Code:   ZoneLO,
		Region: i18n.String{i18n.ES: "La Rioja"},
		Meta:   cbc.Meta{KeyAddressCode: "26"},
	},
	{
		Code:   ZoneLU,
		Region: i18n.String{i18n.ES: "Lugo"},
		Meta:   cbc.Meta{KeyAddressCode: "27"},
	},
	{
		Code:   ZoneM,
		Region: i18n.String{i18n.ES: "Madrid"},
		Meta:   cbc.Meta{KeyAddressCode: "28"},
	},
	{
		Code:   ZoneMA,
		Region: i18n.String{i18n.ES: "Málaga"},
		Meta:   cbc.Meta{KeyAddressCode: "29"},
	},
	{
		Code:   ZoneMU,
		Region: i18n.String{i18n.ES: "Murcia"},
		Meta:   cbc.Meta{KeyAddressCode: "30"},
	},
	{
		Code:   ZoneNA,
		Region: i18n.String{i18n.ES: "Navarra"},
		Meta:   cbc.Meta{KeyAddressCode: "31"},
	},
	{
		Code:   ZoneOR,
		Region: i18n.String{i18n.ES: "Orense"},
		Meta:   cbc.Meta{KeyAddressCode: "32"},
	},
	{
		Code:   ZoneOU,
		Region: i18n.String{i18n.ES: "Orense"},
		Meta:   cbc.Meta{KeyAddressCode: "32"},
	},
	{
		Code:   ZoneO,
		Region: i18n.String{i18n.ES: "Asturias"},
		Meta:   cbc.Meta{KeyAddressCode: "33"},
	},
	{
		Code:   ZoneP,
		Region: i18n.String{i18n.ES: "Palencia"},
		Meta:   cbc.Meta{KeyAddressCode: "34"},
	},
	{
		Code:   ZoneGC,
		Region: i18n.String{i18n.ES: "Las Palmas"},
		Meta:   cbc.Meta{KeyAddressCode: "35"},
	},
	{
		Code:   ZonePO,
		Region: i18n.String{i18n.ES: "Pontevedra"},
		Meta:   cbc.Meta{KeyAddressCode: "36"},
	},
	{
		Code:   ZoneSA,
		Region: i18n.String{i18n.ES: "Salamanca"},
		Meta:   cbc.Meta{KeyAddressCode: "37"},
	},
	{
		Code:   ZoneTF,
		Region: i18n.String{i18n.ES: "Santa Cruz de Tenerife"},
		Meta:   cbc.Meta{KeyAddressCode: "38"},
	},
	{
		Code:   ZoneS,
		Region: i18n.String{i18n.ES: "Cantabria"},
		Meta:   cbc.Meta{KeyAddressCode: "39"},
	},
	{
		Code:   ZoneSG,
		Region: i18n.String{i18n.ES: "Segovia"},
		Meta:   cbc.Meta{KeyAddressCode: "40"},
	},
	{
		Code:   ZoneSE,
		Region: i18n.String{i18n.ES: "Sevilla"},
		Meta:   cbc.Meta{KeyAddressCode: "41"},
	},
	{
		Code:   ZoneSO,
		Region: i18n.String{i18n.ES: "Soria"},
		Meta:   cbc.Meta{KeyAddressCode: "42"},
	},
	{
		Code:   ZoneT,
		Region: i18n.String{i18n.ES: "Tarragona"},
		Meta:   cbc.Meta{KeyAddressCode: "43"},
	},
	{
		Code:   ZoneTE,
		Region: i18n.String{i18n.ES: "Teruel"},
		Meta:   cbc.Meta{KeyAddressCode: "44"},
	},
	{
		Code:   ZoneTO,
		Region: i18n.String{i18n.ES: "Toledo"},
		Meta:   cbc.Meta{KeyAddressCode: "45"},
	},
	{
		Code:   ZoneV,
		Region: i18n.String{i18n.ES: "Valencia"},
		Meta:   cbc.Meta{KeyAddressCode: "46"},
	},
	{
		Code:   ZoneVA,
		Region: i18n.String{i18n.ES: "Valladolid"},
		Meta:   cbc.Meta{KeyAddressCode: "47"},
	},
	{
		Code:   ZoneBI,
		Region: i18n.String{i18n.ES: "Vizcaya"},
		Meta:   cbc.Meta{KeyAddressCode: "48"},
	},
	{
		Code:   ZoneZA,
		Region: i18n.String{i18n.ES: "Zamora"},
		Meta:   cbc.Meta{KeyAddressCode: "49"},
	},
	{
		Code:   ZoneZ,
		Region: i18n.String{i18n.ES: "Zaragoza"},
		Meta:   cbc.Meta{KeyAddressCode: "50"},
	},
	{
		Code:   ZoneCE,
		Region: i18n.String{i18n.ES: "Ceuta"},
		Meta:   cbc.Meta{KeyAddressCode: "51"},
	},
	{
		Code:   ZoneML,
		Region: i18n.String{i18n.ES: "Melilla"},
		Meta:   cbc.Meta{KeyAddressCode: "52"},
	},
}
