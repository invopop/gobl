// Package zugferd adds validations and checks for successful conversion to the ZUGFeRD format.
package zugferd

import (
	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

const (
	// V2 is the key for ZUGFeRD V2.2 and upwards.
	V2 cbc.Key = "de-zugferd-v2"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V2,
		Name: i18n.String{
			i18n.EN: "German ZUGFeRD 2.X",
		},
		Requires: []cbc.Key{
			en16931.V2017,
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Support for the German ZUGFeRD version 2.2 and upwards standard for electronic
				invoicing. ZUGFeRD 2.2 and Factur-X 1.0 are fully compatible and technically
				identical formats the use the Factur-X identifier.
				
				Currently this is just a placeholder addon as the EN16931 addon will provide
				all validation requirements.
				
				For more information, visit [www.ferd-net.de](https://www.ferd-net.de/).
			`),
		},
	}
}
