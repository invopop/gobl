package tax_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type normalizeSimple struct {
	name string
}

func (s *normalizeSimple) Normalize() {
	s.name = "normalized"
}

type normalizeWithoutFunc struct {
	name string
}

func TestExtractNormalizers(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.NotPanics(t, func() {
			tax.ExtractNormalizers(nil)
		})
	})
	// Test handled here by regime and addon defs.
}

func TestExtractNormalizersForNew(t *testing.T) {
	t.Run("nil object", func(t *testing.T) {
		assert.Nil(t, tax.ExtractNormalizersForNew(nil, map[cbc.Key]bool{}))
	})

	t.Run("object without addons yields no normalizers", func(t *testing.T) {
		ns := tax.ExtractNormalizersForNew(&normalizeSimple{}, map[cbc.Key]bool{})
		assert.Empty(t, ns)
	})

	t.Run("returns normalizers for unseen addon keys and records them", func(t *testing.T) {
		as := tax.WithAddons("mx-cfdi-v4")
		seen := map[cbc.Key]bool{}
		ns := tax.ExtractNormalizersForNew(&as, seen)
		assert.Len(t, ns, 1)
		assert.True(t, seen["mx-cfdi-v4"])
	})

	t.Run("skips keys already seen", func(t *testing.T) {
		as := tax.WithAddons("mx-cfdi-v4")
		seen := map[cbc.Key]bool{"mx-cfdi-v4": true}
		ns := tax.ExtractNormalizersForNew(&as, seen)
		assert.Empty(t, ns)
	})

	t.Run("second call over the same object returns nothing new", func(t *testing.T) {
		as := tax.WithAddons("mx-cfdi-v4")
		seen := map[cbc.Key]bool{}
		first := tax.ExtractNormalizersForNew(&as, seen)
		require.Len(t, first, 1)
		second := tax.ExtractNormalizersForNew(&as, seen)
		assert.Empty(t, second)
	})

	t.Run("only the newly-added addon is returned across passes", func(t *testing.T) {
		// Simulates a meta-addon appending a further addon between passes:
		// the first pass records the initial addon, and a second pass after
		// AddAddons returns only the newly-introduced one.
		as := tax.WithAddons("mx-cfdi-v4")
		seen := map[cbc.Key]bool{}
		require.Len(t, tax.ExtractNormalizersForNew(&as, seen), 1)

		as.AddAddons("es-verifactu-v1")
		ns := tax.ExtractNormalizersForNew(&as, seen)
		assert.Len(t, ns, 1) // only es-verifactu-v1, mx-cfdi-v4 already seen
		assert.True(t, seen["es-verifactu-v1"])
	})
}

func TestNormalize(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.NotPanics(t, func() {
			tax.Normalize(nil, nil)
		})
	})
	t.Run("simple implementation", func(t *testing.T) {
		s := &normalizeSimple{name: "original"}
		assert.Equal(t, "original", s.name)
		tax.Normalize(nil, s)
		assert.Equal(t, "normalized", s.name)
	})

	t.Run("simple implementation with list", func(t *testing.T) {
		s := &normalizeSimple{name: "original"}
		nl := []tax.Normalizer{
			func(in any) {
				t.Log("setting name")
				if s, ok := in.(*normalizeSimple); ok {
					s.name = "normalized by func"
				}
			},
		}
		tax.Normalize(nl, s)
		assert.Equal(t, "normalized by func", s.name)
	})

	t.Run("with normalizers", func(t *testing.T) {
		nl := []tax.Normalizer{
			func(in any) {
				if s, ok := in.(*normalizeWithoutFunc); ok {
					s.name = "normalized by func"
				}
			},
		}
		s := &normalizeWithoutFunc{name: "original"}
		tax.Normalize(nl, s)
		assert.Equal(t, "normalized by func", s.name)
	})
}
