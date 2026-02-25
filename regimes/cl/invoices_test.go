package cl_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/cl"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Code:     "TEST-001",
		Currency: "CLP",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "CL",
				Code:    "713254975",
			},
			Addresses: []*org.Address{
				{
					Street:   "Av. Libertador Bernardo O'Higgins",
					Number:   "1234",
					Locality: "Santiago",
					Region:   "Región Metropolitana",
					Code:     "8320000",
					Country:  "CL",
				},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "CL",
				Code:    "77668208K", // Valid customer RUT
			},
			Addresses: []*org.Address{
				{
					Street:   "Av. Providencia",
					Number:   "567",
					Locality: "Providencia",
					Region:   "Región Metropolitana",
					Code:     "7500000",
					Country:  "CL",
				},
			},
		},
		IssueDate: cal.MakeDate(2024, 1, 1),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(10000, 0),
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
	t.Parallel()

	tests := []struct {
		name        string
		setupInv    func(*bill.Invoice)
		expectError string
	}{
		{
			name:     "valid invoice",
			setupInv: func(_ *bill.Invoice) {},
		},
		{
			name: "missing supplier",
			setupInv: func(inv *bill.Invoice) {
				inv.SetRegime("CL")
				inv.Supplier = nil
			},
			expectError: "supplier: cannot be blank",
		},
		{
			name: "supplier without tax ID",
			setupInv: func(inv *bill.Invoice) {
				inv.SetRegime("CL")
				inv.Supplier.TaxID = nil
			},
			expectError: "tax_id: cannot be blank",
		},
		{
			name: "supplier with tax ID but no code",
			setupInv: func(inv *bill.Invoice) {
				inv.SetRegime("CL")
				inv.Supplier.TaxID.Code = ""
			},
			expectError: "code: cannot be blank",
		},
		{
			name: "invalid RUT check digit",
			setupInv: func(inv *bill.Invoice) {
				inv.SetRegime("CL")
				inv.Supplier.TaxID.Code = "713254976" // Invalid check digit
			},
			expectError: "invalid RUT check digit",
		},
		{
			name: "supplier without address",
			setupInv: func(inv *bill.Invoice) {
				inv.SetRegime("CL")
				inv.Supplier.Addresses = nil
			},
			expectError: "addresses",
		},
		{
			name: "supplier with empty addresses array",
			setupInv: func(inv *bill.Invoice) {
				inv.SetRegime("CL")
				inv.Supplier.Addresses = []*org.Address{}
			},
			expectError: "addresses",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			inv := testInvoiceStandard(t)
			tt.setupInv(inv)
			require.NoError(t, inv.Calculate())

			err := inv.Validate()
			if tt.expectError == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.expectError)
				}
			}
		})
	}

	// Separate test for regime-specific validation via cl.Validate
	t.Run("regime-specific validation via cl.Validate", func(t *testing.T) {
		t.Parallel()
		inv := &bill.Invoice{
			Supplier: &org.Party{
				Name: "Test Supplier",
				TaxID: &tax.Identity{
					Country: "CL",
					Code:    "713254975",
				},
				Addresses: []*org.Address{
					{
						Street:   "Test Street",
						Locality: "Santiago",
					},
				},
			},
		}
		err := cl.Validate(inv)
		assert.NoError(t, err)
	})
}

func TestInvoiceCustomerValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupInv    func(*bill.Invoice)
		expectError string
	}{
		{
			name: "B2B invoice - customer without RUT",
			setupInv: func(inv *bill.Invoice) {
				inv.SetRegime("CL")
				inv.Customer.TaxID = nil
			},
			expectError: "tax_id: cannot be blank",
		},
		{
			name: "B2B invoice - customer with RUT but no code",
			setupInv: func(inv *bill.Invoice) {
				inv.SetRegime("CL")
				inv.Customer.TaxID = &tax.Identity{
					Country: "CL",
				}
			},
			expectError: "code: cannot be blank",
		},
		{
			name: "B2B invoice - customer with valid RUT",
			setupInv: func(inv *bill.Invoice) {
				inv.SetRegime("CL")
				inv.Customer.TaxID = &tax.Identity{
					Country: "CL",
					Code:    "77668208K",
				}
			},
		},
		{
			name: "B2B invoice - customer with invalid RUT check digit",
			setupInv: func(inv *bill.Invoice) {
				inv.SetRegime("CL")
				inv.Customer.TaxID = &tax.Identity{
					Country: "CL",
					Code:    "77668208X", // Invalid check digit
				}
			},
			expectError: "invalid RUT",
		},
		{
			name: "B2B invoice - customer without address",
			setupInv: func(inv *bill.Invoice) {
				inv.SetRegime("CL")
				inv.Customer.Addresses = nil
			},
			expectError: "addresses",
		},
		{
			name: "B2B invoice - customer with empty addresses array",
			setupInv: func(inv *bill.Invoice) {
				inv.SetRegime("CL")
				inv.Customer.Addresses = []*org.Address{}
			},
			expectError: "addresses",
		},
		{
			name: "B2C boleta - no customer required",
			setupInv: func(inv *bill.Invoice) {
				inv.SetRegime("CL")
				inv.SetTags(tax.TagSimplified)
				inv.Customer = nil
			},
		},
		{
			name: "B2C boleta - customer without RUT is valid",
			setupInv: func(inv *bill.Invoice) {
				inv.SetRegime("CL")
				inv.SetTags(tax.TagSimplified)
				inv.Customer.TaxID = nil
			},
		},
		{
			name: "B2C boleta - customer with RUT is also valid",
			setupInv: func(inv *bill.Invoice) {
				inv.SetRegime("CL")
				inv.SetTags(tax.TagSimplified)
				inv.Customer.TaxID = &tax.Identity{
					Country: "CL",
					Code:    "77668208K",
				}
			},
		},
		{
			name: "B2C boleta - customer without address is valid",
			setupInv: func(inv *bill.Invoice) {
				inv.SetRegime("CL")
				inv.SetTags(tax.TagSimplified)
				inv.Customer.Addresses = nil
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			inv := testInvoiceStandard(t)
			tt.setupInv(inv)
			require.NoError(t, inv.Calculate())

			err := inv.Validate()
			if tt.expectError == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.expectError)
				}
			}
		})
	}
}

func TestInvoiceNilSafety(t *testing.T) {
	t.Parallel()

	t.Run("nil invoice should not panic", func(t *testing.T) {
		t.Parallel()
		err := cl.Validate((*bill.Invoice)(nil))
		assert.NoError(t, err)
	})
}
