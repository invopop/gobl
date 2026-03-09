package sa_test

import (
	"testing"

	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceValidation(t *testing.T) {
	t.Run("valid standard invoice", func(t *testing.T) {
		i := testStandardInvoice(t)
		require.NoError(t, i.Calculate())
		assert.NoError(t, i.Validate())
	})

	t.Run("valid simplified invoice without customer", func(t *testing.T) {
		i := testSimplifiedInvoice(t)
		require.NoError(t, i.Calculate())
		assert.NoError(t, i.Validate())
	})

	t.Run("standard invoice missing supplier tax ID code", func(t *testing.T) {
		i := testStandardInvoice(t)
		i.Supplier.TaxID.Code = ""
		require.NoError(t, i.Calculate())
		err := i.Validate()
		assert.ErrorContains(t, err, "supplier")
		assert.ErrorContains(t, err, "tax_id")
	})

	t.Run("simplified invoice missing supplier tax ID code", func(t *testing.T) {
		i := testSimplifiedInvoice(t)
		i.Supplier.TaxID.Code = ""
		require.NoError(t, i.Calculate())
		err := i.Validate()
		assert.ErrorContains(t, err, "supplier")
		assert.ErrorContains(t, err, "tax_id")
	})

	t.Run("standard invoice missing customer", func(t *testing.T) {
		i := testStandardInvoice(t)
		i.Customer = nil
		require.NoError(t, i.Calculate())
		err := i.Validate()
		assert.ErrorContains(t, err, "customer")
	})

	t.Run("simplified invoice allows no customer", func(t *testing.T) {
		i := testSimplifiedInvoice(t)
		i.Customer = nil
		require.NoError(t, i.Calculate())
		assert.NoError(t, i.Validate())
	})

	t.Run("reverse charge requires customer tax ID", func(t *testing.T) {
		i := testStandardInvoice(t)
		i.Tags = tax.WithTags(tax.TagReverseCharge)
		require.NoError(t, i.Calculate())
		assert.NoError(t, i.Validate())

		// Remove customer tax ID code
		i = testStandardInvoice(t)
		i.Tags = tax.WithTags(tax.TagReverseCharge)
		i.Customer.TaxID.Code = ""
		require.NoError(t, i.Calculate())
		err := i.Validate()
		assert.ErrorContains(t, err, "customer")
		assert.ErrorContains(t, err, "tax_id")
	})

	t.Run("reverse charge without customer tax ID object", func(t *testing.T) {
		i := testStandardInvoice(t)
		i.Tags = tax.WithTags(tax.TagReverseCharge)
		i.Customer.TaxID = nil
		require.NoError(t, i.Calculate())
		err := i.Validate()
		assert.ErrorContains(t, err, "customer")
		assert.ErrorContains(t, err, "tax_id")
	})

	t.Run("standard invoice without reverse charge does not require customer tax ID", func(t *testing.T) {
		i := testStandardInvoice(t)
		i.Customer.TaxID = nil
		require.NoError(t, i.Calculate())
		assert.NoError(t, i.Validate())
	})

	t.Run("zero-rated line item", func(t *testing.T) {
		i := testStandardInvoice(t)
		i.Lines[0].Taxes[0].Rate = tax.RateZero
		require.NoError(t, i.Calculate())
		assert.NoError(t, i.Validate())
		assert.Equal(t, tax.KeyZero, i.Lines[0].Taxes[0].Key)
	})

	t.Run("exempt line item", func(t *testing.T) {
		i := testStandardInvoice(t)
		i.Lines[0].Taxes[0].Rate = ""
		i.Lines[0].Taxes[0].Key = tax.KeyExempt
		require.NoError(t, i.Calculate())
		assert.NoError(t, i.Validate())
		assert.Equal(t, tax.KeyExempt, i.Lines[0].Taxes[0].Key)
	})
}
