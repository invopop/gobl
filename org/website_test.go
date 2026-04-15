package org_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

func TestWebsiteNormalize(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var w *org.Website
		assert.NotPanics(t, func() {
			w.Normalize()
		})
	})
	w := &org.Website{
		Label: "  My Site\t",
		Title: "\nMain  Site  ",
		URL:   "   https://example.com/path?q=1   ",
	}

	w.Normalize()

	assert.Equal(t, "My Site", w.Label)
	assert.Equal(t, "Main  Site", w.Title)
	assert.Equal(t, "https://example.com/path?q=1", w.URL)
}

func TestWebsiteValidate(t *testing.T) {
	tests := []struct {
		name string
		url  string
		ok   bool
	}{
		{name: "valid https", url: "https://example.org", ok: true},
		{name: "valid http with path and query", url: "http://example.org/path?x=1", ok: true},
		{name: "empty", url: "", ok: false},
		{name: "no scheme", url: "www.example.org", ok: true},
		{name: "invalid", url: "not a url", ok: false},
		{name: "whitespace needs normalize", url: "   https://example.org  ", ok: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := &org.Website{
				URL: tc.url,
			}

			// Ensure Normalize helps when there is surrounding whitespace.
			w.Normalize()

			err := rules.Validate(w)

			if tc.ok {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
