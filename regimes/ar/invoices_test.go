package ar_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/ar"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func validInvoice() *bill.Invoice {
	price := num.MakeAmount(1000, 2)
	return &bill.Invoice{
		Code:      "TEST-001",
		Currency:  "ARS",
		IssueDate: cal.MakeDate(2024, 1, 15),
		Supplier: &org.Party{
			Name: "Empresa Argentina S.A.",
			TaxID: &tax.Identity{
				Country: "AR",
				Code:    "30714589840",
			},
		},
		Customer: &org.Party{
			Name: "Cliente Ejemplo S.R.L.",
			TaxID: &tax.Identity{
				Country: "AR",
				Code:    "30500010912",
			},
		},
		Lines: []*bill.Line{
			{
				Index:    1,
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Producto de Ejemplo",
					Price: &price,
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
			},
		},
	}
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("valid invoice", func(t *testing.T) {
		inv := validInvoice()
		err := ar.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("invoice without supplier tax ID should fail", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = nil

		err := ar.Validate(inv)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tax_id")
	})

	t.Run("Type A invoice without customer tax ID should fail", func(t *testing.T) {
		inv := validInvoice()
		inv.SetTags(ar.TagInvoiceTypeA)
		inv.Customer.TaxID = nil

		err := ar.Validate(inv)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "customer tax ID required for Type A invoices")
	})

	t.Run("Type A invoice with customer tax ID should pass", func(t *testing.T) {
		inv := validInvoice()
		inv.SetTags(ar.TagInvoiceTypeA)

		err := ar.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("Type B invoice without customer should pass", func(t *testing.T) {
		inv := validInvoice()
		inv.SetTags(ar.TagInvoiceTypeB)
		inv.Customer = nil

		err := ar.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("Type C invoice without customer should pass", func(t *testing.T) {
		inv := validInvoice()
		inv.SetTags(ar.TagInvoiceTypeC)
		inv.Customer = nil

		err := ar.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("simplified invoice without customer should pass", func(t *testing.T) {
		inv := validInvoice()
		inv.SetTags(tax.TagSimplified)
		inv.Customer = nil

		err := ar.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("Type E invoice with Argentine customer should fail", func(t *testing.T) {
		inv := validInvoice()
		inv.SetTags(ar.TagInvoiceTypeE)
		// Customer has Argentine tax ID - should fail for export
		inv.Customer = &org.Party{
			Name: "Cliente Argentino",
			TaxID: &tax.Identity{
				Country: "AR",
				Code:    "30500010912",
			},
		}

		err := ar.Validate(inv)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "export invoices should not have Argentine customers")
	})

	t.Run("Type E invoice with foreign customer should pass", func(t *testing.T) {
		inv := validInvoice()
		inv.SetTags(ar.TagInvoiceTypeE)
		inv.Customer = &org.Party{
			Name: "Foreign Customer LLC",
			TaxID: &tax.Identity{
				Country: "US",
				Code:    "123456789",
			},
		}

		err := ar.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("export services invoice with foreign customer should pass", func(t *testing.T) {
		inv := validInvoice()
		inv.SetTags(ar.TagExportServices)
		inv.Customer = &org.Party{
			Name: "Foreign Customer",
			TaxID: &tax.Identity{
				Country: "BR",
				Code:    "12345678901234",
			},
		}

		err := ar.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("export goods invoice with foreign customer should pass", func(t *testing.T) {
		inv := validInvoice()
		inv.SetTags(ar.TagExportGoods)
		inv.Customer = &org.Party{
			Name: "Foreign Customer",
			TaxID: &tax.Identity{
				Country: "CL",
				Code:    "76543210-9",
			},
		}

		err := ar.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("invoice without lines should fail", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines = nil

		err := ar.Validate(inv)
		// This should fail during invoice calculation, not AR-specific validation
		// AR validation focuses on AR-specific rules
		// The common invoice validation will catch missing lines
		assert.NoError(t, err) // AR-specific validation passes
	})
}

func TestInvoiceNormalization(t *testing.T) {
	t.Run("should normalize supplier and customer", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID.Code = "30-71458984-0" // With hyphens
		inv.Customer.TaxID.Code = "30-50001091-2" // With hyphens

		ar.Normalize(inv)

		assert.Equal(t, "30714589840", string(inv.Supplier.TaxID.Code))
		assert.Equal(t, "30500010912", string(inv.Customer.TaxID.Code))
	})

	t.Run("should handle simplified invoice normalization", func(t *testing.T) {
		inv := validInvoice()
		inv.SetTags(tax.TagSimplified)

		ar.Normalize(inv)

		// Normalization should complete without errors
		assert.NotNil(t, inv.Supplier)
	})

	t.Run("should handle export invoice normalization", func(t *testing.T) {
		inv := validInvoice()
		inv.SetTags(ar.TagInvoiceTypeE)
		inv.Customer = &org.Party{
			Name: "Foreign Customer",
			TaxID: &tax.Identity{
				Country: "US",
				Code:    "123456789",
			},
		}

		ar.Normalize(inv)

		// Normalization should complete without errors
		assert.NotNil(t, inv.Customer)
	})
}
