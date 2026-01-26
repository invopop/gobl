package favat_test

import (
	"testing"

	"github.com/invopop/gobl/addons/pl/favat"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceScenarios(t *testing.T) {
	tests := []struct {
		name             string
		invoiceType      cbc.Key
		tags             []cbc.Key
		expectedExtValue string
	}{
		{
			name:             "regular invoice",
			invoiceType:      bill.InvoiceTypeStandard,
			tags:             []cbc.Key{},
			expectedExtValue: "VAT",
		},
		{
			name:             "prepayment invoice",
			invoiceType:      bill.InvoiceTypeStandard,
			tags:             []cbc.Key{tax.TagPartial},
			expectedExtValue: "ZAL",
		},
		{
			name:             "settlement invoice",
			invoiceType:      bill.InvoiceTypeStandard,
			tags:             []cbc.Key{favat.TagSettlement},
			expectedExtValue: "ROZ",
		},
		{
			name:             "simplified invoice",
			invoiceType:      bill.InvoiceTypeStandard,
			tags:             []cbc.Key{tax.TagSimplified},
			expectedExtValue: "UPR",
		},
		{
			name:             "credit note",
			invoiceType:      bill.InvoiceTypeCreditNote,
			tags:             []cbc.Key{},
			expectedExtValue: "KOR",
		},
		{
			name:             "prepayment credit note",
			invoiceType:      bill.InvoiceTypeCreditNote,
			tags:             []cbc.Key{tax.TagPartial},
			expectedExtValue: "KOR_ZAL",
		},
		{
			name:             "settlement credit note",
			invoiceType:      bill.InvoiceTypeCreditNote,
			tags:             []cbc.Key{favat.TagSettlement},
			expectedExtValue: "KOR_ROZ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := &bill.Invoice{
				Regime:    tax.WithRegime("PL"),
				Addons:    tax.WithAddons(favat.V3),
				Currency:  currency.PLN,
				Code:      "TEST",
				Type:      tt.invoiceType,
				IssueDate: cal.MakeDate(2022, 12, 29),
				Supplier: &org.Party{
					Name: "Test Party",
					TaxID: &tax.Identity{
						Country: "PL",
						Code:    "1111111111",
					},
					Addresses: []*org.Address{
						{
							Street:   "ul. Testowa 1",
							Locality: "Warsaw",
							Country:  "PL",
						},
					},
				},
				Customer: &org.Party{
					Name: "Test Customer",
					TaxID: &tax.Identity{
						Country: "PL",
						Code:    "2222222222",
					},
				},
				Lines: []*bill.Line{
					{
						Quantity: num.MakeAmount(1, 3),
						Item: &org.Item{
							Name:  "test-item",
							Price: num.NewAmount(1000, 3),
						},
					},
				},
			}

			if tt.invoiceType == bill.InvoiceTypeCreditNote {
				inv.Preceding = []*org.DocumentRef{
					{
						Code:      "ORIG",
						IssueDate: cal.NewDate(2022, 12, 27),
					},
				}
			}

			if len(tt.tags) > 0 {
				inv.SetTags(tt.tags...)
			}

			require.NoError(t, inv.Calculate())
			require.NoError(t, inv.Validate())

			assert.Equal(t, tt.expectedExtValue, inv.Tax.Ext.Get(favat.ExtKeyInvoiceType).String(),
				"invoice type extension should be set correctly")
		})
	}
}

func TestSettlementTag(t *testing.T) {
	t.Run("settlement tag is available", func(t *testing.T) {
		inv := &bill.Invoice{
			Regime:    tax.WithRegime("PL"),
			Addons:    tax.WithAddons(favat.V3),
			Currency:  currency.PLN,
			Code:      "TEST",
			Type:      bill.InvoiceTypeStandard,
			IssueDate: cal.MakeDate(2022, 12, 29),
			Supplier: &org.Party{
				Name: "Test Party",
				TaxID: &tax.Identity{
					Country: "PL",
					Code:    "1111111111",
				},
				Addresses: []*org.Address{
					{
						Street:   "ul. Testowa 1",
						Locality: "Warsaw",
						Country:  "PL",
					},
				},
			},
			Customer: &org.Party{
				Name: "Test Customer",
				TaxID: &tax.Identity{
					Country: "PL",
					Code:    "2222222222",
				},
			},
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(1, 3),
					Item: &org.Item{
						Name:  "test-item",
						Price: num.NewAmount(1000, 3),
					},
				},
			},
		}

		inv.SetTags(favat.TagSettlement)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())

		assert.True(t, inv.HasTags(favat.TagSettlement))
		assert.Equal(t, "ROZ", inv.Tax.Ext.Get(favat.ExtKeyInvoiceType).String())
	})
}

func TestScenariosSummary(t *testing.T) {
	t.Run("addon has scenarios configured", func(t *testing.T) {
		ad := tax.AddonForKey(favat.V3)
		require.NotNil(t, ad)
		require.NotNil(t, ad.Scenarios)
		assert.Greater(t, len(ad.Scenarios), 0, "should have at least one scenario set")

		// Check that invoice scenarios exist
		var invoiceScenarios *tax.ScenarioSet
		for _, ss := range ad.Scenarios {
			if ss.Schema == bill.ShortSchemaInvoice {
				invoiceScenarios = ss
				break
			}
		}
		require.NotNil(t, invoiceScenarios, "should have invoice scenarios")
		assert.Equal(t, 7, len(invoiceScenarios.List), "should have 7 invoice scenarios")
	})
}
