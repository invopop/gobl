package arca_test

import (
	"testing"

	"github.com/invopop/gobl/addons/ar/arca"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
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

	// Test cases for all document type scenarios
	testCases := []struct {
		name        string
		invType     cbc.Key
		tags        []cbc.Key
		expectedDoc string
	}{
		// Type A - Standard (B2B)
		{
			name:        "standard invoice type A",
			invType:     bill.InvoiceTypeStandard,
			tags:        nil,
			expectedDoc: "001",
		},
		{
			name:        "debit note type A",
			invType:     bill.InvoiceTypeDebitNote,
			tags:        nil,
			expectedDoc: "002",
		},
		{
			name:        "credit note type A",
			invType:     bill.InvoiceTypeCreditNote,
			tags:        nil,
			expectedDoc: "003",
		},
		// Type B - Simplified (B2C)
		{
			name:        "simplified invoice type B",
			invType:     bill.InvoiceTypeStandard,
			tags:        []cbc.Key{tax.TagSimplified},
			expectedDoc: "006",
		},
		{
			name:        "simplified debit note type B",
			invType:     bill.InvoiceTypeDebitNote,
			tags:        []cbc.Key{tax.TagSimplified},
			expectedDoc: "007",
		},
		{
			name:        "simplified credit note type B",
			invType:     bill.InvoiceTypeCreditNote,
			tags:        []cbc.Key{tax.TagSimplified},
			expectedDoc: "008",
		},
		// Export invoices
		{
			name:        "export invoice",
			invType:     bill.InvoiceTypeStandard,
			tags:        []cbc.Key{tax.TagExport},
			expectedDoc: "019",
		},
		{
			name:        "export debit note",
			invType:     bill.InvoiceTypeDebitNote,
			tags:        []cbc.Key{tax.TagExport},
			expectedDoc: "020",
		},
		{
			name:        "export credit note",
			invType:     bill.InvoiceTypeCreditNote,
			tags:        []cbc.Key{tax.TagExport},
			expectedDoc: "021",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inv := testInvoiceForScenario(t, tc.invType, tc.tags)
			require.NoError(t, inv.Calculate())
			assert.Equal(t, tc.expectedDoc, inv.Tax.Ext[arca.ExtKeyDocType].String())
		})
	}
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
			"001": false, // Invoice A
			"002": false, // Debit Note A
			"003": false, // Credit Note A
			"006": false, // Invoice B
			"007": false, // Debit Note B
			"008": false, // Credit Note B
			"019": false, // Export Invoice
			"020": false, // Export Debit Note
			"021": false, // Export Credit Note
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

func testInvoiceForScenario(t *testing.T, invType cbc.Key, tags []cbc.Key) *bill.Invoice {
	t.Helper()
	inv := testInvoiceWithGoods(t)
	inv.Type = invType
	if len(tags) > 0 {
		inv.SetTags(tags...)
	}
	// Clear existing doc type to let scenarios set it
	delete(inv.Tax.Ext, arca.ExtKeyDocType)
	return inv
}
