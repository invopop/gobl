package arca_test

import (
	"testing"

	"github.com/invopop/gobl/addons/ar/arca"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceDocumentScenarios(t *testing.T) {
	ad := tax.AddonForKey(arca.V4)

	t.Run("addon has scenarios defined", func(t *testing.T) {
		require.NotNil(t, ad.Scenarios)
		require.NotEmpty(t, ad.Scenarios)

		// Find the invoice scenario set
		var invoiceScenarios *tax.ScenarioSet
		for _, ss := range ad.Scenarios {
			if ss.Schema == bill.ShortSchemaInvoice {
				invoiceScenarios = ss
				break
			}
		}
		require.NotNil(t, invoiceScenarios, "invoice scenarios should be defined")
		assert.NotEmpty(t, invoiceScenarios.List)
	})

	t.Run("invoice type A - B2B with AR customer", func(t *testing.T) {
		inv := testInvoiceB2B(t)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "1", inv.Tax.Ext[arca.ExtKeyDocType].String())
	})

	t.Run("debit note type A - B2B with AR customer", func(t *testing.T) {
		inv := testInvoiceB2B(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "2", inv.Tax.Ext[arca.ExtKeyDocType].String())
	})

	t.Run("credit note type A - B2B with AR customer", func(t *testing.T) {
		inv := testInvoiceB2B(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "3", inv.Tax.Ext[arca.ExtKeyDocType].String())
	})

	t.Run("invoice type B - B2C no customer", func(t *testing.T) {
		inv := testInvoiceB2C(t)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "6", inv.Tax.Ext[arca.ExtKeyDocType].String())
	})

	t.Run("invoice type B - B2C no tax ID", func(t *testing.T) {
		inv := testInvoiceB2C(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "12345678",
				Ext: tax.Extensions{
					arca.ExtKeyIdentityType: "96", // DNI
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "6", inv.Tax.Ext[arca.ExtKeyDocType].String())
	})

	t.Run("invoice type B - B2C foreign customer", func(t *testing.T) {
		inv := testInvoiceB2C(t)
		inv.Customer.TaxID = &tax.Identity{
			Country: "US",
			Code:    "123456789",
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "6", inv.Tax.Ext[arca.ExtKeyDocType].String())
	})

	t.Run("debit note type B - B2C", func(t *testing.T) {
		inv := testInvoiceB2C(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Customer = nil
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "7", inv.Tax.Ext[arca.ExtKeyDocType].String())
	})

	t.Run("credit note type B - B2C", func(t *testing.T) {
		inv := testInvoiceB2C(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Customer = nil
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "8", inv.Tax.Ext[arca.ExtKeyDocType].String())
	})

	t.Run("invoice type C - simplified regime", func(t *testing.T) {
		inv := testInvoiceTypeC(t)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "11", inv.Tax.Ext[arca.ExtKeyDocType].String())
	})

	t.Run("debit note type C - simplified regime", func(t *testing.T) {
		inv := testInvoiceTypeC(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "12", inv.Tax.Ext[arca.ExtKeyDocType].String())
	})

	t.Run("credit note type C - simplified regime", func(t *testing.T) {
		inv := testInvoiceTypeC(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "13", inv.Tax.Ext[arca.ExtKeyDocType].String())
	})

	t.Run("type C takes precedence over B2B", func(t *testing.T) {
		// Even with an AR customer, simplified-regime tag should result in type C
		inv := testInvoiceTypeC(t)
		inv.Customer = &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "AR",
				Code:    "30500010912",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "11", inv.Tax.Ext[arca.ExtKeyDocType].String())
	})

	t.Run("type C takes precedence over B2C", func(t *testing.T) {
		// Even without a customer, simplified-regime tag should result in type C
		inv := testInvoiceTypeC(t)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "11", inv.Tax.Ext[arca.ExtKeyDocType].String())
	})

	t.Run("type A for monotributo customer", func(t *testing.T) {
		inv := testInvoiceB2B(t)
		// Set customer as Monotributo (VAT status 6) - should be Type A
		inv.Customer.Ext = tax.Extensions{
			arca.ExtKeyVATStatus: "6",
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "1", inv.Tax.Ext[arca.ExtKeyDocType].String())
	})

	t.Run("type B for exempt customer", func(t *testing.T) {
		inv := testInvoiceB2B(t)
		// Set customer as VAT Exempt (VAT status 4) - should be Type B
		inv.Customer.Ext = tax.Extensions{
			arca.ExtKeyVATStatus: "4",
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "6", inv.Tax.Ext[arca.ExtKeyDocType].String())
	})

	t.Run("type B for VAT exempt customer with AR tax ID", func(t *testing.T) {
		inv := testInvoiceB2B(t)
		// Set customer as VAT Exempt (VAT status 4) - should be Type B
		// They have an AR tax ID but are exempt from VAT
		inv.Customer.Ext = tax.Extensions{
			arca.ExtKeyVATStatus: "4",
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "6", inv.Tax.Ext[arca.ExtKeyDocType].String())
	})
}

func TestScenarioSummary(t *testing.T) {
	ad := tax.AddonForKey(arca.V4)

	t.Run("all scenarios have required fields", func(t *testing.T) {
		for _, ss := range ad.Scenarios {
			require.NotEmpty(t, ss.Schema, "scenario set should have schema")
			for i, s := range ss.List {
				assert.NotEmpty(t, s.Name, "scenario %d should have name", i)
				assert.NotEmpty(t, s.Ext, "scenario %d should have extensions", i)
				assert.NotEmpty(t, s.Ext[arca.ExtKeyDocType], "scenario %d should have doc type", i)
			}
		}
	})

	t.Run("scenarios cover all main document types", func(t *testing.T) {
		expectedDocTypes := map[cbc.Code]bool{
			"1":  false, // Invoice A
			"2":  false, // Debit Note A
			"3":  false, // Credit Note A
			"6":  false, // Invoice B
			"7":  false, // Debit Note B
			"8":  false, // Credit Note B
			"11": false, // Invoice C
			"12": false, // Debit Note C
			"13": false, // Credit Note C
		}

		for _, ss := range ad.Scenarios {
			for _, s := range ss.List {
				docType := s.Ext[arca.ExtKeyDocType]
				if _, ok := expectedDocTypes[docType]; ok {
					expectedDocTypes[docType] = true
				}
			}
		}

		for docType, found := range expectedDocTypes {
			assert.True(t, found, "document type %s should be covered by scenarios", docType)
		}
	})
}

// testInvoiceB2B creates an invoice with an Argentine customer (B2B scenario)
func testInvoiceB2B(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := testInvoiceWithGoods(t)
	// Clear any existing tags and doc type to let scenarios determine type
	inv.Tags = tax.Tags{}
	delete(inv.Tax.Ext, arca.ExtKeyDocType)
	return inv
}

// testInvoiceB2C creates an invoice for B2C scenario (no customer or foreign)
func testInvoiceB2C(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := testInvoiceWithGoods(t)
	// Clear any existing tags and doc type to let scenarios determine type
	inv.Tags = tax.Tags{}
	delete(inv.Tax.Ext, arca.ExtKeyDocType)
	return inv
}
