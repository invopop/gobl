// Package flow10 implements the French CTC Flow 10 ("e-reporting")
// rules for transactions that fall outside the Flow 2 clearance
// perimeter: B2C sales, cross-border B2B invoices, and payment
// receipts subject to e-reporting to the French tax authority.
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
	// Namespace is the rules namespace for the French CTC Flow 10 addon.
	Namespace rules.Code = "FR-CTC-FLOW10"

	// Key identifies the French CTC Flow 10 addon family.
	Key cbc.Key = "fr-ctc-flow10"

	// V1 is the first version of the French CTC Flow 10 addon.
	V1 cbc.Key = Key + "-v1"
)

func init() {
	tax.RegisterAddonDef(newV1Addon())
	rules.RegisterWithGuard(
		Key.String(),
		rules.GOBL.Add(Namespace),
		is.InContext(tax.AddonIn(V1)),
		billInvoiceRules(),
		billPaymentRules(),
		orgPartyRules(),
	)
}

func newV1Addon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "France CTC Flow 10 (E-Reporting)",
			i18n.FR: "France CTC Flux 10 (E-Reporting)",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Support for the French CTC Flow 10 e-reporting reform.
				Covers transactions that fall outside Flow 2 clearance:
				B2C sales, cross-border B2B invoices, and payment receipts
				subject to e-reporting to the DGFiP via the PPF.
			`),
			i18n.FR: here.Doc(`
				Support pour le Flux 10 « e-reporting » de la réforme
				française CTC. Couvre les transactions hors clearance
				Flux 2 : ventes B2C, factures B2B transfrontalières et
				encaissements soumis au e-reporting à la DGFiP via le
				Portail Public de Facturation (PPF).
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
	case *org.Identity:
		normalizeIdentity(obj)
	case *tax.Combo:
		normalizeTaxCombo(obj)
	}
}
