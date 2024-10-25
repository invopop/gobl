package nfse_test

import (
	"testing"

	"github.com/invopop/gobl/addons/br/nfse"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestInvoicesValidation(t *testing.T) {
	tests := []struct {
		name string
		inv  *bill.Invoice
		err  string
	}{
		{
			name: "valid invoice",
			inv:  &bill.Invoice{},
		},
		{
			name: "charges present",
			inv: &bill.Invoice{
				Charges: []*bill.Charge{
					{
						Amount: num.MakeAmount(100, 2),
					},
				},
			},
			err: "charges: not supported by nfse.",
		},
		{
			name: "discounts present",
			inv: &bill.Invoice{
				Discounts: []*bill.Discount{
					{
						Amount: num.MakeAmount(100, 2),
					},
				},
			},
			err: "discounts: not supported by nfse.",
		},
	}

	addon := tax.AddonForKey(nfse.V1)
	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			err := addon.Validator(ts.inv)
			if ts.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), ts.err)
				}
			}
		})
	}
}
