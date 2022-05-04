package org

import (
	"regexp"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Item is used to describe a single product or service. Minimal usage
// implies just adding a name and price, more complete usage consists
// of adding descriptions, supplier IDs, SKUs, dimensions, etc.
//
// A set of additional code, ID, or SKU can be included in the `codes` property.
// Each `ItemCode` can be defined with an optional type agreed upon between the
// supplier and customer.
// For general purpose use, the Item's `Ref` property is much
// easier to use.
//
// We recommend setting prices with the item's "net" value, without tax,
// unless the document you're building supports the `price_includes_tax`
// option included in the `bill.Invoice` definition for example.
type Item struct {
	// Unique identify of this item independent of the Supplier IDs
	UUID string `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Primary reference code that identifies this item. Additional codes can be provided in the 'codes' field.
	Ref string `json:"ref,omitempty" jsonschema:"title=Ref"`
	// Brief name of the item
	Name string `json:"name"`
	// Detailed description
	Description string `json:"desc,omitempty"`
	// Currency used for the item's price.
	Currency string `json:"currency,omitempty" jsonschema:"title=Currency"`
	// Base price of a single unit to be sold.
	Price num.Amount `json:"price" jsonschema:"title=Price"`
	// Free-text unit of measure.
	Unit Unit `json:"unit,omitempty" jsonschema:"title=Unit"`
	//	List of additional codes, IDs, or SKUs which can be used to identify the item. The should be agreed upon between supplier and customer.
	Codes []*ItemCode `json:"codes,omitempty" jsonschema:"title=Codes"`
	// Country code of where this item was from originally.
	Origin l10n.Code `json:"origin,omitempty" jsonschema:"title=Country of Origin"`
	// Additional meta information that may be useful
	Meta Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

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

// ItemCode contains a value and optional label property that means additional
// codes can be added to an item.
type ItemCode struct {
	// Local or human reference for the type of code the value
	// represents.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// The item code's value.
	Value string `json:"value" jsonschema:"title=Value"`
}

var unitCodeRegexp = regexp.MustCompile(`^[a-z0-9]+$`)

// Validate checks that an address looks okay.
func (i *Item) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.Name, validation.Required),
		validation.Field(&i.Price, validation.Required),
		validation.Field(&i.Origin, l10n.IsCountry),
		validation.Field(&i.Unit, validation.Match(unitCodeRegexp)),
	)
}
