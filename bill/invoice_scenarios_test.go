package bill_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenarios(t *testing.T) {
	t.Run("basic flow", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		require.NoError(t, inv.Calculate())
		ss := inv.ScenarioSummary() //nolint:staticcheck
		require.NotNil(t, ss)
	})

	t.Run("invalid tags", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Tax = &bill.Tax{
			Tags: []cbc.Key{
				"random",
			},
		}
		require.ErrorContains(t, inv.Calculate(), "tax: (tags: invalid tag 'random'.)")
	})

	t.Run("scenario for new note", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Supplier.TaxID = &tax.Identity{
			Country: "IT",
			Code:    "12345678903",
		}
		inv.Tax = &bill.Tax{
			Tags: []cbc.Key{
				tax.TagReverseCharge,
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Len(t, inv.Notes, 1)
		assert.Equal(t, "Reverse Charge / Inversione del soggetto passivo", inv.Notes[0].Text)
	})

	t.Run("scenario for existing note", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Supplier.TaxID = &tax.Identity{
			Country: "IT",
			Code:    "12345678903",
		}
		inv.Tax = &bill.Tax{
			Tags: []cbc.Key{
				tax.TagReverseCharge,
			},
		}
		inv.Notes = append(inv.Notes, &cbc.Note{
			Key:  cbc.NoteKeyLegal,
			Src:  tax.TagReverseCharge,
			Text: "Random to replace",
		})
		require.NoError(t, inv.Calculate())
		assert.Len(t, inv.Notes, 1)
		assert.Equal(t, "Reverse Charge / Inversione del soggetto passivo", inv.Notes[0].Text)
	})

	t.Run("without tax defined", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Supplier.TaxID = &tax.Identity{
			Country: "IT",
			Code:    "12345678903",
		}
		inv.Tax = nil
		require.NoError(t, inv.Calculate())
		assert.Len(t, inv.Tax.Ext, 2)
		assert.Equal(t, "TD01", inv.Tax.Ext[it.ExtKeySDIDocumentType].String())
	})
}
