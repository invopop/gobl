package gobl_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/schema"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	// basic error
	err := gobl.ErrNoDocument

	assert.Equal(t, "no-document", err.Error())
	assert.Equal(t, cbc.Key("no-document"), err.Key())
	assert.Equal(t, "", err.Message())
	assert.Nil(t, err.Faults())
	data, _ := json.Marshal(err)
	assert.JSONEq(t, `{"key":"no-document"}`, string(data))
	assert.True(t, err.Is(gobl.ErrNoDocument))

	se := errors.New("simple error message")
	err = gobl.ErrValidation.WithCause(se)
	assert.Equal(t, "validation: simple error message", err.Error())
	assert.Equal(t, "simple error message", err.Message())
	data, _ = json.Marshal(err)
	assert.JSONEq(t, `{"key":"validation","message":"simple error message"}`, string(data))

	err2 := err.WithReason("overwrite message")
	assert.Equal(t, "validation: simple error message", err.Error(), "do not modify original")
	assert.Equal(t, "validation: overwrite message", err2.Error())
	assert.Nil(t, err2.Faults())
	assert.Equal(t, "overwrite message", err2.Message())
	data, _ = json.Marshal(err2)
	assert.JSONEq(t, `{"key":"validation","message":"overwrite message"}`, string(data))

	err = gobl.ErrCalculation.WithCause(err2)
	assert.Equal(t, "validation: overwrite message", err.Error())

	fe := gobl.FieldErrors{
		"field": errors.New("field error"),
	}
	err = gobl.ErrValidation.WithCause(fe)
	assert.Equal(t, "validation: (field: field error.).", err.Error())
	data, _ = json.Marshal(err)
	assert.JSONEq(t, `{"key":"validation","fields":{"field":"field error"}}`, string(data))

	// check nested error with Is
	err = gobl.ErrValidation.WithCause(schema.ErrUnknownSchema)
	assert.True(t, errors.Is(err, schema.ErrUnknownSchema))
}
