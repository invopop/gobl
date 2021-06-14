package region

import (
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/region/es"
	"github.com/invopop/gobl/tax"
)

// Code defines the code used to identify a region.
type Code string

// Defined region codes.
const (
	ES Code = es.Code // Spain
)

var regions = map[Code]Region{
	ES: es.New(), // Spain
}

// Region represents the methods we expect to be available from a region.
type Region interface {
	// Taxes provides a region's tax definition
	Taxes() *tax.Region

	// Currency provides the regions default currency definition.
	Currency() *currency.Def

	// ValidateTaxID checks the tax ID objects contents to see if they
	// are considered valid for the region.
	ValidateTaxID(id *org.TaxID) error
}

// Codes provides a list of region IDs that we know about.
func Codes() []Code {
	codes := make([]Code, len(regions))
	i := 0
	for c := range regions {
		codes[i] = c
		i++
	}
	return codes
}

// For returns the region definition for the document or nil if the
// region code is invalid.
func For(code Code) Region {
	return regions[code]
}

// List provides the list of regions and their definitions. Only really meant
// for exporting data.
func List() map[Code]Region {
	return regions
}
