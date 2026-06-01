package net

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterAuthority(t *testing.T) {
	original := Authorities
	t.Cleanup(func() { Authorities = original })
	Authorities = nil

	RegisterAuthority("kyc.example.com")
	RegisterAuthority("auth.example.org")

	assert.Equal(t, []Address{"kyc.example.com", "auth.example.org"}, Authorities)
}

func TestNewClientAuthoritiesIndependent(t *testing.T) {
	original := Authorities
	t.Cleanup(func() { Authorities = original })
	Authorities = []Address{"kyc.example.com"}

	c := NewClient()
	assert.Equal(t, []Address{"kyc.example.com"}, c.authorities)

	// Mutating the global after construction must not affect the client.
	Authorities = append(Authorities, "auth.example.org")
	assert.Equal(t, []Address{"kyc.example.com"}, c.authorities)
}

func TestWithAuthorities(t *testing.T) {
	original := Authorities
	t.Cleanup(func() { Authorities = original })
	Authorities = nil

	c := NewClient(WithAuthorities("kyc.example.com", "auth.example.org"))
	assert.Equal(t, []Address{"kyc.example.com", "auth.example.org"}, c.authorities)
}
