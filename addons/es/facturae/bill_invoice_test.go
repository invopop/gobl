package facturae_test

import (
	"testing"

	"github.com/invopop/gobl/addons/es/facturae"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceValidation(t *testing.T) {
	t.Run("standard invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("missing supplier tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-ES-BILL-INVOICE-02] ($.supplier.tax_id) invoice supplier tax ID in Spain is required")
	})

	t.Run("missing customer tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.NoError(t, err, "if there is no country, customer tax ID code is not required")
	})

	t.Run("missing ext key doc type", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		delete(inv.Tax.Ext, facturae.ExtKeyDocType)
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-ES-FACTURAE-V3-BILL-INVOICE-03] ($.tax.ext) tax ext require 'es-facturae-doc-type' and 'es-facturae-invoice-class' extensions")
	})

}

func TestInvoicePrecedingValidation(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Type = bill.InvoiceTypeCreditNote

	require.NoError(t, inv.Calculate())
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "[GOBL-ES-FACTURAE-V3-BILL-INVOICE-04] ($.preceding) preceding document reference is required for credit-note, corrective, debit-note invoices")

	inv.Preceding = []*org.DocumentRef{
		{
			Code: "123TEST",
		},
	}
	require.NoError(t, inv.Calculate())
	err = rules.Validate(inv)
	assert.ErrorContains(t, err, "[GOBL-ES-FACTURAE-V3-BILL-INVOICE-05] ($.preceding[0].issue_date) preceding document issue date is required; [GOBL-ES-FACTURAE-V3-BILL-INVOICE-06] ($.preceding[0].ext) preceding document ext require 'es-facturae-correction' extension")

	inv.Preceding[0].Ext = tax.Extensions{
		facturae.ExtKeyCorrection: "01",
	}
	inv.Preceding[0].IssueDate = cal.NewDate(2022, 6, 13)
	require.NoError(t, inv.Calculate())
	err = rules.Validate(inv)
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
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "general",
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
