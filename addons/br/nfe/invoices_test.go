package nfe_test

import (
	"fmt"
	"testing"

	"github.com/invopop/gobl/addons/br/nfe"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestInvoicesValidation(t *testing.T) {
	addon := tax.AddonForKey(nfe.V4)

	t.Run("validates tax extensions", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = nil
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "tax: cannot be blank")

		inv.Tax = &bill.Tax{}
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "br-nfe-model: required")
		assert.ErrorContains(t, err, "br-nfe-presence: required")

		inv.Tax.Ext = tax.Extensions{
			nfe.ExtKeyModel:    nfe.ModelNFe,
			nfe.ExtKeyPresence: nfe.PresenceDelivery,
		}
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "br-nfe-presence: value '4' not allowed")

		inv.Tax.Ext[nfe.ExtKeyPresence] = nfe.PresenceInPerson
		err = addon.Validator(inv)
		assert.NoError(t, err)
	})

	t.Run("validates required notes", func(t *testing.T) {
		inv := validInvoice()
		inv.Notes = nil
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "notes: note with key `reason` required")

		inv.Notes = []*org.Note{nil}
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "notes: note with key `reason` required")

		inv.Notes[0] = &org.Note{
			Key:  org.NoteKeyGeneral,
			Text: "General note",
		}
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "notes: note with key `reason` required")

		inv.Notes[0].Key = org.NoteKeyReason
		inv.Notes[0].Text = "1234567890123456789012345678901234567890123456789012345678901" // 61 chars
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "notes: (0: (text: the length must be between 1 and 60")

		inv.Notes[0].Text = "123456789012345678901234567890123456789012345678901234567890" // 60 chars
		err = addon.Validator(inv)
		assert.NoError(t, err)
	})

	t.Run("validates payment when invoice is due", func(t *testing.T) {
		inv := validInvoice()
		inv.Totals = &bill.Totals{}
		inv.Payment = nil

		inv.Totals.Due = &num.AmountZero
		err := addon.Validator(inv)
		assert.NoError(t, err)

		inv.Totals.Due = nil
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "payment: cannot be blank")

		inv.Totals.Due = num.NewAmount(1, 2)
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "payment: cannot be blank")

		inv.Payment = &bill.PaymentDetails{}
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "instructions: cannot be blank")

		inv.Payment.Instructions = &pay.Instructions{}
		err = addon.Validator(inv)
		assert.NoError(t, err)
	})

	t.Run("validates invoice totals due field", func(t *testing.T) {
		inv := validInvoice()
		inv.Totals = &bill.Totals{}

		inv.Totals.Due = num.NewAmount(-1, 2)
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "due: must be no less than 0.")

		inv.Totals.Due = &num.AmountZero
		err = addon.Validator(inv)
		assert.NoError(t, err)

		inv.Totals.Due = num.NewAmount(1, 2)
		err = addon.Validator(inv)
		assert.NoError(t, err)
	})

	t.Run("validates NFe presence when model is NFe", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax.Ext[nfe.ExtKeyModel] = nfe.ModelNFe
		inv.Tax.Ext[nfe.ExtKeyPresence] = nfe.PresenceDelivery
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "br-nfe-presence: value '4' not allowed")

		inv.Tax.Ext[nfe.ExtKeyPresence] = nfe.PresenceInPerson
		err = addon.Validator(inv)
		assert.NoError(t, err)
	})

	t.Run("validates NFCe presence when model is NFCe", func(t *testing.T) {
		inv := validInvoice()
		inv.Customer = nil // For NFCe, customer is optional, so remove it to avoid other validation errors

		inv.Tax.Ext[nfe.ExtKeyModel] = nfe.ModelNFCe
		inv.Tax.Ext[nfe.ExtKeyPresence] = nfe.PresenceNotApplicable
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "br-nfe-presence: invalid value")

		inv.Tax.Ext[nfe.ExtKeyPresence] = nfe.PresenceInPerson
		err = addon.Validator(inv)
		assert.NoError(t, err)
	})
}

