package gobl_test

import (
	"encoding/json"
	"maps"
	"os"
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sampleInvoice = `{
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

func TestFindCalculationDiscrepancies(t *testing.T) {
	t.Run("omitted calculated values are not discrepancies", func(t *testing.T) {
		data := []byte(sampleInvoice)
		inv := parseAndCalculateInvoice(t, data)
		require.Equal(t, "20.00", inv.Totals.Payable.String())

		discrepancies, err := gobl.FindCalculationDiscrepancies(data, inv)
		require.NoError(t, err)
		assert.Empty(t, discrepancies)
	})

	t.Run("incorrect calculated values are reported with their paths", func(t *testing.T) {
		data := withInvoiceValues(t, map[string]any{
			"lines": []any{map[string]any{
				"quantity": "2",
				"item":     map[string]any{"name": "Tulips", "price": "10.00"},
				"sum":      "19.00",
			}},
			"totals": map[string]any{"sum": "19.00", "payable": "19.00"},
		})
		inv := parseAndCalculateInvoice(t, data)

		discrepancies, err := gobl.FindCalculationDiscrepancies(data, inv)
		require.NoError(t, err)
		require.Len(t, discrepancies, 3)
		assert.Equal(t, "$.lines[0].sum", discrepancies[0].Path)
		assert.JSONEq(t, `"19.00"`, string(discrepancies[0].Provided))
		assert.JSONEq(t, `"20.00"`, string(discrepancies[0].Calculated))
		assert.Equal(t, "$.totals.sum", discrepancies[1].Path)
		assert.Equal(t, "$.totals.payable", discrepancies[2].Path)
	})

	t.Run("semantically equal amounts are not discrepancies", func(t *testing.T) {
		data := withInvoiceValues(t, map[string]any{
			"lines": []any{map[string]any{
				"quantity": "2",
				"item":     map[string]any{"name": "Tulips", "price": "10.00"},
				"sum":      "20.0",
			}},
		})
		inv := parseAndCalculateInvoice(t, data)

		discrepancies, err := gobl.FindCalculationDiscrepancies(data, inv)
		require.NoError(t, err)
		assert.Empty(t, discrepancies)
	})

	t.Run("an explicit zero value for a non-pointer calculated field is not a discrepancy", func(t *testing.T) {
		// "type" is calculated but isn't a pointer, so an explicit ""
		// can't be told apart from having left it out; normalizeInvoice
		// fills it in with "standard" either way.
		data := withInvoiceValues(t, map[string]any{"type": ""})
		inv := parseAndCalculateInvoice(t, data)
		require.Equal(t, "standard", inv.Type.String())

		discrepancies, err := gobl.FindCalculationDiscrepancies(data, inv)
		require.NoError(t, err)
		assert.Empty(t, discrepancies)
	})

	t.Run("a deeply nested tax breakdown reports precise paths under an envelope", func(t *testing.T) {
		data, err := os.ReadFile("examples/ie/out/invoice-b2b.json")
		require.NoError(t, err)

		var root map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(data, &root))
		var doc map[string]any
		require.NoError(t, json.Unmarshal(root["doc"], &doc))
		totals := doc["totals"].(map[string]any)
		taxes := totals["taxes"].(map[string]any)
		categories := taxes["categories"].([]any)
		category := categories[0].(map[string]any)
		category["amount"] = "999.99"
		docBytes, err := json.Marshal(doc)
		require.NoError(t, err)
		root["doc"] = docBytes
		full, err := json.Marshal(root)
		require.NoError(t, err)

		obj, err := gobl.Parse(full)
		require.NoError(t, err)
		env := obj.(*gobl.Envelope)
		require.NoError(t, env.Calculate())

		discrepancies, err := gobl.FindCalculationDiscrepancies(full, env)
		require.NoError(t, err)
		require.Len(t, discrepancies, 1)
		assert.Equal(t, "$.doc.totals.taxes.categories[0].amount", discrepancies[0].Path)
		assert.JSONEq(t, `"999.99"`, string(discrepancies[0].Provided))
	})
}

func parseAndCalculateInvoice(t *testing.T, data []byte) *bill.Invoice {
	t.Helper()
	obj, err := gobl.Parse(data)
	require.NoError(t, err)
	inv, ok := obj.(*bill.Invoice)
	require.True(t, ok)
	require.NoError(t, inv.Calculate())
	return inv
}

func withInvoiceValues(t *testing.T, values map[string]any) []byte {
	t.Helper()
	var doc map[string]any
	require.NoError(t, json.Unmarshal([]byte(sampleInvoice), &doc))
	maps.Copy(doc, values)
	data, err := json.Marshal(doc)
	require.NoError(t, err)
	return data
}
