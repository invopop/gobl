package org_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/invopop/gobl/org"
)

func TestEmailValidation(t *testing.T) {
	valid := org.Email{
		Address: "foobar@invopop.example.com",
	}
	assert.NoError(t, valid.Validate())

	invalid := org.Email{
		Address: "foobar",
	}
	assert.EqualError(t, invalid.Validate(), "addr: must be a valid email address.")
}
