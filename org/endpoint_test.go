package org_test

import (
	"testing"

	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEndpointValidation(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		e := &org.Endpoint{Label: "GOBL Net", URI: "gobl:acme.example.com"}
		assert.NoError(t, rules.Validate(e))
	})

	t.Run("uri required", func(t *testing.T) {
		e := &org.Endpoint{Label: "no uri"}
		err := rules.Validate(e)
		assert.ErrorContains(t, err, "uri")
	})

	t.Run("invalid uri", func(t *testing.T) {
		e := &org.Endpoint{URI: "no scheme"}
		err := rules.Validate(e)
		require.Error(t, err)
	})
}

func TestEndpointNormalize(t *testing.T) {
	e := &org.Endpoint{Label: "  spaced  ", URI: "gobl:acme.example.com"}
	norm.Normalize(e)
	assert.Equal(t, "spaced", e.Label)
}

func TestEndpointNormalizeNil(t *testing.T) {
	// nil pointer must not panic.
	var e *org.Endpoint
	assert.NotPanics(t, func() { norm.Normalize(e) })
}

func TestPartyFirstEndpoint(t *testing.T) {
	t.Run("nil party", func(t *testing.T) {
		var p *org.Party
		assert.Nil(t, p.FirstEndpoint())
	})

	t.Run("no endpoints", func(t *testing.T) {
		p := &org.Party{Name: "Acme"}
		assert.Nil(t, p.FirstEndpoint())
	})

	t.Run("returns first entry", func(t *testing.T) {
		p := &org.Party{
			Endpoints: []*org.Endpoint{
				{URI: "gobl:acme.example.com"},
				{URI: "iso6523-actorid-upis::9920:x3157928m"},
			},
		}
		require.NotNil(t, p.FirstEndpoint())
		assert.Equal(t, "gobl:acme.example.com", p.FirstEndpoint().URI.String())
	})

	t.Run("skips nil entries", func(t *testing.T) {
		p := &org.Party{
			Endpoints: []*org.Endpoint{
				nil,
				{URI: "gobl:acme.example.com"},
			},
		}
		require.NotNil(t, p.FirstEndpoint())
		assert.Equal(t, "gobl:acme.example.com", p.FirstEndpoint().URI.String())
	})
}

func TestPartyEndpointLookup(t *testing.T) {
	p := &org.Party{
		Endpoints: []*org.Endpoint{
			{URI: "iso6523-actorid-upis::9920:x3157928m"},
			{URI: "gobl:acme.example.com"},
			{URI: "mailto:billing@example.com"},
		},
	}
	assert.Equal(t, "gobl:acme.example.com", p.Endpoint("gobl").URI.String())
	assert.Equal(t, "iso6523-actorid-upis::9920:x3157928m", p.Endpoint("iso6523-actorid-upis").URI.String())
	assert.Equal(t, "mailto:billing@example.com", p.Endpoint("mailto").URI.String())
	assert.Nil(t, p.Endpoint("ftp"))
}
