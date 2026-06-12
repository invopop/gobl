package lu_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/lu"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validLUInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime: tax.WithRegime("LU"),
		Series: "2024",
		Code:   "0001",
		Supplier: &org.Party{
			Name: "Acme Luxembourg S.A.",
			TaxID: &tax.Identity{
				Country: "LU",
				// 263752 mod 89 = 45  →  valid
				Code: "26375245",
			},
		},
		Customer: &org.Party{
			Name: "Client International B.V.",
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Software consulting",
					Price: num.NewAmount(100000, 2),
					Unit:  org.UnitService,
				},
				Taxes: tax.Set{
					{Category: "VAT", Rate: "standard"},
				},
			},
		},
	}
}

func TestInvoiceLUValidation(t *testing.T) {
	t.Run("valid invoice with TVA code", func(t *testing.T) {
		inv := validLUInvoice()
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("valid invoice with RCS identity instead of TVA", func(t *testing.T) {
		inv := validLUInvoice()
		inv.Supplier.TaxID = nil
		inv.Supplier.Identities = []*org.Identity{
			{Type: lu.IdentityTypeRCS, Code: "B263475"},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("invalid when supplier has neither TVA nor RCS", func(t *testing.T) {
		inv := validLUInvoice()
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-LU-BILL-INVOICE-01]")
	})

	t.Run("invalid when supplier TVA code is empty and no RCS", func(t *testing.T) {
		inv := validLUInvoice()
		inv.Supplier.TaxID.Code = ""
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-LU-BILL-INVOICE-01]")
	})
}
