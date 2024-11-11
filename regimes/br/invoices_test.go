package br_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/br"
	"github.com/stretchr/testify/assert"
)

func TestInvoiceValidation(t *testing.T) {
	t.Run("supplier required", func(t *testing.T) {
		inv := new(bill.Invoice)
		err := br.Validate(inv)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "supplier: cannot be blank")
		}
	})

	t.Run("valid invoice", func(t *testing.T) {
		inv := &bill.Invoice{
			Supplier: &org.Party{
				Name: "Test Supplier",
			},
		}
		err := br.Validate(inv)
		assert.NoError(t, err)
	})
}
