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

type discrepancyItem struct {
	Label  string `json:"label"`
	Amount string `json:"amount" jsonschema_extras:"calculated=true"`
}

type discrepancyTotals struct {
	Tax string `json:"tax"`
}

type discrepancyDoc struct {
	Label   string             `json:"label"`
	Ignored string             `json:"-" jsonschema_extras:"calculated=true"`
	Ext     map[string]string  `json:"ext,omitempty" jsonschema_extras:"calculated=true"`
	Totals  *discrepancyTotals `json:"totals,omitempty" jsonschema_extras:"calculated=true"`
	Items   []discrepancyItem  `json:"items,omitempty"`
}

func TestFindCalculationDiscrepancies_Fields(t *testing.T) {
	t.Run("fields without the calculated tag are never reported", func(t *testing.T) {
		data := []byte(`{"label": "before"}`)
		calculated := &discrepancyDoc{Label: "after"}

		discrepancies, err := gobl.FindCalculationDiscrepancies(data, calculated)
		require.NoError(t, err)
		assert.Empty(t, discrepancies)
	})

	t.Run("a calculated field with a different value is reported", func(t *testing.T) {
		data := []byte(`{"items": [{"label": "tulips", "amount": "10.00"}]}`)
		calculated := &discrepancyDoc{
			Items: []discrepancyItem{{Label: "tulips", Amount: "20.00"}},
		}

		discrepancies, err := gobl.FindCalculationDiscrepancies(data, calculated)
		require.NoError(t, err)
		require.Len(t, discrepancies, 1)
		assert.Equal(t, "$.items[0].amount", discrepancies[0].Path)
		assert.JSONEq(t, `"10.00"`, string(discrepancies[0].Provided))
		assert.JSONEq(t, `"20.00"`, string(discrepancies[0].Calculated))
	})

	t.Run("a calculated field omitted from the input is not reported", func(t *testing.T) {
		data := []byte(`{"label": "invoice"}`)
		calculated := &discrepancyDoc{
			Label:  "invoice",
			Totals: &discrepancyTotals{Tax: "6.00"},
		}

		discrepancies, err := gobl.FindCalculationDiscrepancies(data, calculated)
		require.NoError(t, err)
		assert.Empty(t, discrepancies)
	})

	t.Run("an explicit null for a calculated field is not reported", func(t *testing.T) {
		data := []byte(`{"label": "invoice", "totals": null}`)
		calculated := &discrepancyDoc{
			Label:  "invoice",
			Totals: &discrepancyTotals{Tax: "6.00"},
		}

		discrepancies, err := gobl.FindCalculationDiscrepancies(data, calculated)
		require.NoError(t, err)
		assert.Empty(t, discrepancies)
	})

	t.Run("fields nested under a calculated field are treated as calculated even when not tagged themselves", func(t *testing.T) {
		data := []byte(`{"totals": {"tax": "5.00"}}`)
		calculated := &discrepancyDoc{Totals: &discrepancyTotals{Tax: "6.00"}}

		discrepancies, err := gobl.FindCalculationDiscrepancies(data, calculated)
		require.NoError(t, err)
		require.Len(t, discrepancies, 1)
		assert.Equal(t, "$.totals.tax", discrepancies[0].Path)
	})

	t.Run("fields excluded from JSON are never reported even if marked calculated", func(t *testing.T) {
		data := []byte(`{"ignored": "x"}`)
		calculated := &discrepancyDoc{Ignored: "y"}

		discrepancies, err := gobl.FindCalculationDiscrepancies(data, calculated)
		require.NoError(t, err)
		assert.Empty(t, discrepancies)
	})

	t.Run("a calculated field that isn't a Go struct or slice is compared as a whole", func(t *testing.T) {
		data := []byte(`{"ext": {"a": "1"}}`)
		calculated := &discrepancyDoc{Ext: map[string]string{"a": "2"}}

		discrepancies, err := gobl.FindCalculationDiscrepancies(data, calculated)
		require.NoError(t, err)
		require.Len(t, discrepancies, 1)
		assert.Equal(t, "$.ext", discrepancies[0].Path)
		assert.JSONEq(t, `{"a":"1"}`, string(discrepancies[0].Provided))
		assert.JSONEq(t, `{"a":"2"}`, string(discrepancies[0].Calculated))
	})
}

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

func TestFindCalculationDiscrepancies_Invoice(t *testing.T) {
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
