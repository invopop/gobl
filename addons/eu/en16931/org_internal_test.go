package en16931

import (
	"testing"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

// These white-box tests exercise the defensive guards in the item identity
// helpers, which the rule engine never reaches because it always passes the
// correctly-typed field value.

func TestItemIdentityHelperGuards(t *testing.T) {
	t.Run("max-one-legal ignores non-slice input", func(t *testing.T) {
		assert.True(t, itemHasMaxOneLegalIdentity("not a slice"))
	})
	t.Run("max-one-legal skips nil entries", func(t *testing.T) {
		ids := []*org.Identity{nil, {Scope: org.IdentityScopeLegal}}
		assert.True(t, itemHasMaxOneLegalIdentity(ids))
	})
	t.Run("legal-scheme ignores non-slice input", func(t *testing.T) {
		assert.True(t, itemLegalIdentitiesHaveScheme(42))
	})
	t.Run("legal-scheme skips nil entries", func(t *testing.T) {
		ids := []*org.Identity{
			nil,
			{Scope: org.IdentityScopeLegal, Ext: tax.ExtensionsOf(cbc.CodeMap{iso.ExtKeySchemeID: "0160"})},
		}
		assert.True(t, itemLegalIdentitiesHaveScheme(ids))
	})
}

func TestIdentityScopeIsGuards(t *testing.T) {
	test := identityScopeIs(org.IdentityScopeClass)
	t.Run("matches pointer", func(t *testing.T) {
		assert.True(t, test.Check(&org.Identity{Scope: org.IdentityScopeClass}))
	})
	t.Run("nil pointer does not match", func(t *testing.T) {
		var id *org.Identity
		assert.False(t, test.Check(id))
	})
	t.Run("matches value", func(t *testing.T) {
		assert.True(t, test.Check(org.Identity{Scope: org.IdentityScopeClass}))
	})
	t.Run("other types do not match", func(t *testing.T) {
		assert.False(t, test.Check("nope"))
	})
}
