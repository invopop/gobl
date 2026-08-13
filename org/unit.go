package org

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/jsonschema"
)

var regexpUNECEUnit = regexp.MustCompile(UnitPatternUNECE)

const (
	// UnitPatternUNECE identifies legacy UN/ECE codes that need to be migrated
	// to the untdid-unit extension.
	UnitPatternUNECE = `^[A-Z0-9]{2,3}$`
	// UnitMetaKeySymbol identifies the display symbol in a unit definition's metadata.
	UnitMetaKeySymbol cbc.Key = `symbol`
)

// Set of common units. Some local formats may define additional non-standard
// codes which may be added.
const (
	// Measurement units
	UnitMilligram        cbc.Key = `mg`
	UnitCentigram        cbc.Key = `cg`
	UnitGram             cbc.Key = `g`
	UnitKilogram         cbc.Key = `kg`
	UnitMetricTon        cbc.Key = `t`
	UnitMillimetre       cbc.Key = `mm`
	UnitCentimetre       cbc.Key = `cm`
	UnitDecimetre        cbc.Key = `dm`
	UnitMetre            cbc.Key = `m`
	UnitLinearMetre      cbc.Key = `lm`
	UnitKilometre        cbc.Key = `km`
	UnitInch             cbc.Key = `in`
	UnitFoot             cbc.Key = `ft`
	UnitLinearFoot       cbc.Key = `lft`
	UnitSquareMilimetre  cbc.Key = `mm2`
	UnitSquareCentimetre cbc.Key = `cm2`
	UnitSquareDecimetre  cbc.Key = `dm2`
	UnitSquareMetre      cbc.Key = `m2`
	UnitHectare          cbc.Key = `ha`
	UnitAcre             cbc.Key = `ac`
	UnitCubicMilimetre   cbc.Key = `mm3`
	UnitCubicCentimetre  cbc.Key = `cm3`
	UnitCubicDecimetre   cbc.Key = `dm3`
	UnitCubicMetre       cbc.Key = `m3`
	UnitMillilitre       cbc.Key = "ml"
	UnitCentilitre       cbc.Key = `cl`
	UnitDecilitre        cbc.Key = `dl`
	UnitLitre            cbc.Key = `l`
	UnitKilolitre        cbc.Key = `kl`
	UnitWatt             cbc.Key = `w`
	UnitKilowatt         cbc.Key = `kw`
	UnitKilowattHour     cbc.Key = `kwh`
	UnitKilojoule        cbc.Key = `kj`
	UnitKilocalorie      cbc.Key = `kcal`
	UnitYear             cbc.Key = `yr`
	UnitMonth            cbc.Key = `mon`
	UnitWeek             cbc.Key = `wk`
	UnitDay              cbc.Key = `day`
	UnitSecond           cbc.Key = `s`
	UnitHour             cbc.Key = `h`
	UnitMinute           cbc.Key = `min`
	UnitRate             cbc.Key = `rate`
	UnitPiece            cbc.Key = `piece`
	UnitItem             cbc.Key = `item`
	UnitActivity         cbc.Key = `activity`
	UnitService          cbc.Key = `service`
	UnitGroup            cbc.Key = `group`
	UnitSet              cbc.Key = `set`
	UnitTrip             cbc.Key = `trip`
	UnitJob              cbc.Key = `job`
	UnitAssortment       cbc.Key = `assortment`
	UnitOutfit           cbc.Key = `outfit`
	UnitKit              cbc.Key = `kit`
	UnitBaseBox          cbc.Key = `basebox`
	UnitBulkPack         cbc.Key = `pk`
	UnitOne              cbc.Key = `one`

	// Presentation Unit Codes
	UnitBag       cbc.Key = `bag`
	UnitBox       cbc.Key = `box`
	UnitBin       cbc.Key = `bin`
	UnitCan       cbc.Key = `can`
	UnitTub       cbc.Key = `tub`
	UnitCase      cbc.Key = `case`
	UnitTray      cbc.Key = `tray`
	UnitPortion   cbc.Key = `portion` // non-standard (src: ES)
	UnitDozen     cbc.Key = `dozen`
	UnitPair      cbc.Key = `pair`
	UnitRoll      cbc.Key = `roll`
	UnitCarton    cbc.Key = `carton`
	UnitCylinder  cbc.Key = `cylinder`
	UnitBarrel    cbc.Key = `barrel`
	UnitJerrican  cbc.Key = `jerrican`
	UnitCarboy    cbc.Key = `carboy`
	UnitDemijohn  cbc.Key = `demijohn`
	UnitBottle    cbc.Key = `bottle`
	UnitSixPack   cbc.Key = `6pack` // non-standard (src: ES)
	UnitCanister  cbc.Key = `canister`
	UnitPackage   cbc.Key = `pkg`
	UnitPacket    cbc.Key = `pkt`
	UnitBunch     cbc.Key = `bunch`
	UnitBundle    cbc.Key = `bdl`
	UnitBlock     cbc.Key = `blk`
	UnitTetraBrik cbc.Key = `tetrabrik` // non-standard (src: ES)
	UnitPallet    cbc.Key = `pallet`
	UnitReel      cbc.Key = `reel`
	UnitSack      cbc.Key = `sack`
	UnitSheet     cbc.Key = `sheet`
	UnitEnvelope  cbc.Key = `envelope`
	UnitUnit      cbc.Key = `unit`
	UnitLot       cbc.Key = `lot`
)

