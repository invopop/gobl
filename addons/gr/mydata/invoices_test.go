package mydata_test

import (
	"testing"

	"github.com/invopop/gobl/addons/gr/mydata"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime: tax.WithRegime("GR"),
		Addons: tax.WithAddons(mydata.V1),
		Tags:   tax.WithTags(mydata.TagServices),
		Series: "TEST",
		Code:   "0002",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "EL",
				Code:    "728089281",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "EL",
				Code:    "841442160",
			},
			Addresses: []*org.Address{
				{
					Locality: "Athens",
					Code:     "11528",
				},
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.MakeAmount(10000, 2),
					Unit:  org.UnitPackage,
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
			},
		},
		Payment: &bill.Payment{
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCreditTransfer,
			},
		},
	}
}

func TestInvoiceValidation(t *testing.T) {
	inv := validInvoice()
	require.NoError(t, inv.Calculate())
	assert.NoError(t, inv.Validate())

	// Make it invalid
	inv.Series = ""
	inv.Supplier.TaxID.Code = ""
	inv.Customer.Addresses = nil
	inv.Lines[0].Quantity = num.MakeAmount(0, 0)

	require.NoError(t, inv.Calculate())

	err := inv.Validate()
	assert.ErrorContains(t, err, "series: cannot be blank")
	assert.ErrorContains(t, err, "supplier: (tax_id: (code: cannot be blank")
	assert.ErrorContains(t, err, "customer: (addresses: cannot be blank")
	assert.ErrorContains(t, err, "lines: (0: (total: must be greater than 0")

	// Go in two parts as the payment errors are independent
	inv.Payment.Instructions.Key = "debit-transfer"
	inv.Payment.Instructions.Ext = nil
	require.NoError(t, inv.Calculate())
	err = inv.Validate()
	assert.ErrorContains(t, err, "payment: (instructions: (ext: (gr-mydata-payment-means: required.).).)")
}

func TestSimplifiedInvoiceValidation(t *testing.T) {
	inv := validInvoice()
	inv.SetTags(tax.TagSimplified)
	inv.Customer.TaxID = nil
	inv.Customer.Addresses = nil

	require.NoError(t, inv.Calculate())
	assert.NoError(t, inv.Validate())
}

func TestPrecedingValidation(t *testing.T) {
	inv := validInvoice()

	inv.Preceding = []*bill.Preceding{
		{
			Code: "123",
			Stamps: []*head.Stamp{
				{
					Provider: "unexpected",
					Value:    "1234",
				},
			},
		},
	}
	inv.Type = bill.InvoiceTypeCreditNote

	require.NoError(t, inv.Calculate())

	err := inv.Validate()
	assert.ErrorContains(t, err, "preceding: (0: (stamps: missing iapr-mark stamp.).)")

	inv.Preceding[0].Stamps[0].Provider = "iapr-mark"
	require.NoError(t, inv.Validate())
}

func TestInvoiceLineItemIncomeExt(t *testing.T) {
	t.Run("no ext", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("income cat, no type", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Item.Ext = tax.Extensions{
			mydata.ExtKeyIncomeCat: "category1_1",
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, inv.Validate(), "lines: (0: (item: (ext: (gr-mydata-income-type: required.).).).)")
	})

	t.Run("income type, no cat", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Item.Ext = tax.Extensions{
			mydata.ExtKeyIncomeType: "E3_106",
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, inv.Validate(), "lines: (0: (item: (ext: (gr-mydata-income-cat: required.).).).)")
	})

	t.Run("income cat with type", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Item.Ext = tax.Extensions{
			mydata.ExtKeyIncomeType: "E3_106",
			mydata.ExtKeyIncomeCat:  "category1_1",
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})
}
