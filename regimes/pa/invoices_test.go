package pa_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/pa"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/require"
)

func TestInvoiceValidation(t *testing.T) {
	tests := []struct {
		name string
		inv  *bill.Invoice
		err  string
	}{
		{
			name: "valid standard invoice",
			inv: &bill.Invoice{
				Supplier: &org.Party{
					TaxID: &tax.Identity{
						Country: "PA",
						Code:    "8-442-445-90",
					},
				},
				Customer: &org.Party{
					TaxID: &tax.Identity{
						Country: "PA",
						Code:    "155596713-2-2015-59",
					},
				},
			},
		},
		{
			name: "valid simplified invoice without customer",
			inv: &bill.Invoice{
				Tags: tax.WithTags(tax.TagSimplified),
				Supplier: &org.Party{
					TaxID: &tax.Identity{
						Country: "PA",
						Code:    "8-442-445-90",
					},
				},
			},
		},
		{
			name: "valid simplified invoice with customer",
			inv: &bill.Invoice{
				Tags: tax.WithTags(tax.TagSimplified),
				Supplier: &org.Party{
					TaxID: &tax.Identity{
						Country: "PA",
						Code:    "8-442-445-90",
					},
				},
				Customer: &org.Party{
					TaxID: &tax.Identity{
						Country: "PA",
						Code:    "CIP-000-000-0000",
					},
				},
			},
		},
		{
			name: "missing supplier TaxID",
			inv: &bill.Invoice{
				Supplier: &org.Party{
					Name:  "Test",
					TaxID: nil,
				},
				Customer: &org.Party{
					TaxID: &tax.Identity{
						Country: "PA",
						Code:    "155596713-2-2015-59",
					},
				},
			},
			err: "tax_id: cannot be blank",
		},
		{
			name: "supplier TaxID without code",
			inv: &bill.Invoice{
				Supplier: &org.Party{
					TaxID: &tax.Identity{
						Country: "PA",
						Code:    "",
					},
				},
				Customer: &org.Party{
					TaxID: &tax.Identity{
						Country: "PA",
						Code:    "155596713-2-2015-59",
					},
				},
			},
			err: "code: cannot be blank",
		},
		// Supplier presence is enforced by bill.Invoice core validation, not the regime.
		{
			name: "nil supplier passes",
			inv: &bill.Invoice{
				Supplier: nil,
				Customer: &org.Party{
					TaxID: &tax.Identity{
						Country: "PA",
						Code:    "155596713-2-2015-59",
					},
				},
			},
		},
		{
			name: "standard invoice missing customer",
			inv: &bill.Invoice{
				Supplier: &org.Party{
					TaxID: &tax.Identity{
						Country: "PA",
						Code:    "8-442-445-90",
					},
				},
				Customer: nil,
			},
			err: "customer: cannot be blank",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pa.Validate(tt.inv)
			if tt.err == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tt.err)
		})
	}
}
