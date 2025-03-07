package saft_test

import (
	"testing"

	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocType(t *testing.T) {
	t.Run("invoice", func(t *testing.T) {
		tests := []struct {
			name string
			inv  *bill.Invoice
			want cbc.Code
		}{
			{
				name: "nil",
				inv:  nil,
				want: cbc.CodeEmpty,
			},
			{
				name: "no tax",
				inv:  &bill.Invoice{},
				want: cbc.CodeEmpty,
			},
			{
				name: "no extensions",
				inv: &bill.Invoice{
					Tax: &bill.Tax{},
				},
				want: cbc.CodeEmpty,
			},
			{
				name: "invoice type",
				inv: &bill.Invoice{
					Tax: &bill.Tax{
						Ext: tax.Extensions{
							saft.ExtKeyInvoiceType: saft.InvoiceTypeStandard,
						},
					},
				},
				want: saft.InvoiceTypeStandard,
			},
			{
				name: "work type",
				inv: &bill.Invoice{
					Tax: &bill.Tax{
						Ext: tax.Extensions{
							saft.ExtKeyWorkType: saft.WorkTypeProforma,
						},
					},
				},
				want: saft.WorkTypeProforma,
			},
			{
				name: "invoice type takes precedence",
				inv: &bill.Invoice{
					Tax: &bill.Tax{
						Ext: tax.Extensions{
							saft.ExtKeyInvoiceType: saft.InvoiceTypeStandard,
							saft.ExtKeyWorkType:    saft.WorkTypeProforma,
						},
					},
				},
				want: saft.InvoiceTypeStandard,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := saft.DocType(tt.inv)
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			})
		}
	})

	t.Run("payment", func(t *testing.T) {
		tests := []struct {
			name string
			pmt  *bill.Payment
			want cbc.Code
		}{
			{
				name: "nil",
				pmt:  nil,
				want: cbc.CodeEmpty,
			},
			{
				name: "no extensions",
				pmt:  &bill.Payment{},
				want: cbc.CodeEmpty,
			},
			{
				name: "payment type",
				pmt: &bill.Payment{
					Ext: tax.Extensions{
						saft.ExtKeyPaymentType: saft.PaymentTypeCash,
					},
				},
				want: saft.PaymentTypeCash,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := saft.DocType(tt.pmt)
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			})
		}
	})

	t.Run("delivery", func(t *testing.T) {
		tests := []struct {
			name string
			dlv  *bill.Delivery
			want cbc.Code
		}{
			{
				name: "nil",
				dlv:  nil,
				want: cbc.CodeEmpty,
			},
			{
				name: "no tax",
				dlv:  &bill.Delivery{},
				want: cbc.CodeEmpty,
			},
			{
				name: "no extensions",
				dlv: &bill.Delivery{
					Tax: &bill.Tax{},
				},
				want: cbc.CodeEmpty,
			},
			{
				name: "movement type",
				dlv: &bill.Delivery{
					Tax: &bill.Tax{
						Ext: tax.Extensions{
							saft.ExtKeyMovementType: saft.MovementTypeDeliveryNote,
						},
					},
				},
				want: saft.MovementTypeDeliveryNote,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := saft.DocType(tt.dlv)
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			})
		}
	})

	t.Run("order", func(t *testing.T) {
		tests := []struct {
			name string
			ord  *bill.Order
			want cbc.Code
		}{
			{
				name: "nil",
				ord:  nil,
				want: cbc.CodeEmpty,
			},
			{
				name: "no tax",
				ord:  &bill.Order{},
				want: cbc.CodeEmpty,
			},
			{
				name: "no extensions",
				ord: &bill.Order{
					Tax: &bill.Tax{},
				},
				want: cbc.CodeEmpty,
			},
			{
				name: "work type",
				ord: &bill.Order{
					Tax: &bill.Tax{
						Ext: tax.Extensions{
							saft.ExtKeyWorkType: saft.WorkTypePurchaseOrder,
						},
					},
				},
				want: saft.WorkTypePurchaseOrder,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := saft.DocType(tt.ord)
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			})
		}
	})

	t.Run("unknown type", func(t *testing.T) {
		got, err := saft.DocType("not a document")
		assert.Error(t, err)
		assert.Equal(t, cbc.CodeEmpty, got)
		assert.Contains(t, err.Error(), "unsupported document type")
	})
}
