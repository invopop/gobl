package rules_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFaults(t *testing.T) {
	t.Run("empty faults should not return error", func(t *testing.T) {
		var faults rules.Faults
		assert.NoError(t, faults)
	})
}

func TestFaultMarshalJSON(t *testing.T) {
	t.Run("field-level fault", func(t *testing.T) {
		set := rules.For(new(Email),
			rules.Field("addr",
				rules.Assert("01", "email required", is.Present),
			),
		)
		faults := set.Validate(&Email{Addr: ""})
		require.Error(t, faults)
		f := faults.First()
		data, err := json.Marshal(f)
		require.NoError(t, err)
		var m map[string]string
		require.NoError(t, json.Unmarshal(data, &m))
		assert.Equal(t, "$.addr", m["path"])
		assert.Contains(t, m["code"], "01")
		assert.Equal(t, "email required", m["message"])
	})

	t.Run("root-level fault has $ path", func(t *testing.T) {
		set := rules.For(new(Email),
			rules.Assert("01", "always fails", is.Expr(`false`)),
		)
		faults := set.Validate(&Email{Addr: "ok@test.com"})
		require.Error(t, faults)
		f := faults.First()
		data, err := json.Marshal(f)
		require.NoError(t, err)
		var m map[string]string
		require.NoError(t, json.Unmarshal(data, &m))
		assert.Equal(t, "$", m["path"])
	})
}

func TestFaultListMarshalJSON(t *testing.T) {
	set := rules.For(new(Person),
		rules.Field("name",
			rules.Assert("01", "name required", is.Present),
		),
		rules.Assert("02", "always fails", is.Expr(`false`)),
	)
	faults := set.Validate(&Person{Name: ""})
	require.Error(t, faults)
	data, err := json.Marshal(faults)
	require.NoError(t, err)
	var arr []map[string]string
	require.NoError(t, json.Unmarshal(data, &arr))
	assert.GreaterOrEqual(t, len(arr), 2)
}

func TestFaultListHasCode(t *testing.T) {
	set := rules.For(new(Email),
		rules.Field("addr",
			rules.Assert("01", "email required", is.Present),
		),
	)
	faults := set.Validate(&Email{Addr: ""})
	require.Error(t, faults)

	t.Run("matching code", func(t *testing.T) {
		assert.True(t, faults.HasCode(faults.First().Code()))
	})

	t.Run("non-matching code", func(t *testing.T) {
		assert.False(t, faults.HasCode("NONEXISTENT"))
	})
}

func TestFaultListLast(t *testing.T) {
	set := rules.For(new(Person),
		rules.Field("name",
			rules.Assert("01", "name required", is.Present),
		),
		rules.Assert("02", "always fails", is.Expr(`false`)),
	)
	faults := set.Validate(&Person{Name: ""})
	require.Error(t, faults)
	require.GreaterOrEqual(t, faults.Len(), 2)
	last := faults.Last()
	assert.Equal(t, faults.At(faults.Len()-1), last)
}

func TestFaultError(t *testing.T) {
	t.Run("single fault with path", func(t *testing.T) {
		set := rules.For(new(Email),
			rules.Field("addr",
				rules.Assert("01", "email required", is.Present),
			),
		)
		faults := set.Validate(&Email{Addr: ""})
		require.Error(t, faults)
		assert.Contains(t, faults.Error(), "($.addr)")
		assert.Contains(t, faults.Error(), "email required")
	})

	t.Run("single fault without path (root level)", func(t *testing.T) {
		set := rules.For(new(Email),
			rules.Assert("01", "always fails", is.Expr(`false`)),
		)
		faults := set.Validate(&Email{Addr: "ok@test.com"})
		require.Error(t, faults)
		assert.Contains(t, faults.Error(), "always fails")
		// Root-level faults have no path prefix in Error()
		assert.NotContains(t, faults.Error(), "($.)")
	})

	t.Run("multi-fault semicolon-joined", func(t *testing.T) {
		set := rules.For(new(Person),
			rules.Field("name",
				rules.Assert("01", "name required", is.Present),
			),
			rules.Assert("02", "always fails", is.Expr(`false`)),
		)
		faults := set.Validate(&Person{Name: ""})
		require.Error(t, faults)
		assert.Contains(t, faults.Error(), "; ")
	})
}
