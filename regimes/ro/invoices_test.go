package ro_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Series: "TEST",
		Code:   "0002",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "RO",
				Code:    "14399840",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "RO",
				Code:    "18547290",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.NewAmount(10000, 2),
					Unit:  org.UnitPackage,
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "general",
					},
				},
			},
		},
	}
}

func simplifiedInvoice() *bill.Invoice {
	return &bill.Invoice{
		Series: "TEST",
		Code:   "0003",
		Tags: tax.Tags{
			List: []cbc.Key{tax.TagSimplified},
		},
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "RO",
				Code:    "14399840",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.NewAmount(5000, 2),
					Unit:  org.UnitPackage,
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "general",
					},
				},
			},
		},
	}
}

func TestInvoiceValidation(t *testing.T) {
	tests := []struct {
		name string
		inv  *bill.Invoice
		err  string
	}{
		{
			name: "valid standard invoice",
			inv:  validInvoice(),
		},
		{
			name: "missing supplier tax ID code",
			inv: func() *bill.Invoice {
				inv := validInvoice()
				inv.Supplier.TaxID.Code = ""
				return inv
			}(),
			err: "supplier: (tax_id: (code: cannot be blank.).)",
		},
		{
			name: "missing customer",
			inv: func() *bill.Invoice {
				inv := validInvoice()
				inv.Customer = nil
				return inv
			}(),
			err: "customer: cannot be blank.",
		},
		{
			name: "simplified invoice without customer",
			inv:  simplifiedInvoice(),
		},
		{
			name: "simplified invoice with customer",
			inv: func() *bill.Invoice {
				inv := simplifiedInvoice()
				inv.Customer = &org.Party{
					Name: "Test Customer",
				}
				return inv
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, tt.inv.Calculate())
			if tt.err == "" {
				assert.NoError(t, tt.inv.Validate())
			} else {
				assert.ErrorContains(t, tt.inv.Validate(), tt.err)
			}
		})
	}
}
