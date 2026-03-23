package is_test

import (
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/stretchr/testify/assert"
)

func TestHasContext(t *testing.T) {
	inner := is.In("ES", "PT")

	t.Run("String includes inner test description", func(t *testing.T) {
		h := is.HasContext(inner)
		assert.Equal(t, "context: one of [ES, PT]", h.String())
	})

	t.Run("Check fallback passes when value itself matches", func(t *testing.T) {
		h := is.HasContext(inner)
		assert.True(t, h.Check("ES"))
	})

	t.Run("Check fallback fails when value does not match", func(t *testing.T) {
		h := is.HasContext(inner)
		assert.False(t, h.Check("FR"))
	})

	t.Run("CheckWithContext matching context value", func(t *testing.T) {
		h := is.HasContext(inner)
		rc := &rules.Context{}
		rc.Set("regime", "PT")
		ct := h.(rules.ContextualTest)
		assert.True(t, ct.CheckWithContext(rc, nil))
	})

	t.Run("CheckWithContext non-matching context", func(t *testing.T) {
		h := is.HasContext(inner)
		rc := &rules.Context{}
		rc.Set("regime", "FR")
		ct := h.(rules.ContextualTest)
		assert.False(t, ct.CheckWithContext(rc, nil))
	})

	t.Run("CheckWithContext empty context", func(t *testing.T) {
		h := is.HasContext(inner)
		rc := &rules.Context{}
		ct := h.(rules.ContextualTest)
		assert.False(t, ct.CheckWithContext(rc, nil))
	})
}
