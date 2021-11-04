package gb

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/region"
	"github.com/invopop/gobl/tax"
)

// regionDef holds everything related to validation and taxes in Great Britain
type regionDef struct{}

// New provides the region definition
func New() region.Region {
	return new(regionDef)
}

// Code provides this region's code
func (regionDef) Code() region.Code {
	return region.GB
}

// Taxes provides all of this region's tax definitions.
func (regionDef) Taxes() *tax.Region {
	return &taxRegion
}

// Currency provides this region's main currency.
func (regionDef) Currency() *currency.Def {
	d, ok := currency.Get(currency.Code("GBP"))
	if !ok {
		return nil
	}
	return &d
}

// Validate checks the document type and determines if it can be validated.
func (r regionDef) Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	}
	return nil
}
