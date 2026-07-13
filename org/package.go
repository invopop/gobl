package org

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
)

// UNECERec21Key is used in package key definition maps to provide the
// equivalent UN/ECE Recommendation 21 package type code.
const UNECERec21Key cbc.Key = "unece-rec21"

// UNECERec21MutuallyDefined is the UN/ECE Rec 21 code used when a package
// key cannot be mapped to a standard code.
const UNECERec21MutuallyDefined cbc.Code = "ZZ"

// Standard package keys with mappings to UN/ECE Recommendation 21 package
// type codes. Keys may be extended with sub-keys, e.g. "box+gift", and will
// map to the code of their base key.
const (
	PackageKeyBag      cbc.Key = "bag"
	PackageKeyBale     cbc.Key = "bale"
	PackageKeyBarrel   cbc.Key = "barrel"
	PackageKeyBin      cbc.Key = "bin"
	PackageKeyBox      cbc.Key = "box"
	PackageKeyBundle   cbc.Key = "bundle"
	PackageKeyCan      cbc.Key = "can"
	PackageKeyCarton   cbc.Key = "carton"
	PackageKeyCase     cbc.Key = "case"
	PackageKeyCrate    cbc.Key = "crate"
	PackageKeyDrum     cbc.Key = "drum"
	PackageKeyEnvelope cbc.Key = "envelope"
	PackageKeyPallet   cbc.Key = "pallet"
	PackageKeyReel     cbc.Key = "reel"
	PackageKeyRoll     cbc.Key = "roll"
	PackageKeySack     cbc.Key = "sack"
	PackageKeyTray     cbc.Key = "tray"
	PackageKeyTub      cbc.Key = "tub"
	PackageKeyTube     cbc.Key = "tube"
	PackageKeyUnpacked cbc.Key = "unpacked"
)

// PackageKeyDefinitions describes each of the standard package keys with
// their UN/ECE Recommendation 21 mappings.
var PackageKeyDefinitions = []*cbc.Definition{
	{
		Key:  PackageKeyBag,
		Name: i18n.NewString("Bag"),
		Map:  cbc.CodeMap{UNECERec21Key: "BG"},
	},
	{
		Key:  PackageKeyBale,
		Name: i18n.NewString("Bale"),
		Map:  cbc.CodeMap{UNECERec21Key: "BL"},
	},
	{
		Key:  PackageKeyBarrel,
		Name: i18n.NewString("Barrel"),
		Map:  cbc.CodeMap{UNECERec21Key: "BA"},
	},
	{
		Key:  PackageKeyBin,
		Name: i18n.NewString("Bin"),
		Map:  cbc.CodeMap{UNECERec21Key: "BI"},
	},
	{
		Key:  PackageKeyBox,
		Name: i18n.NewString("Box"),
		Map:  cbc.CodeMap{UNECERec21Key: "BX"},
	},
	{
		Key:  PackageKeyBundle,
		Name: i18n.NewString("Bundle"),
		Map:  cbc.CodeMap{UNECERec21Key: "BE"},
	},
	{
		Key:  PackageKeyCan,
		Name: i18n.NewString("Can"),
		Map:  cbc.CodeMap{UNECERec21Key: "CA"},
	},
	{
		Key:  PackageKeyCarton,
		Name: i18n.NewString("Carton"),
		Map:  cbc.CodeMap{UNECERec21Key: "CT"},
	},
	{
		Key:  PackageKeyCase,
		Name: i18n.NewString("Case"),
		Map:  cbc.CodeMap{UNECERec21Key: "CS"},
	},
	{
		Key:  PackageKeyCrate,
		Name: i18n.NewString("Crate"),
		Map:  cbc.CodeMap{UNECERec21Key: "CR"},
	},
	{
		Key:  PackageKeyDrum,
		Name: i18n.NewString("Drum"),
		Map:  cbc.CodeMap{UNECERec21Key: "DR"},
	},
	{
		Key:  PackageKeyEnvelope,
		Name: i18n.NewString("Envelope"),
		Map:  cbc.CodeMap{UNECERec21Key: "EN"},
	},
	{
		Key:  PackageKeyPallet,
		Name: i18n.NewString("Pallet"),
		Map:  cbc.CodeMap{UNECERec21Key: "PX"},
	},
	{
		Key:  PackageKeyReel,
		Name: i18n.NewString("Reel"),
		Map:  cbc.CodeMap{UNECERec21Key: "RL"},
	},
	{
		Key:  PackageKeyRoll,
		Name: i18n.NewString("Roll"),
		Map:  cbc.CodeMap{UNECERec21Key: "RO"},
	},
	{
		Key:  PackageKeySack,
		Name: i18n.NewString("Sack"),
		Map:  cbc.CodeMap{UNECERec21Key: "SA"},
	},
	{
		Key:  PackageKeyTray,
		Name: i18n.NewString("Tray"),
		Map:  cbc.CodeMap{UNECERec21Key: "DS"}, // plastic, matches UnitTray
	},
	{
		Key:  PackageKeyTub,
		Name: i18n.NewString("Tub"),
		Map:  cbc.CodeMap{UNECERec21Key: "TB"},
	},
	{
		Key:  PackageKeyTube,
		Name: i18n.NewString("Tube"),
		Map:  cbc.CodeMap{UNECERec21Key: "TU"},
	},
	{
		Key:  PackageKeyUnpacked,
		Name: i18n.NewString("Unpacked"),
		Desc: i18n.NewString("Unpacked or unpackaged goods."),
		Map:  cbc.CodeMap{UNECERec21Key: "NE"},
	},
}

