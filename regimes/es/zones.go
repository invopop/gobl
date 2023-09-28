package es

import (
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

var zones = tax.NewZoneStoreEmbedded(Data, "data/zones.json")
