// Package iso is used to define ISO/IEC extensions and codes that may be used
// in documents.
package iso

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterCatalogueDef(newCatalogue())
}

func newCatalogue() *tax.CatalogueDef {
	return &tax.CatalogueDef{
		Key:  "iso",
		Name: i18n.NewString("ISO/IEC Data Elements"),
		Extensions: []*cbc.KeyDefinition{
			extSchemeID,
		},
	}
}
