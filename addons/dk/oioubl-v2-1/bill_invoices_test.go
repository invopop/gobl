package oioubl_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	oioubl "github.com/invopop/gobl/addons/dk/oioubl-v2-1"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Regime:    tax.WithRegime("DK"),
		Addons:    tax.WithAddons(oioubl.V2_1),
		IssueDate: cal.MakeDate(2026, 1, 1),
		Type:      "standard",
		Currency:  "DKK",
		Series:    "2026",
		Code:      "1000",
		Supplier: &org.Party{
			Name: "Eksempel A/S",
			TaxID: &tax.Identity{
				Country: "DK",
				Code:    "12345674",
			},
			Inboxes: []*org.Inbox{
				{Scheme: "0184", Code: "12345674"},
			},
			Addresses: []*org.Address{
				{Street: "Hovedgaden 1", Locality: "København", Code: "1000", Country: "DK"},
			},
		},
		Customer: &org.Party{
			Name: "Kunde ApS",
			TaxID: &tax.Identity{
				Country: "DK",
				Code:    "88146328",
			},
			Inboxes: []*org.Inbox{
				{Scheme: "0184", Code: "88146328"},
			},
			People: []*org.Person{
				{Name: &org.Name{Given: "Anders", Surname: "Jensen"}},
			},
			Addresses: []*org.Address{
				{Street: "Bygaden 5", Locality: "Aarhus", Code: "8000", Country: "DK"},
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Produkt",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{Category: "VAT", Rate: "standard"},
				},
			},
		},
	}
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("standard invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("missing supplier inboxes (F-INV031)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Inboxes = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "F-INV031")
	})

	t.Run("missing customer inboxes (F-INV044)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Inboxes = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "F-INV044")
	})

	t.Run("missing customer people (F-INV046)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.People = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "F-INV046")
	})

	t.Run("customer with two people is allowed (loose vs F-INV046)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.People = append(inv.Customer.People,
			&org.Person{Name: &org.Name{Given: "Mette", Surname: "Hansen"}},
		)
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("ordering absent is allowed", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Ordering = nil
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("ordering present without code fails (F-INV024)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Ordering = &bill.Ordering{}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "F-INV024")
	})

	t.Run("ordering present with code passes", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Ordering = &bill.Ordering{Code: "PO-2026-001"}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})
}
