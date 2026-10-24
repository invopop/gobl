package cl_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/regimes/cl"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("should create a new Chile regime", func(t *testing.T) {
		regime := cl.New()
		require.NotNil(t, regime)

		assert.Equal(t, "CL", regime.Country.String())
		assert.Equal(t, currency.CLP, regime.Currency)
		assert.Equal(t, "America/Santiago", regime.TimeZone)
		assert.Equal(t, tax.CategoryVAT, regime.TaxScheme)
		assert.Equal(t, "Chile", regime.Name["en"])
		assert.NotEmpty(t, regime.Sources)
		require.Len(t, regime.Corrections, 1)
		assert.Equal(t, bill.ShortSchemaInvoice, regime.Corrections[0].Schema)
		assert.Contains(t, regime.Corrections[0].Types, bill.InvoiceTypeCreditNote)
		assert.Contains(t, regime.Corrections[0].Types, bill.InvoiceTypeDebitNote)
	})
}
