package org

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/jsonschema"
)

// Unit is used to represent standard unit types.
type Unit string

// Set of common units based on UN/ECE recommendation 20 and 21. Some local formats
// may define additional non-standard codes which may be added. There are so
// many different unit codes in the world, that it's impractical to try and define them
// all, this is thus a selection of which we think are the most useful.
const (
	UnitEmpty Unit = `` // No unit defined

	// Measurement units
	UnitGram         Unit = `g`
	UnitKilogram     Unit = `kg`
	UnitMetricTon    Unit = `t`
	UnitMillimetre   Unit = `mm`
	UnitCentimetre   Unit = `cm`
	UnitMetre        Unit = `m`
	UnitKilometre    Unit = `km`
	UnitInch         Unit = `in`
	UnitFoot         Unit = `ft`
	UnitSquareMetre  Unit = `m2`
	UnitCubicMetre   Unit = `m3`
	UnitCentilitre   Unit = `cl`
	UnitLitre        Unit = `l`
	UnitWatt         Unit = `w`
	UnitKilowatt     Unit = `kw`
	UnitKilowattHour Unit = `kwh`
	UnitDay          Unit = `day`
	UnitSecond       Unit = `s`
	UnitHour         Unit = `h`
	UnitMinute       Unit = `min`
	UnitPiece        Unit = `piece`

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

// DefUnit serves to define unit keys.
type DefUnit struct {
	// Key for the Unit
	Unit Unit `json:"unit" jsonschema:"title=Unit"`
	// Description of the unit
	Description string `json:"description" jsonschema:"title=Description"`
	// Standard UN/ECE code
	UNECE cbc.Code `json:"unece" jsonschema:"title=UN/ECE Unit Code"`
}

// UnitDefinitions describes each of the unit constants.
// Order is important.
var UnitDefinitions = []DefUnit{
	// Recommendations Nº 20
	// source: https://unece.org/trade/documents/2021/06/uncefact-rec20-0
	{UnitGram, "Metric grams", "GRM"},
	{UnitKilogram, "Metric kilograms", "KGM"},
	{UnitMetricTon, "Metric tons", "TNE"},
	{UnitMillimetre, "Milimetres", "MMT"},
	{UnitCentimetre, "Centimetres", "CMT"},
	{UnitMetre, "Metres", "MTR"},
	{UnitKilometre, "Kilometers", "KMT"},
	{UnitInch, "Inches", "INH"},
	{UnitFoot, "Feet", "FOT"},
	{UnitSquareMetre, "Square metres", "MTK"},
	{UnitCubicMetre, "Cubic metres", "MTQ"},
	{UnitCentilitre, "Centilitres", "CLT"},
	{UnitLitre, "Litres", "LTR"},
	{UnitWatt, "Watts", "WTT"},
	{UnitKilowatt, "Kilowatts", "KWT"},
	{UnitKilowattHour, "Kilowatt Hours", "KWH"},
	{UnitDay, "Days", "DAY"},
	{UnitSecond, "Seconds", "SEC"},
	{UnitHour, "Hours", "HUR"},
	{UnitMinute, "Minutes", "MIN"},
	{UnitPiece, "Pieces", "H87"},

	// Recommendations Nº 21
	// source: https://unece.org/trade/documents/2021/06/uncefact-rec21
	{UnitBag, "Bags", "XBG"},
	{UnitBox, "Boxes", "XBX"},
	{UnitBin, "Bins", "XBI"},
	{UnitCan, "Cans", "XCA"},
	{UnitTub, "Tubs", "XTB"},
	{UnitCase, "Cases", "XCS"},
	{UnitTray, "Trays", "XDS"},    // plastic
	{UnitPortion, "Portions", ""}, // non-standard (src: ES)
	{UnitDozen, "Dozens", ""},     // non-standard (src: ES)
	{UnitRoll, "Rolls", "XRO"},
	{UnitCarton, "Cartons", "XCT"},
	{UnitCylinder, "Cylinders", "XCY"},
	{UnitBarrel, "Barrels", "XBA"},
	{UnitJerrican, "Jerricans", "XJY"}, // cylindrical
	{UnitCarboy, "Carboys", "XCO"},     // non-protected
	{UnitDemijohn, "Demijohn", "XDJ"},  // non-protected
	{UnitBottle, "Bottles", "XBO"},     // non-protected, cylindrical
	{UnitSixPack, "Six Packs", ""},     // non-standard (src: ES)
	{UnitCanister, "Canisters", "XCI"},
	{UnitPackage, "Packages", "XPK"},
	{UnitBunch, "Bunches", "XBH"},
	{UnitTetraBrik, "Tetra-Briks", ""}, // non-standard (src: ES)
	{UnitPallet, "Pallets", "XPX"},
	{UnitReel, "Reels", "XRL"},
	{UnitSack, "Sacks", "XSA"},
	{UnitSheet, "Sheets", "XST"},
	{UnitEnvelope, "Envelopes", "XEN"},
}

// Validate ensures the unit looks correct
func (u Unit) Validate() error {
	return validation.Validate(string(u), validation.Match(cbc.KeyValidationRegexp))
}

// UNECE provides the unit's UN/ECE equivalent
// value. If not available, returns CodeEmpty.
func (u Unit) UNECE() cbc.Code {
	for _, def := range UnitDefinitions {
		if def.Unit == u {
			return def.UNECE
		}
	}
	return cbc.CodeEmpty
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (u Unit) JSONSchema() *jsonschema.Schema {
	s := &jsonschema.Schema{
		Title:       "Unit",
		Type:        "string",
		AnyOf:       make([]*jsonschema.Schema, len(UnitDefinitions)),
		Description: "Unit describes how the quantity of the product should be interpreted.",
	}
	for i, v := range UnitDefinitions {
		s.AnyOf[i] = &jsonschema.Schema{
			Const:       v.Unit,
			Description: v.Description,
		}
	}
	// Add the custom unit to the end
	s.AnyOf = append(s.AnyOf, &jsonschema.Schema{
		Pattern:     cbc.KeyPattern,
		Description: "Custom unit definition",
	})
	return s
}
