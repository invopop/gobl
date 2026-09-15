package untdid

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// unitCodes maps GOBL unit keys to the UN/ECE Recommendation 20 and 21 codes
// published in the UNTDID. It lives here, and not in org, so the generic unit
// definitions remain independent of any external code list, and so that every
// format built on these codes can share one table.
var unitCodes = map[cbc.Key]cbc.Code{
	org.UnitMilligram:        "MGM",
	org.UnitCentigram:        "CGM",
	org.UnitGram:             "GRM",
	org.UnitKilogram:         "KGM",
	org.UnitMetricTon:        "TNE",
	org.UnitMillimetre:       "MMT",
	org.UnitCentimetre:       "CMT",
	org.UnitDecimetre:        "DMT",
	org.UnitMetre:            "MTR",
	org.UnitLinearMetre:      "LM",
	org.UnitKilometre:        "KMT",
	org.UnitInch:             "INH",
	org.UnitFoot:             "FOT",
	org.UnitLinearFoot:       "LF",
	org.UnitSquareMilimetre:  "MMK",
	org.UnitSquareCentimetre: "CMK",
	org.UnitSquareDecimetre:  "DMK",
	org.UnitSquareMetre:      "MTK",
	org.UnitAcre:             "ACR",
	org.UnitHectare:          "HAR",
	org.UnitCubicMilimetre:   "MMQ",
	org.UnitCubicCentimetre:  "CMQ",
	org.UnitCubicDecimetre:   "DMQ",
	org.UnitCubicMetre:       "MTQ",
	org.UnitMillilitre:       "MLT",
	org.UnitCentilitre:       "CLT",
	org.UnitDecilitre:        "DLT",
	org.UnitLitre:            "LTR",
	org.UnitKilolitre:        "K6",
	org.UnitWatt:             "WTT",
	org.UnitKilowatt:         "KWT",
	org.UnitKilowattHour:     "KWH",
	org.UnitKilojoule:        "KJO",
	org.UnitKilocalorie:      "E14",
	org.UnitRate:             "A9",
	org.UnitYear:             "ANN",
	org.UnitMonth:            "MON",
	org.UnitWeek:             "WEE",
	org.UnitDay:              "DAY",
	org.UnitSecond:           "SEC",
	org.UnitHour:             "HUR",
	org.UnitMinute:           "MIN",
	org.UnitPiece:            "H87",
	org.UnitItem:             "EA",
	org.UnitPair:             "PR",
	org.UnitDozen:            "DZN",
	org.UnitAssortment:       "AS",
	org.UnitService:          "E48",
	org.UnitJob:              "E51",
	org.UnitActivity:         "ACT",
	org.UnitTrip:             "E54",
	org.UnitGroup:            "10",
	org.UnitOutfit:           "11",
	org.UnitKit:              "KT",
	org.UnitBaseBox:          "BB",
	org.UnitBulkPack:         "AB",
	org.UnitOne:              "C62",
	org.UnitBag:              "XBG",
	org.UnitBox:              "XBX",
	org.UnitBin:              "XBI",
	org.UnitCan:              "XCA",
	org.UnitTub:              "XTB",
	org.UnitCase:             "XCS",
	org.UnitTray:             "XDS",
	org.UnitPortion:          "13",
	org.UnitSet:              "SET",
	org.UnitRoll:             "XRO",
	org.UnitCarton:           "XCT",
	org.UnitCylinder:         "XCY",
	org.UnitBarrel:           "XBA",
	org.UnitJerrican:         "XJY",
	org.UnitCarboy:           "XCO",
	org.UnitDemijohn:         "XDJ",
	org.UnitBottle:           "XBO",
	org.UnitCanister:         "XCI",
	org.UnitPackage:          "XPK",
	org.UnitPacket:           "XPA",
	org.UnitBunch:            "XBH",
	org.UnitBundle:           "XBE",
	org.UnitBlock:            "XOK",
	org.UnitPallet:           "XPX",
	org.UnitReel:             "XRL",
	org.UnitSack:             "XSA",
	org.UnitSheet:            "XST",
	org.UnitEnvelope:         "XEN",
	org.UnitLot:              "XLT",
	org.UnitUnit:             "XUN",
}

// unitKeys reverses unitCodes for constant time lookups. The forward map is
// one-to-one, so each code resolves to a single unit.
var unitKeys = func() map[cbc.Code]cbc.Key {
	m := make(map[cbc.Code]cbc.Key, len(unitCodes))
	for unit, code := range unitCodes {
		m[code] = unit
	}
	return m
}()

// UnitCode converts a GOBL unit key into its corresponding UNTDID unit code.
// It returns an empty code when the unit has no standard mapping.
func UnitCode(unit cbc.Key) cbc.Code {
	return unitCodes[unit]
}

// UnitKey converts a UNTDID unit code into its corresponding GOBL unit key.
// It returns an empty key when GOBL has no standard mapping.
func UnitKey(code cbc.Code) cbc.Key {
	return unitKeys[code]
}

// NormalizeUnit returns the unit a document should carry alongside its
// extensions. The untdid-unit extension determines the unit whenever it is
// set, as it comes from the document format itself rather than the GOBL
// vocabulary, and leaves the unit empty for a code GOBL cannot express.
//
// The extension is never added nor removed: a document that carries one keeps
// it alongside the unit it resolves to, and a document without one is not
// given one just because its unit has a code.
func NormalizeUnit(unit cbc.Key, ext tax.Extensions) cbc.Key {
	if code := ext.Get(ExtKeyUnit); code != cbc.CodeEmpty {
		return UnitKey(code)
	}
	return unit
}
