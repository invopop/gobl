package ro_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/ro"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestInvoiceValidation(t *testing.T) {
	t.Run("valid supplier with CUI", func(t *testing.T) {
		inv := &bill.Invoice{
			Supplier: &org.Party{
				Name: "Romanian Tech SRL",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    "13547272",
				},
			},
		}
		assert.NoError(t, ro.Validate(inv))
	})

	t.Run("nil supplier is allowed", func(t *testing.T) {
		inv := &bill.Invoice{Supplier: nil}
		assert.NoError(t, ro.Validate(inv))
	})

	t.Run("supplier without tax ID is rejected", func(t *testing.T) {
		inv := &bill.Invoice{
			Supplier: &org.Party{
				Name:  "Romanian Tech SRL",
				TaxID: nil,
			},
		}
		err := ro.Validate(inv)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "tax_id: cannot be blank")
		}
	})

	t.Run("supplier with tax ID but no code is rejected", func(t *testing.T) {
		inv := &bill.Invoice{
			Supplier: &org.Party{
				Name: "Romanian Tech SRL",
				TaxID: &tax.Identity{
					Country: "RO",
					Code:    "",
				},
			},
		}
		err := ro.Validate(inv)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "code: cannot be blank")
		}
	})
}
