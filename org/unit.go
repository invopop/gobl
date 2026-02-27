package org

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Unit represents either a unit key defined by GOBL *or* a two to three letter code
// defined by the UN/ECE.
type Unit string

const (
	// UnitPatternUNECE is a regular expression for UN/ECE unit codes when a unit is not covered by GOBL.
	UnitPatternUNECE = `^[A-Z0-9]{2,3}$`
	// UnitUNECEMutuallyDefined is the UN/ECE code for mutually defined units.
	UnitUNECEMutuallyDefined cbc.Code = `ZZ`
)

var regexpUNECEUnit = regexp.MustCompile(UnitPatternUNECE)

// Set of common units based on UN/ECE recommendation 20 and 21 extensions. Some local formats
// may define additional non-standard codes which may be added.
//
// The UN/ECE defines a very large set of units which would be impractical to support
// here in GOBL, so the Unit type will also accept any UN/ECE unit code instead of
// one of the keys defined here.
const (
	UnitEmpty Unit = `` // No unit defined

	// Measurement units
	UnitMilligram        Unit = `mg`
	UnitCentigram        Unit = `cg`
	UnitGram             Unit = `g`
	UnitKilogram         Unit = `kg`
	UnitMetricTon        Unit = `t`
	UnitMillimetre       Unit = `mm`
	UnitCentimetre       Unit = `cm`
	UnitDecimetre        Unit = `dm`
	UnitMetre            Unit = `m`
	UnitLinearMetre      Unit = `lm`
	UnitKilometre        Unit = `km`
	UnitInch             Unit = `in`
	UnitFoot             Unit = `ft`
	UnitLinearFoot       Unit = `lft`
	UnitSquareMilimetre  Unit = `mm2`
	UnitSquareCentimetre Unit = `cm2`
	UnitSquareDecimetre  Unit = `dm2`
	UnitSquareMetre      Unit = `m2`
	UnitHectare          Unit = `ha`
	UnitAcre             Unit = `ac`
	UnitCubicMilimetre   Unit = `mm3`
	UnitCubicCentimetre  Unit = `cm3`
	UnitCubicDecimetre   Unit = `dm3`
	UnitCubicMetre       Unit = `m3`
	UnitMillilitre       Unit = "ml"
	UnitCentilitre       Unit = `cl`
	UnitDecilitre        Unit = `dl`
	UnitLitre            Unit = `l`
	UnitKilolitre        Unit = `kl`
	UnitWatt             Unit = `w`
	UnitKilowatt         Unit = `kw`
	UnitKilowattHour     Unit = `kwh`
	UnitYear             Unit = `yr`
	UnitMonth            Unit = `mon`
	UnitWeek             Unit = `wk`
	UnitDay              Unit = `day`
	UnitSecond           Unit = `s`
	UnitHour             Unit = `h`
	UnitMinute           Unit = `min`
	UnitRate             Unit = `rate`
	UnitPiece            Unit = `piece`
	UnitItem             Unit = `item`
	UnitActivity         Unit = `activity`
	UnitService          Unit = `service`
	UnitGroup            Unit = `group`
	UnitSet              Unit = `set`
	UnitTrip             Unit = `trip`
	UnitJob              Unit = `job`
	UnitAssortment       Unit = `assortment`
	UnitOutfit           Unit = `outfit`
	UnitKit              Unit = `kit`
	UnitBaseBox          Unit = `basebox`
	UnitBulkPack         Unit = `pk`
	UnitOne              Unit = `one`
	UnitFlatFee          Unit = `ff`

	// Presentation Unit Codes
	UnitBag       Unit = `bag`
	UnitBox       Unit = `box`
	UnitBin       Unit = `bin`
	UnitCan       Unit = `can`
	UnitTub       Unit = `tub`
	UnitCase      Unit = `case`
	UnitTray      Unit = `tray`
	UnitPortion   Unit = `portion` // non-standard (src: ES)
	UnitDozen     Unit = `dozen`
	UnitPair      Unit = `pair`
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
	UnitPacket    Unit = `pkt`
	UnitBunch     Unit = `bunch`
	UnitBundle    Unit = `bdl`
	UnitBlock     Unit = `blk`
	UnitTetraBrik Unit = `tetrabrik` // non-standard (src: ES)
	UnitPallet    Unit = `pallet`
	UnitReel      Unit = `reel`
	UnitSack      Unit = `sack`
	UnitSheet     Unit = `sheet`
	UnitEnvelope  Unit = `envelope`
	UnitUnit      Unit = `unit`
	UnitLot       Unit = `lot`
)

