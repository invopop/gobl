package pt_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithholdingTaxCategories(t *testing.T) {
	r := pt.New()

	for _, code := range []cbc.Code{"IRS", "IRC"} {
		t.Run(code.String(), func(t *testing.T) {
			cat := r.CategoryDef(code)
			require.NotNil(t, cat, "category %s should be defined in the PT regime", code)
			assert.Equal(t, code, cat.Code)
			assert.True(t, cat.Retained, "category %s should be retained", code)
			assert.Empty(t, cat.Rates, "category %s should be rate-less", code)
		})
	}
}
