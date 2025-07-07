package org_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestPersonValidation(t *testing.T) {
	t.Run("valid person", func(t *testing.T) {
		p := &org.Person{
			Label: "Attn",
			Name: &org.Name{
				Given:   "John",
				Surname: "Doe",
			},
			Identities: []*org.Identity{
				{Type: "SSN", Code: "123-45-6789"},
			},
			Addresses: []*org.Address{
				{Street: "123 Main St", Locality: "Anytown", Country: "US"},
			},
			Emails: []*org.Email{
				{
					Address: "foo@example.com",
				},
			},
		}
		assert.NoError(t, p.Validate())
	})
}
