package gobl_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testDoc is a minimal struct used to generate rule faults in error tests.
type testDoc struct {
	Name string `json:"name"`
}

var testDocRules = rules.For(new(testDoc),
	rules.Field("name",
		rules.Assert("01", "name required", is.Present),
	),
)

func TestError(t *testing.T) {
	t.Run("basic error with no cause", func(t *testing.T) {
		err := gobl.ErrNoDocument
		assert.Equal(t, "no-document", err.Error())
		assert.Equal(t, cbc.Key("no-document"), err.Key())
		assert.Equal(t, "", err.Message())
		assert.Nil(t, err.Faults())
		data, _ := json.Marshal(err)
		assert.JSONEq(t, `{"key":"no-document"}`, string(data))
		assert.True(t, err.Is(gobl.ErrNoDocument))
	})

	t.Run("with cause from simple error", func(t *testing.T) {
		err := gobl.ErrValidation.WithCause(errors.New("simple error message"))
		assert.Equal(t, "validation: simple error message", err.Error())
		assert.Equal(t, "simple error message", err.Message())
		assert.Nil(t, err.Faults())
		data, _ := json.Marshal(err)
		assert.JSONEq(t, `{"key":"validation","message":"simple error message"}`, string(data))
	})

	t.Run("with reason overwrites message without modifying original", func(t *testing.T) {
		err := gobl.ErrValidation.WithCause(errors.New("original message"))
		err2 := err.WithReason("overwrite message")
		assert.Equal(t, "validation: original message", err.Error(), "original unchanged")
		assert.Equal(t, "validation: overwrite message", err2.Error())
		assert.Equal(t, "overwrite message", err2.Message())
		assert.Nil(t, err2.Faults())
		data, _ := json.Marshal(err2)
		assert.JSONEq(t, `{"key":"validation","message":"overwrite message"}`, string(data))
	})

	t.Run("with cause from *Error returns inner error unchanged", func(t *testing.T) {
		inner := gobl.ErrValidation.WithReason("inner message")
		err := gobl.ErrCalculation.WithCause(inner)
		assert.Equal(t, "validation: inner message", err.Error())
	})

	t.Run("with cause from rules.Faults", func(t *testing.T) {
		faults := testDocRules.Validate(&testDoc{Name: ""})
		require.Error(t, faults)

		err := gobl.ErrValidation.WithCause(faults)
		require.NotNil(t, err.Faults())
		assert.Equal(t, 1, err.Faults().Len())
		assert.Equal(t, "$.name", err.Faults().First().Path())

		data, _ := json.Marshal(err)
		var out struct {
			Key    string           `json:"key"`
			Faults []map[string]any `json:"faults"`
		}
		require.NoError(t, json.Unmarshal(data, &out))
		assert.Equal(t, "validation", out.Key)
		assert.Len(t, out.Faults, 1)
		assert.Equal(t, "$.name", out.Faults[0]["path"])
	})

	t.Run("Is checks nested cause chain", func(t *testing.T) {
		err := gobl.ErrValidation.WithCause(schema.ErrUnknownSchema)
		assert.True(t, errors.Is(err, schema.ErrUnknownSchema))
	})
}
