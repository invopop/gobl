package sa_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func testStandardInvoice(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Code:     "INV-001",
		Currency: "SAR",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "SA",
				Code:    "300075588700003",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "SA",
				Code:    "310122393500003",
			},
		},
		IssueDate: cal.MakeDate(2024, 1, 15),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(5, 0),
				Item: &org.Item{
					Name:  "Service Item",
					Price: num.NewAmount(200, 0),
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

func testSimplifiedInvoice(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Code:     "INV-002",
		Currency: "SAR",
		Tags:     tax.WithTags(tax.TagSimplified),
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "SA",
				Code:    "300075588700003",
			},
		},
		IssueDate: cal.MakeDate(2024, 1, 15),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(3, 0),
				Item: &org.Item{
					Name:  "Product Item",
					Price: num.NewAmount(100, 0),
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
