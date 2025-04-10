// Package choruspro implements the Chorus-Pro add-on for GOBL.
package choruspro

import (
	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

const (
	// V1 is the key for the Chorus-Pro add-on version 1.x
	V1 cbc.Key = "fr-choruspro-v1"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "French Chorus-Pro v1",
		},
		Requires: []cbc.Key{
			en16931.V2017,
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Support for the Chorus-Pro add-on for GOBL. This addon ensures that all fields required by Chorus-Pro are present and valid.
				This addon is does not enforce any specific format for the data, but will ensure that all required fields are present.
				This addon extends the EN16931 addon to ensure that the data is compliant with the Chorus-Pro specifications.


				This addon can then be used in addition to any format specific addon such as Factur-X, UBL, CII...
			`),
			// List of formats: https://portail.chorus-pro.gouv.fr/aife_documentation?id=kb_article_view&sys_kb_id=03feffcec333d65043b12775e00131fa
		},
		Normalizer: normalize,
		Validator:  validate,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeInvoice(obj)
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	}
	return nil
}
