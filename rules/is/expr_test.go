package is_test

import (
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type exprPerson struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// exprCode is a named string type for testing Expr with value types.
type exprCode string

// exprInt is a named int type.
type exprInt int

// exprUint is a named uint type.
type exprUint uint

// exprFloat is a named float type.
type exprFloat float64

// exprBool is a named bool type.
type exprBool bool

func TestExpr(t *testing.T) {
	t.Run("struct-based expr passes", func(t *testing.T) {
		set := rules.For(new(exprPerson),
			rules.Assert("01", "name required",
				is.Expr(`name != ""`),
			),
		)
		faults := set.Validate(&exprPerson{Name: "Alice", Age: 30})
		assert.Nil(t, faults)
	})

	t.Run("struct-based expr fails", func(t *testing.T) {
		set := rules.For(new(exprPerson),
			rules.Assert("01", "name required",
				is.Expr(`name != ""`),
			),
		)
		faults := set.Validate(&exprPerson{Name: "", Age: 30})
		require.Error(t, faults)
	})

	t.Run("named string type via this", func(t *testing.T) {
		set := rules.For(exprCode(""),
			rules.Assert("01", "code required",
				is.Expr(`this != ""`),
			),
		)
		faults := set.Validate(exprCode("ABC"))
		assert.Nil(t, faults)

		faults = set.Validate(exprCode(""))
		require.Error(t, faults)
	})

	t.Run("named int type via this", func(t *testing.T) {
		set := rules.For(exprInt(0),
			rules.Assert("01", "must be positive",
				is.Expr(`this > 0`),
			),
		)
		faults := set.Validate(exprInt(5))
		assert.Nil(t, faults)

		faults = set.Validate(exprInt(0))
		require.Error(t, faults)
	})

	t.Run("named uint type via this", func(t *testing.T) {
		set := rules.For(exprUint(0),
			rules.Assert("01", "must be positive",
				is.Expr(`this > 0`),
			),
		)
		faults := set.Validate(exprUint(5))
		assert.Nil(t, faults)

		faults = set.Validate(exprUint(0))
		require.Error(t, faults)
	})

	t.Run("named float type via this", func(t *testing.T) {
		set := rules.For(exprFloat(0),
			rules.Assert("01", "must be positive",
				is.Expr(`this > 0`),
			),
		)
		faults := set.Validate(exprFloat(1.5))
		assert.Nil(t, faults)

		faults = set.Validate(exprFloat(0))
		require.Error(t, faults)
	})

	t.Run("named bool type via this", func(t *testing.T) {
		set := rules.For(exprBool(false),
			rules.Assert("01", "must be true",
				is.Expr(`this == true`),
			),
		)
		faults := set.Validate(exprBool(true))
		assert.Nil(t, faults)

		faults = set.Validate(exprBool(false))
		require.Error(t, faults)
	})

	t.Run("String returns expression text", func(t *testing.T) {
		e := is.Expr(`this > 5`)
		assert.Equal(t, "this > 5", e.String())
	})

	t.Run("format args", func(t *testing.T) {
		e := is.Expr("this > %d", 5)
		assert.Equal(t, "this > 5", e.String())
	})

	t.Run("uncompiled expr panics", func(t *testing.T) {
		e := is.Expr(`this > 0`)
		assert.Panics(t, func() {
			e.Check(42)
		})
	})

	t.Run("invalid expression panics at compile time", func(t *testing.T) {
		assert.Panics(t, func() {
			rules.For(new(exprPerson),
				rules.Assert("01", "bad", is.Expr(`=== invalid`)),
			)
		})
	})
}
