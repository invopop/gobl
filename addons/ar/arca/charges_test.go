package arca_test

import (
	"testing"

	"github.com/invopop/gobl/addons/ar/arca"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChargeValidation(t *testing.T) {
	ad := tax.AddonForKey(arca.V4)

	t.Run("valid tax charge", func(t *testing.T) {
		charge := &bill.Charge{
			Key:     bill.ChargeKeyTax,
			Percent: num.NewPercentage(10, 2),
			Ext: tax.Extensions{
				arca.ExtKeyTaxType: "1", // National Taxes
			},
		}
		err := ad.Validator(charge)
		require.NoError(t, err)
	})

	t.Run("tax charge missing ext", func(t *testing.T) {
		charge := &bill.Charge{
			Key:     bill.ChargeKeyTax,
			Percent: num.NewPercentage(10, 2),
			// No ext
		}
		err := ad.Validator(charge)
		assert.ErrorContains(t, err, "ar-arca-tax-type: required")
	})

	t.Run("tax charge with ext missing tax type", func(t *testing.T) {
		charge := &bill.Charge{
			Key:     bill.ChargeKeyTax,
			Percent: num.NewPercentage(10, 2),
			Ext: tax.Extensions{
				"other-ext": "value",
			},
		}
		err := ad.Validator(charge)
		assert.ErrorContains(t, err, "ar-arca-tax-type: required")
	})

	t.Run("tax type present but missing percent", func(t *testing.T) {
		charge := &bill.Charge{
			Key: bill.ChargeKeyTax,
			Ext: tax.Extensions{
				arca.ExtKeyTaxType: "1",
			},
			// No percent
		}
		err := ad.Validator(charge)
		assert.ErrorContains(t, err, "percent: cannot be blank")
	})

	t.Run("non-tax charge does not require ext or percent", func(t *testing.T) {
		charge := &bill.Charge{
			Key:    "custom-charge",
			Reason: "Some custom charge",
			// No percent and no ext
		}
		err := ad.Validator(charge)
		require.NoError(t, err)
	})

	t.Run("tax type other requires reason", func(t *testing.T) {
		charge := &bill.Charge{
			Key:     bill.ChargeKeyTax,
			Percent: num.NewPercentage(10, 2),
			Ext: tax.Extensions{
				arca.ExtKeyTaxType: "99", // TaxTypeOther
			},
			// No reason
		}
		err := ad.Validator(charge)
		assert.ErrorContains(t, err, "reason is required when tax type is 'other'")
	})

	t.Run("tax type other with reason", func(t *testing.T) {
		charge := &bill.Charge{
			Key:     bill.ChargeKeyTax,
			Percent: num.NewPercentage(10, 2),
			Reason:  "Custom tax description",
			Ext: tax.Extensions{
				arca.ExtKeyTaxType: "99", // TaxTypeOther
			},
		}
		err := ad.Validator(charge)
		require.NoError(t, err)
	})

	t.Run("non-other tax types do not require reason", func(t *testing.T) {
		charge := &bill.Charge{
			Key:     bill.ChargeKeyTax,
			Percent: num.NewPercentage(10, 2),
			Ext: tax.Extensions{
				arca.ExtKeyTaxType: "1", // TaxTypeNationalTaxes
			},
			// No reason provided
		}
		err := ad.Validator(charge)
		require.NoError(t, err)
	})
}

func TestChargeIntegration(t *testing.T) {
	t.Run("charge on invoice is validated", func(t *testing.T) {
		inv := testInvoiceWithCharge(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())

		require.Len(t, inv.Charges, 1)
		assert.Equal(t, "5", inv.Charges[0].Ext[arca.ExtKeyTaxType].String())
	})

	t.Run("invoice with tax charge missing ext fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Charges = []*bill.Charge{
			{
				Key:     bill.ChargeKeyTax,
				Percent: num.NewPercentage(3, 2),
				// No ext
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "ar-arca-tax-type: required")
	})

	t.Run("invoice with charge missing percent fails", func(t *testing.T) {
		inv := testInvoiceWithCharge(t)
		require.NoError(t, inv.Calculate())

		inv.Charges[0].Percent = nil

		err := inv.Validate()
		assert.ErrorContains(t, err, "percent: cannot be blank")
	})

	t.Run("invoice with non-tax charge passes", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Charges = []*bill.Charge{
			{
				Key:    "custom-charge",
				Reason: "Some custom charge without tax type",
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})
}

func testInvoiceWithCharge(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := testInvoiceWithGoods(t)
	inv.Charges = []*bill.Charge{
		{
			Key:     bill.ChargeKeyTax,
			Percent: num.NewPercentage(3, 2), // 3%
			Ext: tax.Extensions{
				arca.ExtKeyTaxType: "5", // Gross Income Tax
			},
		},
	}
	return inv
}