// HasValidPackageKey provides a validator to ensure a package key is at
// least *based* on one of the standard keys, allowing extensions such as
// "box+gift".
var HasValidPackageKey = cbc.HasValidKeyIn(validBasePackageKeys()...)

func validBasePackageKeys() []cbc.Key {
	list := make([]cbc.Key, len(PackageKeyDefinitions))
	for i, v := range PackageKeyDefinitions {
		list[i] = v.Key
	}
	return list
}

// Package describes a physical package or container used to pack and ship
// an item, based on the UBL "Package" class.
type Package struct {
	uuid.Identify
	// Label for the package for presentation in output documents.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// Key for the type of package, based on the list pre-defined by GOBL
	// with mappings to UN/ECE Recommendation 21 codes.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// Count of packages of this same type, when more than one.
	Count int `json:"count,omitempty" jsonschema:"title=Count"`
	// Identities used to identify the individual package, such as a
	// GS1 Serial Shipping Container Code (SSCC).
	Identities []*Identity `json:"identities,omitempty" jsonschema:"title=Identities"`
	// Attributes describe features of the package, such as its dimensions
	// and gross weight.
	Attributes []*Attribute `json:"attributes,omitempty" jsonschema:"title=Attributes"`
}

func packageRules() *rules.Set {
	return rules.For(new(Package),
		rules.Field("key",
			rules.AssertIfPresent("01", "package key must be or extend one of the pre-defined keys",
				HasValidPackageKey,
			),
		),
		rules.Field("count",
			rules.Assert("02", "package count must be zero or positive", is.Min(0)),
		),
		rules.Field("attributes",
			rules.Assert("03", "package attributes must not contain duplicate keys",
				AttributesHaveUniqueKeys(),
			),
		),
	)
}

func normalizePackage(p *Package) {
	uuid.Normalize(&p.UUID)
	p.Label = cbc.NormalizeString(p.Label)
	p.Attributes = CleanAttributes(p.Attributes)
}

// UNECERec21 provides the UN/ECE Recommendation 21 package type code for
// the package's key, matching extended keys to their base definition.
// Returns "ZZ" (mutually defined) for unmatched keys, or an empty code
// when no key is set.
func (p *Package) UNECERec21() cbc.Code {
	if p == nil || p.Key == "" {
		return cbc.CodeEmpty
	}
	for _, def := range PackageKeyDefinitions {
		if p.Key == def.Key || p.Key.HasPrefix(def.Key) {
			return def.Map[UNECERec21Key]
		}
	}
	return UNECERec21MutuallyDefined
}

// IsEmpty returns true if the package has no meaningful content.
func (p *Package) IsEmpty() bool {
	return p == nil || (p.Label == "" &&
		p.Key == "" &&
		p.Count == 0 &&
		len(p.Identities) == 0 &&
		len(p.Attributes) == 0)
}

// CleanPackages removes any nil or empty packages from the list.
func CleanPackages(pkgs []*Package) []*Package {
	var cleaned []*Package
	for _, p := range pkgs {
		if p.IsEmpty() {
			continue
		}
		cleaned = append(cleaned, p)
	}
	return cleaned
}

// JSONSchemaExtend adds extra details to the schema.
func (Package) JSONSchemaExtend(js *jsonschema.Schema) {
	prop, ok := js.Properties.Get("key")
	if !ok {
		return
	}
	anyOf := make([]*jsonschema.Schema, 0, len(PackageKeyDefinitions)+1)
	for _, def := range PackageKeyDefinitions {
		anyOf = append(anyOf, &jsonschema.Schema{
			Const:       def.Key,
			Title:       def.Name.String(),
			Description: def.Desc.String(),
		})
	}
	anyOf = append(anyOf, &jsonschema.Schema{
		Title:   "Other",
		Pattern: cbc.KeyPattern,
	})
	prop.AnyOf = anyOf
}
