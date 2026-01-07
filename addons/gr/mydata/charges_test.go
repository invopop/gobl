package mydata_test

import (
	"testing"

	"github.com/invopop/gobl/addons/gr/mydata"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
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
			Ext: tax.Extensions{
				mydata.ExtKeyFee: "13",
			},
		}
		ad.Normalizer(c)
		assert.Equal(t, mydata.TaxTypeFee, c.Ext[mydata.ExtKeyTaxType].String())
		assert.Equal(t, "13", c.Ext[mydata.ExtKeyFee].String())
	})

	t.Run("with other tax extension", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(300, 2),
			Ext: tax.Extensions{
				mydata.ExtKeyOtherTax: "8",
			},
		}
		ad.Normalizer(c)
		assert.Equal(t, mydata.TaxTypeOtherTax, c.Ext[mydata.ExtKeyTaxType].String())
		assert.Equal(t, "8", c.Ext[mydata.ExtKeyOtherTax].String())
	})

	t.Run("with stamp duty extension", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(1200, 2),
			Ext: tax.Extensions{
				mydata.ExtKeyStampDuty: "1",
			},
		}
		ad.Normalizer(c)
		assert.Equal(t, mydata.TaxTypeStampDuty, c.Ext[mydata.ExtKeyTaxType].String())
		assert.Equal(t, "1", c.Ext[mydata.ExtKeyStampDuty].String())
	})

	t.Run("without any tax extension", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(1000, 2),
			Reason: "Some charge",
		}
		ad.Normalizer(c)
		assert.Nil(t, c.Ext)
	})

	t.Run("with type extension set", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(1000, 2),
			Ext: tax.Extensions{
				mydata.ExtKeyTaxType: mydata.TaxTypeOtherTax,
				mydata.ExtKeyFee:     "13",
			},
		}
		ad.Normalizer(c)
		assert.Equal(t, mydata.TaxTypeOtherTax, c.Ext[mydata.ExtKeyTaxType].String())
	})
}

func TestValidateCharge(t *testing.T) {
	ad := tax.AddonForKey(mydata.V1)

	t.Run("nil charge", func(t *testing.T) {
		var c *bill.Charge
		err := ad.Validator(c)
		assert.NoError(t, err)
	})

	t.Run("valid fee", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(1000, 2),
			Ext: tax.Extensions{
				mydata.ExtKeyTaxType: "2",
				mydata.ExtKeyFee:     "13",
			},
		}
		err := ad.Validator(c)
		assert.NoError(t, err)
	})

	t.Run("valid other tax", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(300, 2),
			Ext: tax.Extensions{
				mydata.ExtKeyTaxType:  "3",
				mydata.ExtKeyOtherTax: "8",
			},
		}
		err := ad.Validator(c)
		assert.NoError(t, err)
	})

	t.Run("valid stamp duty", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(1200, 2),
			Ext: tax.Extensions{
				mydata.ExtKeyTaxType:   "4",
				mydata.ExtKeyStampDuty: "1",
			},
		}
		err := ad.Validator(c)
		assert.NoError(t, err)
	})

	t.Run("missing fee extension for fee type", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(1000, 2),
			Ext: tax.Extensions{
				mydata.ExtKeyTaxType: mydata.TaxTypeFee,
			},
		}
		err := ad.Validator(c)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "gr-mydata-fee: required")
	})

	t.Run("missing other tax extension for other tax type", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(300, 2),
			Ext: tax.Extensions{
				mydata.ExtKeyTaxType: mydata.TaxTypeOtherTax,
			},
		}
		err := ad.Validator(c)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "gr-mydata-other-tax: required")
	})

	t.Run("missing stamp duty extension for stamp duty type", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(1200, 2),
			Ext: tax.Extensions{
				mydata.ExtKeyTaxType: mydata.TaxTypeStampDuty,
			},
		}
		err := ad.Validator(c)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "gr-mydata-stamp-duty: required")
	})

	t.Run("multiple specific extensions", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(1500, 2),
			Ext: tax.Extensions{
				mydata.ExtKeyTaxType:  mydata.TaxTypeFee,
				mydata.ExtKeyFee:      "13",
				mydata.ExtKeyOtherTax: "8",
			},
		}
		err := ad.Validator(c)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "gr-mydata-other-tax: only one allowed")
		assert.Contains(t, err.Error(), "gr-mydata-fee: only one allowed")
	})
}

func TestChargeNormalizationAndValidation(t *testing.T) {
	ad := tax.AddonForKey(mydata.V1)

	t.Run("normalize then validate - withholding", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(1500, 2),
			Ext: tax.Extensions{
				mydata.ExtKeyFee: "5",
			},
		}
		ad.Normalizer(c)
		err := ad.Validator(c)
		assert.NoError(t, err)
		assert.Equal(t, mydata.TaxTypeFee, c.Ext[mydata.ExtKeyTaxType].String())
	})

	t.Run("normalize then validate - fee", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(1000, 2),
			Ext: tax.Extensions{
				mydata.ExtKeyFee: "13",
			},
		}
		ad.Normalizer(c)
		err := ad.Validator(c)
		assert.NoError(t, err)
		assert.Equal(t, "2", c.Ext[mydata.ExtKeyTaxType].String())
	})

	t.Run("normalize then validate - other tax", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(300, 2),
			Ext: tax.Extensions{
				mydata.ExtKeyOtherTax: "8",
			},
		}
		ad.Normalizer(c)
		err := ad.Validator(c)
		assert.NoError(t, err)
		assert.Equal(t, "3", c.Ext[mydata.ExtKeyTaxType].String())
	})

	t.Run("normalize then validate - stamp duty", func(t *testing.T) {
		c := &bill.Charge{
			Amount: num.MakeAmount(1200, 2),
			Ext: tax.Extensions{
				mydata.ExtKeyStampDuty: "1",
			},
		}
		ad.Normalizer(c)
		err := ad.Validator(c)
		assert.NoError(t, err)
		assert.Equal(t, "4", c.Ext[mydata.ExtKeyTaxType].String())
	})
}
