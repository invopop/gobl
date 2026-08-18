package gobl

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const calculationInvoice = `{
  "$schema": "https://gobl.org/draft-0/bill/invoice",
  "$regime": "NL",
  "currency": "EUR",
  "issue_date": "2022-07-12",
  "supplier": {
    "tax_id": {"country": "NL", "code": "000099995B57"},
    "name": "Foobar BV"
  },
  "lines": [{
    "quantity": "2",
    "item": {"name": "Tulips", "price": "10.00"}
  }]
}`

func TestCalculateWithDiscrepancies(t *testing.T) {
	t.Run("omitted calculated values", func(t *testing.T) {
		result, discrepancies, err := calculateWithDiscrepancies([]byte(calculationInvoice))
		require.NoError(t, err)
		assert.Empty(t, discrepancies)
		inv := result.(*bill.Invoice)
		require.NotNil(t, inv.Totals)
		assert.Equal(t, "20.00", inv.Totals.Payable.String())
	})

	t.Run("incorrect calculated values", func(t *testing.T) {
		data := withInvoiceValues(t, calculationInvoice, map[string]any{
			"lines": []any{map[string]any{
				"quantity": "2",
				"item":     map[string]any{"name": "Tulips", "price": "10.00"},
				"sum":      "19.00",
			}},
			"totals": map[string]any{"sum": "19.00", "payable": "19.00"},
		})

		_, discrepancies, err := calculateWithDiscrepancies(data)
		require.NoError(t, err)
		require.Len(t, discrepancies, 3)
		assert.Equal(t, "$.lines[0].sum", discrepancies[0].Path)
		assert.JSONEq(t, `"19.00"`, string(discrepancies[0].Provided))
		assert.JSONEq(t, `"20.00"`, string(discrepancies[0].Calculated))
		assert.Equal(t, "$.totals.sum", discrepancies[1].Path)
		assert.Equal(t, "$.totals.payable", discrepancies[2].Path)
	})

	t.Run("semantically equal amounts", func(t *testing.T) {
		data := withInvoiceValues(t, calculationInvoice, map[string]any{
			"lines": []any{map[string]any{
				"quantity": "2",
				"item":     map[string]any{"name": "Tulips", "price": "10.00"},
				"sum":      "20.0",
			}},
		})

		_, discrepancies, err := calculateWithDiscrepancies(data)
		require.NoError(t, err)
		assert.Empty(t, discrepancies)
	})

	t.Run("envelope paths and result", func(t *testing.T) {
		var doc map[string]any
		require.NoError(t, json.Unmarshal([]byte(calculationInvoice), &doc))
		doc["totals"] = map[string]any{"payable": "19.00"}
		data, err := json.Marshal(map[string]any{
			"$schema": "https://gobl.org/draft-0/envelope",
			"head": map[string]any{
				"uuid": "0198f976-812a-7c64-92a2-640e467159f9",
			},
			"doc": doc,
		})
		require.NoError(t, err)

		result, discrepancies, err := calculateWithDiscrepancies(data)
		require.NoError(t, err)
		assert.IsType(t, &Envelope{}, result)
		require.Len(t, discrepancies, 1)
		assert.Equal(t, "$.doc.totals.payable", discrepancies[0].Path)
	})

	t.Run("normalization-only changes are ignored", func(t *testing.T) {
		data := withInvoiceValues(t, calculationInvoice, map[string]any{
			"supplier": map[string]any{
				"tax_id": map[string]any{"country": "NL", "code": "000099995B57"},
				"name":   "  Foobar BV  ",
			},
		})

		_, discrepancies, err := calculateWithDiscrepancies(data)
		require.NoError(t, err)
		assert.Empty(t, discrepancies)
	})
}

func withInvoiceValues(t *testing.T, source string, values map[string]any) []byte {
	t.Helper()
	var doc map[string]any
	require.NoError(t, json.Unmarshal([]byte(source), &doc))
	for key, value := range values {
		doc[key] = value
	}
	data, err := json.Marshal(doc)
	require.NoError(t, err)
	return data
}
func TestCalculateMain(t *testing.T) {
	t.Run("returns the calculated value", func(t *testing.T) {
		value, err := Calculate([]byte(calculationInvoice))
		require.NoError(t, err)

		inv := value.(*bill.Invoice)
		require.NotNil(t, inv.Totals)
		assert.Equal(t, "20.00", inv.Totals.Payable.String())
	})

	t.Run("returns discrepancies as an error", func(t *testing.T) {
		data := withInvoiceValues(t, calculationInvoice, map[string]any{
			"totals": map[string]any{"payable": "19.00"},
		})

		value, err := Calculate(data, WithDiscrepancies())
		assert.Nil(t, value)
		var discrepancies CalculationDiscrepancies
		require.ErrorAs(t, err, &discrepancies)
		require.Len(t, discrepancies, 1)
		assert.Equal(t, "$.totals.payable", discrepancies[0].Path)
	})
}
