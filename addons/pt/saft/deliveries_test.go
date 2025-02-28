package saft_test

import (
	"testing"

	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeliveryValidation(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)
	require.NotNil(t, addon)

	t.Run("valid delivery", func(t *testing.T) {
		dlv := validDelivery()
		require.NoError(t, addon.Validator(dlv))
	})

	t.Run("missing movement type", func(t *testing.T) {
		dlv := validDelivery()

		dlv.Tax = nil
		assert.ErrorContains(t, addon.Validator(dlv), "tax: (ext: (pt-saft-movement-type: required")

		dlv.Tax = new(bill.Tax)
		assert.ErrorContains(t, addon.Validator(dlv), "tax: (ext: (pt-saft-movement-type: required")
	})
}

func TestDeliveryNormalization(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)
	require.NotNil(t, addon)

	t.Run("note type", func(t *testing.T) {
		dlv := &bill.Delivery{
			Type: bill.DeliveryTypeNote,
		}
		addon.Normalizer(dlv)
		require.NotNil(t, dlv.Tax)
		require.NotNil(t, dlv.Tax.Ext)
		assert.Equal(t, saft.MovementTypeDeliveryNote, dlv.Tax.Ext[saft.ExtKeyMovementType])
	})

	t.Run("waybill type", func(t *testing.T) {
		dlv := &bill.Delivery{
			Type: bill.DeliveryTypeWaybill,
		}
		addon.Normalizer(dlv)
		require.NotNil(t, dlv.Tax)
		require.NotNil(t, dlv.Tax.Ext)
		assert.Equal(t, saft.MovementTypeWaybill, dlv.Tax.Ext[saft.ExtKeyMovementType])
	})

	t.Run("respect existing value", func(t *testing.T) {
		dlv := &bill.Delivery{
			Type: bill.DeliveryTypeNote,
			Tax: &bill.Tax{
				Ext: tax.Extensions{
					saft.ExtKeyMovementType: saft.MovementTypeFixedAssets,
				},
			},
		}
		addon.Normalizer(dlv)
		assert.Equal(t, saft.MovementTypeFixedAssets, dlv.Tax.Ext[saft.ExtKeyMovementType])
	})
}

func validDelivery() *bill.Delivery {
	date := cal.NewDate(2023, 1, 1)

	price, err := num.AmountFromString("100.00")
	if err != nil {
		panic(err)
	}

	quantity, err := num.AmountFromString("1")
	if err != nil {
		panic(err)
	}

	return &bill.Delivery{
		Type:      bill.DeliveryTypeNote,
		IssueDate: *date,
		Series:    "GR SERIES-A",
		Code:      "123",
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: "PT",
				Code:    "123456789",
			},
			Name: "Test Supplier",
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: "PT",
				Code:    "987654321",
			},
			Name: "Test Customer",
		},
		Lines: []*bill.Line{
			{
				Item: &org.Item{
					Name:  "Test Item",
					Price: &price,
				},
				Quantity: quantity,
			},
		},
		Tax: &bill.Tax{
			Ext: tax.Extensions{
				saft.ExtKeyMovementType: saft.MovementTypeDeliveryNote,
			},
		},
	}
}
