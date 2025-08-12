package ticket_test

import (
	"encoding/json"
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/it/ticket"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func exampleStandardInvoice(t *testing.T) *bill.Invoice {
	t.Helper()
	i := &bill.Invoice{
		Regime:   tax.WithRegime("IT"),
		Addons:   tax.WithAddons(ticket.V1),
		Code:     "123TEST",
		Currency: "EUR",
		Tax: &bill.Tax{
			PricesInclude: tax.CategoryVAT,
			Ext: tax.Extensions{
				ticket.ExtKeyLottery: "12345678",
			},
		},
		Type: bill.InvoiceTypeStandard,
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "IT",
				Code:    "12345678903",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item 0",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "general",
					},
				},
				Discounts: []*bill.LineDiscount{
					{
						Reason:  "Testing",
						Percent: num.NewPercentage(10, 2),
					},
				},
			},
			{
				Quantity: num.MakeAmount(13, 0),
				Item: &org.Item{
					Name:  "Test Item 1",
					Price: num.NewAmount(1000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Ext: tax.Extensions{
							ticket.ExtKeyExempt: "N4",
						},
					},
				},
				Discounts: []*bill.LineDiscount{
					{
						Reason:  "Testing",
						Percent: num.NewPercentage(10, 2),
					},
				},
			},
		},
	}
	return i
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("standard invoice", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("test correction", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		inv.Lines[0].Ext = tax.Extensions{
			ticket.ExtKeyLine: "1234567890",
		}
		inv.Lines[1].Ext = tax.Extensions{
			ticket.ExtKeyLine: "1234567890",
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Correct(bill.Corrective, bill.WithStamps([]*head.Stamp{
			{
				Provider: ticket.StampRef,
				Value:    "1234567890",
			},
		})))
		require.NoError(t, inv.Validate())

		json, err := json.MarshalIndent(inv, "", "  ")
		require.NoError(t, err)
		t.Log(string(json))
	})
}

func TestSupplierValidation(t *testing.T) {
	t.Run("invalid Tax ID", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		inv.Supplier.TaxID = &tax.Identity{
			Country: "IT",
			Code:    "RSSGNN60R30H501U",
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "code: contains invalid characters")
	})

	t.Run("missing supplier", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		inv.Supplier = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "supplier: cannot be blank.")
	})
}

func TestInvoiceLineTaxes(t *testing.T) {
	t.Run("item with no taxes", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		inv.Lines = append(inv.Lines, &bill.Line{
			Quantity: num.MakeAmount(10, 0),
			Item: &org.Item{
				Name:  "Test Item 2",
				Price: num.NewAmount(10000, 2),
			},
		})
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		require.EqualError(t, err, "lines: (2: (taxes: missing category VAT.).).")
	})

	t.Run("item with no tax rate nor key", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		inv.Lines = append(inv.Lines, &bill.Line{
			Quantity: num.MakeAmount(10, 0),
			Item: &org.Item{
				Name:  "Test Item 2",
				Price: num.NewAmount(10000, 2),
			},
			Taxes: tax.Set{
				{
					Category: "VAT",
				},
			},
		})
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		require.ErrorContains(t, err, "lines: (2: (taxes: (0: (percent: required for 'standard' in 'VAT'.).).).)")
	})

	t.Run("normalization when exempt key provided", func(t *testing.T) {
		list := []struct {
			Code        cbc.Code
			ExpectedKey cbc.Key
		}{
			{
				Code:        "N1",
				ExpectedKey: tax.KeyOutsideScope,
			},
			{
				Code:        "N2",
				ExpectedKey: tax.KeyOutsideScope,
			},
			{
				Code:        "N3",
				ExpectedKey: tax.KeyExport,
			},
			{
				Code:        "N4",
				ExpectedKey: tax.KeyExempt,
			},
			{
				Code:        "N5",
				ExpectedKey: tax.KeyExempt,
			},
			{
				Code:        "N6",
				ExpectedKey: tax.KeyReverseCharge,
			},
		}
		for _, row := range list {
			inv := exampleStandardInvoice(t)
			inv.Lines = append(inv.Lines, &bill.Line{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item 2",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Ext: tax.Extensions{
							ticket.ExtKeyExempt: row.Code,
						},
					},
				},
			})
			require.NoError(t, inv.Calculate())
			assert.Equal(t, row.ExpectedKey, inv.Lines[2].Taxes[0].Key)
			assert.Equal(t, row.Code, inv.Lines[2].Taxes[0].Ext.Get(ticket.ExtKeyExempt))
			assert.NoError(t, inv.Validate())
		}
	})

	t.Run("item with zero percent", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		inv.Lines = append(inv.Lines, &bill.Line{
			Quantity: num.MakeAmount(10, 0),
			Item: &org.Item{
				Name:  "Test Item 2",
				Price: num.NewAmount(10000, 2),
			},
			Taxes: tax.Set{
				{
					Category: "VAT",
					Percent:  num.NewPercentage(0, 2),
				},
			},
		})
		require.NoError(t, inv.Calculate())
		assert.Equal(t, tax.KeyZero, inv.Lines[2].Taxes[0].Key)
	})

	t.Run("sale with foreign tax combo country", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		inv.Lines[1] = &bill.Line{
			Quantity: num.MakeAmount(10, 0),
			Item: &org.Item{
				Name:  "Test Item 2",
				Price: num.NewAmount(10000, 2),
			},
			Taxes: tax.Set{
				{
					Country:  "FR",
					Category: "VAT",
					Rate:     "general",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "FR", inv.Lines[1].Taxes[0].Country.String())
		assert.Equal(t, tax.KeyStandard, inv.Lines[1].Taxes[0].Key)
		assert.Equal(t, "N2", inv.Lines[1].Taxes[0].Ext.Get(ticket.ExtKeyExempt).String())
	})

	t.Run("item with Invalid Percentage", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		inv.Lines = append(inv.Lines, &bill.Line{
			Quantity: num.MakeAmount(10, 0),
			Item: &org.Item{
				Name:  "Test Item 2",
				Price: num.NewAmount(10000, 2),
			},
			Taxes: tax.Set{
				{
					Category: "VAT",
					Percent:  num.NewPercentage(24, 2),
				},
			},
		})
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		require.EqualError(t, err, "lines: (2: (taxes: (0: (percent: must be a valid value.).).).).")
	})
}

func TestInvoiceTax(t *testing.T) {
	t.Run("invalid PricesInclude", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		inv.Tax.PricesInclude = tax.CategoryGST
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		require.EqualError(t, err, "tax: (prices_include: must be a valid value.).")
	})

	t.Run("missing PricesInclude", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		inv.Tax.PricesInclude = ""
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("missing Tax", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		inv.Tax = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("lottery code length", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		inv.Tax.Ext[ticket.ExtKeyLottery] = "1234567"
		require.NoError(t, inv.Calculate())
		require.EqualError(t, inv.Validate(), "tax: (ext: (it-ticket-lottery: does not match pattern.).).")
	})

	t.Run("lottery code empty", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		inv.Tax.Ext[ticket.ExtKeyLottery] = ""
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("lottery code uppercase", func(t *testing.T) {
		inv := exampleStandardInvoice(t)
		inv.Tax.Ext[ticket.ExtKeyLottery] = "1234567a"
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "1234567A", string(inv.Tax.Ext[ticket.ExtKeyLottery]))
	})

}
