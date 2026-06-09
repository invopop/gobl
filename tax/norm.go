package tax

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

func init() {
	norm.Register(
		// Global cleanup applied to every cbc.Code and tax.Extensions wherever
		// they appear in a document, so individual normalizers don't need to
		// trim codes or clean extensions by hand. NormalizeCode only trims
		// whitespace, so it is safe to apply before a normalizer that inspects
		// a raw code (e.g. org.Inbox).
		norm.For(func(c *cbc.Code) { *c = cbc.NormalizeCode(*c) }),
		norm.For(func(e *Extensions) { *e = e.Clean() }),
		norm.For(normalizeCombo),
		// Tax identities are an exception to the normal normalization rules:
		// they are normalized by their own country's regime (not the document
		// regime or any addon). The country-specific normalizers register
		// themselves with a tax.IdentityIn guard; here we only provide the
		// fallback used when the identity's country has no regime at all.
		norm.When(identityWithoutRegime,
			norm.For(func(id *Identity) { NormalizeIdentity(id) }),
		),
	)
}

// identityWithoutRegime passes for a tax identity whose country has no
// registered tax regime, mirroring the previous Identity.Normalize fallback.
var identityWithoutRegime rules.Test = is.Func("identity country has no regime",
	func(value any) bool {
		id, ok := value.(*Identity)
		return ok && id != nil && regimes.For(id.Country.Code()) == nil
	},
)
