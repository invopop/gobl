package org

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Unit is used to represent standard unit types.
type Unit string

// Set of common units based on UN/ECE recommendation 20 and 21. Some local formats
// may define additional non-standard codes which may be added. There are so
// many different unit codes in the world, that it's impractical to try and define them
// all, this is thus a selection of which we think are the most useful.
const (
	// Default empty value is a "piece"
	UnitPiece Unit = ``

	// Measurement units
	UnitGram         Unit = `g`
	UnitKilogram     Unit = `kg`
	UnitMetricTon    Unit = `t`
	UnitMetre        Unit = `m`
	UnitCentimetre   Unit = `cm`
	UnitMillimetre   Unit = `mm`
	UnitKilometre    Unit = `km`
	UnitInch         Unit = `in`
	UnitSquareMetre  Unit = `m2`
	UnitCubicMetre   Unit = `m3`
	UnitCentilitres  Unit = `cl`
	UnitLitre        Unit = `l`
	UnitWatt         Unit = `w`
	UnitKilowatt     Unit = `kw`
	UnitKilowattHour Unit = `kwh`
	UnitDay          Unit = `day`
	UnitSecond       Unit = `s`
	UnitHour         Unit = `h`
	UnitMinute       Unit = `min`

	// Presentation Unit Codes
	UnitBag       Unit = `bag`
	UnitBox       Unit = `box`
	UnitBin       Unit = `bin`
	UnitCan       Unit = `can`
	UnitTub       Unit = `tub`
	UnitCase      Unit = `case`
	UnitTray      Unit = `tray`
	UnitPortion   Unit = `portion` // non-standard (src: ES)
	UnitDozen     Unit = `dozen`   // non-standard (src: ES)
	UnitRoll      Unit = `roll`
	UnitCarton    Unit = `carton`
	UnitCylinder  Unit = `cylinder`
	UnitBarrel    Unit = `barrel`
	UnitJerrican  Unit = `jerrican`
	UnitCarboy    Unit = `carboy`
	UnitDemijohn  Unit = `demijohn`
	UnitBottle    Unit = `bottle`
	UnitSixPack   Unit = `6pack` // non-standard (src: ES)
	UnitCanister  Unit = `canister`
	UnitPackage   Unit = `pkg`
	UnitBunch     Unit = `bunch`
	UnitTetraBrik Unit = `tetrabrik` // non-standard (src: ES)
	UnitPallet    Unit = `pallet`
	UnitReel      Unit = `reel`
	UnitSack      Unit = `sack`
	UnitSheet     Unit = `sheet`
	UnitEnvelope  Unit = `envelope`
)

// UNECEUnitMap defines the conversion of GOBL Unit Codes into
// UN/ECE recommended units.
var UNECEUnitMap = map[Unit]string{
	// Recommendations Nº 20
	// source: https://unece.org/trade/documents/2021/06/uncefact-rec20-0
	UnitPiece:        "H87",
	UnitGram:         "GRM",
	UnitKilogram:     "KGM",
	UnitMetricTon:    "TNE",
	UnitMetre:        "MTR",
	UnitCentimetre:   "CMT",
	UnitMillimetre:   "MMT",
	UnitCentilitres:  "CLT",
	UnitLitre:        "LTR",
	UnitSquareMetre:  "MTK",
	UnitCubicMetre:   "MTQ",
	UnitKilometre:    "KMT", // "KTM" is no longer used
	UnitWatt:         "WTT",
	UnitKilowatt:     "KWT",
	UnitKilowattHour: "KWH",
	UnitDay:          "DAY",
	UnitHour:         "HUR",
	UnitMinute:       "MIN",
	UnitSecond:       "SEC",

	// Recommendations Nº 21
	// source: https://unece.org/trade/documents/2021/06/uncefact-rec21
	UnitBag:      "XBG",
	UnitBox:      "XBX",
	UnitBin:      "XBI",
	UnitCase:     "XCS",
	UnitTub:      "XTB",
	UnitRoll:     "XRO",
	UnitCan:      `XCA`,
	UnitTray:     "XDS", // plastic
	UnitCarton:   "XCT",
	UnitCylinder: "XCY",
	UnitBarrel:   "XBA",
	UnitJerrican: "XJY", // cylindrical
	UnitBottle:   "XBO", // non-protected, cylindrical
	UnitCarboy:   "XCO", // non-protected
	UnitDemijohn: "XDJ", // non-protected
	UnitCanister: "XCI",
	UnitPackage:  "XPK",
	UnitBunch:    "XBH",
	UnitPallet:   "XPX",
	UnitReel:     "XRL",
	UnitSack:     "XSA",
	UnitSheet:    "XST",
	UnitEnvelope: "XEN",
}

var unitCodeRegexp = regexp.MustCompile(`^[a-z0-9]+$`)

// Validate ensures the unit looks correct
func (u Unit) Validate() error {
	return validation.Validate(u, validation.Match(unitCodeRegexp))
}
