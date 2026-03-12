package es_test

import (
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceValidation(t *testing.T) {
	t.Run("standard invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("missing supplier tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetRegime("ES")
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.Equal(t, "[GOBL-ES-BILL-INVOICE-02] supplier.tax_id: supplier tax ID is required", err.Error())
	})

}
