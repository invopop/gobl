package en16931_test

import (
	"testing"

	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenarios(t *testing.T) {
	t.Run("provides list of scenarios", func(t *testing.T) {
		scenarios := en16931.Scenarios()
		assert.NotEmpty(t, scenarios)
	})

	t.Run("prepayment invoice (386)", func(t *testing.T) {
		inv := &bill.Invoice{
			Addons:   tax.WithAddons(en16931.V2017),
			Code:     "TEST-001",
			Currency: "EUR",
			Type:     bill.InvoiceTypeStandard,
			Supplier: &org.Party{Name: "Test Supplier"},
			Customer: &org.Party{Name: "Test Customer"},
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Test",
						Price: num.NewAmount(10000, 2),
					},
				},
			},
		}
		inv.SetTags(tax.TagPrepayment)

		err := inv.Calculate()
		require.NoError(t, err)

		assert.Equal(t, "386", inv.Tax.Ext[untdid.ExtKeyDocumentType].String())
	})

	t.Run("factored invoice (393)", func(t *testing.T) {
		inv := &bill.Invoice{
			Addons:   tax.WithAddons(en16931.V2017),
			Code:     "TEST-001",
			Currency: "EUR",
			Type:     bill.InvoiceTypeStandard,
			Supplier: &org.Party{Name: "Test Supplier"},
			Customer: &org.Party{Name: "Test Customer"},
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Test",
						Price: num.NewAmount(10000, 2),
					},
				},
			},
		}
		inv.SetTags(tax.TagFactoring)

		err := inv.Calculate()
		require.NoError(t, err)

		assert.Equal(t, "393", inv.Tax.Ext[untdid.ExtKeyDocumentType].String())
	})

	t.Run("factored credit note (396)", func(t *testing.T) {
		inv := &bill.Invoice{
			Addons:   tax.WithAddons(en16931.V2017),
			Code:     "TEST-001",
			Currency: "EUR",
			Type:     bill.InvoiceTypeCreditNote,
			Supplier: &org.Party{Name: "Test Supplier"},
			Customer: &org.Party{Name: "Test Customer"},
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Test",
						Price: num.NewAmount(10000, 2),
					},
				},
			},
		}
		inv.SetTags(tax.TagFactoring)

		err := inv.Calculate()
		require.NoError(t, err)

		assert.Equal(t, "396", inv.Tax.Ext[untdid.ExtKeyDocumentType].String())
	})
}
