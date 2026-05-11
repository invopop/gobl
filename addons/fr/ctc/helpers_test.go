package ctc

import (
	"testing"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// addonContext activates the fr-ctc rule guard so the addon's
// validators fire on standalone objects (bill.Reason / org.Party /
// org.Identity) that do not themselves carry the addon.
func addonContext() rules.WithContext {
	return func(rc *rules.Context) {
		rc.Set(rules.ContextKey(V1), tax.AddonForKey(V1))
	}
}

// runNormalize invokes the addon's registered normalizer on the given
// object, matching what tax.Normalize would do during Calculate.
func runNormalize(t *testing.T, doc any) {
	t.Helper()
	tax.Normalize([]tax.Normalizer{tax.AddonForKey(V1).Normalizer}, doc)
}

// frPartyWithSIREN returns a French supplier party with a SIREN
// identity, suitable for invoice-supplier and payment-supplier slots.
func frPartyWithSIREN() *org.Party {
	return &org.Party{
		Name: "Supplier SARL",
		TaxID: &tax.Identity{
			Country: "FR",
			Code:    "39356000000",
		},
		Identities: []*org.Identity{
			{
				Type:  fr.IdentityTypeSIREN,
				Code:  "356000000",
				Scope: org.IdentityScopeLegal,
				Ext: tax.ExtensionsOf(tax.ExtMap{
					iso.ExtKeySchemeID: identitySchemeIDSIREN,
				}),
			},
		},
		Addresses: []*org.Address{{Country: "FR"}},
	}
}

// frCustomerWithSIREN returns a French customer party with a SIREN
// identity. Used to trigger the Flow 2 dispatcher (both parties
// resolve as French).
func frCustomerWithSIREN() *org.Party {
	return &org.Party{
		Name: "Customer SAS",
		TaxID: &tax.Identity{
			Country: "FR",
			Code:    "44732829320",
		},
		Identities: []*org.Identity{
			{
				Type:  fr.IdentityTypeSIREN,
				Code:  "732829320",
				Scope: org.IdentityScopeLegal,
				Ext: tax.ExtensionsOf(tax.ExtMap{
					iso.ExtKeySchemeID: identitySchemeIDSIREN,
				}),
			},
		},
		Addresses: []*org.Address{{Country: "FR"}},
	}
}

// deCustomerWithVATID returns a German customer party with an EU-VAT
// identity (ICD scheme 0223). Used to trigger the Flow 10 dispatcher
// branch (at least one party is non-French).
func deCustomerWithVATID() *org.Party {
	return &org.Party{
		Name: "Kunde Deutschland GmbH",
		TaxID: &tax.Identity{
			Country: "DE",
			Code:    "111111125",
		},
		Identities: []*org.Identity{
			{
				Code:  "DE111111125",
				Scope: org.IdentityScopeLegal,
				Ext: tax.ExtensionsOf(tax.ExtMap{
					iso.ExtKeySchemeID: identitySchemeIDEUVAT,
				}),
			},
		},
		Addresses: []*org.Address{{Country: "DE"}},
	}
}
