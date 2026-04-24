package bill_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaxValidation(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		tx := &bill.Tax{}
		assert.NoError(t, rules.Validate(tx))
	})
	t.Run("with rounding", func(t *testing.T) {
		tx := &bill.Tax{
			Rounding: "precise",
		}
		assert.NoError(t, rules.Validate(tx))
	})
	t.Run("with invalid rounding", func(t *testing.T) {
		tx := &bill.Tax{
			Rounding: "currency-foo",
		}
		err := rules.Validate(tx)
		assert.ErrorContains(t, err, "rounding model is not valid")
	})
	t.Run("with tax point", func(t *testing.T) {
		tx := &bill.Tax{
			Point: "delivery",
		}
		assert.NoError(t, rules.Validate(tx))
	})
	t.Run("with invalid tax point", func(t *testing.T) {
		tx := &bill.Tax{
			Point: "invalid",
		}
		err := rules.Validate(tx)
		assert.ErrorContains(t, err, "tax point is not valid")
	})
}

func TestTaxNormalize(t *testing.T) {
	t.Run("switch rounding, sum-then-round", func(t *testing.T) {
		tx := &bill.Tax{
			Rounding: "sum-then-round",
		}
		tx.Normalize(tax.Normalizers{})
		assert.Equal(t, "precise", tx.Rounding.String())
	})
	t.Run("switch rounding, round-then-sum", func(t *testing.T) {
		tx := &bill.Tax{
			Rounding: "round-then-sum",
		}
		tx.Normalize(tax.Normalizers{})
		assert.Equal(t, "currency", tx.Rounding.String())
	})
}

func TestInvoiceTaxTagsMigration(t *testing.T) {
	// Sample document taken from spanish examples.
	in := `{
		"$schema": "https://gobl.org/draft-0/bill/invoice",
		"uuid": "3aea7b56-59d8-4beb-90bd-f8f280d852a0",
		"type": "standard",
		"code": "SAMPLE-001",
		"issue_date": "2022-02-01",
		"currency": "EUR",
		"tax": {
			"tags": [
				"simplified"
			]
		},
		"supplier": {
			"name": "Provide One S.L.",
			"tax_id": {
				"country": "ES",
				"code": "B98602642"
			},
			"addresses": [
				{
					"num": "42",
					"street": "Calle Pradillo",
					"locality": "Madrid",
					"region": "Madrid",
					"code": "28002",
					"country": "ES"
				}
			],
			"emails": [
				{
					"addr": "billing@example.com"
				}
			]
		},
		"lines": [
			{
				"i": 1,
				"quantity": "20",
				"item": {
					"name": "Main product",
					"price": "90.00"
				},
				"sum": "1800.00",
				"discounts": [
					{
						"percent": "10%",
						"amount": "180.00",
						"reason": "Special discount"
					}
				],
				"taxes": [
					{
						"cat": "VAT",
						"rate": "standard",
						"percent": "21.0%"
					}
				],
				"total": "1620.00"
			},
			{
				"i": 2,
				"quantity": "1",
				"item": {
					"name": "Something else",
					"price": "10.00"
				},
				"sum": "10.00",
				"taxes": [
					{
						"cat": "VAT",
						"rate": "standard",
						"percent": "21.0%"
					}
				],
				"total": "10.00"
			}
		],
		"totals": {
			"sum": "1630.00",
			"total": "1630.00",
			"taxes": {
				"categories": [
					{
						"code": "VAT",
						"rates": [
							{
								"key": "standard",
								"base": "1630.00",
								"percent": "21.0%",
								"amount": "342.30"
							}
						],
						"amount": "342.30"
					}
				],
				"sum": "342.30"
			},
			"tax": "342.30",
			"total_with_tax": "1972.30",
			"payable": "1972.30"
		}
	}`
	inv := new(bill.Invoice)
	require.NoError(t, json.Unmarshal([]byte(in), inv))

	assert.Equal(t, "simplified", inv.GetTags()[0].String())
}

