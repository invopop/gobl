// Package choruspro provides extensions and validations for the Chorus Pro standard.
package choruspro

import (
	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

const (
	// V1 is the key for the Chorus Pro standard
	V1 cbc.Key = "fr-choruspro-v1"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "Chorus Pro",
			i18n.FR: "Chorus Pro",
		},
		Requires: []cbc.Key{
			en16931.V2017,
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Support for the French Chorus Pro platform for B2G (Business-to-Government) electronic invoicing.
				This addon provides the necessary structures and validations to ensure compliance with the 
				Chorus Pro specifications.

				It requires the EN16931 addon as it expands on the European standard with French-specific
				requirements for public sector invoicing.
			`),
			i18n.FR: here.Doc(`
				Support pour la plateforme française Chorus Pro pour la facturation électronique B2G 
				(Business-to-Government). Cet addon fournit les structures et validations nécessaires 
				pour assurer la conformité avec les spécifications Chorus Pro.

				Il nécessite l'addon EN16931 car il étend le standard européen avec des exigences 
				spécifiques françaises pour la facturation du secteur public.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Chorus Pro Specifications",
					i18n.FR: "Spécifications Chorus Pro",
				},
				URL: "https://communaute.chorus-pro.gouv.fr/wp-content/uploads/2018/11/External_Specifications_EDI_Appendix_V4.10.pdf",
			},
		},
		Extensions: extensions,
		Validator:  validate,
		Normalizer: normalize,
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *org.Party:
		return validateParty(obj)
	}
	return nil
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeInvoice(obj)
	case *org.Party:
		normalizeParty(obj)
	}
}
