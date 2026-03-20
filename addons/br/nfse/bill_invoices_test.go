package nfse_test

import (
	"testing"

	"github.com/invopop/gobl/addons/br/nfse"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func withAddonContext() rules.WithContext {
	return func(rc *rules.Context) {
		rc.Set(rules.ContextKey(nfse.V1), tax.AddonForKey(nfse.V1))
	}
}

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
			err:  "series is required",
		},
		{
			name: "invalid code (non-digits)",
			inv: &bill.Invoice{
				Code: "ABC-123",
			},
			err: "code must be a positive integer",
		},
		{
			name: "invalid code (padding zeroes)",
			inv: &bill.Invoice{
				Code: "000123",
			},
			err: "code must be a positive integer",
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
			err: "not supported by NFS-e",
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
			err: "not supported by NFS-e",
		},
		{
			name: "series missing",
			inv:  &bill.Invoice{},
			err:  "series is required",
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			err := rules.Validate(ts.inv, withAddonContext())
			if ts.err != "" {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), ts.err)
				}
			}
		})
	}
}

func TestSuppliersValidation(t *testing.T) {
	t.Run("validates supplier", func(t *testing.T) {
		sup := new(org.Party)
		inv := &bill.Invoice{
			Supplier: sup,
		}
		err := rules.Validate(inv, withAddonContext())
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "supplier name is required")
		}

		sup.Name = "Test"
		err = rules.Validate(inv, withAddonContext())
		if assert.Error(t, err) {
			assert.NotContains(t, err.Error(), "supplier name is required")
		}
	})

	t.Run("validates tax ID", func(t *testing.T) {
		sup := new(org.Party)
		inv := &bill.Invoice{
			Supplier: sup,
		}
		err := rules.Validate(inv, withAddonContext())
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "supplier tax ID is required")
		}

		sup.TaxID = new(tax.Identity)
		err = rules.Validate(inv, withAddonContext())
		if assert.Error(t, err) {
			assert.NotContains(t, err.Error(), "supplier tax ID is required")
			assert.Contains(t, err.Error(), "supplier tax ID code is required")
		}

		sup.TaxID.Code = "123"
		err = rules.Validate(inv, withAddonContext())
		if assert.Error(t, err) {
			assert.NotContains(t, err.Error(), "supplier tax ID code is required")
		}
	})

	t.Run("validates addresses", func(t *testing.T) {
		sup := new(org.Party)
		inv := &bill.Invoice{
			Supplier: sup,
		}
		err := rules.Validate(inv, withAddonContext())
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "supplier must have at least one address")
		}

		sup.Addresses = []*org.Address{nil}
		err = rules.Validate(inv, withAddonContext())
		if assert.Error(t, err) {
			assert.NotContains(t, err.Error(), "supplier must have at least one address")
			assert.Contains(t, err.Error(), "supplier address must not be empty")
		}

		sup.Addresses[0] = new(org.Address)
		err = rules.Validate(inv, withAddonContext())
		if assert.Error(t, err) {
			assert.NotContains(t, err.Error(), "supplier address must not be empty")
			assert.Contains(t, err.Error(), "supplier address requires a street")
			assert.Contains(t, err.Error(), "supplier address requires a number")
			assert.Contains(t, err.Error(), "supplier address requires a locality")
			assert.Contains(t, err.Error(), "supplier address requires a state")
			assert.Contains(t, err.Error(), "supplier address requires a postal code")
		}

		sup.Addresses[0] = &org.Address{
			Code:     "12345678",
			Locality: "Test",
			Number:   "123",
			State:    "RJ",
			Street:   "Test",
		}
		err = rules.Validate(inv, withAddonContext())
		if assert.Error(t, err) {
			assert.NotContains(t, err.Error(), "supplier address requires a street")
			assert.NotContains(t, err.Error(), "supplier address requires a number")
			assert.NotContains(t, err.Error(), "supplier address requires a locality")
			assert.NotContains(t, err.Error(), "supplier address requires a state")
			assert.NotContains(t, err.Error(), "supplier address requires a postal code")
		}
	})

	t.Run("validates extensions", func(t *testing.T) {
		sup := new(org.Party)
		inv := &bill.Invoice{
			Supplier: sup,
		}
		err := rules.Validate(inv, withAddonContext())
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "supplier requires 'br-ibge-municipality', 'br-nfse-simples', and 'br-nfse-fiscal-incentive' extensions")
		}

		sup.Ext = tax.Extensions{
			nfse.ExtKeySimples:         "1",
			"br-ibge-municipality":     "12345678",
			nfse.ExtKeyFiscalIncentive: "2",
		}
		err = rules.Validate(inv, withAddonContext())
		if assert.Error(t, err) {
			assert.NotContains(t, err.Error(), "supplier requires 'br-ibge-municipality', 'br-nfse-simples', and 'br-nfse-fiscal-incentive' extensions")
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
