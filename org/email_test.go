package org_test

import (
	"context"
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestEmailNormalize_TrimsFields(t *testing.T) {
	e := &org.Email{
		Label:   "  Work  ",
		Address: "  john.doe@example.com  ",
	}
	e.Normalize()

	assert.Equal(t, "Work", e.Label)
	assert.Equal(t, "john.doe@example.com", e.Address)
}

func TestEmailNormalize_NilReceiver(t *testing.T) {
	assert.NotPanics(t, func() {
		var e *org.Email
		e.Normalize()
	})
}

func TestEmailValidate(t *testing.T) {
	t.Run("valid email after normalize", func(t *testing.T) {
		e := &org.Email{
			Address: "  jane.doe@example.com  ",
		}
		e.Normalize()
		assert.NoError(t, e.Validate())
	})

	t.Run("valid email with context", func(t *testing.T) {
		e := &org.Email{
			Address: "john.doe@example.com",
		}
		assert.NoError(t, e.ValidateWithContext(context.Background()))
	})

	t.Run("empty address is invalid", func(t *testing.T) {
		e := &org.Email{}
		assert.Error(t, e.Validate())
	})

	t.Run("invalid format is rejected", func(t *testing.T) {
		e := &org.Email{Address: "not-an-email"}
		assert.Error(t, e.Validate())
	})

	t.Run("accepts uppercase", func(t *testing.T) {
		e := &org.Email{
			Address: "John.Doe+tag@Example.COM",
		}
		assert.NoError(t, e.Validate())
	})

	t.Run("invalid with whitespace after normalize", func(t *testing.T) {
		e := &org.Email{
			Address: "   ",
		}
		e.Normalize()
		assert.Error(t, e.Validate())
	})

	t.Run("invalid missing @", func(t *testing.T) {
		e := &org.Email{
			Address: "johndoe.example.com",
		}
		assert.Error(t, e.Validate())
	})

}
