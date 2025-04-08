package cpp_test

import (
	"testing"

	"github.com/invopop/gobl/addons/fr/cpp"
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
		Addons:   tax.WithAddons(cpp.V1),
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

func TestInvoiceValidation(t *testing.T) {
	// test with siren and normal tax id
	inv := testInvoiceStandard(t)
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	// test with siret and normal tax id
	inv = testInvoiceStandard(t)
	inv.Supplier.TaxID.Code = "39183804200000"
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	// test with siret and siren
	inv = testInvoiceStandard(t)
	inv.Customer.TaxID.Code = "39183804200000"
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	// test with no siren or siret
	inv = testInvoiceStandard(t)
	inv.Supplier.TaxID.Code = "44732829320"
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	// test with extension set
	inv = testInvoiceStandard(t)
	inv.Supplier.TaxID.Code = "44732829320"
	inv.Supplier.Identities = []*org.Identity{
		{
			Key:  fr.IdentityKeySiren,
			Code: cbc.Code(inv.Supplier.TaxID.Code.String()[2:]),
		},
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestInvoiceNormalization(t *testing.T) {
	inv := testInvoiceStandard(t)
	ad := tax.AddonForKey(cpp.V1)
	ad.Normalizer(inv)
}
