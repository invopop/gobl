package cbc_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

func TestSourceValidation(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		src := &cbc.Source{
			Title: i18n.NewString("Test"),
			URL:   "http://example.com",
		}
		assert.NoError(t, rules.Validate(src))
	})
	t.Run("missing URL", func(t *testing.T) {
		src := &cbc.Source{}
		assert.ErrorContains(t, rules.Validate(src),
			"[GOBL-CBC-SOURCE-01] ($.url) url is required and must be a URL")
	})

	t.Run("invalid URL", func(t *testing.T) {
		src := &cbc.Source{
			Title: i18n.NewString("Test"),
			URL:   "http:\\example",
		}
		assert.ErrorContains(t, rules.Validate(src),
			"[GOBL-CBC-SOURCE-01] ($.url) url is required and must be a URL")
	})
}
