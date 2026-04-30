// Package flow6 handles the extensions, validations and normalization
// for the French CTC Flow 6 — CDV (cycle de vie) lifecycle statuses
// exchanged between PAs (plateformes agréées) for B2B invoices.
//
// The addon is standalone: it does not require fr-ctc-flow2-v1. It
// operates on bill.Status documents, carries the codebooks needed for
// the gobl.cii CDAR round-trip, and validates the subset of (key, type)
// / reason / action / role combinations that Flow 6 accepts.
package flow6

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
	// Key identifies the French CTC Flow 6 addon family.
	Key cbc.Key = "fr-ctc-flow6"

	// V1 is the key for the French CTC Flow 6 addon.
	V1 cbc.Key = Key + "-v1"
)

func init() {
	tax.RegisterAddonDef(newAddon())
	schema.Register(schema.GOBL.Add("addons/fr/ctc/flow6"),
		Characteristic{},
	)
	rules.RegisterWithGuard(
		Key.String(),
		rules.GOBL.Add("FR-CTC-FLOW6"),
		is.InContext(tax.AddonIn(V1)),
		billStatusRules(),
		billReasonRules(),
		billActionRules(),
		orgPartyRules(),
	)
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "France CTC Flow 6",
			i18n.FR: "France CTC Flux 6",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Support for the French CTC (Continuous Transaction Control)
				Flow 6 lifecycle messages (Cycle de Vie) exchanged between
				registered platforms (plateformes agréées) for B2B invoices.

				This addon operates on bill.Status documents. It carries the
				code tables (ProcessConditionCode, ReasonCode, RequestedAction,
				RoleCode) that the gobl.cii CDAR converter reads to round-trip
				to and from the French PPF XML, and validates the subset of
				(key, type) / reason / action / role combinations that Flow 6
				accepts.

				It does not depend on Flow 2: a platform may report lifecycle
				events for any compliant invoice, whether or not the invoice
				itself went through the Flow 2 B2B clearance path.
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
	case *bill.Reason:
		normalizeReason(obj)
	case *org.Party:
		// Party-level normalization handles the case where a party is
		// processed in isolation (e.g. through a direct tax.Normalize call);
		// the status-level normalizer applies the contextual role defaults
		// because those depend on which slot the party occupies.
		_ = obj
	}
}
