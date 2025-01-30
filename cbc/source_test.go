package cbc_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/stretchr/testify/assert"
)

func TestSourceValidation(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		src := &cbc.Source{
			Title: i18n.NewString("Test"),
			URL:   "http://example.com",
		}
		assert.NoError(t, src.Validate())
	})
	t.Run("missing URL", func(t *testing.T) {
		src := &cbc.Source{}
		assert.ErrorContains(t, src.Validate(), "url: cannot be blank.")
	})

	t.Run("invalid URL", func(t *testing.T) {
		src := &cbc.Source{
			Title: i18n.NewString("Test"),
			URL:   "http:\\example",
		}
		assert.ErrorContains(t, src.Validate(), "url: must be a valid URL.")
	})
}
