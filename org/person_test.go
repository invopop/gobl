package org_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestPersonNormalize(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var p *org.Person
		assert.NotPanics(t, func() {
			p.Normalize(nil)
		})
	})
	t.Run("complete", func(t *testing.T) {
		p := &org.Person{
			Label: " Attn ",
			Name: &org.Name{
				Given:   " John ",
				Surname: " Doe ",
			},
			Identities: []*org.Identity{
				{Type: "SSN", Code: " 123-45-6789 "},
			},
			Addresses: []*org.Address{
				{Street: " 123 Main St ", Locality: "Anytown", Country: "US"},
			},
			Emails: []*org.Email{
				{
					Address: " foo@example.com ",
				},
			},
		}
		p.Normalize(nil)
		assert.Equal(t, "Attn", p.Label)
		assert.Equal(t, "John", p.Name.Given)
		assert.Equal(t, "Doe", p.Name.Surname)
		assert.Equal(t, "foo@example.com", p.Emails[0].Address)
		assert.Equal(t, "123-45-6789", p.Identities[0].Code.String())
		assert.Equal(t, "123 Main St", p.Addresses[0].Street)
	})
}

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
