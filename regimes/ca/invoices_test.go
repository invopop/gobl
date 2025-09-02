package ca_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/ca"
	"github.com/stretchr/testify/assert"
)

func TestInvoiceValidation(t *testing.T) {
	t.Run("should pass with a valid invoice", func(t *testing.T) {
		inv := &bill.Invoice{
			Supplier: &org.Party{
				Name: "Test Supplier",
			},
		}
		err := ca.New().Validator(inv)
		assert.NoError(t, err)
	})

	t.Run("should fail without a supplier", func(t *testing.T) {
		inv := &bill.Invoice{
			Customer: &org.Party{
				Name: "Test Customer",
			},
		}
		err := ca.New().Validator(inv)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "supplier: cannot be blank")
	})
}
