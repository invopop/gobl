package tax_test

import (
	"testing"

	"github.com/invopop/gobl/addons/es/facturae"
	"github.com/invopop/gobl/addons/es/tbai"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCorrectionSetDef(t *testing.T) {
	cs := tax.CorrectionSet{
		{
			Schema:     bill.ShortSchemaInvoice,
			Extensions: []cbc.Key{facturae.ExtKeyCorrection},
		},
		{
			Schema:  "note/message",
			CopyTax: true,
		},
	}

	t.Run("finds matching schema", func(t *testing.T) {
		cd := cs.Def(bill.ShortSchemaInvoice)
		require.NotNil(t, cd)
		assert.Equal(t, facturae.ExtKeyCorrection, cd.Extensions[0])
	})

	t.Run("finds by suffix", func(t *testing.T) {
		assert.NotNil(t, cs.Def("note/message"))
	})

	t.Run("returns nil for unknown schema", func(t *testing.T) {
		assert.Nil(t, cs.Def("unknown/schema"))
	})

	t.Run("returns nil for nil set", func(t *testing.T) {
		var nilSet tax.CorrectionSet
		assert.Nil(t, nilSet.Def(bill.ShortSchemaInvoice))
	})
}

func TestCorrectionDefinitionMerge(t *testing.T) {
	t.Run("merges extensions and flags", func(t *testing.T) {
		cd1 := &tax.CorrectionDefinition{
			Schema:     bill.ShortSchemaInvoice,
			Extensions: []cbc.Key{facturae.ExtKeyCorrection},
		}
		cd2 := &tax.CorrectionDefinition{
			Schema:     bill.ShortSchemaInvoice,
			Extensions: []cbc.Key{tbai.ExtKeyCorrection},
			CopyTax:    true,
		}
		merged := cd1.Merge(cd2)
		assert.Len(t, merged.Extensions, 2)
		assert.Contains(t, merged.Extensions, facturae.ExtKeyCorrection)
		assert.Contains(t, merged.Extensions, tbai.ExtKeyCorrection)
		assert.True(t, merged.CopyTax)
	})

	t.Run("nil receiver returns other", func(t *testing.T) {
		other := &tax.CorrectionDefinition{Schema: bill.ShortSchemaInvoice}
		var cd *tax.CorrectionDefinition
		assert.Equal(t, other, cd.Merge(other))
	})

	t.Run("nil other returns receiver", func(t *testing.T) {
		cd := &tax.CorrectionDefinition{Schema: bill.ShortSchemaInvoice}
		assert.Equal(t, cd, cd.Merge(nil))
	})

	t.Run("different schemas returns receiver", func(t *testing.T) {
		cd := &tax.CorrectionDefinition{Schema: bill.ShortSchemaInvoice}
		other := &tax.CorrectionDefinition{Schema: "note/message"}
		assert.Equal(t, cd, cd.Merge(other))
	})

	t.Run("does not mutate inputs", func(t *testing.T) {
		cd := &tax.CorrectionDefinition{
			Schema:     bill.ShortSchemaInvoice,
			Types:      []cbc.Key{"a"},
			Extensions: []cbc.Key{facturae.ExtKeyCorrection},
			Stamps:     []cbc.Key{"s1"},
			CopyTax:    false,
		}
		other := &tax.CorrectionDefinition{
			Schema:     bill.ShortSchemaInvoice,
			Types:      []cbc.Key{"b"},
			Extensions: []cbc.Key{tbai.ExtKeyCorrection},
			Stamps:     []cbc.Key{"s2"},
			CopyTax:    true,
		}
		_ = cd.Merge(other)

		// Receiver must remain unchanged
		assert.False(t, cd.CopyTax, "receiver CopyTax must not be mutated")
		assert.Equal(t, []cbc.Key{"a"}, cd.Types)
		assert.Equal(t, []cbc.Key{facturae.ExtKeyCorrection}, cd.Extensions)
		assert.Equal(t, []cbc.Key{"s1"}, cd.Stamps)

		// Other must remain unchanged
		assert.True(t, other.CopyTax)
		assert.Equal(t, []cbc.Key{"b"}, other.Types)
	})
}

func TestCorrectionDefinitionHasType(t *testing.T) {
	t.Run("returns true for matching type", func(t *testing.T) {
		cd := &tax.CorrectionDefinition{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote},
		}
		assert.True(t, cd.HasType(bill.InvoiceTypeCreditNote))
		assert.True(t, cd.HasType(bill.InvoiceTypeDebitNote))
	})

	t.Run("returns false for non-matching type", func(t *testing.T) {
		cd := &tax.CorrectionDefinition{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
		}
		assert.False(t, cd.HasType(bill.InvoiceTypeCorrective))
	})

	t.Run("returns false for nil definition", func(t *testing.T) {
		var cd *tax.CorrectionDefinition
		assert.False(t, cd.HasType(bill.InvoiceTypeCreditNote))
	})
}

func TestCorrectionDefinitionHasExtension(t *testing.T) {
	t.Run("returns true for matching key", func(t *testing.T) {
		cd := &tax.CorrectionDefinition{
			Extensions: []cbc.Key{facturae.ExtKeyCorrection},
		}
		assert.True(t, cd.HasExtension(facturae.ExtKeyCorrection))
	})

	t.Run("returns false for non-matching key", func(t *testing.T) {
		cd := &tax.CorrectionDefinition{
			Extensions: []cbc.Key{facturae.ExtKeyCorrection},
		}
		assert.False(t, cd.HasExtension("unknown-key"))
	})

	t.Run("returns false for nil definition", func(t *testing.T) {
		var cd *tax.CorrectionDefinition
		assert.False(t, cd.HasExtension(facturae.ExtKeyCorrection))
	})
}

func TestCorrectionNormalizeMerge(t *testing.T) {
	var calls []string

	cd1 := &tax.CorrectionDefinition{
		Schema: bill.ShortSchemaInvoice,
		Normalize: func(_ any) {
			calls = append(calls, "first")
		},
	}
	cd2 := &tax.CorrectionDefinition{
		Schema: bill.ShortSchemaInvoice,
		Normalize: func(_ any) {
			calls = append(calls, "second")
		},
	}

	t.Run("chains both normalizers", func(t *testing.T) {
		calls = nil
		merged := cd1.Merge(cd2)
		assert.NotNil(t, merged.Normalize)
		merged.Normalize(nil)
		assert.Equal(t, []string{"first", "second"}, calls)
	})

	t.Run("keeps single normalizer", func(t *testing.T) {
		calls = nil
		noNorm := &tax.CorrectionDefinition{Schema: bill.ShortSchemaInvoice}
		merged := cd1.Merge(noNorm)
		assert.NotNil(t, merged.Normalize)
		merged.Normalize(nil)
		assert.Equal(t, []string{"first"}, calls)
	})

	t.Run("adopts other normalizer", func(t *testing.T) {
		calls = nil
		noNorm := &tax.CorrectionDefinition{Schema: bill.ShortSchemaInvoice}
		merged := noNorm.Merge(cd2)
		assert.NotNil(t, merged.Normalize)
		merged.Normalize(nil)
		assert.Equal(t, []string{"second"}, calls)
	})

	t.Run("nil when both nil", func(t *testing.T) {
		noNorm1 := &tax.CorrectionDefinition{Schema: bill.ShortSchemaInvoice}
		noNorm2 := &tax.CorrectionDefinition{Schema: bill.ShortSchemaInvoice}
		merged := noNorm1.Merge(noNorm2)
		assert.Nil(t, merged.Normalize)
	})
}
