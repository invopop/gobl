package tax_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type scenarioTestDocument struct {
	typ        cbc.Key
	tags       []cbc.Key
	exts       []tax.Extensions
	categories []cbc.Code
}

func (d *scenarioTestDocument) GetType() cbc.Key {
	return d.typ
}
func (d *scenarioTestDocument) GetTags() []cbc.Key {
	return d.tags
}
func (d *scenarioTestDocument) GetExtensions() []tax.Extensions {
	return d.exts
}
func (d *scenarioTestDocument) GetTaxCategories() []cbc.Code {
	return d.categories
}

func TestScenarioSetSummary(t *testing.T) {
	ss := &tax.ScenarioSet{
		Schema: bill.ShortSchemaInvoice,
		List: []*tax.Scenario{
			{
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Ext: tax.ExtensionsOf(tax.ExtMap{
					"xx-test": "normal",
				}),
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Filter: func(doc any) bool {
					inv, ok := doc.(*bill.Invoice)
					if !ok {
						return false
					}
					return inv.Totals.Paid()
				},
				Ext: tax.ExtensionsOf(tax.ExtMap{
					"xx-test": "paid",
				}),
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Tags:  []cbc.Key{tax.TagSimplified},
				Ext: tax.ExtensionsOf(tax.ExtMap{
					"xx-test": "simple",
				}),
				Note: &tax.Note{
					Category: tax.CategoryVAT,
					Key:      "note1",
					Text:     "This is a note1",
				},
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Tags:  []cbc.Key{tax.TagSimplified, tax.TagPartial},
				Note: &tax.Note{
					Category: tax.CategoryVAT,
					Key:      "note1",
					Text:     "This will replace previous note1",
				},
			},
			{
				Types:   []cbc.Key{bill.InvoiceTypeStandard},
				ExtKey:  "yy-test",
				ExtCode: "BAR",
				Note: &tax.Note{
					Category: tax.CategoryVAT,
					Key:      "note2",
					Text:     "This is a note 2",
				},
			},
			{
				Types:  []cbc.Key{bill.InvoiceTypeStandard},
				ExtKey: "zz-test",
				Note: &tax.Note{
					Category: tax.CategoryVAT,
					Key:      "note3",
					Text:     "This is a note 3",
				},
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeDebitNote},
				Note: &tax.Note{
					Category: tax.CategoryVAT,
					Key:      "note4",
					Text:     "This should not be used",
				},
			},
			{
				Tags: []cbc.Key{"with-code"},
				Codes: cbc.CodeMap{
					"code1": "value1",
				},
			},
		},
	}
	t.Run("standard invoice", func(t *testing.T) {
		doc := &scenarioTestDocument{
			typ: bill.InvoiceTypeStandard,
		}
		sum := ss.SummaryFor(doc)
		require.NotNil(t, sum)
		assert.Equal(t, "normal", sum.Ext.Get("xx-test").String())
	})
	t.Run("standard invoice", func(t *testing.T) {
		inv := scenariosInvoiceExample()
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Advance{
				{
					Percent:     num.NewPercentage(1, 0),
					Description: "prepaid",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		sum := ss.SummaryFor(inv)
		require.NotNil(t, sum)
		assert.Equal(t, "paid", sum.Ext.Get("xx-test").String())
	})
	t.Run("simplified invoice", func(t *testing.T) {
		doc := &scenarioTestDocument{
			typ:  bill.InvoiceTypeStandard,
			tags: []cbc.Key{tax.TagSimplified},
		}
		sum := ss.SummaryFor(doc)
		require.NotNil(t, sum)
		assert.Equal(t, "simple", sum.Ext.Get("xx-test").String())
	})
	t.Run("simplified partial invoice", func(t *testing.T) {
		doc := &scenarioTestDocument{
			typ:  bill.InvoiceTypeStandard,
			tags: []cbc.Key{tax.TagSimplified, tax.TagPartial},
		}
		sum := ss.SummaryFor(doc)
		require.NotNil(t, sum)
		assert.Equal(t, "simple", sum.Ext.Get("xx-test").String())
		assert.Equal(t, "This will replace previous note1", sum.Notes[0].Text)
	})
	t.Run("invoice with extensions", func(t *testing.T) {
		doc := &scenarioTestDocument{
			typ:  bill.InvoiceTypeStandard,
			exts: []tax.Extensions{tax.ExtensionsOf(tax.ExtMap{"yy-test": "BAR"})},
		}
		sum := ss.SummaryFor(doc)
		require.NotNil(t, sum)
		assert.Equal(t, "normal", sum.Ext.Get("xx-test").String())
		assert.Equal(t, "This is a note 2", sum.Notes[0].Text)
	})
	t.Run("invoice with extensions and no value", func(t *testing.T) {
		doc := &scenarioTestDocument{
			typ: bill.InvoiceTypeStandard,
			exts: []tax.Extensions{
				tax.ExtensionsOf(tax.ExtMap{"yy-test": "BAR"}),
				tax.ExtensionsOf(tax.ExtMap{"zz-test": "FOO"}),
			},
		}
		sum := ss.SummaryFor(doc)
		require.NotNil(t, sum)
		assert.Equal(t, "normal", sum.Ext.Get("xx-test").String())
		assert.Equal(t, "This is a note 3", sum.Notes[1].Text)
	})
	t.Run("extension keys", func(t *testing.T) {
		keys := ss.ExtensionKeys()
		require.Len(t, keys, 1)
		assert.Contains(t, keys, cbc.Key("xx-test"))
	})
	t.Run("notes extraction", func(t *testing.T) {
		notes := ss.Notes()
		require.Len(t, notes, 5)
		assert.Equal(t, "This is a note1", notes[0].Text)
	})
	t.Run("summary with note added", func(t *testing.T) {
		doc := &scenarioTestDocument{
			typ:  bill.InvoiceTypeStandard,
			tags: []cbc.Key{tax.TagSimplified},
		}
		sum := ss.SummaryFor(doc)
		require.NotNil(t, sum)
		require.Len(t, sum.Notes, 1)
		assert.Equal(t, "This is a note1", sum.Notes[0].Text)
		assert.Equal(t, tax.CategoryVAT, sum.Notes[0].Category)
	})
}

func TestScenarioSetCategoryFilter(t *testing.T) {
	ss := &tax.ScenarioSet{
		Schema: bill.ShortSchemaInvoice,
		List: []*tax.Scenario{
			{
				Tags:       []cbc.Key{tax.TagReverseCharge},
				Categories: []cbc.Code{tax.CategoryVAT},
				Note: &tax.Note{
					Category: tax.CategoryVAT,
					Key:      tax.KeyReverseCharge,
					Text:     "Reverse charge VAT",
				},
			},
			{
				Tags:       []cbc.Key{tax.TagReverseCharge},
				Categories: []cbc.Code{"IGIC"},
				Note: &tax.Note{
					Category: "IGIC",
					Key:      tax.KeyReverseCharge,
					Text:     "Reverse charge IGIC",
				},
			},
		},
	}

	t.Run("matches VAT category", func(t *testing.T) {
		doc := &scenarioTestDocument{
			typ:        bill.InvoiceTypeStandard,
			tags:       []cbc.Key{tax.TagReverseCharge},
			categories: []cbc.Code{tax.CategoryVAT},
		}
		sum := ss.SummaryFor(doc)
		require.NotNil(t, sum)
		require.Len(t, sum.Notes, 1)
		assert.Equal(t, "Reverse charge VAT", sum.Notes[0].Text)
	})

	t.Run("matches IGIC category", func(t *testing.T) {
		doc := &scenarioTestDocument{
			typ:        bill.InvoiceTypeStandard,
			tags:       []cbc.Key{tax.TagReverseCharge},
			categories: []cbc.Code{"IGIC"},
		}
		sum := ss.SummaryFor(doc)
		require.NotNil(t, sum)
		require.Len(t, sum.Notes, 1)
		assert.Equal(t, "Reverse charge IGIC", sum.Notes[0].Text)
	})

	t.Run("no match without category", func(t *testing.T) {
		doc := &scenarioTestDocument{
			typ:  bill.InvoiceTypeStandard,
			tags: []cbc.Key{tax.TagReverseCharge},
		}
		sum := ss.SummaryFor(doc)
		require.NotNil(t, sum)
		assert.Empty(t, sum.Notes)
	})

	t.Run("no match with wrong category", func(t *testing.T) {
		doc := &scenarioTestDocument{
			typ:        bill.InvoiceTypeStandard,
			tags:       []cbc.Key{tax.TagReverseCharge},
			categories: []cbc.Code{"IRPF"},
		}
		sum := ss.SummaryFor(doc)
		require.NotNil(t, sum)
		assert.Empty(t, sum.Notes)
	})

	t.Run("both categories present matches both", func(t *testing.T) {
		doc := &scenarioTestDocument{
			typ:        bill.InvoiceTypeStandard,
			tags:       []cbc.Key{tax.TagReverseCharge},
			categories: []cbc.Code{tax.CategoryVAT, "IGIC"},
		}
		sum := ss.SummaryFor(doc)
		require.NotNil(t, sum)
		assert.Len(t, sum.Notes, 2)
	})
}

func scenariosInvoiceExample() *bill.Invoice {
	i := &bill.Invoice{
		Series:    "TEST",
		Code:      "00123",
		IssueDate: cal.MakeDate(2022, 6, 13),
		Tax: &bill.Tax{
			PricesInclude: tax.CategoryVAT,
		},
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "54387763P",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Key:      "standard",
					},
				},
			},
		},
	}
	return i
}
