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
	typ  cbc.Key
	tags []cbc.Key
	exts []tax.Extensions
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

func TestScenarioSetSummary(t *testing.T) {
	ss := &tax.ScenarioSet{
		Schema: bill.ShortSchemaInvoice,
		List: []*tax.Scenario{
			{
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Ext: tax.Extensions{
					"xx-test": "normal",
				},
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
				Ext: tax.Extensions{
					"xx-test": "paid",
				},
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Tags:  []cbc.Key{tax.TagSimplified},
				Ext: tax.Extensions{
					"xx-test": "simple",
				},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Code: "note1",
					Text: "This is a note1",
				},
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Tags:  []cbc.Key{tax.TagSimplified, tax.TagPartial},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Code: "note1",
					Text: "This will replace previous note1",
				},
			},
			{
				Types:   []cbc.Key{bill.InvoiceTypeStandard},
				ExtKey:  "yy-test",
				ExtCode: "BAR",
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Code: "note2",
					Text: "This is a note 2",
				},
			},
			{
				Types:  []cbc.Key{bill.InvoiceTypeStandard},
				ExtKey: "zz-test",
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Code: "note2",
					Text: "This is a note 3",
				},
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeDebitNote},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Code: "note3",
					Text: "This should not be used",
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
		assert.Equal(t, "normal", sum.Ext["xx-test"].String())
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
		assert.Equal(t, "paid", sum.Ext["xx-test"].String())
	})
	t.Run("simplified invoice", func(t *testing.T) {
		doc := &scenarioTestDocument{
			typ:  bill.InvoiceTypeStandard,
			tags: []cbc.Key{tax.TagSimplified},
		}
		sum := ss.SummaryFor(doc)
		require.NotNil(t, sum)
		assert.Equal(t, "simple", sum.Ext["xx-test"].String())
	})
	t.Run("simplified partial invoice", func(t *testing.T) {
		doc := &scenarioTestDocument{
			typ:  bill.InvoiceTypeStandard,
			tags: []cbc.Key{tax.TagSimplified, tax.TagPartial},
		}
		sum := ss.SummaryFor(doc)
		require.NotNil(t, sum)
		assert.Equal(t, "simple", sum.Ext["xx-test"].String())
		assert.Equal(t, "This will replace previous note1", sum.Notes[0].Text)
	})
	t.Run("invoice with extensions", func(t *testing.T) {
		doc := &scenarioTestDocument{
			typ:  bill.InvoiceTypeStandard,
			exts: []tax.Extensions{{"yy-test": "BAR"}},
		}
		sum := ss.SummaryFor(doc)
		require.NotNil(t, sum)
		assert.Equal(t, "normal", sum.Ext["xx-test"].String())
		assert.Equal(t, "This is a note 2", sum.Notes[0].Text)
	})
	t.Run("invoice with extensions and no value", func(t *testing.T) {
		doc := &scenarioTestDocument{
			typ: bill.InvoiceTypeStandard,
			exts: []tax.Extensions{
				{"yy-test": "BAR"},
				{"zz-test": "FOO"},
			},
		}
		sum := ss.SummaryFor(doc)
		require.NotNil(t, sum)
		assert.Equal(t, "normal", sum.Ext["xx-test"].String())
		assert.Equal(t, "This is a note 3", sum.Notes[1].Text)
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
					Price: num.MakeAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
			},
		},
	}
	return i
}