// DefUnit serves to define unit keys.
type DefUnit struct {
	// Key for the Unit
	Unit Unit `json:"unit" jsonschema:"title=Unit"`
	// Name of the Unit
	Name string `json:"name" jsonschema:"title=Name"`
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
	{UnitMilligram, "Milligrams", "", "MGM"},
	{UnitCentigram, "Centigrams", "", "CGM"},
	{UnitGram, "Metric grams", "", "GRM"},
	{UnitKilogram, "Metric kilograms", "", "KGM"},
	{UnitMetricTon, "Metric tons", "", "TNE"},
	{UnitMillimetre, "Milimetres", "", "MMT"},
	{UnitCentimetre, "Centimetres", "", "CMT"},
	{UnitDecimetre, "Decimetres", "A unit of length equal to one-tenth of a metre.", "DMT"},
	{UnitMetre, "Metres", "", "MTR"},
	{UnitLinearMetre, "Linear metres", "A unit of length defining the number of metres of a linear item.", "LM"},
	{UnitKilometre, "Kilometers", "", "KMT"},
	{UnitInch, "Inches", "", "INH"},
	{UnitFoot, "Feet", "", "FOT"},
	{UnitLinearFoot, "Linear feet", "A unit of length defining the number of feet of a linear item.", "LF"},
	{UnitSquareMilimetre, "Square millimetres", "", "MMK"},
	{UnitSquareCentimetre, "Square centimetres", "", "CMK"},
	{UnitSquareDecimetre, "Square decimetres", "", "DMK"},
	{UnitSquareMetre, "Square metres", "", "MTK"},
	{UnitAcre, "Acres", "A unit of area equal to 43,560 square feet.", "ACR"},
	{UnitHectare, "Hectares", "A unit of area equal to 10,000 square metres.", "HAR"},
	{UnitCubicMilimetre, "Cubic millimetres", "", "MMQ"},
	{UnitCubicCentimetre, "Cubic centimetres", "", "CMQ"},
	{UnitCubicDecimetre, "Cubic decimetres", "", "DMQ"},
	{UnitCubicMetre, "Cubic metres", "", "MTQ"},
	{UnitMillilitre, "Millilitres", "", "MLT"},
	{UnitCentilitre, "Centilitres", "", "CLT"},
	{UnitDecilitre, "Decilitres", "", "DLT"},
	{UnitLitre, "Litres", "", "LTR"},
	{UnitKilolitre, "Kilolitres", "", ""},
	{UnitWatt, "Watts", "", "WTT"},
	{UnitKilowatt, "Kilowatts", "", "KWT"},
	{UnitKilowattHour, "Kilowatt Hours", "", "KWH"},
	{UnitRate, "Rate", "A unit of quantity expressed as a rate for usage of a facility or service.", "A9"},
	{UnitYear, "Years", "A unit of time equal to twelve months.", "ANN"},
	{UnitMonth, "Months", "Unit of time equal to 1/12 of a year of 365,25 days.", "MON"},
	{UnitWeek, "Weeks", "A unit of time equal to seven days.", "WEE"},
	{UnitDay, "Days", "", "DAY"},
	{UnitSecond, "Seconds", "", "SEC"},
	{UnitHour, "Hours", "", "HUR"},
	{UnitMinute, "Minutes", "", "MIN"},
	{UnitPiece, "Pieces", "A unit of count defining the number of pieces (piece: a single item, article or exemplar).", "H87"},
	{UnitItem, "Items", " A unit of count defining the number of items regarded as separate units.", "EA"},
	{UnitPair, "Pairs", "A unit of count defining the number of pairs (pair: item described by two's).", "PR"},
	{UnitDozen, "Dozens", "A unit of count defining the number of units in multiples of 12.", "DZN"},
	{UnitAssortment, "Assortments", "A unit of count defining the number of assortments (assortment: a collection of items or components of a single product packaged together).", "AS"},
	{UnitService, "Service Units", "A unit of count defining the number of service units (service unit: defined period / property / facility / utility of supply).", "E48"},
	{UnitJob, "Jobs", "A unit of count defining the number of jobs.", "E51"},
	{UnitActivity, "Activities", "A unit of count defining the number of activities (activity: a unit of work or action).", "ACT"},
	{UnitTrip, "Trips", "A unit of count defining the number of trips (trip: a journey to a place and back again).", "E54"},
	{UnitGroup, "Groups", "A unit of count defining the number of groups (group: set of items classified together).", "10"},
	{UnitOutfit, "Outfits", "A unit of count defining the number of outfits (outfit: a complete set of equipment / materials / objects used for a specific purpose).", "11"},
	{UnitKit, "Kits", "A unit of count defining the number of kits (kit: tub, barrel or pail).", "KT"},
	{UnitBaseBox, "Base Boxes", "A unit of area of 112 sheets of tin mil products (tin plate, tin free steel or black plate) 14 by 20 inches, or 31,360 square inches.", "BB"},
	{UnitBulkPack, "Bulk Packs", "A unit of count defining the number of items per bulk pack.", "AB"},
	{UnitOne, "One", "A single generic unit of a service or product.", "C62"},
	{UnitFlatFee, "Flat Fee", "A fixed one-time charge.", ""},

	// Recommendations Nº 21
	// source: https://unece.org/trade/documents/2021/06/uncefact-rec21
	{UnitBag, "Bags", "", "XBG"},
	{UnitBox, "Boxes", "", "XBX"},
	{UnitBin, "Bins", "", "XBI"},
	{UnitCan, "Cans", "", "XCA"},
	{UnitTub, "Tubs", "", "XTB"},
	{UnitCase, "Cases", "", "XCS"},
	{UnitTray, "Trays", "", "XDS"},    // plastic
	{UnitPortion, "Portions", "", ""}, // non-standard (src: ES)
	{UnitSet, "Sets", "A unit of count defining the number of sets (set: a number of objects grouped together).", "SET"},
	{UnitRoll, "Rolls", "", "XRO"},
	{UnitCarton, "Cartons", "", "XCT"},
	{UnitCylinder, "Cylinders", "", "XCY"},
	{UnitBarrel, "Barrels", "", "XBA"},
	{UnitJerrican, "Jerricans", "Jerrican, cylindrical", "XJY"},
	{UnitCarboy, "Carboys", "", "XCO"},    // non-protected
	{UnitDemijohn, "Demijohn", "", "XDJ"}, // non-protected
	{UnitBottle, "Bottles", "", "XBO"},    // non-protected, cylindrical
	{UnitSixPack, "Six Packs", "", ""},    // non-standard (src: ES)
	{UnitCanister, "Canisters", "", "XCI"},
	{UnitPackage, "Packages", "Standard packaging unit.", "XPK"},
	{UnitPacket, "Packets", "", "PA"},
	{UnitBunch, "Bunches", "", "XBH"},
	{UnitBundle, "Bundles", "", "BE"},
	{UnitBlock, "Blocks", "", ""},
	{UnitTetraBrik, "Tetra-Briks", "", ""}, // non-standard (src: ES)
	{UnitPallet, "Pallets", "", "XPX"},
	{UnitReel, "Reels", "", "XRL"},
	{UnitSack, "Sacks", "", "XSA"},
	{UnitSheet, "Sheets", "", "XST"},
	{UnitEnvelope, "Envelopes", "", "XEN"},
	{UnitLot, "Lot", "", "XLT"},
	{UnitUnit, "Unit", "A type of package composed of a single item or object, not otherwise specified as a unit of transport equipment.", "XUN"},
}

