package facturae_test

import (
	"testing"

	"github.com/invopop/gobl/addons/es/facturae"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceValidation(t *testing.T) {
	inv := testInvoiceStandard(t)
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	delete(inv.Tax.Ext, facturae.ExtKeyDocType)
	err := inv.Validate()
	assert.ErrorContains(t, err, "tax: (ext: (es-facturae-doc-type: required.).)")
}

func TestInvoicePrecedingValidation(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Type = bill.InvoiceTypeCreditNote

	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	assert.ErrorContains(t, err, "preceding: cannot be blank.")

	inv.Preceding = []*org.DocumentRef{
		{
			Code: "123TEST",
		},
	}
	require.NoError(t, inv.Calculate())
	err = inv.Validate()
	assert.ErrorContains(t, err, "preceding: (0: (ext: (es-facturae-correction: required.); issue_date: cannot be blank.).)")

	inv.Preceding[0].Ext = tax.Extensions{
		facturae.ExtKeyCorrection: "01",
	}
	inv.Preceding[0].IssueDate = cal.NewDate(2022, 6, 13)
	require.NoError(t, inv.Calculate())
	err = inv.Validate()
	assert.NoError(t, err)
}

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	i := &bill.Invoice{
		Regime: tax.WithRegime("ES"),
		Addons: tax.WithAddons(facturae.V3),
		// Tags:     tax.WithTags(tax.TagSelfBilled),
		Code:     "123TEST",
		Currency: "EUR",
		Tax: &bill.Tax{
			// Addons:        []cbc.Key{facturae.KeyV3},
			PricesInclude: tax.CategoryVAT,
		},
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "54387763P",
			},
			Name: "Test Customer",
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
				Discounts: []*bill.LineDiscount{
					{
						Reason:  "Testing",
						Percent: num.NewPercentage(10, 2),
					},
				},
			},
		},
	}
	return i
}
