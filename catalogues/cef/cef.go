// Package cef provides codes issue by the "Connecting Europe Facility"
// (CEF Digital) initiative.
package cef

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
		Key:  "cef",
		Name: i18n.NewString("Connecting Europe Facility (CEF)"),
		Extensions: []*cbc.Definition{
			extVATEX,
		},
	}
}
