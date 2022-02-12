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

var regions = newCollection()

// collection holds a set of regions.
type collection struct {
	regions map[Code]Region
}

// newCollection expects an array of regions, from which
func newCollection() *collection {
	c := new(collection)
	c.regions = make(map[Code]Region)
	return c
}

func (c *collection) add(r Region) {
	c.regions[r.Code()] = r
}

func (c *collection) codes() []Code {
	codes := make([]Code, len(c.regions))
	i := 0
	for code := range c.regions {
		codes[i] = code
		i++
	}
	return codes
}

func (c *collection) forCode(code Code) Region {
	return c.regions[code]
}

// Register adds a new region to the shared global list of regions.
func Register(region Region) {
	regions.add(region)
}

// For returns the region definition for the document or nil if the
// region code is invalid or has not been registered.
func For(code Code) Region {
	return regions.forCode(code)
}

// Codes provides a list of region codes contained in the collection
// of registered regions.
func Codes() []Code {
	return regions.codes()
}
