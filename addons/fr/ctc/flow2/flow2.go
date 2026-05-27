// Package flow2 implements the French CTC Flow 2 ("facturation")
// rules for domestic B2B clearance invoices exchanged between two
// French parties. Built on the EU EN16931 CIUS profile (which is
// declared via the addon's Requires).
package flow2

import (
	"github.com/invopop/gobl/addons/eu/en16931"
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
	// Namespace is the rules namespace for the French CTC Flow 2 addon.
	Namespace rules.Code = "FR-CTC-FLOW2"

	// Key identifies the French CTC Flow 2 addon family.
	Key cbc.Key = "fr-ctc-flow2"

	// V1 is the first version of the French CTC Flow 2 addon.
	V1 cbc.Key = Key + "-v1"
)

func init() {
	tax.RegisterAddonDef(newV1Addon())
	rules.RegisterWithGuard(
		Key.String(),
		rules.GOBL.Add(Namespace),
		is.InContext(tax.AddonIn(V1)),
		billInvoiceRules(),
		orgPartyRules(),
		orgIdentityRules(),
		orgInboxRules(),
		orgItemRules(),
	)
}

func newV1Addon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Requires: []cbc.Key{
			en16931.V2017,
		},
		Name: i18n.String{
			i18n.EN: "France CTC Flow 2 (B2B Clearance)",
			i18n.FR: "France CTC Flux 2 (Clearance B2B)",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Support for Flow 2 ("facturation") of the French CTC
				reform: domestic B2B clearance, applied to invoices
				exchanged between two parties identifiable as French
				(SIREN or French VAT ID on both sides). Built on the
				EU EN16931 CIUS profile (declared as a dependency).
			`),
			i18n.FR: here.Doc(`
				Support du Flux 2 (« facturation ») de la réforme
				française CTC : clearance B2B domestique appliquée aux
				factures émises entre deux parties identifiables comme
				françaises (SIREN ou TVA française des deux côtés).
				Bâti sur le profil CIUS EU EN16931 (déclaré comme
				dépendance).
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
	case *org.Identity:
		normalizeIdentity(obj)
	}
}
