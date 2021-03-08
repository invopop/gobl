package pack

import (
	"github.com/invopop/gobl"
	"github.com/invopop/gobl/tax"
)

// Pack represents the methods available for use inside a
// region.
type Pack interface {
	TaxDefs() *tax.Defs
	Validate(gobl.Document) error
}
