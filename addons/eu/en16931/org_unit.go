package en16931

import (
	"fmt"

	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// unitUNTDIDMap maps GOBL unit keys to UNTDID unit codes used by EN 16931.
// It deliberately lives in this package so the generic org unit definitions
// remain independent of any external code list.
var unitUNTDIDMap = map[cbc.Key]cbc.Code{
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

// untdidUnitMap reverses unitUNTDIDMap for constant time lookups. The forward
// map is one-to-one, so each code resolves to a single unit.
var untdidUnitMap = func() map[cbc.Code]cbc.Key {
	m := make(map[cbc.Code]cbc.Key, len(unitUNTDIDMap))
	for unit, code := range unitUNTDIDMap {
		m[code] = unit
	}
	return m
}()

// UnitToUNTDID converts a GOBL unit key into its corresponding UNTDID unit
// code. It returns an empty code when the unit has no standard mapping.
func UnitToUNTDID(unit cbc.Key) cbc.Code {
	return unitUNTDIDMap[unit]
}

// UnitFromUNTDID converts a UNTDID unit code into its corresponding GOBL unit
// key. It returns an empty unit when GOBL has no standard mapping.
func UnitFromUNTDID(code cbc.Code) cbc.Key {
	return untdidUnitMap[code]
}

func normalizeOrgItem(item *org.Item) {
	code := item.Ext.Get(untdid.ExtKeyUnit)
	if unit := UnitFromUNTDID(code); unit != cbc.KeyEmpty {
		item.Unit = unit
	}
	if item.Unit == cbc.KeyEmpty {
		item.Unit = org.UnitOne
	}
	if code == cbc.CodeEmpty {
		if code = UnitToUNTDID(item.Unit); code != cbc.CodeEmpty {
			item.Ext = item.Ext.Set(untdid.ExtKeyUnit, code)
		}
	}
}

func orgItemRules() *rules.Set {
	return rules.For(new(org.Item),
		rules.Assert("02", fmt.Sprintf("UNTDID unit `%s` must be present and valid (BR-23)", untdid.ExtKeyUnit),
			is.Func("required valid UNTDID unit", func(value any) bool {
				item, ok := value.(*org.Item)
				return ok && item != nil &&
					item.Ext.Has(untdid.ExtKeyUnit) &&
					tax.ExtensionHasValidCode(untdid.ExtKeyUnit).Check(item.Ext)
			}),
		),
	)
}
