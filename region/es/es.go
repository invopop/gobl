package es

import (
	"github.com/invopop/gobl"
	"github.com/invopop/gobl/region/pack"
	"github.com/invopop/gobl/tax"
)

type spain struct{}

// New returns a new instance of the spanish region pack
func New() pack.Pack {
	s := new(spain)
	return s
}

// TaxDefs provides all of this regions tax definitions.
func (spain) TaxDefs() *tax.Defs {
	return &taxDefs
}

func (spain) Validate(doc gobl.Document) error {

	return nil
}
