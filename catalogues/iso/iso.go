// Package iso is used to define ISO/IEC extensions and codes that may be used
// in documents.
package iso

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterCatalogueDef("iso.json")
}

const (
	// ExtKeySchemeID is used by the ISO 6523 scheme identifier.
	ExtKeySchemeID cbc.Key = "iso-scheme-id"
)
