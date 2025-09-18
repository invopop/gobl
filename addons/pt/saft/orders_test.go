package saft_test

import (
	"testing"
	"time"

	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderValidation(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)

	t.Run("valid order", func(t *testing.T) {
		ord := validOrder()
		require.NoError(t, addon.Validator(ord))
	})

	t.Run("missing work type", func(t *testing.T) {
		ord := validOrder()

		ord.Tax = nil
		assert.ErrorContains(t, addon.Validator(ord), "tax: (ext: (pt-saft-work-type: required")

		ord.Tax = new(bill.Tax)
		assert.ErrorContains(t, addon.Validator(ord), "tax: (ext: (pt-saft-work-type: required")
	})

	t.Run("invalid work type", func(t *testing.T) {
		ord := validOrder()

		ord.Tax.Ext = tax.Extensions{
			saft.ExtKeyWorkType: saft.WorkTypeProforma, // Proforma is not valid in orders, only in invoices
		}

		assert.ErrorContains(t, addon.Validator(ord), "value 'PF' invalid")
	})

	t.Run("missing VAT category in lines", func(t *testing.T) {
		ord := validOrder()

		ord.Lines[0].Taxes = nil
		assert.ErrorContains(t, addon.Validator(ord), "lines: (0: (taxes: missing category VAT")
	})

	t.Run("missing value date", func(t *testing.T) {
		ord := validOrder()
		ord.ValueDate = nil
		assert.ErrorContains(t, addon.Validator(ord), "value_date: cannot be blank")
	})
}

func TestOrderNormalization(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)

	t.Run("nil order", func(t *testing.T) {
		assert.NotPanics(t, func() {
			var inv *bill.Order
			addon.Normalizer(inv)
		})
	})

	t.Run("purchase order type", func(t *testing.T) {
		ord := &bill.Order{
			Type: bill.OrderTypePurchase,
		}
		addon.Normalizer(ord)
		require.NotNil(t, ord.Tax)
		require.NotNil(t, ord.Tax.Ext)
		assert.Equal(t, saft.WorkTypePurchaseOrder, ord.Tax.Ext[saft.ExtKeyWorkType])
	})

	t.Run("quote order type", func(t *testing.T) {
		ord := &bill.Order{
			Type: bill.OrderTypeQuote,
		}
		addon.Normalizer(ord)
		require.NotNil(t, ord.Tax)
		require.NotNil(t, ord.Tax.Ext)
		assert.Equal(t, saft.WorkTypeBudgets, ord.Tax.Ext[saft.ExtKeyWorkType])
	})

	t.Run("respect existing value", func(t *testing.T) {
		ord := &bill.Order{
			Type: bill.OrderTypePurchase,
			Tax: &bill.Tax{
				Ext: tax.Extensions{
					saft.ExtKeyWorkType: saft.WorkTypeOther,
				},
			},
		}
		addon.Normalizer(ord)
		assert.Equal(t, saft.WorkTypeOther, ord.Tax.Ext[saft.ExtKeyWorkType])
	})

	t.Run("sets default value date from issue date", func(t *testing.T) {
		ord := validOrder()
		ord.ValueDate = nil
		addon.Normalizer(ord)
		assert.Equal(t, &ord.IssueDate, ord.ValueDate)
	})

	t.Run("sets default value date from operation date", func(t *testing.T) {
		ord := validOrder()
		ord.OperationDate = cal.NewDate(2024, 12, 2)
		ord.ValueDate = nil
		addon.Normalizer(ord)
		assert.Equal(t, ord.OperationDate, ord.ValueDate)
	})

	t.Run("keeps existing value date", func(t *testing.T) {
		ord := validOrder()
		ord.ValueDate = cal.NewDate(2024, 12, 2)
		addon.Normalizer(ord)
		assert.Equal(t, cal.NewDate(2024, 12, 2), ord.ValueDate)
	})

	t.Run("sets today as value date when no issue date is set", func(t *testing.T) {
		ord := validOrder()
		ord.IssueDate = cal.Date{}
		ord.ValueDate = nil

		addon.Normalizer(ord)

		loc, err := time.LoadLocation("Europe/Lisbon")
		require.NoError(t, err)
		today := cal.TodayIn(loc)
		assert.Equal(t, &today, ord.ValueDate)
	})

}

func validOrder() *bill.Order {
	return &bill.Order{
		Regime: tax.WithRegime("PT"),
		Addons: tax.WithAddons(saft.V1),
		Type:   bill.OrderTypePurchase,
		Tax: &bill.Tax{
			Ext: tax.Extensions{
				saft.ExtKeyWorkType: saft.WorkTypePurchaseOrder,
			},
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Code:    "123456789",
				Country: "PT",
			},
			Name: "Test Supplier",
		},
		Customer: &org.Party{
			Name: "Test Customer",
		},
		Series:    "NE SERIES-A",
		Code:      "123",
		Currency:  "EUR",
		IssueDate: cal.MakeDate(2023, 1, 1),
		ValueDate: cal.NewDate(2022, 12, 31),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(100, 0),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Key:      "standard",
					},
				},
			},
		},
	}
}
