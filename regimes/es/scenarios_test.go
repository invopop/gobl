package es_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	_ "github.com/invopop/gobl/regimes"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	i := &bill.Invoice{
		Code:     "123TEST",
		Currency: "EUR",
		Tax: &bill.Tax{
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

func testInvoiceSimplified(t *testing.T) *bill.Invoice {
	t.Helper()
	i := &bill.Invoice{
		Currency: "EUR",
		Code:     "123TEST",
		Tax: &bill.Tax{
			Tags: []cbc.Key{tax.TagSimplified},
		},
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
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

func TestInvoiceDocumentScenarios(t *testing.T) {
	i := testInvoiceStandard(t)
	require.NoError(t, i.Calculate())
	assert.Len(t, i.Notes, 0)

	// TODO: refactor this to have the scenarios add extensions,
	// or perform these checks in the conversion module.
	ss := i.ScenarioSummary() //nolint:staticcheck
	assert.Contains(t, ss.Codes, es.KeyFacturaEInvoiceDocumentType)
	assert.Equal(t, ss.Codes[es.KeyFacturaEInvoiceDocumentType], cbc.Code("FC"))

	i = testInvoiceStandard(t)
	i.Tax.Tags = []cbc.Key{es.TagTravelAgency}
	require.NoError(t, i.Calculate())
	assert.Len(t, i.Notes, 1)
	assert.Equal(t, i.Notes[0].Src, es.TagTravelAgency)
	assert.Equal(t, i.Notes[0].Text, "Régimen especial de las agencias de viajes.")

	i = testInvoiceSimplified(t)
	require.NoError(t, i.Calculate())
	require.NoError(t, i.Validate())
}
