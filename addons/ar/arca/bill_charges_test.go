package arca_test

import (
	"testing"

	"github.com/invopop/gobl/addons/ar/arca"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChargeValidation(t *testing.T) {
	t.Run("valid tax charge", func(t *testing.T) {
		charge := &bill.Charge{
			Key:     bill.ChargeKeyTax,
			Percent: num.NewPercentage(10, 2),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				arca.ExtKeyTaxType: "1", // National Taxes
			}),
		}
		err := rules.Validate(charge, withAddonContext())
		require.NoError(t, err)
	})

	t.Run("tax charge missing ext", func(t *testing.T) {
		charge := &bill.Charge{
			Key:     bill.ChargeKeyTax,
			Percent: num.NewPercentage(10, 2),
			// No ext
		}
		err := rules.Validate(charge, withAddonContext())
		assert.ErrorContains(t, err, "tax charge requires 'ar-arca-tax-type' extension")
	})

	t.Run("tax charge with ext missing tax type", func(t *testing.T) {
		charge := &bill.Charge{
			Key:     bill.ChargeKeyTax,
			Percent: num.NewPercentage(10, 2),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				"other-ext": "value",
			}),
		}
		err := rules.Validate(charge, withAddonContext())
		assert.ErrorContains(t, err, "tax charge requires 'ar-arca-tax-type' extension")
	})

	t.Run("tax type present but missing percent", func(t *testing.T) {
		charge := &bill.Charge{
			Key: bill.ChargeKeyTax,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				arca.ExtKeyTaxType: "1",
			}),
			// No percent
		}
		err := rules.Validate(charge, withAddonContext())
		assert.ErrorContains(t, err, "percent is required when tax type is set")
	})

	t.Run("non-tax charge does not require ext or percent", func(t *testing.T) {
		charge := &bill.Charge{
			Key:    "custom-charge",
			Reason: "Some custom charge",
			// No percent and no ext
		}
		err := rules.Validate(charge, withAddonContext())
		require.NoError(t, err)
	})

	t.Run("tax type other requires reason", func(t *testing.T) {
		charge := &bill.Charge{
			Key:     bill.ChargeKeyTax,
			Percent: num.NewPercentage(10, 2),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				arca.ExtKeyTaxType: "99", // TaxTypeOther
			}),
			// No reason
		}
		err := rules.Validate(charge, withAddonContext())
		assert.ErrorContains(t, err, "reason is required when tax type is 'other'")
	})

	t.Run("tax type other with reason", func(t *testing.T) {
		charge := &bill.Charge{
			Key:     bill.ChargeKeyTax,
			Percent: num.NewPercentage(10, 2),
			Reason:  "Custom tax description",
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				arca.ExtKeyTaxType: "99", // TaxTypeOther
			}),
		}
		err := rules.Validate(charge, withAddonContext())
		require.NoError(t, err)
	})

	t.Run("non-other tax types do not require reason", func(t *testing.T) {
		charge := &bill.Charge{
			Key:     bill.ChargeKeyTax,
			Percent: num.NewPercentage(10, 2),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				arca.ExtKeyTaxType: "1", // TaxTypeNationalTaxes
			}),
			// No reason provided
		}
		err := rules.Validate(charge, withAddonContext())
		require.NoError(t, err)
	})
}

func withAddonContext() rules.WithContext {
	return func(rc *rules.Context) {
		rc.Set(rules.ContextKey(arca.V4), tax.AddonForKey(arca.V4))
	}
}

func TestChargeIntegration(t *testing.T) {
	t.Run("charge on invoice is validated", func(t *testing.T) {
		inv := testInvoiceWithCharge(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))

		require.Len(t, inv.Charges, 1)
		assert.Equal(t, "5", inv.Charges[0].Ext.Get(arca.ExtKeyTaxType).String())
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
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "tax charge requires 'ar-arca-tax-type' extension")
	})

	t.Run("invoice with charge missing percent fails", func(t *testing.T) {
		inv := testInvoiceWithCharge(t)
		require.NoError(t, inv.Calculate())

		inv.Charges[0].Percent = nil

		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "percent is required when tax type is set")
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
		require.NoError(t, rules.Validate(inv))
	})
}

func testInvoiceWithCharge(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := testInvoiceWithGoods(t)
	inv.Charges = []*bill.Charge{
		{
			Key:     bill.ChargeKeyTax,
			Percent: num.NewPercentage(3, 2), // 3%
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				arca.ExtKeyTaxType: "5", // Gross Income Tax
			}),
		},
	}
	return inv
}
