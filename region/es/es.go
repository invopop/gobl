package es

import (
	"github.com/invopop/gobl"
	"github.com/invopop/gobl/tax"
)

type Spain struct{}

// TaxRegion provides all of this regions tax definitions.
func (Spain) TaxRegion() *tax.Region {
	return &taxRegion
}

func (Spain) Validate(doc gobl.Document) error {

	return nil
}
