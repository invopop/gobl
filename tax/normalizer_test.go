package tax_test

import (
	"testing"

	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
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
