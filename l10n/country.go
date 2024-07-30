package l10n

import (
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// ISOCountryCode defines an ISO 3166-2 country code.
type ISOCountryCode Code

// TaxCountryCode defines a code that may coincide with a country code,
// but may not always coincide.
type TaxCountryCode Code

// CountryDef provides the structure use to define a Country Code
// definition.
type CountryDef struct {
	// ISO 3166-2 Country code
	Code Code `json:"code" jsonschema:"title=ISO 3166-2 Country Code"`
	// ISO 3166-1 alpha-3 Country code
	Alpha3 string `json:"alpha3" jsonschema:"title=ISO 3166-1 Alpha-3 Country Code"`
	// English name of the country
	Name string `json:"name" jsonschema:"title=Name"`
	// Internet Top-Level-Domain
	TLD string `json:"tld" jsonschema:"title=Top level domain"`
	// When true, the country is an ISO code.
	ISO bool `json:"iso" jsonschema:"title=ISO"`
	// When true, the code can be used for tax purposes.
	Tax bool `json:"tax" jsonschema:"title=Tax"`
	// Indicates that a different code can be used for lookups.
	AltCode Code `json:"alt_code" jsonschema:"title=Tax Code"`
}

// CountryDefs provides helps for managing the list of countries
type CountryDefs []*CountryDef

func validISOCountryCodes() []ISOCountryCode {
	list := make([]ISOCountryCode, 0, CountryDefinitions.Len())
	for _, d := range CountryDefinitions {
		if !d.ISO {
			continue
		}
		list = append(list, ISOCountryCode(d.Code))
	}
	return list
}

func validTaxCountryCodes() []TaxCountryCode {
	list := make([]TaxCountryCode, 0, CountryDefinitions.Len())
	for _, d := range CountryDefinitions {
		if !d.Tax {
			continue
		}
		list = append(list, TaxCountryCode(d.Code))
	}
	return list
}

var (
	isISOCountry = validation.In(validISOCountryCodes()...).Error("must be a valid ISO country code")
	isTaxCountry = validation.In(validTaxCountryCodes()...).Error("must be a valid tax country code")
)

// Validate ensures the ISO country code is inside the known and valid
// list of countries.
func (c ISOCountryCode) Validate() error {
	return isISOCountry.Validate(c)
}

// Validate ensures the tax country code is inside the known and valid
// list of country codes for taxes.
func (c TaxCountryCode) Validate() error {
	return isTaxCountry.Validate(c)
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (ISOCountryCode) JSONSchema() *jsonschema.Schema {
	s := &jsonschema.Schema{
		Title:       "ISO Country Code",
		Type:        "string",
		OneOf:       make([]*jsonschema.Schema, 0, CountryDefinitions.Len()),
		Description: `Defines an ISO 3166-2 country code`,
	}
	for _, d := range CountryDefinitions {
		if !d.ISO {
			continue
		}
		s.OneOf = append(s.OneOf, &jsonschema.Schema{
			Const: d.Code,
			Title: d.Name,
		})
	}
	return s
}

// In returns true if the country code is contained inside the provided set
func (c ISOCountryCode) In(set ...ISOCountryCode) bool {
	for _, x := range set {
		if c == x {
			return true
		}
	}
	return false
}

// String provides string representation of the ISO country code
func (c ISOCountryCode) String() string {
	return string(c)
}

// Name provides the ISO Country Name for the code
func (c ISOCountryCode) Name() string {
	if d := CountryDefinitions.Code(Code(c)); d != nil {
		return d.Name
	}
	return ""
}

// Alpha3 provides the ISO 3166-1 alpha-3 country code
func (c ISOCountryCode) Alpha3() string {
	if d := CountryDefinitions.Code(Code(c)); d != nil {
		return d.Alpha3
	}
	return ""
}

// In returns true if the tax country code is contained inside the provided set
func (c TaxCountryCode) In(set ...TaxCountryCode) bool {
	for _, x := range set {
		if c == x {
			return true
		}
	}
	return false
}

// String provides string representation of the tax country code
func (c TaxCountryCode) String() string {
	return string(c)
}

// Name provides the ISO Country Name for the code
func (c TaxCountryCode) Name() string {
	if d := CountryDefinitions.Code(Code(c)); d != nil {
		return d.Name
	}
	return ""
}

// JSONSchema provides a representation of the type for usage in Schema.
func (TaxCountryCode) JSONSchema() *jsonschema.Schema {
	s := &jsonschema.Schema{
		Title:       "Tax Country Code",
		Type:        "string",
		OneOf:       make([]*jsonschema.Schema, 0, CountryDefinitions.Len()),
		Description: `Defines an ISO base country code used for tax purposes`,
	}
	for _, d := range CountryDefinitions {
		if !d.Tax {
			continue
		}
		s.OneOf = append(s.OneOf, &jsonschema.Schema{
			Const: d.Code,
			Title: d.Name,
		})
	}
	return s
}

// Len provides the length of the country definitions
func (cdfs CountryDefs) Len() int {
	return len(cdfs)
}

// Code finds the country definition for the given country code
func (cdfs CountryDefs) Code(c Code) *CountryDef {
	for _, v := range cdfs {
		if v.Code == c {
			return v
		}
	}
	return nil
}
