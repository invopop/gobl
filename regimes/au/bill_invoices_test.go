package au_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/au"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoice(price int64, combo *tax.Combo) *bill.Invoice {
	return &bill.Invoice{
		Regime:   tax.WithRegime(au.CountryCode),
		Code:     "0001",
		Currency: "AUD",
		Supplier: &org.Party{
			Name: "Windows to Fit Pty Ltd",
			TaxID: &tax.Identity{
				Country: "AU",
				Code:    "51824753556",
			},
		},
		Customer: &org.Party{
			Name: "Building Company",
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Window frame",
					Price: num.NewAmount(price, 0),
				},
				Taxes: tax.Set{combo},
			},
		},
	}
}

func standardGST() *tax.Combo {
	return &tax.Combo{Category: tax.CategoryGST, Rate: tax.RateGeneral}
}

func gstFree() *tax.Combo {
	return &tax.Combo{Category: tax.CategoryGST, Key: tax.KeyZero, Percent: num.NewPercentage(0, 3)}
}

func inputTaxed() *tax.Combo {
	return &tax.Combo{Category: tax.CategoryGST, Key: tax.KeyExempt}
}

func validateInvoice(t *testing.T, inv *bill.Invoice) error {
	t.Helper()
	require.NoError(t, inv.Calculate())
	return rules.Validate(inv)
}

func TestInvoiceSupplierABN(t *testing.T) {
	t.Run("valid tax invoice", func(t *testing.T) {
		assert.NoError(t, validateInvoice(t, testInvoice(100, standardGST())))
	})

	t.Run("missing ABN when GST is charged", func(t *testing.T) {
		inv := testInvoice(100, standardGST())
		inv.Supplier.TaxID = nil
		assert.ErrorContains(t, validateInvoice(t, inv),
			"[GOBL-AU-BILL-INVOICE-01] ($.supplier) invoice supplier must have an ABN when GST is charged")
	})

	t.Run("missing ABN on a GST-free sale", func(t *testing.T) {
		inv := testInvoice(100, gstFree())
		inv.Supplier.TaxID = nil
		assert.NoError(t, validateInvoice(t, inv))
	})

	t.Run("missing ABN on an input-taxed sale", func(t *testing.T) {
		inv := testInvoice(100, inputTaxed())
		inv.Supplier.TaxID = nil
		assert.NoError(t, validateInvoice(t, inv))
	})

	t.Run("uncalculated invoice without totals", func(t *testing.T) {
		inv := testInvoice(100, standardGST())
		inv.Type = bill.InvoiceTypeStandard
		inv.Supplier.TaxID = nil
		assert.NotContains(t, errorText(rules.Validate(inv)), "GOBL-AU-BILL-INVOICE")
	})

	t.Run("missing ABN on a credit note", func(t *testing.T) {
		inv := testInvoice(100, standardGST())
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Preceding = []*org.DocumentRef{{Code: "0000"}}
		inv.Supplier.TaxID = nil
		assert.NotContains(t, errorText(validateInvoice(t, inv)), "GOBL-AU-BILL-INVOICE-01")
	})
}

func TestInvoiceCustomerThreshold(t *testing.T) {
	t.Run("no customer below A$1,000", func(t *testing.T) {
		inv := testInvoice(900, standardGST()) // 990 with GST
		inv.Customer = nil
		assert.NoError(t, validateInvoice(t, inv))
	})

	t.Run("no customer at A$1,000 with GST", func(t *testing.T) {
		inv := testInvoice(1000, standardGST()) // 1,100 with GST
		inv.Customer = nil
		assert.ErrorContains(t, validateInvoice(t, inv),
			"[GOBL-AU-BILL-INVOICE-02] ($.customer) invoice customer is required when GST is charged on a total of A$1,000 or more")
	})

	t.Run("customer present at A$1,000 with GST", func(t *testing.T) {
		assert.NoError(t, validateInvoice(t, testInvoice(1000, standardGST())))
	})

	t.Run("no customer on a large GST-free sale", func(t *testing.T) {
		inv := testInvoice(5000, gstFree())
		inv.Customer = nil
		assert.NoError(t, validateInvoice(t, inv))
	})

	t.Run("no customer on a large sale in another currency", func(t *testing.T) {
		inv := testInvoice(5000, standardGST())
		inv.Currency = "USD"
		inv.ExchangeRates = []*currency.ExchangeRate{{From: "USD", To: "AUD", Amount: num.MakeAmount(15, 1)}}
		inv.Customer = nil
		assert.NoError(t, validateInvoice(t, inv))
	})
}

func errorText(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
