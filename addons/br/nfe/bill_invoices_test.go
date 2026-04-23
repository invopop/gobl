package nfe_test

import (
	"fmt"
	"testing"

	"github.com/invopop/gobl/addons/br/nfe"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/br"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoicesValidation(t *testing.T) {
	t.Run("validates tax extensions", func(t *testing.T) {
		inv := validCalculatedInvoice(t)
		inv.Tax = nil
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "tax is required")

		inv.Tax = &bill.Tax{}
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "tax requires 'br-nfe-model' and 'br-nfe-presence' extensions")

		inv.Tax.Ext = tax.ExtensionsOf(tax.ExtMap{
			nfe.ExtKeyModel:    nfe.ModelNFe,
			nfe.ExtKeyPresence: nfe.PresenceDelivery,
		})
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "NF-e invoices do not support '4' for 'br-nfe-presence'")

		inv.Tax.Ext = inv.Tax.Ext.Set(nfe.ExtKeyPresence, nfe.PresenceInPerson)
		err = rules.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("validates required notes", func(t *testing.T) {
		inv := validCalculatedInvoice(t)
		inv.Notes = nil
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "a note with key 'reason' is required to describe the nature of the operation (natOp)")

		inv.Notes = []*org.Note{nil}
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "a note with key 'reason' is required to describe the nature of the operation (natOp)")

		inv.Notes[0] = &org.Note{
			Key:  org.NoteKeyGeneral,
			Text: "General note",
		}
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "a note with key 'reason' is required to describe the nature of the operation (natOp)")

		inv.Notes[0].Key = org.NoteKeyReason
		inv.Notes[0].Text = "1234567890123456789012345678901234567890123456789012345678901" // 61 chars
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "reason note text must be between 1 and 60 characters")

		inv.Notes[0].Text = "123456789012345678901234567890123456789012345678901234567890" // 60 chars
		err = rules.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("validates payment when invoice is due", func(t *testing.T) {
		inv := validCalculatedInvoice(t)
		inv.Payment = nil
		inv.Totals.Due = &num.AmountZero
		err := rules.Validate(inv)
		assert.NoError(t, err)

		inv.Totals.Due = nil
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "payment is required")

		inv.Totals.Due = num.NewAmount(1, 2)
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "payment is required")

		inv.Payment = &bill.PaymentDetails{}
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "payment instructions are required")

		inv.Payment.Instructions = &pay.Instructions{
			Key: pay.MeansKeyCash,
			Ext: tax.ExtensionsOf(tax.ExtMap{
				nfe.ExtKeyPaymentMeans: "01",
			}),
		}
		err = rules.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("validates invoice totals due field", func(t *testing.T) {
		inv := validCalculatedInvoice(t)

		inv.Totals.Due = num.NewAmount(-1, 2)
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "due amount must not be negative")

		inv.Totals.Due = &num.AmountZero
		err = rules.Validate(inv)
		assert.NoError(t, err)

		inv.Totals.Due = num.NewAmount(1, 2)
		err = rules.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("validates NFe presence when model is NFe", func(t *testing.T) {
		inv := validCalculatedInvoice(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(nfe.ExtKeyModel, nfe.ModelNFe)
		inv.Tax.Ext = inv.Tax.Ext.Set(nfe.ExtKeyPresence, nfe.PresenceDelivery)
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "NF-e invoices do not support '4' for 'br-nfe-presence'")

		inv.Tax.Ext = inv.Tax.Ext.Set(nfe.ExtKeyPresence, nfe.PresenceInPerson)
		err = rules.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("validates NFCe presence when model is NFCe", func(t *testing.T) {
		inv := validCalculatedInvoice(t)
		inv.Customer = nil // For NFCe, customer is optional

		inv.Tax.Ext = inv.Tax.Ext.Set(nfe.ExtKeyModel, nfe.ModelNFCe)
		inv.Tax.Ext = inv.Tax.Ext.Set(nfe.ExtKeyPresence, nfe.PresenceNotApplicable)
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "NFC-e invoices require in-person or delivery for 'br-nfe-presence'")

		inv.Tax.Ext = inv.Tax.Ext.Set(nfe.ExtKeyPresence, nfe.PresenceInPerson)
		err = rules.Validate(inv)
		assert.NoError(t, err)
	})
}

func TestInvoiceSeriesValidation(t *testing.T) {
	tests := []struct {
		series cbc.Code
		err    string
	}{
		{series: "0"},
		{series: "1"},
		{series: "12"},
		{series: "123"},
		{series: "999"},
		{series: "", err: "series is required"},
		{series: "1000", err: "series format is invalid; must be 0 or 1-999"},
		{series: "abc", err: "series format is invalid; must be 0 or 1-999"},
		{series: "012", err: "series format is invalid; must be 0 or 1-999"},
		{series: "00", err: "series format is invalid; must be 0 or 1-999"},
		{series: "-3", err: "series format is invalid; must be 0 or 1-999"},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("validates series %s", tt.series)
		t.Run(name, func(t *testing.T) {
			inv := validCalculatedInvoice(t)
			inv.Series = tt.series
			err := rules.Validate(inv)
			if tt.err != "" {
				assert.ErrorContains(t, err, tt.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSupplierValidation(t *testing.T) {
	t.Run("nil supplier", func(t *testing.T) {
		inv := validCalculatedInvoice(t)
		inv.Supplier = nil
		err := rules.Validate(inv)
		// supplier presence is validated at GOBL level, but our rules
		// will still produce errors for nested nil fields - check no panic
		assert.Error(t, err) // GOBL core requires supplier
	})

	t.Run("validates supplier name", func(t *testing.T) {
		inv := validCalculatedInvoice(t)
		inv.Supplier.Name = ""
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier name is required")

		inv.Supplier.Name = "Test Company"
		err = rules.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("validates supplier addresses required", func(t *testing.T) {
		inv := validCalculatedInvoice(t)
		inv.Supplier.Addresses = nil
		inv.Supplier.Ext = tax.Extensions{} // remove ext to avoid municipality check on nil addresses
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier must have at least one address")

		inv.Supplier.Addresses = []*org.Address{}
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier must have at least one address")

		inv.Supplier.Addresses = []*org.Address{nil}
		inv.Supplier.Ext = tax.ExtensionsOf(tax.ExtMap{
			"br-ibge-municipality": "3304557",
		})
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier address must not be empty")

		inv.Supplier.Addresses = []*org.Address{
			{
				Street:   "Rua Test",
				Number:   "100",
				Locality: "São Paulo",
				State:    "SP",
				Code:     "01310100",
			},
		}
		err = rules.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("validates supplier tax ID required", func(t *testing.T) {
		inv := validCalculatedInvoice(t)
		inv.Supplier.TaxID = nil
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier tax ID is required")

		inv.Supplier.TaxID = &tax.Identity{
			Country: "BR",
		}
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier tax ID code is required")

		inv.Supplier.TaxID.Code = "55263640000186"
		err = rules.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("validates supplier municipality extension when addresses exist", func(t *testing.T) {
		inv := validCalculatedInvoice(t)
		inv.Supplier.Ext = tax.Extensions{}
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "requires 'br-ibge-municipality' extension when addresses are present")

		inv.Supplier.Ext = tax.ExtensionsOf(tax.ExtMap{
			"br-ibge-municipality": "3304557",
		})
		err = rules.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("validates supplier address fields", func(t *testing.T) {
		inv := validCalculatedInvoice(t)
		inv.Supplier.Addresses = []*org.Address{
			{
				Street:   "",
				Number:   "100",
				Locality: "São Paulo",
				State:    "SP",
				Code:     "01310100",
			},
		}
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier address requires a street")

		inv.Supplier.Addresses[0].Street = "Rua Test"
		inv.Supplier.Addresses[0].Number = ""
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier address requires a number")

		inv.Supplier.Addresses[0].Number = "100"
		inv.Supplier.Addresses[0].Locality = ""
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier address requires a locality")

		inv.Supplier.Addresses[0].Locality = "São Paulo"
		inv.Supplier.Addresses[0].State = ""
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier address requires a state")

		inv.Supplier.Addresses[0].State = "SP"
		inv.Supplier.Addresses[0].Code = ""
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier address requires a postal code")
	})
}

func TestCustomerValidation(t *testing.T) {
	t.Run("validates customer required for NFe", func(t *testing.T) {
		inv := validCalculatedInvoice(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(nfe.ExtKeyModel, nfe.ModelNFe)
		inv.Customer = nil
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "customer is required for NF-e")

		inv.Customer = &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "BR",
				Code:    "05700736000196",
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
			Ext: tax.ExtensionsOf(tax.ExtMap{
				"br-ibge-municipality": "3550308",
			}),
		}
		err = rules.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("validates customer addresses required for NFe", func(t *testing.T) {
		inv := validCalculatedInvoice(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(nfe.ExtKeyModel, nfe.ModelNFe)
		inv.Customer.Addresses = nil
		inv.Customer.Ext = tax.Extensions{} // avoid municipality check
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "customer must have at least one address for NF-e")

		inv.Customer.Addresses = []*org.Address{}
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "customer must have at least one address for NF-e")

		inv.Customer.Addresses = []*org.Address{nil}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			"br-ibge-municipality": "3550308",
		})
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "customer address must not be empty")

		inv.Customer.Addresses = []*org.Address{
			{
				Street:   "Rua das Flores",
				Number:   "123",
				Locality: "São Paulo",
				State:    "SP",
				Code:     "01310000",
			},
		}
		err = rules.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("customer not required for NFCe", func(t *testing.T) {
		inv := validCalculatedInvoice(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(nfe.ExtKeyModel, nfe.ModelNFCe)
		inv.Tax.Ext = inv.Tax.Ext.Set(nfe.ExtKeyPresence, nfe.PresenceInPerson)
		inv.Customer = nil
		err := rules.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("validates customer tax ID required", func(t *testing.T) {
		inv := validCalculatedInvoice(t)
		inv.Customer.TaxID = nil
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "customer tax ID is required")

		inv.Customer.TaxID = &tax.Identity{
			Country: "BR",
		}
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "customer tax ID code is required")

		inv.Customer.TaxID.Code = "05700736000196"
		err = rules.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("validates customer municipality when addresses exist", func(t *testing.T) {
		inv := validCalculatedInvoice(t)
		inv.Customer.Ext = tax.Extensions{}
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "requires 'br-ibge-municipality' extension when addresses are present")

		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			"br-ibge-municipality": "3550308",
		})
		err = rules.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("validates customer address fields", func(t *testing.T) {
		inv := validCalculatedInvoice(t)
		inv.Customer.Addresses = []*org.Address{
			{
				Street:   "",
				Number:   "123",
				Locality: "São Paulo",
				State:    "SP",
				Code:     "01310000",
			},
		}
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "customer address requires a street")

		inv.Customer.Addresses[0].Street = "Rua das Flores"
		inv.Customer.Addresses[0].Number = ""
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "customer address requires a number")

		inv.Customer.Addresses[0].Number = "123"
		inv.Customer.Addresses[0].Locality = ""
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "customer address requires a locality")

		inv.Customer.Addresses[0].Locality = "São Paulo"
		inv.Customer.Addresses[0].State = ""
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "customer address requires a state")

		inv.Customer.Addresses[0].State = "SP"
		inv.Customer.Addresses[0].Code = ""
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "customer address requires a postal code")
	})
}

func TestInvoiceCurrencyValidation(t *testing.T) {
	t.Run("non-BRL currency without exchange rates", func(t *testing.T) {
		inv := validCalculatedInvoice(t)
		inv.Currency = "USD"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-BR-NFE-BILL-INVOICE-34] invoice must be in BRL or provide exchange rate for conversion")
	})

	t.Run("non-BRL currency with exchange rates", func(t *testing.T) {
		inv := validCalculatedInvoice(t)
		inv.Currency = "USD"
		inv.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   "USD",
				To:     "BRL",
				Amount: num.MakeAmount(500, 2),
			},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.NoError(t, err)
	})
}

// validInvoice creates a raw invoice suitable for scenario tests that call Calculate() themselves.
func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Addons:   tax.WithAddons(nfe.V4),
		Currency: "BRL",
		Series:   cbc.Code("123"),
		Supplier: &org.Party{
			Name: "Test Supplier LTDA",
			TaxID: &tax.Identity{
				Country: "BR",
				Code:    "55263640000186",
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
			Ext: tax.ExtensionsOf(tax.ExtMap{
				"br-ibge-municipality": "3304557",
			}),
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "BR",
				Code:    "05700736000196",
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
			Ext: tax.ExtensionsOf(tax.ExtMap{
				"br-ibge-municipality": "3550308",
			}),
		},
		Tax: &bill.Tax{
			Ext: tax.ExtensionsOf(tax.ExtMap{
				nfe.ExtKeyModel:    nfe.ModelNFe,
				nfe.ExtKeyPresence: nfe.PresenceInPerson,
			}),
		},
		Notes: []*org.Note{
			{
				Key:  org.NoteKeyReason,
				Text: "VENDA DE MERCADORIA",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Product",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: br.TaxCategoryICMS,
						Percent:  num.NewPercentage(18, 2),
					},
					{
						Category: br.TaxCategoryPIS,
						Percent:  num.NewPercentage(165, 4),
					},
					{
						Category: br.TaxCategoryCOFINS,
						Percent:  num.NewPercentage(760, 4),
					},
				},
			},
		},
		Payment: &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCash,
			},
		},
	}
}

// validCalculatedInvoice creates a fully valid invoice with Calculate() applied,
// suitable for post-modification testing of specific rule violations.
func validCalculatedInvoice(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := &bill.Invoice{
		Addons:   tax.WithAddons(nfe.V4),
		Currency: "BRL",
		Series:   cbc.Code("123"),
		Supplier: &org.Party{
			Name: "Test Supplier LTDA",
			TaxID: &tax.Identity{
				Country: "BR",
				Code:    "55263640000186",
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
			Ext: tax.ExtensionsOf(tax.ExtMap{
				"br-ibge-municipality": "3304557",
			}),
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "BR",
				Code:    "05700736000196",
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
			Ext: tax.ExtensionsOf(tax.ExtMap{
				"br-ibge-municipality": "3550308",
			}),
		},
		Notes: []*org.Note{
			{
				Key:  org.NoteKeyReason,
				Text: "VENDA DE MERCADORIA",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Product",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: br.TaxCategoryICMS,
						Percent:  num.NewPercentage(18, 2),
					},
					{
						Category: br.TaxCategoryPIS,
						Percent:  num.NewPercentage(165, 4),
					},
					{
						Category: br.TaxCategoryCOFINS,
						Percent:  num.NewPercentage(760, 4),
					},
				},
			},
		},
		Payment: &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCash,
			},
		},
	}
	require.NoError(t, inv.Calculate())
	// Presence is not set by scenarios, set it manually after Calculate
	inv.Tax.Ext = inv.Tax.Ext.Set(nfe.ExtKeyPresence, nfe.PresenceInPerson)
	require.NoError(t, rules.Validate(inv))
	return inv
}
