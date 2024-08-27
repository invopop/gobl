package tax_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenarioSetSummary(t *testing.T) {
	ss := &tax.ScenarioSet{
		Schema: bill.ShortSchemaInvoice,
		List: []*tax.Scenario{
			{
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Ext: tax.Extensions{
					"xx-test": "100",
				},
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Tags:  []cbc.Key{tax.TagSimplified},
				Ext: tax.Extensions{
					"xx-test": "200",
				},
				Note: &cbc.Note{
					Key:  cbc.NoteKeyLegal,
					Code: "note1",
					Text: "This is a note1",
				},
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Tags:  []cbc.Key{tax.TagSimplified, tax.TagPartial},
				Note: &cbc.Note{
					Key:  cbc.NoteKeyLegal,
					Code: "note1",
					Text: "This will replace previous note1",
				},
			},
			{
				Types:    []cbc.Key{bill.InvoiceTypeStandard},
				ExtKey:   "yy-test",
				ExtValue: "BAR",
				Note: &cbc.Note{
					Key:  cbc.NoteKeyLegal,
					Code: "note2",
					Text: "This is a note 2",
				},
			},
			{
				Types:  []cbc.Key{bill.InvoiceTypeStandard},
				ExtKey: "zz-test",
				Note: &cbc.Note{
					Key:  cbc.NoteKeyLegal,
					Code: "note2",
					Text: "This is a note 3",
				},
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeDebitNote},
				Note: &cbc.Note{
					Key:  cbc.NoteKeyLegal,
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
		sum := ss.SummaryFor(bill.InvoiceTypeStandard, nil, nil)
		require.NotNil(t, sum)
		assert.Equal(t, "100", sum.Ext["xx-test"].String())
	})
	t.Run("simplified invoice", func(t *testing.T) {
		sum := ss.SummaryFor(bill.InvoiceTypeStandard, []cbc.Key{tax.TagSimplified}, nil)
		require.NotNil(t, sum)
		assert.Equal(t, "200", sum.Ext["xx-test"].String())
	})
	t.Run("simplified partial invoice", func(t *testing.T) {
		sum := ss.SummaryFor(bill.InvoiceTypeStandard, []cbc.Key{tax.TagSimplified, tax.TagPartial}, nil)
		require.NotNil(t, sum)
		assert.Equal(t, "200", sum.Ext["xx-test"].String())
		assert.Equal(t, "This will replace previous note1", sum.Notes[0].Text)
	})
	t.Run("invoice with extensions", func(t *testing.T) {
		sum := ss.SummaryFor(bill.InvoiceTypeStandard, []cbc.Key{}, []tax.Extensions{{"yy-test": "BAR"}})
		require.NotNil(t, sum)
		assert.Equal(t, "100", sum.Ext["xx-test"].String())
		assert.Equal(t, "This is a note 2", sum.Notes[0].Text)
	})
	t.Run("invoice with extensions and no value", func(t *testing.T) {
		sum := ss.SummaryFor(bill.InvoiceTypeStandard, []cbc.Key{}, []tax.Extensions{{"yy-test": "BAR", "zz-test": "FOO"}})
		require.NotNil(t, sum)
		assert.Equal(t, "100", sum.Ext["xx-test"].String())
		assert.Equal(t, "This is a note 3", sum.Notes[1].Text)
	})

}
