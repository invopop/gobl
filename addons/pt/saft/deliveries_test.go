package saft_test

import (
	"testing"

	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeliveryValidation(t *testing.T) {
	t.Run("valid delivery", func(t *testing.T) {
		dlv := validDelivery()
		require.NoError(t, rules.Validate(dlv, withAddonContext()))
	})

	t.Run("missing movement type", func(t *testing.T) {
		dlv := validDelivery()

		dlv.Tax = nil
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "tax requires 'pt-saft-movement-type' extension")

		dlv.Tax = new(bill.Tax)
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "tax requires 'pt-saft-movement-type' extension")
	})

	t.Run("missing despatch date", func(t *testing.T) {
		dlv := validDelivery()
		dlv.DespatchDate = nil
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "cannot be blank")
	})

	t.Run("invalid series format", func(t *testing.T) {
		dlv := validDelivery()

		dlv.Series = "SERIES-A"
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "series format must be valid")
	})

	t.Run("invalid code format", func(t *testing.T) {
		dlv := validDelivery()

		dlv.Code = "ABCD"
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "code format must be valid")
	})

	t.Run("valid full code", func(t *testing.T) {
		dlv := validDelivery()

		dlv.Series = ""
		dlv.Code = "GR SERIES-A/123"
		assert.NoError(t, rules.Validate(dlv, withAddonContext()))
	})

	t.Run("invalid full code", func(t *testing.T) {
		dlv := validDelivery()

		dlv.Series = ""
		dlv.Code = "ABCDEF"
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "code format must be valid")
	})

	t.Run("missing supplier tax ID", func(t *testing.T) {
		dlv := validDelivery()

		dlv.Supplier.TaxID = nil
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "supplier tax ID is required")

		dlv.Supplier.TaxID = &tax.Identity{
			Country: "PT",
			Code:    "",
		}
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "supplier tax ID code is required")

		// dlv.Supplier = nil is caught by core GOBL rules (supplier is required)
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
		assert.Equal(t, saft.MovementTypeDeliveryNote, dlv.Tax.Ext.Get(saft.ExtKeyMovementType))
	})

	t.Run("waybill type", func(t *testing.T) {
		dlv := &bill.Delivery{
			Type: bill.DeliveryTypeWaybill,
		}
		addon.Normalizer(dlv)
		require.NotNil(t, dlv.Tax)
		require.NotNil(t, dlv.Tax.Ext)
		assert.Equal(t, saft.MovementTypeWaybill, dlv.Tax.Ext.Get(saft.ExtKeyMovementType))
	})

	t.Run("respect existing value", func(t *testing.T) {
		dlv := &bill.Delivery{
			Type: bill.DeliveryTypeNote,
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(tax.ExtMap{
					saft.ExtKeyMovementType: saft.MovementTypeFixedAssets,
				}),
			},
		}
		addon.Normalizer(dlv)
		assert.Equal(t, saft.MovementTypeFixedAssets, dlv.Tax.Ext.Get(saft.ExtKeyMovementType))
	})
}

func validDelivery() *bill.Delivery {
	date := cal.NewDate(2023, 1, 1)

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
			Name: "Test Customer",
		},
		DespatchDate: date,
		Lines: []*bill.Line{
			{
				Index:    1,
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name: "Test Item",
					Unit: "one",
					Ext: tax.ExtensionsOf(tax.ExtMap{
						saft.ExtKeyProductType: saft.ProductTypeService,
					}),
				},
			},
		},
		Tax: &bill.Tax{
			Ext: tax.ExtensionsOf(tax.ExtMap{
				saft.ExtKeyMovementType: saft.MovementTypeDeliveryNote,
			}),
		},
		Totals: &bill.Totals{},
	}
}
