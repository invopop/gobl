package mydata_test

import (
	"testing"

	"github.com/invopop/gobl/addons/gr/mydata"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeCharge(t *testing.T) {
	ad := tax.AddonForKey(mydata.V1)

	t.Run("nil charge", func(t *testing.T) {
		var c *bill.Charge
		ad.Normalizer(c)
		assert.Nil(t, c)
	})

	t.Run("with fee extension", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(1000, 2),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				mydata.ExtKeyFee: "13",
			}),
		}
		ad.Normalizer(c)
		assert.Equal(t, mydata.TaxTypeFee, c.Ext.Get(mydata.ExtKeyTaxType).String())
		assert.Equal(t, "13", c.Ext.Get(mydata.ExtKeyFee).String())
	})

	t.Run("with other tax extension", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(300, 2),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				mydata.ExtKeyOtherTax: "8",
			}),
		}
		ad.Normalizer(c)
		assert.Equal(t, mydata.TaxTypeOtherTax, c.Ext.Get(mydata.ExtKeyTaxType).String())
		assert.Equal(t, "8", c.Ext.Get(mydata.ExtKeyOtherTax).String())
	})

	t.Run("with stamp duty extension", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(1200, 2),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				mydata.ExtKeyStampDuty: "1",
			}),
		}
		ad.Normalizer(c)
		assert.Equal(t, mydata.TaxTypeStampDuty, c.Ext.Get(mydata.ExtKeyTaxType).String())
		assert.Equal(t, "1", c.Ext.Get(mydata.ExtKeyStampDuty).String())
	})

	t.Run("without any tax extension", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(1000, 2),
			Reason: "Some charge",
		}
		ad.Normalizer(c)
		assert.True(t, c.Ext.IsZero())
	})

	t.Run("with type extension set", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(1000, 2),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				mydata.ExtKeyTaxType: mydata.TaxTypeOtherTax,
				mydata.ExtKeyFee:     "13",
			}),
		}
		ad.Normalizer(c)
		assert.Equal(t, mydata.TaxTypeOtherTax, c.Ext.Get(mydata.ExtKeyTaxType).String())
	})

	t.Run("with stamp duty charge key and existing extension", func(t *testing.T) {
		c := &bill.Charge{
			Key:    bill.ChargeKeyStampDuty,
			Amount: num.MakeAmount(1200, 2),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				mydata.ExtKeyFee: "13",
			}),
		}
		ad.Normalizer(c)
		assert.Equal(t, mydata.TaxTypeStampDuty, c.Ext.Get(mydata.ExtKeyTaxType).String())
		assert.Equal(t, "13", c.Ext.Get(mydata.ExtKeyFee).String())
	})

	t.Run("with tax charge key and existing extension", func(t *testing.T) {
		c := &bill.Charge{
			Key:    bill.ChargeKeyTax,
			Amount: num.MakeAmount(500, 2),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				mydata.ExtKeyOtherTax: "8",
			}),
		}
		ad.Normalizer(c)
		assert.Equal(t, mydata.TaxTypeOtherTax, c.Ext.Get(mydata.ExtKeyTaxType).String())
		assert.Equal(t, "8", c.Ext.Get(mydata.ExtKeyOtherTax).String())
	})

	t.Run("charge key overrides existing tax type", func(t *testing.T) {
		c := &bill.Charge{
			Key:    bill.ChargeKeyStampDuty,
			Amount: num.MakeAmount(1200, 2),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				mydata.ExtKeyTaxType: mydata.TaxTypeFee,
			}),
		}
		ad.Normalizer(c)
		assert.Equal(t, mydata.TaxTypeStampDuty, c.Ext.Get(mydata.ExtKeyTaxType).String())
	})
}

func TestValidateCharge(t *testing.T) {
	t.Run("nil charge", func(t *testing.T) {
		var c *bill.Charge
		err := rules.Validate(c, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("valid fee", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(1000, 2),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				mydata.ExtKeyTaxType: "2",
				mydata.ExtKeyFee:     "13",
			}),
		}
		err := rules.Validate(c, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("valid other tax", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(300, 2),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				mydata.ExtKeyTaxType:  "3",
				mydata.ExtKeyOtherTax: "8",
			}),
		}
		err := rules.Validate(c, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("valid stamp duty", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(1200, 2),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				mydata.ExtKeyTaxType:   "4",
				mydata.ExtKeyStampDuty: "1",
			}),
		}
		err := rules.Validate(c, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("missing fee extension for fee type", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(1000, 2),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				mydata.ExtKeyTaxType: mydata.TaxTypeFee,
			}),
		}
		err := rules.Validate(c, withAddonContext())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "charge with fee tax type requires 'gr-mydata-fee' extension")
	})

	t.Run("missing other tax extension for other tax type", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(300, 2),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				mydata.ExtKeyTaxType: mydata.TaxTypeOtherTax,
			}),
		}
		err := rules.Validate(c, withAddonContext())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "charge with other-tax type requires 'gr-mydata-other-tax' extension")
	})

	t.Run("missing stamp duty extension for stamp duty type", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(1200, 2),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				mydata.ExtKeyTaxType: mydata.TaxTypeStampDuty,
			}),
		}
		err := rules.Validate(c, withAddonContext())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "charge with stamp-duty type requires 'gr-mydata-stamp-duty' extension")
	})

	t.Run("multiple specific extensions", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(1500, 2),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				mydata.ExtKeyTaxType:  mydata.TaxTypeFee,
				mydata.ExtKeyFee:      "13",
				mydata.ExtKeyOtherTax: "8",
			}),
		}
		err := rules.Validate(c, withAddonContext())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "only one of fee, other-tax, or stamp-duty allowed")
		// The "only one" rule covers all three extensions together
	})
}

func withAddonContext() rules.WithContext {
	return func(rc *rules.Context) {
		rc.Set(rules.ContextKey(mydata.V1), tax.AddonForKey(mydata.V1))
	}
}
