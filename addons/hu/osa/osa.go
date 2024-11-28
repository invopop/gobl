// Package osa handles the extensions and validation rules in order to use GOBL with the
// Hungarian OSA format.
package osa

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

const (
	// V3 for Hungary OSA XML v3.x
	V3 cbc.Key = "hu-osa-v3"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V3,
		Name: i18n.String{
			i18n.EN: "Hungary OSA v3.x",
		},
		Extensions: extensions,
	}
}
