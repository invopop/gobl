package nfse_test

import (
	"testing"

	"github.com/invopop/gobl/addons/br/nfse"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
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
			inv: &bill.Invoice{
				Series: "SAMPLE",
			},
		},
		{
			name: "nil invoice",
			inv:  nil,
		},
		{
			name: "missing series",
			inv:  &bill.Invoice{},
			err:  "series: cannot be blank",
		},
		{
			name: "invalid code (non-digits)",
			inv: &bill.Invoice{
				Code: "ABC-123",
			},
			err: "code: must be in a valid format",
		},
		{
			name: "invalid code (padding zeroes)",
			inv: &bill.Invoice{
				Code: "000123",
			},
			err: "code: must be in a valid format",
		},
		{
			name: "valid code",
			inv: &bill.Invoice{
				Series: "SAMPLE",
				Code:   "123000",
			},
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
			err: "charges: not supported by nfse",
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
			err: "discounts: not supported by nfse",
		},
		{
			name: "series missing",
			inv:  &bill.Invoice{},
			err:  "series: cannot be blank",
		},
	}

	addon := tax.AddonForKey(nfse.V1)
	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			err := addon.Validator(ts.inv)
			if ts.err == "" {
				if err != nil {
					assert.NotContains(t, err.Error(), ts.err)
				}
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
			assert.Contains(t, err.Error(), "identities: missing key 'br-nfse-municipal-reg';")
		}

		sup.Identities = append(sup.Identities, &org.Identity{
			Key:  nfse.IdentityKeyMunicipalReg,
			Code: "12345678",
		})
		err = addon.Validator(inv)
		if assert.Error(t, err) {
			assert.NotContains(t, err.Error(), "identities: missing key 'br-nfse-municipal-reg';")
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

	t.Run("validates extensions", func(t *testing.T) {
		sup := new(org.Party)
		inv := &bill.Invoice{
			Supplier: sup,
		}
		err := addon.Validator(inv)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "br-nfse-simples: required")
			assert.Contains(t, err.Error(), "br-nfse-municipality: required")
			assert.Contains(t, err.Error(), "br-nfse-fiscal-incentive: required")
		}

		sup.Ext = tax.Extensions{
			nfse.ExtKeySimples:         "1",
			nfse.ExtKeyMunicipality:    "12345678",
			nfse.ExtKeyFiscalIncentive: "2",
		}
		err = addon.Validator(inv)
		if assert.Error(t, err) {
			assert.NotContains(t, err.Error(), "br-nfse-simples: required")
			assert.NotContains(t, err.Error(), "br-nfse-municipality: required")
			assert.NotContains(t, err.Error(), "br-nfse-fiscal-incentive: required")
		}
	})
}

func TestSuppliersNormalization(t *testing.T) {
	addon := tax.AddonForKey(nfse.V1)

	tests := []struct {
		name     string
		supplier *org.Party
		out      cbc.Code
	}{
		{
			name:     "no supplier",
			supplier: nil,
		},
		{
			name:     "sets default fiscal incentive",
			supplier: &org.Party{},
			out:      "2",
		},
		{
			name: "does not override fiscal incentive",
			supplier: &org.Party{
				Ext: tax.Extensions{
					nfse.ExtKeyFiscalIncentive: "1",
				},
			},
			out: "1",
		},
	}
	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			inv := &bill.Invoice{Supplier: ts.supplier}
			addon.Normalizer(inv)
			if ts.supplier == nil {
				assert.Nil(t, inv.Supplier)
			} else {
				assert.Equal(t, ts.out, inv.Supplier.Ext[nfse.ExtKeyFiscalIncentive])
			}
		})
	}
}
