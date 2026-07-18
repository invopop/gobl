package pint_test

import (
	"testing"

	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/addons/peppol/pint"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

// Both en16931 (required by PINT) and the PINT addon are applied, mirroring how
// they resolve on a real document, so we can confirm PINT's normalizer runs
// after en16931's and wins.

func TestTaxComboNormalization(t *testing.T) {
	t.Run("standard GST rate maps to S (not O)", func(t *testing.T) {
		p := num.MakePercentage(10, 2)
		c := &tax.Combo{
			Category: tax.CategoryGST,
			Key:      tax.KeyStandard,
			Percent:  &p,
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017, pint.V1))
		assert.Equal(t, "S", c.Ext.Get(untdid.ExtKeyTaxCategory).String())
	})

	t.Run("zero-rated GST maps to Z", func(t *testing.T) {
		p := num.MakePercentage(0, 2)
		c := &tax.Combo{
			Category: tax.CategoryGST,
			Key:      tax.KeyZero,
			Percent:  &p,
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017, pint.V1))
		assert.Equal(t, "Z", c.Ext.Get(untdid.ExtKeyTaxCategory).String())
	})

	t.Run("exempt GST maps to E", func(t *testing.T) {
		c := &tax.Combo{
			Category: tax.CategoryGST,
			Key:      tax.KeyExempt,
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017, pint.V1))
		assert.Equal(t, "E", c.Ext.Get(untdid.ExtKeyTaxCategory).String())
	})

	t.Run("standard VAT still maps to S", func(t *testing.T) {
		p := num.MakePercentage(20, 2)
		c := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyStandard,
			Percent:  &p,
		}
		norm.Normalize(c, tax.AddonContext(en16931.V2017, pint.V1))
		assert.Equal(t, "S", c.Ext.Get(untdid.ExtKeyTaxCategory).String())
	})
}
