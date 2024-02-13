package gobl_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	// basic error
	err := gobl.ErrNoDocument

	assert.Equal(t, "no-document", err.Error())
	assert.Equal(t, cbc.Key("no-document"), err.Key())
	assert.Equal(t, "", err.Message())
	assert.Nil(t, err.Fields())
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
	assert.Nil(t, err2.Fields())
	assert.Equal(t, "overwrite message", err2.Message())
	data, _ = json.Marshal(err2)
	assert.JSONEq(t, `{"key":"validation","message":"overwrite message"}`, string(data))

	err = gobl.ErrCalculation.WithCause(err2)
	assert.Equal(t, "validation: overwrite message", err.Error())

	ve := validation.Errors{
		"field": errors.New("field error"),
	}
	err = gobl.ErrValidation.WithCause(ve)
	assert.Equal(t, "validation: (field: field error.).", err.Error())
	data, _ = json.Marshal(err)
	assert.JSONEq(t, `{"key":"validation","fields":{"field":"field error"}}`, string(data))

	// check nested error with Is
	err = gobl.ErrValidation.WithCause(schema.ErrUnknownSchema)
	assert.True(t, errors.Is(err, schema.ErrUnknownSchema))
}

func TestFieldErrors_Error(t *testing.T) {
	errs := gobl.FieldErrors{
		"B": errors.New("B1"),
		"C": errors.New("C1"),
		"A": errors.New("A1"),
	}
	assert.Equal(t, "A: A1; B: B1; C: C1.", errs.Error())

	errs = gobl.FieldErrors{
		"C": gobl.FieldErrors{
			"B": errors.New("B1"),
		},
		"A": errors.New("A1"),
	}
	assert.Equal(t, "A: A1; C: (B: B1.).", errs.Error())

	errs = gobl.FieldErrors{
		"B": errors.New("B1"),
	}
	assert.Equal(t, "B: B1.", errs.Error())

	errs = gobl.FieldErrors{}
	assert.Equal(t, "", errs.Error())
}

func TestFieldErrors_MarshalJSON(t *testing.T) {
	errs := gobl.FieldErrors{
		"A": errors.New("A1"),
		"B": gobl.FieldErrors{
			"2": errors.New("B1"),
		},
		"C": validation.Errors{
			"3": errors.New("C1"),
		},
	}
	data, err := errs.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, `{"A":"A1","B":{"2":"B1"},"C":{"3":"C1"}}`, string(data))
}
