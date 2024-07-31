package l10n

// CountryDefs provides helps for managing the list of countries
type CountryDefs []*CountryDef

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

// ISO provides a list of all the ISO countries.
func (cdfs CountryDefs) ISO() []*CountryDef {
	list := make([]*CountryDef, 0, countries.Len())
	for _, d := range countries {
		if !d.ISO {
			continue
		}
		list = append(list, d)
	}
	return list
}

// Tax provides a list of all the country definitions that can be used
// for tax purposes.
func (cdfs CountryDefs) Tax() []*CountryDef {
	list := make([]*CountryDef, 0, countries.Len())
	for _, d := range countries {
		if !d.Tax {
			continue
		}
		list = append(list, d)
	}
	return list
}
