package fr_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
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
		inv.SetRegime("FR")
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-FR-BILL-INVOICE-01]")
	})

	t.Run("empty supplier tax ID code", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID.Code = ""
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-FR-BILL-INVOICE-01]")
	})

	t.Run("valid with SIREN identity instead of tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID = nil
		inv.Supplier.Identities = []*org.Identity{
			{
				Type: fr.IdentityTypeSIREN,
				Code: "732829320",
			},
		}
		inv.SetRegime("FR")
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("valid with SIRET identity instead of tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID = nil
		inv.Supplier.Identities = []*org.Identity{
			{
				Type: fr.IdentityTypeSIRET,
				Code: "73282932000015",
			},
		}
		inv.SetRegime("FR")
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})
}

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Currency:  "EUR",
		IssueDate: cal.MakeDate(2024, 2, 13),
		Code:      "INV-001",
		Supplier: &org.Party{
			Name: "Fournisseur Exemple SARL",
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "44732829320",
			},
		},
		Customer: &org.Party{
			Name: "Client Exemple SAS",
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "44391838042",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Consultation informatique",
					Price: num.NewAmount(150, 0),
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
