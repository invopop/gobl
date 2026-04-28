package choruspro

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
)

func normalizeOrgParty(party *org.Party) {
	if party == nil {
		return
	}

	if !party.Ext.IsZero() {
		if party.Ext.Get(ExtKeyScheme) != "" {
			return
		}
	}

	if party.TaxID != nil && party.TaxID.Country != "FR" {
		if party.Ext.IsZero() {
			party.Ext = tax.MakeExtensions()
		}
		if l10n.Unions().Code(l10n.EU).HasMember(l10n.Code(party.TaxID.Country)) {
			party.Ext = party.Ext.Merge(
				tax.ExtensionsOf(tax.ExtMap{
					ExtKeyScheme: "2",
				}),
			)
		} else {
			party.Ext = party.Ext.Merge(
				tax.ExtensionsOf(tax.ExtMap{
					ExtKeyScheme: "3",
				}),
			)
		}
		return
	}

	// If FR or no tax ID we search for a SIRET identity and set the scheme to 1
	for _, identity := range party.Identities {
		if identity.Type == fr.IdentityTypeSIRET {
			if party.Ext.IsZero() {
				party.Ext = tax.MakeExtensions()
			}
			party.Ext = party.Ext.Merge(
				tax.ExtensionsOf(tax.ExtMap{
					ExtKeyScheme: "1",
				}),
			)
			return
		}
	}
}
