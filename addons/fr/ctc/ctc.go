// Package ctc registers the French CTC meta-addon (fr-ctc-v1) which
// inspects each GOBL document at normalize-time and appends the
// appropriate flow-specific addon (fr-ctc-flow2-v1, fr-ctc-flow6-v1,
// or fr-ctc-flow10-v1) so the right validation rules fire without the
// caller having to know which flow applies.
//
// The flow-specific addons live in subpackages and are independent:
// callers can declare any of them directly if they prefer.
package ctc

import (
	"github.com/invopop/gobl/addons/fr/ctc/flow10"
	"github.com/invopop/gobl/addons/fr/ctc/flow2"
	"github.com/invopop/gobl/addons/fr/ctc/flow6"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
)

const (
	// Key identifies the French CTC meta-addon family.
	Key cbc.Key = "fr-ctc"

	// V1 is the first version of the French CTC meta-addon.
	V1 cbc.Key = Key + "-v1"
)

func init() {
	tax.RegisterAddonDef(newV1Addon())
}

func newV1Addon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "France CTC (auto-dispatch)",
			i18n.FR: "France CTC (dispatch automatique)",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Meta-addon for the French CTC (Continuous Transaction
				Control) reform. Inspects each GOBL document at
				normalize-time and appends the appropriate flow-specific
				addon — fr-ctc-flow2-v1 (domestic B2B clearance),
				fr-ctc-flow6-v1 (lifecycle status messages) or
				fr-ctc-flow10-v1 (B2C / cross-border B2B e-reporting and
				payment receipts) — so the right validation rules fire
				without the caller having to know which flow applies.

				Callers who prefer explicit control can declare any of
				the flow-specific addons directly.
			`),
			i18n.FR: here.Doc(`
				Méta-addon de la réforme française CTC. Inspecte chaque
				document GOBL au moment de la normalisation et ajoute
				automatiquement l'addon de flux approprié — fr-ctc-flow2-v1
				(clearance B2B domestique), fr-ctc-flow6-v1 (cycle de vie)
				ou fr-ctc-flow10-v1 (e-reporting B2C / B2B transfrontalier
				et encaissements) — afin que les règles de validation
				adaptées se déclenchent sans que l'appelant n'ait à
				connaître le flux applicable.
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
		Normalizer: normalize,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		obj.Addons.AddAddons(dispatchInvoice(obj))
	case *bill.Status:
		obj.Addons.AddAddons(flow6.V1)
	case *bill.Payment:
		obj.Addons.AddAddons(dispatchPayment(obj))
	}
}

// dispatchPayment picks the flow addon for a payment. A payment of
// type advice (CDAR 211) or receipt (CDAR 212) exchanged between two
// French parties is a Flow 6 CDV message. Everything else — payment
// requests, B2C settlements, cross-border B2B receipts — is e-reporting
// to the DGFiP via Flow 10.
func dispatchPayment(pmt *bill.Payment) cbc.Key {
	if pmt == nil {
		return flow10.V1
	}
	if pmt.Type != bill.PaymentTypeAdvice && pmt.Type != bill.PaymentTypeReceipt {
		return flow10.V1
	}
	if partyIsFrench(pmt.Supplier) && partyIsFrench(pmt.Customer) {
		return flow6.V1
	}
	return flow10.V1
}

// dispatchInvoice picks the flow addon for an invoice based on its
// parties: two French parties → Flow 2 clearance; otherwise Flow 10
// e-reporting (B2C if no customer, cross-border B2B if customer is
// non-French).
func dispatchInvoice(inv *bill.Invoice) cbc.Key {
	if partyIsFrench(inv.Supplier) && partyIsFrench(inv.Customer) {
		return flow2.V1
	}
	return flow10.V1
}

// partyIsFrench mirrors the helper of the same name in flow2/org.go.
// Duplicated here to keep the meta-addon's import graph trivial and
// the dispatcher easy to reason about: SIREN identity OR French TaxID
// counts as French.
func partyIsFrench(party *org.Party) bool {
	if party == nil {
		return false
	}
	if party.TaxID != nil && l10n.Code(party.TaxID.Country) == l10n.FR {
		return true
	}
	for _, id := range party.Identities {
		if id == nil {
			continue
		}
		if id.Type == fr.IdentityTypeSIREN {
			return true
		}
		if !id.Ext.IsZero() && id.Ext.Get(iso.ExtKeySchemeID).String() == "0002" {
			return true
		}
	}
	return false
}