func TestTaxMergeExtensions(t *testing.T) {
	t.Run("nil tax", func(t *testing.T) {
		var tx *bill.Tax
		ext := tax.ExtensionsOf(tax.ExtMap{
			"vat-cat": "standard",
		})
		tx = tx.MergeExtensions(ext)
		assert.Equal(t, "standard", tx.Ext.Get("vat-cat").String())
	})
	t.Run("zero extensions", func(t *testing.T) {
		tx := &bill.Tax{}
		tx = tx.MergeExtensions(tax.Extensions{})
		assert.True(t, tx.Ext.IsZero())
	})
	t.Run("with extensions", func(t *testing.T) {
		tx := &bill.Tax{
			Ext: tax.ExtensionsOf(tax.ExtMap{
				"vat-cat": "standard",
			}),
		}
		tx = tx.MergeExtensions(tax.ExtensionsOf(tax.ExtMap{
			"vat-cat": "reduced",
		}))
		assert.Equal(t, "reduced", tx.Ext.Get("vat-cat").String())
	})
	t.Run("new extensions", func(t *testing.T) {
		tx := &bill.Tax{
			Ext: tax.ExtensionsOf(tax.ExtMap{
				"vat-test": "bar",
			}),
		}
		tx = tx.MergeExtensions(tax.ExtensionsOf(tax.ExtMap{
			"vat-cat": "reduced",
		}))
		assert.Equal(t, "reduced", tx.Ext.Get("vat-cat").String())
		assert.Equal(t, "bar", tx.Ext.Get("vat-test").String())
	})
}

func TestTaxMergeNotes(t *testing.T) {
	t.Run("nil tax", func(t *testing.T) {
		var tx *bill.Tax
		n := &tax.Note{Category: tax.CategoryVAT, Key: tax.KeyExempt, Text: "Exempt"}
		tx = tx.MergeNotes(n)
		require.NotNil(t, tx)
		assert.Len(t, tx.Notes, 1)
		assert.Equal(t, "Exempt", tx.Notes[0].Text)
	})
	t.Run("no notes", func(t *testing.T) {
		tx := &bill.Tax{}
		tx = tx.MergeNotes()
		assert.Nil(t, tx.Notes)
	})
	t.Run("with existing notes", func(t *testing.T) {
		tx := &bill.Tax{
			Notes: []*tax.Note{
				{Category: tax.CategoryVAT, Key: tax.KeyExempt, Text: "Existing"},
			},
		}
		tx = tx.MergeNotes(&tax.Note{Category: tax.CategoryVAT, Key: tax.KeyReverseCharge, Text: "New"})
		assert.Len(t, tx.Notes, 2)
		assert.Equal(t, "Existing", tx.Notes[0].Text)
		assert.Equal(t, "New", tx.Notes[1].Text)
	})
}

func TestTaxJSONSchemaExtend(t *testing.T) {
	eg := `{
		"properties": {
			"rounding": {
				"type": "string",
				"title": "Rounding"
			},
			"point": {
				"type": "string",
				"title": "Point"
			}
		}
	}`
	schema := new(jsonschema.Schema)
	require.NoError(t, json.Unmarshal([]byte(eg), schema))

	tx := new(bill.Tax)
	tx.JSONSchemaExtend(schema)

	prop, ok := schema.Properties.Get("rounding")
	require.True(t, ok)
	assert.Len(t, prop.OneOf, 2)
	assert.Equal(t, "precise", prop.OneOf[0].Const)
	assert.Equal(t, "currency", prop.OneOf[1].Const)

	prop, ok = schema.Properties.Get("point")
	require.True(t, ok)
	assert.Len(t, prop.OneOf, 3)
	assert.Equal(t, "issue", prop.OneOf[0].Const)
	assert.Equal(t, "delivery", prop.OneOf[1].Const)
	assert.Equal(t, "payment", prop.OneOf[2].Const)
}

func TestTaxGetExt(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var tx *bill.Tax
		assert.Empty(t, tx.GetExt("any-ext"))
	})
	t.Run("empty", func(t *testing.T) {
		tx := &bill.Tax{}
		assert.Empty(t, tx.GetExt("any-ext"))
	})
	t.Run("with extensions", func(t *testing.T) {
		tx := &bill.Tax{
			Ext: tax.ExtensionsOf(tax.ExtMap{
				"vat-cat":  "standard",
				"vat-rate": "21.0%",
			}),
		}
		assert.Equal(t, "standard", tx.GetExt("vat-cat").String())
		assert.Equal(t, "21.0%", tx.GetExt("vat-rate").String())
		assert.Empty(t, tx.GetExt("non-existent"))
		assert.Empty(t, tx.GetExt(""))
	})
}

func TestTaxHasExt(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var tx *bill.Tax
		assert.False(t, tx.HasExt("any-ext"))
	})
	t.Run("empty", func(t *testing.T) {
		tx := &bill.Tax{}
		assert.False(t, tx.HasExt("any-ext"))
	})
	t.Run("with extensions", func(t *testing.T) {
		tx := &bill.Tax{
			Ext: tax.ExtensionsOf(tax.ExtMap{
				"vat-cat":  "standard",
				"vat-rate": "21.0%",
			}),
		}
		assert.True(t, tx.HasExt("vat-cat"))
		assert.True(t, tx.HasExt("vat-rate"))
		assert.False(t, tx.HasExt("non-existent"))
		assert.False(t, tx.HasExt(""))
	})
}
