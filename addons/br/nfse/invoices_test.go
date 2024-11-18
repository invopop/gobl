package nfse_test

import (
	"testing"

	"github.com/invopop/gobl/addons/br/nfse"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
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
			name: "nil invoice",
			inv:  nil,
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

func TestSuppliersValidation(t *testing.T) {
	addon := tax.AddonForKey(nfse.V1)

	t.Run("validates supplier", func(t *testing.T) {
		sup := new(org.Party)
		inv := &bill.Invoice{
			Supplier: sup,
		}
		err := addon.Validator(inv)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "name: cannot be blank")
		}

		sup.Name = "Test"
		err = addon.Validator(inv)
		if assert.Error(t, err) {
			assert.NotContains(t, err.Error(), "name: cannot be blank")
		}
	})

	t.Run("validates tax ID", func(t *testing.T) {
		sup := new(org.Party)
		inv := &bill.Invoice{
			Supplier: sup,
		}
		err := addon.Validator(inv)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "tax_id: cannot be blank")
		}

		sup.TaxID = new(tax.Identity)
		err = addon.Validator(inv)
		if assert.Error(t, err) {
			assert.NotContains(t, err.Error(), "tax_id: cannot be blank")
			assert.Contains(t, err.Error(), "tax_id: (code: cannot be blank")
		}

		sup.TaxID.Code = "123"
		err = addon.Validator(inv)
		if assert.Error(t, err) {
			assert.NotContains(t, err.Error(), "tax_id: (code: cannot be blank")
		}
	})

	t.Run("validates identities", func(t *testing.T) {
		sup := new(org.Party)
		inv := &bill.Invoice{
			Supplier: sup,
		}
		err := addon.Validator(inv)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "identities: missing key br-nfse-municipal-reg;")
		}

		sup.Identities = append(sup.Identities, &org.Identity{
			Key:  nfse.IdentityKeyMunicipalReg,
			Code: "12345678",
		})
		err = addon.Validator(inv)
		if assert.Error(t, err) {
			assert.NotContains(t, err.Error(), "identities: missing key br-nfse-municipal-reg;")
		}
	})

	t.Run("validates addresses", func(t *testing.T) {
		sup := new(org.Party)
		inv := &bill.Invoice{
			Supplier: sup,
		}
		err := addon.Validator(inv)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "addresses: cannot be blank")
		}

		sup.Addresses = []*org.Address{nil}
		err = addon.Validator(inv)
		if assert.Error(t, err) {
			assert.NotContains(t, err.Error(), "addresses: cannot be blank")
			assert.Contains(t, err.Error(), "addresses: (0: cannot be blank.)")
		}

		sup.Addresses[0] = new(org.Address)
		err = addon.Validator(inv)
		if assert.Error(t, err) {
			assert.NotContains(t, err.Error(), "addresses: (0: cannot be blank.)")
			assert.Contains(t, err.Error(), "addresses: (0: (code: cannot be blank; locality: cannot be blank; num: cannot be blank; state: cannot be blank; street: cannot be blank.).)")
		}

		sup.Addresses[0] = &org.Address{
			Code:     "12345678",
			Locality: "Test",
			Number:   "123",
			State:    "RJ",
			Street:   "Test",
		}
		err = addon.Validator(inv)
		if assert.Error(t, err) {
			assert.NotContains(t, err.Error(), "addresses: (0:")
		}
	})
}
