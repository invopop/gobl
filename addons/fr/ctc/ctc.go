// Package ctc bundles the French CTC (Continuous Transaction Control)
// e-invoicing and e-reporting addon. It covers:
//
//   - Flow 2: domestic B2B clearance (cleared between two French parties).
//   - Flow 10: e-reporting for B2C, cross-border or other transactions
//     that fall outside the Flow 2 clearance perimeter.
//   - Flow 6: lifecycle status messages (Cycle de Vie) on bill.Status
//     documents exchanged between registered platforms.
//
// The invoice rule set is dispatched at validation time: an invoice
// whose supplier and customer both resolve as French (SIREN identity or
// French tax ID) runs the Flow 2 rule set; everything else runs the
// Flow 10 reporting rule set. Flow 6 operates on a separate document
// type (bill.Status) and does not need a predicate.
package ctc

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
)

const (
	// Key identifies the French CTC addon family.
	Key cbc.Key = "fr-ctc"

	// V1 is the key for the first version of the French CTC addon.
	V1 cbc.Key = Key + "-v1"
)

func init() {
	tax.RegisterAddonDef(newAddon())
	schema.Register(schema.GOBL.Add("addons/fr/ctc"),
		Characteristic{},
	)
	rules.RegisterWithGuard(
		Key.String(),
		rules.GOBL.Add("FR-CTC"),
		is.InContext(tax.AddonIn(V1)),
		billInvoiceRules(),
		billPaymentRules(),
		billStatusRules(),
		billReasonRules(),
		billActionRules(),
		orgPartyRules(),
		orgIdentityRules(),
		orgInboxRules(),
		orgItemRules(),
	)
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "France CTC",
			i18n.FR: "France CTC",
		},
		// eu-en16931-v2017 is required only when the dispatcher selects
		// Flow 2 (domestic French B2B clearance); the Flow 2 ruleset
		// enforces it. Flow 10 and Flow 6 work standalone.
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Support for the French CTC (Continuous Transaction Control)
				e-invoicing and e-reporting reform.

				The addon covers three of the flows defined by the French
				specification:

				- Flow 2 ("facturation"): domestic B2B clearance, applied
				  to invoices issued between two parties identifiable as
				  French (SIREN or French VAT ID on both sides).
				- Flow 10 ("e-reporting"): reporting of transactions that
				  fall outside Flow 2 clearance — B2C sales, cross-border
				  B2B, and payment receipts subject to e-reporting.
				- Flow 6 ("cycle de vie"): lifecycle status messages
				  (bill.Status) exchanged between registered platforms.

				The invoice ruleset is dispatched at validation time based
				on whether both parties resolve as French. There is no
				caller-facing switch: identify the parties correctly and
				the right flow runs.
			`),
			i18n.FR: here.Doc(`
				Support pour la réforme française CTC (Contrôle Continu
				des Transactions) de la facturation et du e-reporting.

				L'addon couvre trois flux du cahier des charges :

				- Flux 2 (« facturation ») : clearance B2B domestique,
				  appliqué aux factures émises entre deux parties
				  identifiables comme françaises (SIREN ou numéro de TVA
				  français des deux côtés).
				- Flux 10 (« e-reporting ») : déclaration des transactions
				  hors flux 2 — ventes B2C, B2B transfrontalières et
				  encaissements soumis au e-reporting.
				- Flux 6 (« cycle de vie ») : statuts cycle de vie
				  (bill.Status) échangés entre plateformes agréées.

				Le jeu de règles applicable aux factures est sélectionné
				au moment de la validation selon que les deux parties
				sont françaises ou non. Aucun commutateur explicite n'est
				exposé : il suffit d'identifier correctement les parties.
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
	case *bill.Status:
		normalizeStatus(obj)
	case *bill.Reason:
		normalizeReason(obj)
	case *org.Party:
		normalizeParty(obj)
	case *org.Identity:
		normalizeIdentity(obj)
	}
}
