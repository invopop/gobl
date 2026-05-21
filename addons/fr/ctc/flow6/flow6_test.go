package flow6

import (
	"testing"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// addonContext activates the fr-ctc-flow6 rule guard so the addon's
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

// quiet linter — keep iso import alive for fixtures defined in
// bill_status_test.go.
var _ = iso.ExtKeySchemeID
