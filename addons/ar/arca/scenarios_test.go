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
		// Type A - VAT registered customer
		{
			name:        "invoice type A",
			invType:     bill.InvoiceTypeStandard,
			tags:        []cbc.Key{arca.TagVATRegistered},
			expectedDoc: "1",
		},
		{
			name:        "debit note type A",
			invType:     bill.InvoiceTypeDebitNote,
			tags:        []cbc.Key{arca.TagVATRegistered},
			expectedDoc: "2",
		},
		{
			name:        "credit note type A",
			invType:     bill.InvoiceTypeCreditNote,
			tags:        []cbc.Key{arca.TagVATRegistered},
			expectedDoc: "3",
		},
		// Type B - Final consumers and VAT exempt
		{
			name:        "invoice type B",
			invType:     bill.InvoiceTypeStandard,
			tags:        []cbc.Key{tax.TagSimplified},
			expectedDoc: "6",
		},
		{
			name:        "debit note type B",
			invType:     bill.InvoiceTypeDebitNote,
			tags:        []cbc.Key{tax.TagSimplified},
			expectedDoc: "7",
		},
		{
			name:        "credit note type B",
			invType:     bill.InvoiceTypeCreditNote,
			tags:        []cbc.Key{tax.TagSimplified},
			expectedDoc: "8",
		},
		// Type C - Simplified Regime (Monotributista)
		{
			name:        "invoice type C",
			invType:     bill.InvoiceTypeStandard,
			tags:        []cbc.Key{arca.TagSimplifiedRegime},
			expectedDoc: "11",
		},
		{
			name:        "debit note type C",
			invType:     bill.InvoiceTypeDebitNote,
			tags:        []cbc.Key{arca.TagSimplifiedRegime},
			expectedDoc: "12",
		},
		{
			name:        "credit note type C",
			invType:     bill.InvoiceTypeCreditNote,
			tags:        []cbc.Key{arca.TagSimplifiedRegime},
			expectedDoc: "13",
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
