package region

import (
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/tax"
)

// Code defines the code used to identify a region.
type Code string

// Defined region codes.
const (
	ES Code = "ES" // Spain
	NL Code = "NL" // Netherlands
	GB Code = "GB" // Great Britain (not Northern Ireland)
)

// Region represents the methods we expect to be available from a region.
type Region interface {
	// Taxes provides a region's tax definition
	Taxes() *tax.Region

	// Currency provides the regions default currency definition.
	Currency() *currency.Def

	// Code provides this region's official code.
	Code() Code

	// Validate determines the type of GOBL document provided and attempts
	// to check the contents for errors.
	Validate(obj interface{}) error
}

// Collection holds a set of regions.
type Collection struct {
	regions map[Code]Region
}

// NewCollection expects an array of regions, from which
func NewCollection(regions ...Region) *Collection {
	c := new(Collection)
	c.regions = make(map[Code]Region)
	for _, r := range regions {
		c.regions[r.Code()] = r
	}
	return c
}

// Codes provides a list of region codes contained in the collection.
func (c *Collection) Codes() []Code {
	codes := make([]Code, len(c.regions))
	i := 0
	for code := range c.regions {
		codes[i] = code
		i++
	}
	return codes
}

// For returns the region definition for the document or nil if the
// region code is invalid.
func (c *Collection) For(code Code) Region {
	return c.regions[code]
}

// List provides the list of regions and their definitions. Only really meant
// for exporting data.
func (c *Collection) List() map[Code]Region {
	return c.regions
}