// UnitDefinitions describes each of the unit constants.
// Order is important.
var UnitDefinitions = []*cbc.Definition{
	// Measurement and count units.
	{Key: UnitMilligram, Name: i18n.NewString("Milligrams"), Meta: cbc.Meta{UnitMetaKeySymbol: "mg"}},
	{Key: UnitCentigram, Name: i18n.NewString("Centigrams"), Meta: cbc.Meta{UnitMetaKeySymbol: "cg"}},
	{Key: UnitGram, Name: i18n.NewString("Metric grams"), Meta: cbc.Meta{UnitMetaKeySymbol: "g"}},
	{Key: UnitKilogram, Name: i18n.NewString("Metric kilograms"), Meta: cbc.Meta{UnitMetaKeySymbol: "kg"}},
	{Key: UnitMetricTon, Name: i18n.NewString("Metric tons"), Meta: cbc.Meta{UnitMetaKeySymbol: "t"}},
	{Key: UnitMillimetre, Name: i18n.NewString("Millimetres"), Meta: cbc.Meta{UnitMetaKeySymbol: "mm"}},
	{Key: UnitCentimetre, Name: i18n.NewString("Centimetres"), Meta: cbc.Meta{UnitMetaKeySymbol: "cm"}},
	{Key: UnitDecimetre, Name: i18n.NewString("Decimetres"), Meta: cbc.Meta{UnitMetaKeySymbol: "dm"}, Desc: i18n.NewString("A unit of length equal to one-tenth of a metre.")},
	{Key: UnitMetre, Name: i18n.NewString("Metres"), Meta: cbc.Meta{UnitMetaKeySymbol: "m"}},
	{Key: UnitLinearMetre, Name: i18n.NewString("Linear metres"), Meta: cbc.Meta{UnitMetaKeySymbol: "lm"}, Desc: i18n.NewString("The running length in metres of a uniform-width product (e.g. carpet, fabric, cable), billed per metre regardless of width.")},
	{Key: UnitKilometre, Name: i18n.NewString("Kilometres"), Meta: cbc.Meta{UnitMetaKeySymbol: "km"}},
	{Key: UnitInch, Name: i18n.NewString("Inches"), Meta: cbc.Meta{UnitMetaKeySymbol: "in"}},
	{Key: UnitFoot, Name: i18n.NewString("Feet"), Meta: cbc.Meta{UnitMetaKeySymbol: "ft"}},
	{Key: UnitLinearFoot, Name: i18n.NewString("Linear feet"), Meta: cbc.Meta{UnitMetaKeySymbol: "lft"}, Desc: i18n.NewString("The running length in feet of a uniform-width product (e.g. lumber, trim, cable), billed per foot regardless of width.")},
	{Key: UnitSquareMilimetre, Name: i18n.NewString("Square millimetres"), Meta: cbc.Meta{UnitMetaKeySymbol: "mm²"}},
	{Key: UnitSquareCentimetre, Name: i18n.NewString("Square centimetres"), Meta: cbc.Meta{UnitMetaKeySymbol: "cm²"}},
	{Key: UnitSquareDecimetre, Name: i18n.NewString("Square decimetres"), Meta: cbc.Meta{UnitMetaKeySymbol: "dm²"}},
	{Key: UnitSquareMetre, Name: i18n.NewString("Square metres"), Meta: cbc.Meta{UnitMetaKeySymbol: "m²"}},
	{Key: UnitAcre, Name: i18n.NewString("Acres"), Desc: i18n.NewString("A unit of area equal to 43,560 square feet.")},
	{Key: UnitHectare, Name: i18n.NewString("Hectares"), Desc: i18n.NewString("A unit of area equal to 10,000 square metres.")},
	{Key: UnitCubicMilimetre, Name: i18n.NewString("Cubic millimetres"), Meta: cbc.Meta{UnitMetaKeySymbol: "mm³"}},
	{Key: UnitCubicCentimetre, Name: i18n.NewString("Cubic centimetres"), Meta: cbc.Meta{UnitMetaKeySymbol: "cm³"}},
	{Key: UnitCubicDecimetre, Name: i18n.NewString("Cubic decimetres"), Meta: cbc.Meta{UnitMetaKeySymbol: "dm³"}},
	{Key: UnitCubicMetre, Name: i18n.NewString("Cubic metres"), Meta: cbc.Meta{UnitMetaKeySymbol: "m³"}},
	{Key: UnitMillilitre, Name: i18n.NewString("Millilitres"), Meta: cbc.Meta{UnitMetaKeySymbol: "ml"}},
	{Key: UnitCentilitre, Name: i18n.NewString("Centilitres"), Meta: cbc.Meta{UnitMetaKeySymbol: "cl"}},
	{Key: UnitDecilitre, Name: i18n.NewString("Decilitres"), Meta: cbc.Meta{UnitMetaKeySymbol: "dl"}},
	{Key: UnitLitre, Name: i18n.NewString("Litres"), Meta: cbc.Meta{UnitMetaKeySymbol: "l"}},
	{Key: UnitKilolitre, Name: i18n.NewString("Kilolitres"), Meta: cbc.Meta{UnitMetaKeySymbol: "kl"}},
	{Key: UnitWatt, Name: i18n.NewString("Watts"), Meta: cbc.Meta{UnitMetaKeySymbol: "W"}},
	{Key: UnitKilowatt, Name: i18n.NewString("Kilowatts"), Meta: cbc.Meta{UnitMetaKeySymbol: "kW"}},
	{Key: UnitKilowattHour, Name: i18n.NewString("Kilowatt Hours"), Meta: cbc.Meta{UnitMetaKeySymbol: "kWh"}},
	{Key: UnitKilojoule, Name: i18n.NewString("Kilojoules"), Meta: cbc.Meta{UnitMetaKeySymbol: "kJ"}},
	{Key: UnitKilocalorie, Name: i18n.NewString("Kilocalories"), Meta: cbc.Meta{UnitMetaKeySymbol: "kcal"}},
	{Key: UnitRate, Name: i18n.NewString("Rate"), Desc: i18n.NewString("A unit of quantity expressed as a rate for usage of a facility or service.")},
	{Key: UnitYear, Name: i18n.NewString("Years"), Desc: i18n.NewString("A unit of time equal to twelve months.")},
	{Key: UnitMonth, Name: i18n.NewString("Months"), Desc: i18n.NewString("Unit of time equal to 1/12 of a year of 365,25 days.")},
	{Key: UnitWeek, Name: i18n.NewString("Weeks"), Desc: i18n.NewString("A unit of time equal to seven days.")},
	{Key: UnitDay, Name: i18n.NewString("Days")},
	{Key: UnitSecond, Name: i18n.NewString("Seconds")},
	{Key: UnitHour, Name: i18n.NewString("Hours")},
	{Key: UnitMinute, Name: i18n.NewString("Minutes")},
	{Key: UnitPiece, Name: i18n.NewString("Pieces"), Desc: i18n.NewString("A unit of count defining the number of pieces (piece: a single item, article or exemplar).")},
	{Key: UnitItem, Name: i18n.NewString("Items"), Desc: i18n.NewString("A unit of count defining the number of items regarded as separate units.")},
	{Key: UnitPair, Name: i18n.NewString("Pairs"), Desc: i18n.NewString("A unit of count defining the number of pairs (pair: item described by two's).")},
	{Key: UnitDozen, Name: i18n.NewString("Dozens"), Desc: i18n.NewString("A unit of count defining the number of units in multiples of 12.")},
	{Key: UnitAssortment, Name: i18n.NewString("Assortments"), Desc: i18n.NewString("A unit of count defining the number of assortments (assortment: a collection of items or components of a single product packaged together).")},
	{Key: UnitService, Name: i18n.NewString("Service Units"), Desc: i18n.NewString("A unit of count defining the number of service units (service unit: defined period / property / facility / utility of supply).")},
	{Key: UnitJob, Name: i18n.NewString("Jobs"), Desc: i18n.NewString("A unit of count defining the number of jobs.")},
	{Key: UnitActivity, Name: i18n.NewString("Activities"), Desc: i18n.NewString("A unit of count defining the number of activities (activity: a unit of work or action).")},
	{Key: UnitTrip, Name: i18n.NewString("Trips"), Desc: i18n.NewString("A unit of count defining the number of trips (trip: a journey to a place and back again).")},
	{Key: UnitGroup, Name: i18n.NewString("Groups"), Desc: i18n.NewString("A unit of count defining the number of groups (group: set of items classified together).")},
	{Key: UnitOutfit, Name: i18n.NewString("Outfits"), Desc: i18n.NewString("A unit of count defining the number of outfits (outfit: a complete set of equipment / materials / objects used for a specific purpose).")},
	{Key: UnitKit, Name: i18n.NewString("Kits"), Desc: i18n.NewString("A unit of count defining the number of kits (kit: tub, barrel or pail).")},
	{Key: UnitBaseBox, Name: i18n.NewString("Base Boxes"), Desc: i18n.NewString("A unit of area of 112 sheets of tin mil products (tin plate, tin free steel or black plate) 14 by 20 inches, or 31,360 square inches.")},
	{Key: UnitBulkPack, Name: i18n.NewString("Bulk Packs"), Desc: i18n.NewString("A unit of count defining the number of items per bulk pack.")},
	{Key: UnitOne, Name: i18n.NewString("One"), Desc: i18n.NewString("A single generic unit of a service or product.")},

	// Presentation units.
	{Key: UnitBag, Name: i18n.NewString("Bags")},
	{Key: UnitBox, Name: i18n.NewString("Boxes")},
	{Key: UnitBin, Name: i18n.NewString("Bins")},
	{Key: UnitCan, Name: i18n.NewString("Cans")},
	{Key: UnitTub, Name: i18n.NewString("Tubs")},
	{Key: UnitCase, Name: i18n.NewString("Cases")},
	{Key: UnitTray, Name: i18n.NewString("Trays")},       // plastic
	{Key: UnitPortion, Name: i18n.NewString("Portions")}, // non-standard (src: ES)
	{Key: UnitSet, Name: i18n.NewString("Sets"), Desc: i18n.NewString("A unit of count defining the number of sets (set: a number of objects grouped together).")},
	{Key: UnitRoll, Name: i18n.NewString("Rolls")},
	{Key: UnitCarton, Name: i18n.NewString("Cartons")},
	{Key: UnitCylinder, Name: i18n.NewString("Cylinders")},
	{Key: UnitBarrel, Name: i18n.NewString("Barrels")},
	{Key: UnitJerrican, Name: i18n.NewString("Jerricans"), Desc: i18n.NewString("Jerrican, cylindrical")},
	{Key: UnitCarboy, Name: i18n.NewString("Carboys")},     // non-protected
	{Key: UnitDemijohn, Name: i18n.NewString("Demijohns")}, // non-protected
	{Key: UnitBottle, Name: i18n.NewString("Bottles")},     // non-protected, cylindrical
	{Key: UnitSixPack, Name: i18n.NewString("Six Packs")},  // non-standard (src: ES)
	{Key: UnitCanister, Name: i18n.NewString("Canisters")},
	{Key: UnitPackage, Name: i18n.NewString("Packages"), Desc: i18n.NewString("Standard packaging unit.")},
	{Key: UnitPacket, Name: i18n.NewString("Packets")},
	{Key: UnitBunch, Name: i18n.NewString("Bunches")},
	{Key: UnitBundle, Name: i18n.NewString("Bundles")},
	{Key: UnitBlock, Name: i18n.NewString("Blocks")},
	{Key: UnitTetraBrik, Name: i18n.NewString("Tetra-Briks")}, // non-standard (src: ES)
	{Key: UnitPallet, Name: i18n.NewString("Pallets")},
	{Key: UnitReel, Name: i18n.NewString("Reels")},
	{Key: UnitSack, Name: i18n.NewString("Sacks")},
	{Key: UnitSheet, Name: i18n.NewString("Sheets")},
	{Key: UnitEnvelope, Name: i18n.NewString("Envelopes")},
	{Key: UnitLot, Name: i18n.NewString("Lot")},
	{Key: UnitUnit, Name: i18n.NewString("Unit"), Desc: i18n.NewString("A type of package composed of a single item or object, not otherwise specified as a unit of transport equipment.")},
}

// HasValidUnitKey validates that a key is one of the units defined by GOBL.
var HasValidUnitKey = cbc.InKeyDefs(UnitDefinitions)

// ExtendUnitKeySchema adds the available GOBL units to a cbc.Key property in
// the provided schema. Each model that uses a unit calls this helper so all
// schemas are derived from the same UnitDefinitions source.
func ExtendUnitKeySchema(schema *jsonschema.Schema, property string) {
	prop, ok := schema.Properties.Get(property)
	if !ok {
		return
	}
	prop.OneOf = make([]*jsonschema.Schema, len(UnitDefinitions))
	for i, def := range UnitDefinitions {
		prop.OneOf[i] = &jsonschema.Schema{
			Const:       def.Key,
			Title:       def.Name.String(),
			Description: def.Desc.String(),
		}
	}
}
