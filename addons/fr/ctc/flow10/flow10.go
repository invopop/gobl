// Package flow10 handles the extensions and validation rules for the French
// CTC (Cycle de Traitement de la Commande) Flow 10 e-reporting requirements
// for transactions not subject to Flow 2 domestic B2B clearance.
package flow10

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// Key identifies the French CTC Flow 10 addon family.
	Key cbc.Key = "fr-ctc-flow10"

	// V1 is the key for the French CTC Flow 10 addon
	V1 cbc.Key = Key + "-v1"
)

func init() {
	tax.RegisterAddonDef(newAddon())
	rules.RegisterWithGuard(
		Key.String(),
		rules.GOBL.Add("FR-CTC-FLOW10"),
		is.InContext(tax.AddonIn(V1)),
		billInvoiceRules(),
		billPaymentRules(),
	)
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "France CTC Flow 10",
			i18n.FR: "France CTC Flux 10",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Support for the French CTC (Continuous Transaction Control) Flow 10
				e-reporting requirements from the French electronic invoicing reform.

				Flow 10 covers transactions that must be reported to the tax authority
				but are not subject to the domestic B2B clearance flow (Flow 2). This
				includes B2C, cross-border, and out-of-scope transactions where VAT
				data and payment data must still be transmitted to the PPF.
			`),
			i18n.FR: here.Doc(`
				Support pour le CTC (Contrôle Continu des Transactions) français Flux 10
				pour les exigences de e-reporting de la réforme française de la
				facturation électronique.

				Le Flux 10 couvre les transactions qui doivent être déclarées à
				l'administration fiscale mais qui ne sont pas soumises au flux B2B
				domestique (Flux 2). Cela inclut les transactions B2C, transfrontalières
				et hors champ pour lesquelles les données de TVA et de paiement doivent
				tout de même être transmises au PPF.
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
		Scenarios:  scenarios,
		Normalizer: normalize,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeInvoice(obj)
	case *org.Party:
		normalizeParty(obj)
	}
}
