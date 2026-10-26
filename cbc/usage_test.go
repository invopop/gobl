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
		assert.ErrorContains(t, rules.Validate(u), "[GOBL-CBC-USAGE-02] ($.paths[0]) path must be a JSON path")
	})
	t.Run("index instead of wildcard", func(t *testing.T) {
		u := &cbc.Usage{
			Paths: []string{"$.lines[0].ext"},
		}
		assert.ErrorContains(t, rules.Validate(u), "[GOBL-CBC-USAGE-02]")
	})
}

func TestUsageSchemaValidation(t *testing.T) {
	t.Run("valid schemas", func(t *testing.T) {
		u := &cbc.Usage{
			Schemas: []string{"tax/combo", "org/document-ref", "bill/tax"},
		}
		assert.NoError(t, rules.Validate(u))
	})
	t.Run("invalid schemas", func(t *testing.T) {
		for _, s := range []string{"combo", "Tax/Combo", "https://gobl.org/draft-0/tax/combo", "tax/combo/ext"} {
			u := &cbc.Usage{Schemas: []string{s}}
			assert.ErrorContains(t, rules.Validate(u), "[GOBL-CBC-USAGE-01] ($.schemas[0]) schema must be a relative schema ID", s)
		}
	})
}

func TestUsageMatchesSchema(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var u *cbc.Usage
		assert.True(t, u.MatchesSchema("tax/combo"))
	})
	t.Run("no schemas", func(t *testing.T) {
		u := &cbc.Usage{Keys: []cbc.Key{"foo"}}
		assert.True(t, u.MatchesSchema("tax/combo"))
	})
	u := &cbc.Usage{Schemas: []string{"tax/combo", "org/party"}}
	assert.True(t, u.MatchesSchema("tax/combo"))
	assert.True(t, u.MatchesSchema("https://gobl.org/draft-0/org/party"))
	assert.False(t, u.MatchesSchema("org/item"))
	assert.False(t, u.MatchesSchema("https://gobl.org/draft-0/org/third-party"))
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

func TestUsageMatchesPathPrefix(t *testing.T) {
	u := &cbc.Usage{
		Paths: []string{
			"$.supplier",
			"$.lines[*].taxes[*]",
		},
	}
	tests := []struct {
		path string
		ok   bool
	}{
		{"$.supplier", true},
		{"$.supplier.ext", true},
		{"$.supplier.identities[0].ext", true},
		{"$.suppliers", false},
		{"$.customer.ext", false},
		{"$.lines[0].taxes[1]", true},
		{"$.lines[0].taxes[1].ext", true},
		{"$.lines[0].taxes", false},
		{"$.lines[0].ext", false},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			assert.Equal(t, tt.ok, u.MatchesPath(tt.path))
		})
	}
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
		{"$.tax.ext.foo", true},
		{"$.tax.extra", false},
		{"$.lines[0].ext", false},
		{"$.lines.taxes[0].ext", false},
		{"$.lines[].taxes[0].ext", false},
		{"$.lines[a].taxes[0].ext", false},
		{"$.lines[0].taxes[0].ext.foo", true},
		{"$.lines[0", false},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			assert.Equal(t, tt.ok, u.MatchesPath(tt.path))
		})
	}
}