func TestInvoiceSeriesValidation(t *testing.T) {
	addon := tax.AddonForKey(nfe.V4)

	tests := []struct {
		series cbc.Code
		err    string
	}{
		{series: "0"},
		{series: "1"},
		{series: "12"},
		{series: "123"},
		{series: "999"},
		{series: "", err: "series: cannot be blank"},
		{series: "1000", err: "series: must be in a valid format"},
		{series: "abc", err: "series: must be in a valid format"},
		{series: "012", err: "series: must be in a valid format"},
		{series: "00", err: "series: must be in a valid format"},
		{series: "-3", err: "series: must be in a valid format"},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("validates series %s", tt.series)
		t.Run(name, func(t *testing.T) {
			inv := validInvoice()
			inv.Series = tt.series
			err := addon.Validator(inv)
			if tt.err != "" {
				assert.ErrorContains(t, err, tt.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSupplierValidation(t *testing.T) {
	addon := tax.AddonForKey(nfe.V4)

	t.Run("nil supplier", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier = nil
		err := addon.Validator(inv)
		assert.NoError(t, err) // supplier presence is validated at GOBL level
	})

	t.Run("validates supplier name", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Name = ""
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "name: cannot be blank")

		inv.Supplier.Name = "Test Company"
		err = addon.Validator(inv)
		if err != nil {
			assert.NotContains(t, err.Error(), "name: cannot be blank")
		}
	})

	t.Run("validates supplier addresses required", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Addresses = nil
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "addresses: cannot be blank")

		inv.Supplier.Addresses = []*org.Address{}
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "addresses: cannot be blank")

		inv.Supplier.Addresses = []*org.Address{nil}
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "addresses: (0: cannot be blank")

		inv.Supplier.Addresses = []*org.Address{
			{
				Street:   "Rua Test",
				Number:   "100",
				Locality: "São Paulo",
				State:    "SP",
				Code:     "01310100",
			},
		}
		err = addon.Validator(inv)
		if err != nil {
			assert.NotContains(t, err.Error(), "addresses: cannot be blank")
		}
	})

	t.Run("validates supplier state registration identity", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Identities = nil
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "identities: missing key 'br-nfe-state-reg'")

		inv.Supplier.Identities = []*org.Identity{
			{
				Key:  nfe.IdentityKeyStateReg,
				Code: "35503304557308",
			},
		}
		err = addon.Validator(inv)
		if err != nil {
			assert.NotContains(t, err.Error(), "identities: missing key 'br-nfse-state-reg'")
		}
	})

	t.Run("validates supplier tax ID required", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = nil
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "tax_id: cannot be blank")

		inv.Supplier.TaxID = &tax.Identity{}
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "tax_id: (code: cannot be blank")

		inv.Supplier.TaxID.Code = "55263640000186"
		err = addon.Validator(inv)
		if err != nil {
			assert.NotContains(t, err.Error(), "tax_id: (code: cannot be blank")
		}
	})

	t.Run("validates supplier municipality extension when addresses exist", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Ext = nil
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "br-ibge-municipality: required")

		inv.Supplier.Ext = tax.Extensions{
			"br-ibge-municipality": "3304557",
		}
		err = addon.Validator(inv)
		if err != nil {
			assert.NotContains(t, err.Error(), "br-ibge-municipality: required")
		}
	})

	t.Run("validates supplier address fields", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Addresses = []*org.Address{
			{
				Street:   "",
				Number:   "100",
				Locality: "São Paulo",
				State:    "SP",
				Code:     "01310100",
			},
		}
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "street: cannot be blank")

		inv.Supplier.Addresses[0].Street = "Rua Test"
		inv.Supplier.Addresses[0].Number = ""
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "num: cannot be blank")

		inv.Supplier.Addresses[0].Number = "100"
		inv.Supplier.Addresses[0].Locality = ""
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "locality: cannot be blank")

		inv.Supplier.Addresses[0].Locality = "São Paulo"
		inv.Supplier.Addresses[0].State = ""
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "state: cannot be blank")

		inv.Supplier.Addresses[0].State = "SP"
		inv.Supplier.Addresses[0].Code = ""
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "code: cannot be blank")
	})
}

