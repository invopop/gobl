package choruspro_test

import (
	"testing"

	"github.com/invopop/gobl/addons/fr/choruspro"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/require"
)

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := &bill.Invoice{
		Regime:   tax.WithRegime("FR"),
		Addons:   tax.WithAddons(choruspro.V1),
		Code:     "123TEST",
		Currency: "EUR",
		Tax: &bill.Tax{
			PricesInclude: tax.CategoryVAT,
		},
		Type: bill.InvoiceTypeStandard,
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "732829320",
			},
			Addresses: []*org.Address{
				{
					Street:   "Via di Test",
					Code:     "12345",
					Locality: "Paris",
					Country:  "FR",
					Number:   "3",
				},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "44391838042",
			},
			Addresses: []*org.Address{
				{
					Street:   "Piazza di Test",
					Code:     "38342",
					Locality: "Paris",
					Country:  "FR",
					Number:   "1",
				},
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Payment: &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCard,
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(10000, 2),
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
	return inv
}

func TestInvoicePartyIdentities(t *testing.T) {
	t.Run("With SIREN and normal tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("With SIRET and normal tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID.Code = "39183804200000"
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("With SIRET and SIREN", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID.Code = "39183804200000"
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("With no SIREN or SIRET", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID.Code = "44732829320"
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("With extension set", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID.Code = "44732829320"
		inv.Supplier.Identities = []*org.Identity{
			{
				Type: fr.IdentityTypeSiren,
				Code: cbc.Code(inv.Supplier.TaxID.Code.String()[2:]),
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("With no identities", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Supplier.Identities = nil
		err := inv.Validate()
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot be blank")
	})

	t.Run("With invalid identity", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Supplier.Identities = []*org.Identity{
			{
				Type: "INVALID",
				Code: cbc.Code("1234567890"),
			},
		}
		err := inv.Validate()
		require.Error(t, err)
		require.Contains(t, err.Error(), "at least one identity must be SIREN or SIRET")
	})
}

func TestInvoicePayment(t *testing.T) {
	t.Run("With payment", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("With no payment", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = nil
		err := inv.Validate()
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot be blank")
	})

	t.Run("With no payment instructions", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment.Instructions = nil
		err := inv.Validate()
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot be blank")
	})
}
