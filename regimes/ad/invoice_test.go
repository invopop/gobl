package ad_test

import (
	"strings"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceB2B(t *testing.T) {
	inv := testInvoiceStandard(t)
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	t.Run("calculates IGI at 4.5%", func(t *testing.T) {
		cat := inv.Totals.Taxes.Category(tax.CategoryVAT)
		require.NotNil(t, cat)
		assert.Equal(t, "4.5%", cat.Rates[0].Percent.String())
	})
}

func TestInvoiceB2C(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Tax = &bill.Tax{PricesInclude: tax.CategoryVAT}
	inv.Customer = &org.Party{
		Name: "Pere Martí",
		Addresses: []*org.Address{
			{
				Street:   "Carrer de la Vall 3",
				Locality: "Ordino",
				Code:     "AD300",
				Country:  "AD",
			},
		},
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	t.Run("prices include IGI", func(t *testing.T) {
		assert.Equal(t, tax.CategoryVAT, inv.Tax.PricesInclude)
	})
}

func TestInvoiceDiplomaticScenario(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.SetTags("diplomatic")
	inv.Lines[0].Taxes[0] = &tax.Combo{
		Category: tax.CategoryVAT,
		Rate:     tax.RateZero,
	}
	require.NoError(t, inv.Calculate())

	t.Run("adds diplomatic legal note", func(t *testing.T) {
		found := false
		for _, n := range inv.Notes {
			if strings.Contains(n.Text, "Art. 15") {
				found = true
			}
		}
		assert.True(t, found, "diplomatic legal note not found")
	})
}

func TestInvoiceExportScenario(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.SetTags(tax.TagExport)
	require.NoError(t, inv.Calculate())

	t.Run("adds export legal note", func(t *testing.T) {
		found := false
		for _, n := range inv.Notes {
			if strings.Contains(n.Text, "980-A") {
				found = true
			}
		}
		assert.True(t, found, "export legal note not found")
	})
}

func TestInvoiceCreditNote(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Type = bill.InvoiceTypeCreditNote
	inv.Preceding = []*org.DocumentRef{
		{
			Type:      bill.InvoiceTypeStandard,
			IssueDate: cal.NewDate(2025, 8, 21),
			Series:    "SAMPLE",
			Code:      "001",
		},
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	t.Run("credit note type accepted", func(t *testing.T) {
		assert.Equal(t, bill.InvoiceTypeCreditNote, inv.Type)
	})
}

func TestInvoiceReducedRate(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Lines[0].Taxes[0] = &tax.Combo{
		Category: tax.CategoryVAT,
		Rate:     tax.RateReduced,
	}
	require.NoError(t, inv.Calculate())

	t.Run("applies 1% reduced rate", func(t *testing.T) {
		cat := inv.Totals.Taxes.Category(tax.CategoryVAT)
		require.NotNil(t, cat)
		assert.Equal(t, "1%", cat.Rates[0].Percent.String())
	})
}

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Regime:    tax.WithRegime("AD"),
		Type:      bill.InvoiceTypeStandard,
		Series:    "SAMPLE",
		Code:      "001",
		IssueDate: cal.MakeDate(2025, 8, 21),
		Currency:  "EUR",
		Supplier: &org.Party{
			Name: "Fusta i Disseny S.L.",
			TaxID: &tax.Identity{
				Country: "AD",
				Code:    "L123456A",
			},
			Addresses: []*org.Address{
				{
					Street:   "Carrer Major 12",
					Locality: "Andorra la Vella",
					Code:     "AD500",
					Country:  "AD",
				},
			},
		},
		Customer: &org.Party{
			Name: "Maquinaria Pirineus S.L.",
			TaxID: &tax.Identity{
				Country: "AD",
				Code:    "F654321Z",
			},
			Addresses: []*org.Address{
				{
					Street:   "Avinguda Meritxell 45",
					Locality: "Escaldes-Engordany",
					Code:     "AD700",
					Country:  "AD",
				},
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Taula de fusta de roure",
					Price: num.NewAmount(250, 0),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     tax.RateGeneral,
					},
				},
			},
		},
	}
}