var isValidUnit = validation.In(validUnits()...)

func validUnits() []interface{} {
	list := make([]interface{}, len(UnitDefinitions))
	for i, d := range UnitDefinitions {
		list[i] = string(d.Unit)
	}
	return list
}

// Validate ensures the unit looks correct
func (u Unit) Validate() error {
	if regexpUNECEUnit.MatchString(string(u)) {
		return nil
	}
	return validation.Validate(string(u), isValidUnit.Error("must be a valid value or UN/ECE code"))
}

// UNECE provides the unit's UN/ECE equivalent value.
func (u Unit) UNECE() cbc.Code {
	if u == UnitEmpty {
		return cbc.CodeEmpty
	}
	// If already a UNECE code, return it.
	if regexpUNECEUnit.MatchString(string(u)) {
		return cbc.Code(string(u))
	}
	for _, def := range UnitDefinitions {
		if def.Unit == u {
			return def.UNECE
		}
	}
	return UnitUNECEMutuallyDefined // Assume something else.
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (u Unit) JSONSchema() *jsonschema.Schema {
	s := &jsonschema.Schema{
		Title:       "Unit",
		Type:        "string",
		OneOf:       make([]*jsonschema.Schema, len(UnitDefinitions)),
		Description: "Unit defines how the quantity of the product should be interpreted either using a GOBL lower-case key (e.g. 'kg'), or UN/ECE code upper-case code (e.g. 'KGM').",
	}
	for i, v := range UnitDefinitions {
		s.OneOf[i] = &jsonschema.Schema{
			Const:       v.Unit,
			Title:       v.Name,
			Description: v.Description,
		}
	}
	// Add the UN/ECE unit code pattern as an alternative to the pre-defined units.
	s.OneOf = append(s.OneOf, &jsonschema.Schema{
		Pattern:     UnitPatternUNECE,
		Description: "UN/ECE Unit Code from Recommendations 20 and 21",
	})
	return s
}
