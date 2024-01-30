package template_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/pkg/template"
)

func TestTemplateExecute(t *testing.T) {
	data, err := os.ReadFile("./examples/invoice.yaml")
	require.NoError(t, err)

	row := map[string]any{
		"code":              "1234",
		"customer_country":  "ES",
		"customer_tax_code": "A27425347",
		"customer_name":     "ACME S.L.",
		"lines": []map[string]any{
			{
				"quantity":   "1",
				"item_name":  "Widgets",
				"item_price": "100.00",
			},
			{
				"quantity":   "12",
				"item_name":  "Gadgets",
				"item_price": "5.23",
			},
		},
	}

	tmpl, err := template.New("invoice", string(data))
	require.NoError(t, err)

	out, err := tmpl.Execute(row)
	require.NoError(t, err)
	require.NotNil(t, out)

	inv, ok := out.(*bill.Invoice)
	require.True(t, ok)

	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	assert.Equal(t, "ACME S.L.", inv.Customer.Name)
	assert.Equal(t, "", inv.Series)
	assert.Equal(t, "196.94", inv.Totals.Payable.String())
	assert.Equal(t, "34.18", inv.Totals.Tax.String())
}
