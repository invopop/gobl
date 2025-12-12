package arca_test

import (
	"testing"

	"github.com/invopop/gobl/addons/ar/arca"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChargeNormalization(t *testing.T) {
	ad := tax.AddonForKey(arca.V4)

	testCases := []struct {
		name         string
		chargeKey    cbc.Key
		expectedCode string
	}{
		{"national taxes", arca.ChargeKeyNationalTaxes, "1"},
		{"provincial taxes", arca.ChargeKeyProvincialTaxes, "2"},
		{"municipal taxes", arca.ChargeKeyMunicipalTaxes, "3"},
		{"internal taxes", arca.ChargeKeyInternalTaxes, "4"},
		{"gross income tax", arca.ChargeKeyGrossIncomeTax, "5"},
		{"vat prepayment", arca.ChargeKeyVATPrepayment, "6"},
		{"gross income tax prepayment", arca.ChargeKeyGrossIncomeTaxPrepayment, "7"},
		{"municipal taxes prepayment", arca.ChargeKeyMunicipalTaxesPrepayment, "8"},
		{"other prepayments", arca.ChargeKeyOtherPrepayments, "9"},
		{"vat not categorized prepayment", arca.ChargeKeyVATNotCategorizedPrepayment, "13"},
		{"other", arca.ChargeKeyOther, "99"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			charge := &bill.Charge{
				Key:     tc.chargeKey,
				Percent: num.NewPercentage(10, 2),
			}
			ad.Normalizer(charge)
			assert.Equal(t, tc.expectedCode, charge.Ext[arca.ExtKeyTributeType].String())
		})
	}

	t.Run("unknown charge key does not set tribute type", func(t *testing.T) {
		charge := &bill.Charge{
			Key:     "unknown-key",
			Percent: num.NewPercentage(10, 2),
		}
		ad.Normalizer(charge)
		assert.Empty(t, charge.Ext[arca.ExtKeyTributeType])
	})

	t.Run("existing extensions are merged", func(t *testing.T) {
		charge := &bill.Charge{
			Key:     arca.ChargeKeyNationalTaxes,
			Percent: num.NewPercentage(10, 2),
			Ext: tax.Extensions{
				"custom-key": "custom-value",
			},
		}
		ad.Normalizer(charge)
		assert.Equal(t, "1", charge.Ext[arca.ExtKeyTributeType].String())
		assert.Equal(t, "custom-value", charge.Ext["custom-key"].String())
	})

	t.Run("existing tribute type is overwritten by normalization", func(t *testing.T) {
		charge := &bill.Charge{
			Key:     arca.ChargeKeyNationalTaxes, // would set "1"
			Percent: num.NewPercentage(10, 2),
			Ext: tax.Extensions{
				arca.ExtKeyTributeType: "5", // already set to gross income tax
			},
		}
		ad.Normalizer(charge)
		// Merge overwrites existing values with new ones
		assert.Equal(t, "1", charge.Ext[arca.ExtKeyTributeType].String())
	})
}

func TestChargeValidation(t *testing.T) {
	ad := tax.AddonForKey(arca.V4)

	t.Run("valid charge with tribute type passes", func(t *testing.T) {
		charge := &bill.Charge{
			Key:     arca.ChargeKeyNationalTaxes,
			Percent: num.NewPercentage(10, 2),
			Ext: tax.Extensions{
				arca.ExtKeyTributeType: "1",
			},
		}
		err := ad.Validator(charge)
		require.NoError(t, err)
	})

	t.Run("charge without tribute type does not require percent", func(t *testing.T) {
		charge := &bill.Charge{
			Key:    "custom-charge",
			Reason: "Some custom charge",
			// No percent and no tribute type extension
		}
		err := ad.Validator(charge)
		require.NoError(t, err)
	})

	t.Run("missing percent fails when tribute type is present", func(t *testing.T) {
		charge := &bill.Charge{
			Key: arca.ChargeKeyNationalTaxes,
			Ext: tax.Extensions{
				arca.ExtKeyTributeType: "1",
			},
		}
		err := ad.Validator(charge)
		assert.ErrorContains(t, err, "percent: cannot be blank")
	})

	t.Run("tribute type 'other' requires reason", func(t *testing.T) {
		charge := &bill.Charge{
			Key:     arca.ChargeKeyOther,
			Percent: num.NewPercentage(10, 2),
			Ext: tax.Extensions{
				arca.ExtKeyTributeType: "99", // TributeTypeOther
			},
		}
		err := ad.Validator(charge)
		assert.ErrorContains(t, err, "reason is required when tribute type is 'other'")
	})

	t.Run("tribute type 'other' with reason passes", func(t *testing.T) {
		charge := &bill.Charge{
			Key:     arca.ChargeKeyOther,
			Percent: num.NewPercentage(10, 2),
			Reason:  "Custom tax description",
			Ext: tax.Extensions{
				arca.ExtKeyTributeType: "99", // TributeTypeOther
			},
		}
		err := ad.Validator(charge)
		require.NoError(t, err)
	})

	t.Run("non-other tribute types do not require reason", func(t *testing.T) {
		charge := &bill.Charge{
			Key:     arca.ChargeKeyNationalTaxes,
			Percent: num.NewPercentage(10, 2),
			Ext: tax.Extensions{
				arca.ExtKeyTributeType: "1", // TributeTypeNationalTaxes
			},
			// No reason provided
		}
		err := ad.Validator(charge)
		require.NoError(t, err)
	})
}

func TestChargeIntegration(t *testing.T) {
	t.Run("charge on invoice is normalized and validated", func(t *testing.T) {
		inv := testInvoiceWithCharge(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())

		// Check that the charge was normalized
		require.Len(t, inv.Charges, 1)
		assert.Equal(t, "5", inv.Charges[0].Ext[arca.ExtKeyTributeType].String())
	})

	t.Run("invoice with charge missing percent when tribute type present fails", func(t *testing.T) {
		inv := testInvoiceWithCharge(t)
		require.NoError(t, inv.Calculate())

		// Remove the percent but keep tribute type
		inv.Charges[0].Percent = nil

		err := inv.Validate()
		assert.ErrorContains(t, err, "percent: cannot be blank")
	})

	t.Run("invoice with charge without tribute type passes", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Charges = []*bill.Charge{
			{
				Key:    "custom-charge",
				Reason: "Some custom charge without tribute type",
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
			Key:     arca.ChargeKeyGrossIncomeTax,
			Percent: num.NewPercentage(3, 2), // 3%
		},
	}
	return inv
}
