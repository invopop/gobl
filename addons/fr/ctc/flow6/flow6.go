// Package flow6 implements the French CTC Flow 6 ("Cycle de Vie")
// lifecycle status messages exchanged between registered platforms
// on a bill.Status document.
package flow6

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
	// Namespace is the rules namespace for the French CTC Flow 6 addon.
	Namespace rules.Code = "FR-CTC-FLOW6"

	// Key identifies the French CTC Flow 6 addon family.
	Key cbc.Key = "fr-ctc-flow6"

	// V1 is the first version of the French CTC Flow 6 addon.
	V1 cbc.Key = Key + "-v1"
)

func init() {
	tax.RegisterAddonDef(newV1Addon())
	rules.RegisterWithGuard(
		Key.String(),
		rules.GOBL.Add(Namespace),
		is.InContext(tax.AddonIn(V1)),
		billStatusRules(),
		billPaymentRules(),
		billReasonRules(),
		billActionRules(),
		orgPartyRules(),
	)
}

func newV1Addon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "France CTC Flow 6 (Cycle de Vie)",
			i18n.FR: "France CTC Flux 6 (Cycle de Vie)",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Support for Flow 6 ("Cycle de Vie") lifecycle status
				messages of the French CTC reform. Carries the CDAR
				ProcessConditionCode on bill.Status documents exchanged
				between registered dematerialisation platforms (PDPs)
				and the Portail Public de Facturation (PPF).
			`),
			i18n.FR: here.Doc(`
				Support du Flux 6 (« Cycle de Vie ») de la réforme
				française CTC. Porte le ProcessConditionCode CDAR sur
				les documents bill.Status échangés entre les plateformes
				de dématérialisation agréées (PDP) et le Portail Public
				de Facturation (PPF).
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
		Normalizer: normalize,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Status:
		normalizeStatus(obj)
	case *bill.Payment:
		normalizePayment(obj)
	case *bill.Reason:
		normalizeReason(obj)
	case *bill.Action:
		normalizeAction(obj)
	case *org.Party:
		normalizeParty(obj)
	case *org.Identity:
		normalizeIdentity(obj)
	}
}
