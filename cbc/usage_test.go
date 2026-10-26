package cbc_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

func TestUsageValidation(t *testing.T) {
	t.Run("valid paths", func(t *testing.T) {
		u := &cbc.Usage{
			Paths: []string{
				"$",
				"$.tax.ext",
				"$.lines[*].taxes[*].ext",
				"$.supplier.identities[*].ext",
				"$.$addons",
			},
		}
		assert.NoError(t, rules.Validate(u))
	})
	t.Run("invalid path", func(t *testing.T) {
		u := &cbc.Usage{
			Paths: []string{"tax.ext"},
		}
		assert.ErrorContains(t, rules.Validate(u), "[GOBL-CBC-USAGE-01] ($.paths[0]) path must be a JSON path")
	})
	t.Run("index instead of wildcard", func(t *testing.T) {
		u := &cbc.Usage{
			Paths: []string{"$.lines[0].ext"},
		}
		assert.ErrorContains(t, rules.Validate(u), "[GOBL-CBC-USAGE-01]")
	})
}

func TestUsageMatches(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var u *cbc.Usage
		assert.True(t, u.Matches(nil, nil))
		assert.True(t, u.Matches([]cbc.Key{"foo"}, []cbc.Code{"ES"}))
	})
	t.Run("empty", func(t *testing.T) {
		u := &cbc.Usage{}
		assert.True(t, u.Matches(nil, nil))
		assert.True(t, u.Matches([]cbc.Key{"foo"}, nil))
	})
	t.Run("keys", func(t *testing.T) {
		u := &cbc.Usage{Keys: []cbc.Key{"es-verifactu-v1"}}
		assert.True(t, u.Matches([]cbc.Key{"foo", "es-verifactu-v1"}, nil))
		assert.False(t, u.Matches([]cbc.Key{"foo"}, []cbc.Code{"ES"}))
		assert.False(t, u.Matches(nil, nil))
	})
	t.Run("codes", func(t *testing.T) {
		u := &cbc.Usage{Codes: []cbc.Code{"PT"}}
		assert.True(t, u.Matches(nil, []cbc.Code{"PT"}))
		assert.False(t, u.Matches([]cbc.Key{"foo"}, []cbc.Code{"ES"}))
	})
	t.Run("keys or codes", func(t *testing.T) {
		u := &cbc.Usage{
			Keys:  []cbc.Key{"es-tbai-v1"},
			Codes: []cbc.Code{"PT"},
		}
		assert.True(t, u.Matches([]cbc.Key{"es-tbai-v1"}, []cbc.Code{"ES"}))
		assert.True(t, u.Matches(nil, []cbc.Code{"PT"}))
		assert.False(t, u.Matches([]cbc.Key{"foo"}, []cbc.Code{"ES"}))
	})
}

func TestUsageMatchesPath(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var u *cbc.Usage
		assert.True(t, u.MatchesPath("$.tax.ext"))
	})
	t.Run("no paths", func(t *testing.T) {
		u := &cbc.Usage{Keys: []cbc.Key{"foo"}}
		assert.True(t, u.MatchesPath("$.tax.ext"))
	})
	u := &cbc.Usage{
		Paths: []string{
			"$.tax.ext",
			"$.lines[*].taxes[*].ext",
		},
	}
	tests := []struct {
		path string
		ok   bool
	}{
		{"$.tax.ext", true},
		{"$.lines[0].taxes[12].ext", true},
		{"$.lines[*].taxes[*].ext", true},
		{"$.lines[3].taxes[*].ext", true},
		{"$.tax", false},
		{"$.tax.ext.foo", false},
		{"$.lines[0].ext", false},
		{"$.lines.taxes[0].ext", false},
		{"$.lines[].taxes[0].ext", false},
		{"$.lines[a].taxes[0].ext", false},
		{"$.lines[0].taxes[0].ext.foo", false},
		{"$.lines[0", false},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			assert.Equal(t, tt.ok, u.MatchesPath(tt.path))
		})
	}
}
