package pl_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/pl"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestInvoiceValidation(t *testing.T) {
	t.Run("valid invoice with supplier and valid TaxID", func(t *testing.T) {
		inv := &bill.Invoice{
			Series: "TEST",
			Code:   "0001",
			Supplier: &org.Party{
				Name: "Test Supplier",
				TaxID: &tax.Identity{
					Country: "PL",
					Code:    "5260300788", // Valid Polish NIP
				},
			},
		}
		err := pl.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("invoice with nil supplier", func(t *testing.T) {
		inv := &bill.Invoice{
			Series:   "TEST",
			Code:     "0002",
			Supplier: nil,
		}
		err := pl.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("invoice with supplier but nil TaxID", func(t *testing.T) {
		inv := &bill.Invoice{
			Series: "TEST",
			Code:   "0003",
			Supplier: &org.Party{
				Name:  "Test Supplier",
				TaxID: nil,
			},
		}
		err := pl.Validate(inv)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "tax_id: cannot be blank")
		}
	})

	t.Run("invoice with supplier and TaxID but empty Code", func(t *testing.T) {
		inv := &bill.Invoice{
			Series: "TEST",
			Code:   "0004",
			Supplier: &org.Party{
				Name: "Test Supplier",
				TaxID: &tax.Identity{
					Country: "PL",
					Code:    "", // Empty code
				},
			},
		}
		err := pl.Validate(inv)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "code: cannot be blank")
		}
	})
}
