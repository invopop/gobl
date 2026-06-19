package en16931_test

import (
	"testing"

	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxComboNormalization(t *testing.T) {
	t.Run("standard VAT rate", func(t *testing.T) {
		p := num.MakePercentage(19, 2)
		c := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyStandard,
			Percent:  &p,
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		assert.Equal(t, "S", c.Ext.Get(untdid.ExtKeyTaxCategory).String())
		assert.Equal(t, "19%", c.Percent.String())
	})

	t.Run("unkown rate", func(t *testing.T) {
		c := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      cbc.Key("unknown"),
			Percent:  num.NewPercentage(19, 2),
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		assert.True(t, c.Ext.IsZero())
	})
	t.Run("IGIC", func(t *testing.T) {
		c := &tax.Combo{
			Category: es.TaxCategoryIGIC,
			Percent:  num.NewPercentage(7, 2),
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		assert.Equal(t, "L", c.Ext.Get(untdid.ExtKeyTaxCategory).String())
		assert.Equal(t, "7%", c.Percent.String())
	})

	t.Run("IPSI", func(t *testing.T) {
		c := &tax.Combo{
			Category: es.TaxCategoryIPSI,
			Percent:  num.NewPercentage(7, 2),
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		assert.Equal(t, "M", c.Ext.Get(untdid.ExtKeyTaxCategory).String())
		assert.Equal(t, "7%", c.Percent.String())
	})
	t.Run("exempt", func(t *testing.T) {
		c := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExempt,
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		assert.Equal(t, "E", c.Ext.Get(untdid.ExtKeyTaxCategory).String())
	})
	t.Run("missing rate, without percent", func(t *testing.T) {
		c := &tax.Combo{
			Category: tax.CategoryVAT,
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		// this will raise validation error later
		assert.Equal(t, "S", c.Ext.Get(untdid.ExtKeyTaxCategory).String())
	})

	t.Run("missing rate, with percent", func(t *testing.T) {
		c := &tax.Combo{
			Category: tax.CategoryVAT,
			Percent:  num.NewPercentage(19, 3),
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		assert.Equal(t, "S", c.Ext.Get(untdid.ExtKeyTaxCategory).String())
	})

	t.Run("missing rate, with zero percent", func(t *testing.T) {
		c := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyZero,
			Percent:  num.NewPercentage(0, 3),
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		assert.Equal(t, "Z", c.Ext.Get(untdid.ExtKeyTaxCategory).String())
	})

	t.Run("sales tax", func(t *testing.T) {
		c := &tax.Combo{
			Category: tax.CategoryGST,
			Percent:  num.NewPercentage(19, 2),
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		assert.Equal(t, "O", c.Ext.Get(untdid.ExtKeyTaxCategory).String())
		assert.Equal(t, "19%", c.Percent.String())
	})
}

func TestTaxComboValidation(t *testing.T) {
	t.Run("standard VAT rate", func(t *testing.T) {
		c := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyStandard,
			Percent:  num.NewPercentage(19, 2),
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		err := rules.Validate(c, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
		assert.Equal(t, "S", c.Ext.Get(untdid.ExtKeyTaxCategory).String())
		assert.Equal(t, "19%", c.Percent.String())
	})

	t.Run("exempt with vatex code", func(t *testing.T) {
		c := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExempt,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				"cef-vatex": "VATEX-EU-132",
			}),
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		assert.NoError(t, rules.Validate(c, tax.AddonContext(en16931.V2017)))
	})

	t.Run("exempt without vatex", func(t *testing.T) {
		// BR-E-10 is no longer enforced at the combo level (the invoice-level
		// exemption-note rule covers the remaining requirement), so an exempt
		// combo without a VATEX extension must validate.
		c := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExempt,
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		assert.NoError(t, rules.Validate(c, tax.AddonContext(en16931.V2017)))
	})

	t.Run("country-extension vatex code fails the EN16931 list", func(t *testing.T) {
		// The cef-vatex catalogue only enforces the code shape, so this
		// French CGI code passes the extension layer (it is NOT in the
		// official CEF list, unlike VATEX-FR-FRANCHISE/CNWVAT which the
		// EU list adopted); the EN16931 BR-CL-22 list rule rejects it.
		// Country profiles extend the list by suppressing this rule and
		// asserting their own.
		c := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExempt,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				"cef-vatex": "VATEX-FR-CGI261-1",
			}),
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		err := rules.Validate(c, tax.AddonContext(en16931.V2017))
		assert.ErrorContains(t, err, "VATEX code must belong to the CEF VATEX code list")
	})

	t.Run("malformed vatex code fails the catalogue pattern", func(t *testing.T) {
		c := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExempt,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				"cef-vatex": "EXEMPT-132",
			}),
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		err := rules.Validate(c, tax.AddonContext(en16931.V2017))
		assert.Error(t, err)
	})

	t.Run("reverse charge without vatex", func(t *testing.T) {
		c := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyReverseCharge,
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		err := rules.Validate(c, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
		assert.Equal(t, "AE", c.Ext.Get(untdid.ExtKeyTaxCategory).String())
		assert.Nil(t, c.Percent)
	})

	t.Run("intra-community", func(t *testing.T) {
		c := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyIntraCommunity,
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		err := rules.Validate(c, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
		assert.Equal(t, "K", c.Ext.Get(untdid.ExtKeyTaxCategory).String())
		assert.Nil(t, c.Percent)
	})

	t.Run("export", func(t *testing.T) {
		c := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExport,
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		err := rules.Validate(c, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
		assert.Equal(t, "G", c.Ext.Get(untdid.ExtKeyTaxCategory).String())
		assert.Nil(t, c.Percent)
	})
	t.Run("outside-scope", func(t *testing.T) {
		c := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyOutsideScope,
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		err := rules.Validate(c, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
		assert.Equal(t, "O", c.Ext.Get(untdid.ExtKeyTaxCategory).String())
		assert.Nil(t, c.Percent)
	})

	t.Run("VAT and IPSI mismatch", func(t *testing.T) {
		c := &tax.Combo{
			Category: tax.CategoryVAT,
			Percent:  num.NewPercentage(7, 2),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				untdid.ExtKeyTaxCategory: en16931.TaxCategoryIGIC,
			}),
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		err := rules.Validate(c, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
		assert.Equal(t, "S", c.Ext.Get(untdid.ExtKeyTaxCategory).String())
	})

	t.Run("zero with vatex code", func(t *testing.T) {
		c := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyZero,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				"cef-vatex": "VATEX-EU-132",
			}),
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		err := rules.Validate(c, tax.AddonContext(en16931.V2017))
		assert.ErrorContains(t, err, "VATEX extension must not be set")
	})

	t.Run("standard with vatex code", func(t *testing.T) {
		c := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyStandard,
			Percent:  num.NewPercentage(19, 2),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				"cef-vatex": "VATEX-EU-132",
			}),
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		err := rules.Validate(c, tax.AddonContext(en16931.V2017))
		assert.ErrorContains(t, err, "VATEX extension must not be set")
	})

	t.Run("nil", func(t *testing.T) {
		var tc *tax.Combo
		err := rules.Validate(tc, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})

	t.Run("missing rate", func(t *testing.T) {
		c := &tax.Combo{
			Category: tax.CategoryVAT,
			Percent:  num.NewPercentage(19, 2),
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017))
		c.Ext = tax.Extensions{} // override
		err := rules.Validate(c, tax.AddonContext(en16931.V2017))
		assert.ErrorContains(t, err, "tax category extension is required")
	})
}
