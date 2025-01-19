package sg_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/sg"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Code:    "199912345A",
				Country: "SG",
			},
			Name: "Test Supplier",
			Addresses: []*org.Address{
				{
					Street:  "Test Street",
					Code:    "123456",
					Country: l10n.SG.ISO(),
				},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			Addresses: []*org.Address{
				{
					Street:  "Test Street",
					Code:    "123456",
					Country: l10n.SG.ISO(),
				},
			},
		},
		Code:     "0001",
		Currency: "SGD",
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(100, 0),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryGST,
						Rate:     tax.RateStandard,
					},
				},
			},
		},
	}
}

func TestValidInvoice(t *testing.T) {
	inv := validInvoice()
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestValidReceiptInvoice(t *testing.T) {
	inv := validInvoice()
	inv.SetTags(sg.TagInvoiceReceipt)
	inv.Customer = nil

	require.NoError(t, inv.Calculate())
	assert.Len(t, inv.Notes, 1)
	assert.Equal(t, inv.Notes[0].Src, sg.TagInvoiceReceipt)
	assert.Equal(t, inv.Notes[0].Text, "Price Payable includes GST")
	require.NoError(t, inv.Validate())

}

func TestValidSimplifiedInvoice(t *testing.T) {
	inv := validInvoice()
	inv.SetTags(tax.TagSimplified)
	inv.Customer = nil

	require.NoError(t, inv.Calculate())
	assert.Len(t, inv.Notes, 1)
	assert.Equal(t, inv.Notes[0].Src, tax.TagSimplified)
	assert.Equal(t, inv.Notes[0].Text, "Price Payable includes GST")
	require.NoError(t, inv.Validate())
}
