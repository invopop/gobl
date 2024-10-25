// Package untdid defines the UN/EDIFACT data elements contained in the UNTDID (United Nations Trade Data Interchange Directory).
package untdid

import (
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterCatalogueDef(newCatalogue())
}

func newCatalogue() *tax.CatalogueDef {
	return &tax.CatalogueDef{
		Key:        "untdid",
		Name:       i18n.NewString("UN/EDIFACT Data Elements"),
		Extensions: extensions,
	}
}
