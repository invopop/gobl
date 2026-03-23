package mydata_test

import (
	"testing"

	"github.com/invopop/gobl/addons/gr/mydata"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
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
					Price: num.NewAmount(10000, 2),
					Unit:  org.UnitPackage,
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "general",
					},
				},
			},
		},
		Payment: &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCreditTransfer,
			},
		},
	}
}

func TestInvoiceValidation(t *testing.T) {
	inv := validInvoice()
	require.NoError(t, inv.Calculate())
	assert.NoError(t, rules.Validate(inv))

	// Make it invalid
	inv.Series = ""
	inv.Supplier.TaxID = nil
	inv.Customer.Addresses = nil
	inv.Lines[0].Quantity = num.MakeAmount(0, 0)

	require.NoError(t, inv.Calculate())

	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "series is required")
	assert.ErrorContains(t, err, "supplier tax ID is required")
	assert.ErrorContains(t, err, "customer addresses are required")
	assert.ErrorContains(t, err, "line total must be positive")

	// Go in two parts as the payment errors are independent
	inv.Payment.Instructions.Key = "debit-transfer"
	inv.Payment.Instructions.Ext = nil
	require.NoError(t, inv.Calculate())
	err = rules.Validate(inv)
	assert.ErrorContains(t, err, "payment instructions require 'gr-mydata-payment-means' extension")
}

func TestSimplifiedInvoiceValidation(t *testing.T) {
	inv := validInvoice()
	inv.SetTags(tax.TagSimplified)
	inv.Customer.TaxID = nil
	inv.Customer.Addresses = nil

	require.NoError(t, inv.Calculate())
	assert.NoError(t, rules.Validate(inv))
}

func TestOtherInvoiceTypeValidation(t *testing.T) {
	inv := validInvoice()
	inv.Type = bill.InvoiceTypeOther
	inv.Tax = &bill.Tax{
		Ext: tax.Extensions{
			mydata.ExtKeyInvoiceType: "8.2",
		},
	}
	inv.Customer.TaxID = nil
	inv.Customer.Addresses = nil

	require.NoError(t, inv.Calculate())
	assert.NoError(t, rules.Validate(inv))
}

func TestPrecedingValidation(t *testing.T) {
	inv := validInvoice()

	inv.Preceding = []*org.DocumentRef{
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

	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "preceding document requires 'iapr-mark' stamp")

	inv.Preceding[0].Stamps[0].Provider = "iapr-mark"
	require.NoError(t, rules.Validate(inv))
}

func TestInvoiceLineItemIncomeExt(t *testing.T) {
	t.Run("no ext", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("income cat, no type", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Item.Ext = tax.Extensions{
			mydata.ExtKeyIncomeCat: "category1_1",
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "income extensions 'gr-mydata-income-cat' and 'gr-mydata-income-type' must both be present")
	})

	t.Run("income type, no cat", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Item.Ext = tax.Extensions{
			mydata.ExtKeyIncomeType: "E3_106",
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "income extensions 'gr-mydata-income-cat' and 'gr-mydata-income-type' must both be present")
	})

	t.Run("income cat with type", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Item.Ext = tax.Extensions{
			mydata.ExtKeyIncomeType: "E3_106",
			mydata.ExtKeyIncomeCat:  "category1_1",
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})
}
