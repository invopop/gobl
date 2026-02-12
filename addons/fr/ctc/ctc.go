// Package ctc handles the extensions and validation rules for the French
// CTC (Cycle de Traitement de la Commande) Flow 2 B2B e-invoicing requirements.
package ctc

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
	// Flow2 is the key for the French CTC addon
	Flow2 cbc.Key = "fr-ctc-flow2-v1"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: Flow2,
		Name: i18n.String{
			i18n.EN: "France CTC Flow 2",
			i18n.FR: "France CTC Flux 2",
		},
		Requires: []cbc.Key{
			en16931.V2017,
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Support for the French CTC (Continuous Transaction Control) Flow 2 B2B
				e-invoicing requirements from the French electronic invoicing reform.

				This addon provides the necessary structures and validations to ensure compliance
				with the French CTC specifications for B2B electronic invoicing.

				It requires the EN16931 addon as it extends the European standard with French-specific
				requirements for the e-invoicing reform.

				This addon is required for regulated invoice. This refers to invoices between two parties 
				registered for VAT in France. This addon should not be used for invoices which should be reported.

				Note on currency conversion (BR-FR-CO-12): When an invoice is issued in a non-EUR
				currency, the gobl.ubl library will automatically handle the conversion to EUR and
				present the invoice with both the original currency and EUR equivalents for tax
				amounts, ensuring compliance with French accounting requirements.
			`),
			i18n.FR: here.Doc(`
				Support pour le CTC (Contrôle Continu des Transactions) français Flux 2
				pour les exigences de facturation électronique B2B de la réforme française.

				Cet addon fournit les structures et validations nécessaires pour assurer la
				conformité avec les spécifications CTC françaises pour la facturation électronique B2B.

				Il nécessite l'addon EN16931 car il étend le standard européen avec des exigences
				spécifiques françaises pour la réforme de la facturation électronique.

				Cet addon est requis pour les factures réglementées. Cela concerne les factures entre                           
      			deux parties assujetties à la TVA en France. Cet addon ne doit pas être utilisé pour                                   
            	les factures qui doivent être déclarées.                                                                                                                               

				Note sur la conversion de devises (BR-FR-CO-12) : Lorsqu'une facture est émise dans
				une devise autre que l'EUR, la bibliothèque gobl.ubl gère automatiquement la conversion
				en EUR et présente la facture avec à la fois la devise d'origine et les équivalents en
				EUR pour les montants de TVA, garantissant la conformité avec les exigences comptables
				françaises.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "External Specifications",
					i18n.FR: "Spécifications Externes",
				},
				URL: "https://www.impots.gouv.fr/specifications-externes-b2b",
			},
		},
		Extensions: extensions,
		Tags: []*tax.TagSet{
			invoiceTags,
		},
		Normalizer: normalize,
		Validator:  validate,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeInvoice(obj)
	case *org.Party:
		normalizeParty(obj)
	case *org.Identity:
		normalizeIdentity(obj)
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *org.Party:
		return validateParty(obj)
	case *org.Identity:
		return validateIdentity(obj)
	case *org.Inbox:
		return validateInbox(obj)
	case *org.Item:
		return validateItem(obj)
	case []*org.Attachment:
		return validateOrgAttachments(obj)
	}
	return nil
}
