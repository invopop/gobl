package is_test

import (
	"errors"
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/stretchr/testify/assert"
)

func TestFunc(t *testing.T) {
	isPositive := func(v any) bool {
		n, ok := v.(int)
		return ok && n > 0
	}

	t.Run("passes when function returns true", func(t *testing.T) {
		assert.True(t, is.Func("positive", isPositive).Check(5))
	})

	t.Run("fails when function returns false", func(t *testing.T) {
		assert.False(t, is.Func("positive", isPositive).Check(-1))
	})

	t.Run("fails for unexpected type", func(t *testing.T) {
		assert.False(t, is.Func("positive", isPositive).Check("hello"))
	})

	t.Run("String returns description", func(t *testing.T) {
		assert.Equal(t, "positive", is.Func("positive", isPositive).String())
	})
}

func TestFuncError(t *testing.T) {
	validate := func(v any) error {
		s, ok := v.(string)
		if !ok || s == "" {
			return errors.New("must be a non-empty string")
		}
		return nil
	}

	t.Run("passes when function returns nil", func(t *testing.T) {
		assert.True(t, is.FuncError("non-empty string", validate).Check("hello"))
	})

	t.Run("fails when function returns error", func(t *testing.T) {
		assert.False(t, is.FuncError("non-empty string", validate).Check(""))
	})

	t.Run("fails for wrong type", func(t *testing.T) {
		assert.False(t, is.FuncError("non-empty string", validate).Check(42))
	})

	t.Run("String returns description", func(t *testing.T) {
		assert.Equal(t, "non-empty string", is.FuncError("non-empty string", validate).String())
	})
}

func TestFuncContext(t *testing.T) {
	fn := func(ctx rules.Context, _ any) bool {
		regime, _ := ctx.Value("regime").(string)
		return regime == "ES"
	}
	fc := is.FuncContext("regime is ES", fn)

	t.Run("String returns description", func(t *testing.T) {
		assert.Equal(t, "regime is ES", fc.String())
	})

	t.Run("CheckWithContext passes with matching context", func(t *testing.T) {
		rc := &rules.Context{}
		rc.Set("regime", "ES")
		assert.True(t, fc.CheckWithContext(rc, nil))
	})

	t.Run("CheckWithContext fails with non-matching context", func(t *testing.T) {
		rc := &rules.Context{}
		rc.Set("regime", "PT")
		assert.False(t, fc.CheckWithContext(rc, nil))
	})

	t.Run("Check fallback with empty context", func(t *testing.T) {
		assert.False(t, fc.Check(nil))
	})
}
