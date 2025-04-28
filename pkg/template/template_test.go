package template_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/pkg/template"
)

func TestTemplateExecute(t *testing.T) {
	data, err := os.ReadFile("./examples/invoice.yaml.tmpl")
	require.NoError(t, err)

	tmpl, err := template.New("invoice", string(data))
	require.NoError(t, err)

	t.Run("basics", func(t *testing.T) {
		row := map[string]any{
			"code":              "1234",
			"customer_country":  "ES",
			"customer_tax_code": "A27425347",
			"customer_name":     "ACME S.L.",
			"customer_meta": here.Doc(`
				foo: "test bar"
			`),
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

		out, err := tmpl.Execute(row)
		require.NoError(t, err)
		require.NotNil(t, out)

		inv, ok := out.(*bill.Invoice)
		require.True(t, ok)

		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())

		assert.Equal(t, "ACME S.L.", inv.Customer.Name)
		assert.Equal(t, "", inv.Series.String())
		assert.Equal(t, "196.94", inv.Totals.Payable.String())
		assert.Equal(t, "34.18", inv.Totals.Tax.String())
		assert.Equal(t, "test bar", inv.Customer.Meta["foo"])
	})
	t.Run("optional fallback", func(t *testing.T) {
		row := map[string]any{
			"code": time.Now(),
		}

		out, err := tmpl.Execute(row)
		require.NoError(t, err)

		inv, ok := out.(*bill.Invoice)
		require.True(t, ok)

		assert.Equal(t, "", inv.Code.String())
	})

	t.Run("invalid YAML", func(t *testing.T) {
		data := here.Doc(`
			$schema: "https://gobl.org/draft-0/bill/invoice"
			series: "{{ .series | optional }}"
			  - "foo"
			code: "{{ .code | optional }}"
		`)
		tmpl, err := template.New("invoice", data)
		require.NoError(t, err)

		row := map[string]any{
			"code": "TEST",
		}
		_, err = tmpl.Execute(row)
		require.ErrorContains(t, err, "parsing input: yaml: line 2: did not find expected ke")
	})
}

func TestTemplateMust(t *testing.T) {
	data, err := os.ReadFile("./examples/invoice.yaml.tmpl")
	require.NoError(t, err)

	assert.NotPanics(t, func() {
		_ = template.Must(template.New("invoice", string(data)))
	})

	assert.Panics(t, func() {
		_ = template.Must(template.New("invoice", `bad {{ baahhh }}`))
	})
}
