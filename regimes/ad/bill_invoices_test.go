package ad_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	adregime "github.com/invopop/gobl/regimes/ad"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/require"
)

// validInvoice returns a minimal but complete Andorran invoice
// that passes all validation rules. Individual tests mutate this
// to test specific failure cases.
func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime:   tax.WithRegime(adregime.CountryCode),
		Currency: "EUR",
		Code:     "0001",
		Supplier: &org.Party{
			Name: "Acme Andorra SL",
			TaxID: &tax.Identity{
				Country: "AD",
				Code:    "L132950X",
			},
			Addresses: []*org.Address{
				{Street: "Carrer Major 1", Locality: "Andorra la Vella", Country: "AD"},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			Addresses: []*org.Address{
				{Street: "Carrer Prat 2", Locality: "Escaldes", Country: "AD"},
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Consulting service",
					Price: num.NewAmount(100, 0), // ← NewAmount returns *num.Amount directly
				},
				Taxes: tax.Set{
					{Category: tax.CategoryVAT, Rate: tax.RateGeneral},
				},
			},
		},
	}
}

func TestValidInvoice(t *testing.T) {
	inv := validInvoice()
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
}

func TestValidInvoiceWithNRTIdentity(t *testing.T) {
	// Supplier below IGI threshold — has NRT org identity but no tax ID code
	inv := validInvoice()
	inv.Supplier.TaxID = nil
	inv.Supplier.Identities = []*org.Identity{
		{Type: adregime.IdentityTypeNRT, Code: "F123456A"},
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
}

func TestMissingSupplierNRT(t *testing.T) {
	inv := validInvoice()
	inv.Supplier.TaxID = nil // no tax ID and no org identity
	require.NoError(t, inv.Calculate())
	require.ErrorContains(t, rules.Validate(inv),
		"invoice supplier in Andorra must have an NRT tax ID code or an NRT identity")
}

func TestSupplierTaxIDWithoutCode(t *testing.T) {
	inv := validInvoice()
	inv.Supplier.TaxID = &tax.Identity{Country: "AD"} // present but no code
	require.NoError(t, inv.Calculate())
	require.ErrorContains(t, rules.Validate(inv),
		"invoice supplier in Andorra must have an NRT tax ID code or an NRT identity")
}