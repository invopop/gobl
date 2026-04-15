package tax_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryDefValidations(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		c := baseCategoryDef()
		err := rules.Validate(c)
		require.NoError(t, err)
	})

	t.Run("informative", func(t *testing.T) {
		c := baseCategoryDef()
		c.Informative = true
		err := rules.Validate(c)
		require.NoError(t, err)
	})

	t.Run("retained", func(t *testing.T) {
		c := baseCategoryDef()
		c.Retained = true
		err := rules.Validate(c)
		require.NoError(t, err)
	})

	t.Run("informative and retained", func(t *testing.T) {
		c := baseCategoryDef()
		c.Informative = true
		c.Retained = true
		err := rules.Validate(c)
		assert.ErrorContains(t, err, "[GOBL-TAX-CATEGORYDEF-04] category def cannot be retained and informative")
	})

	t.Run("with valid extensions", func(t *testing.T) {
		c := baseCategoryDef()
		c.Extensions = []cbc.Key{pt.ExtKeyRegion}

		err := rules.Validate(c)
		require.NoError(t, err)
	})

	t.Run("with invalid extensions", func(t *testing.T) {
		c := baseCategoryDef()
		c.Extensions = []cbc.Key{"INVALID"}
		err := rules.Validate(c)
		assert.ErrorContains(t, err, "[GOBL-CBC-KEY-02] ($.extensions[0]) key must match the required pattern")
	})
}

func baseCategoryDef() *tax.CategoryDef {
	return &tax.CategoryDef{
		Code:  "TEST",
		Name:  i18n.NewString("TEST"),
		Title: i18n.NewString("Test tax"),
	}
}