func TestCustomerValidation(t *testing.T) {
	addon := tax.AddonForKey(nfe.V4)

	t.Run("validates customer required for NFe", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax.Ext[nfe.ExtKeyModel] = nfe.ModelNFe
		inv.Customer = nil
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "customer: cannot be blank")

		inv.Customer = &org.Party{
			TaxID: &tax.Identity{
				Code: "05700736000196",
			},
			Addresses: []*org.Address{
				{
					Street:   "Rua das Flores",
					Number:   "123",
					Locality: "São Paulo",
					State:    "SP",
					Code:     "01310000",
				},
			},
			Ext: tax.Extensions{
				"br-ibge-municipality": "3550308",
			},
		}
		err = addon.Validator(inv)
		assert.NoError(t, err)
	})

	t.Run("validates customer addresses required for NFe", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax.Ext[nfe.ExtKeyModel] = nfe.ModelNFe
		inv.Customer.Addresses = nil
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "addresses: cannot be blank")

		inv.Customer.Addresses = []*org.Address{}
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "addresses: cannot be blank")

		inv.Customer.Addresses = []*org.Address{nil}
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "addresses: (0: cannot be blank")

		inv.Customer.Addresses = []*org.Address{
			{
				Street:   "Rua das Flores",
				Number:   "123",
				Locality: "São Paulo",
				State:    "SP",
				Code:     "01310000",
			},
		}
		err = addon.Validator(inv)
		assert.NoError(t, err)
	})

	t.Run("customer not required for NFCe", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax.Ext[nfe.ExtKeyModel] = nfe.ModelNFCe
		inv.Tax.Ext[nfe.ExtKeyPresence] = nfe.PresenceInPerson
		inv.Customer = nil
		err := addon.Validator(inv)
		assert.NoError(t, err)
	})

	t.Run("validates customer tax ID required", func(t *testing.T) {
		inv := validInvoice()
		inv.Customer.TaxID = nil
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "tax_id: cannot be blank")

		inv.Customer.TaxID = &tax.Identity{}
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "tax_id: (code: cannot be blank")

		inv.Customer.TaxID.Code = "05700736000196"
		err = addon.Validator(inv)
		assert.NoError(t, err)
	})

	t.Run("validates customer municipality when addresses exist", func(t *testing.T) {
		inv := validInvoice()
		inv.Customer.Ext = nil
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "br-ibge-municipality: required")

		inv.Customer.Ext = tax.Extensions{
			"br-ibge-municipality": "3550308",
		}
		err = addon.Validator(inv)
		assert.NoError(t, err)
	})

	t.Run("validates customer address fields", func(t *testing.T) {
		inv := validInvoice()
		inv.Customer.Addresses = []*org.Address{
			{
				Street:   "",
				Number:   "123",
				Locality: "São Paulo",
				State:    "SP",
				Code:     "01310000",
			},
		}
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "street: cannot be blank")

		inv.Customer.Addresses[0].Street = "Rua das Flores"
		inv.Customer.Addresses[0].Number = ""
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "num: cannot be blank")

		inv.Customer.Addresses[0].Number = "123"
		inv.Customer.Addresses[0].Locality = ""
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "locality: cannot be blank")

		inv.Customer.Addresses[0].Locality = "São Paulo"
		inv.Customer.Addresses[0].State = ""
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "state: cannot be blank")

		inv.Customer.Addresses[0].State = "SP"
		inv.Customer.Addresses[0].Code = ""
		err = addon.Validator(inv)
		assert.ErrorContains(t, err, "code: cannot be blank")
	})
}

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime:   tax.WithRegime("BR"),
		Addons:   tax.WithAddons(nfe.V4),
		Currency: "BRL",
		Series:   cbc.Code("123"),
		Supplier: &org.Party{
			Name: "Test Supplier LTDA",
			TaxID: &tax.Identity{
				Code: "55263640000186",
			},
			Identities: []*org.Identity{
				{
					Key:  nfe.IdentityKeyStateReg,
					Code: "35503304557308",
				},
			},
			Addresses: []*org.Address{
				{
					Street:   "Av Paulista",
					Number:   "1578",
					Locality: "São Paulo",
					State:    "SP",
					Code:     "01310100",
				},
			},
			Ext: tax.Extensions{
				"br-ibge-municipality": "3304557",
			},
		},
		Tax: &bill.Tax{
			Ext: tax.Extensions{
				nfe.ExtKeyModel:    nfe.ModelNFe,
				nfe.ExtKeyPresence: nfe.PresenceInPerson,
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Code: "05700736000196",
			},
			Addresses: []*org.Address{
				{
					Street:   "Rua das Flores",
					Number:   "123",
					Locality: "São Paulo",
					State:    "SP",
					Code:     "01310000",
				},
			},
			Ext: tax.Extensions{
				"br-ibge-municipality": "3550308",
			},
		},
		Notes: []*org.Note{
			{
				Key:  org.NoteKeyReason,
				Text: "VENDA DE MERCADORIA",
			},
		},
		Payment: &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCash,
			},
		},
	}
}
