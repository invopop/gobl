package sa_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/sa"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime: tax.WithRegime("SA"),
		Code:   "TEST-001",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "SA",
				Code:    "300000000000003",
			},
			Identities: []*org.Identity{
				{Type: sa.IdentityTypeCRN, Code: "1234567890"},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "SA",
				Code:    "300000000100003",
			},
		},
		Currency:  "SAR",
		IssueDate: cal.MakeDate(2024, 1, 1),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(10000, 2),
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

// --- Rule 01 (BR-KSA-39): supplier must have valid tax ID ---
// --- Rule 02 (BR-KSA-08): supplier must have exactly 1 valid identity ---

func TestSupplierTaxIDRequired(t *testing.T) {
	t.Run("valid invoice passes", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("supplier without tax ID fails", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier must have a valid tax ID code")
	})

	t.Run("supplier with empty tax ID code fails", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = &tax.Identity{Country: "SA", Code: ""}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier must have a valid tax ID code")
	})
}

func TestSupplierIdentities(t *testing.T) {
	t.Run("supplier with no identities fails", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Identities = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier must have a valid identity")
	})

	t.Run("supplier with one valid CRN identity passes", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Identities = []*org.Identity{
			{Type: sa.IdentityTypeCRN, Code: "1234567890"},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("supplier with one valid MOM identity passes", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Identities = []*org.Identity{
			{Type: sa.IdentityTypeMom, Code: "1234567890"},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("supplier with one valid MLS identity passes", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Identities = []*org.Identity{
			{Type: sa.IdentityTypeMLS, Code: "1234567890"},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("supplier with one valid 700 identity passes", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Identities = []*org.Identity{
			{Type: sa.IdentityType700, Code: "1234567890"},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("supplier with one valid SAG identity passes", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Identities = []*org.Identity{
			{Type: sa.IdentityTypeSAG, Code: "1234567890"},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("supplier with one valid OTH identity passes", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Identities = []*org.Identity{
			{Type: sa.IdentityTypeOTH, Code: "1234567890"},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("supplier with invalid identity type fails", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Identities = []*org.Identity{
			{Type: "INVALID", Code: "1234567890"},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier must have a valid identity")
	})

	t.Run("supplier with two identities fails", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Identities = []*org.Identity{
			{Type: sa.IdentityTypeCRN, Code: "1234567890"},
			{Type: sa.IdentityTypeMLS, Code: "1234567890"},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier must have a valid identity")
	})

	t.Run("supplier with NAT identity type fails (not in supplier list)", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Identities = []*org.Identity{
			{Type: sa.IdentityTypeNational, Code: "1234567890"},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier must have a valid identity")
	})

	t.Run("supplier with TIN identity type fails (not in supplier list)", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Identities = []*org.Identity{
			{Type: sa.IdentityTypeTIN, Code: "123456789012345"},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier must have a valid identity")
	})
}
