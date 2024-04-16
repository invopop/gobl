package tax_test

import (
	"context"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/regimes/mx"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeExtensions(t *testing.T) {
	var em tax.Extensions

	em2 := tax.NormalizeExtensions(em)
	assert.Nil(t, em2)

	em = tax.Extensions{
		"key": "",
	}
	em2 = tax.NormalizeExtensions(em)
	assert.Nil(t, em2)

	em = tax.Extensions{
		"key": "foo",
		"bar": "",
	}
	em2 = tax.NormalizeExtensions(em)
	assert.NotNil(t, em2)
	assert.Len(t, em2, 1)
	assert.Equal(t, "foo", em2["key"].String())
}

func TestExtValue(t *testing.T) {
	ev := tax.ExtValue("IT")
	assert.Equal(t, "IT", ev.String())
	assert.Equal(t, cbc.Code("IT"), ev.Code())
	assert.Equal(t, cbc.KeyEmpty, ev.Key())

	ev = tax.ExtValue("testing")
	assert.Equal(t, "testing", ev.String())
	assert.Equal(t, cbc.Key("testing"), ev.Key())
	assert.Equal(t, cbc.CodeEmpty, ev.Code())

	ev = tax.ExtValue("A string")
	assert.Equal(t, cbc.CodeEmpty, ev.Code())
	assert.Equal(t, cbc.KeyEmpty, ev.Key())
	assert.Equal(t, "A string", ev.String())
}

func TestExtValidation(t *testing.T) {
	t.Run("with mexico", func(t *testing.T) {
		// Use mexico for tests as it has more extensions
		mr := mx.New()
		ctx := mr.WithContext(context.Background())

		t.Run("test patterns", func(t *testing.T) {
			em := tax.Extensions{
				mx.ExtKeyCFDIPostCode: "12345",
			}
			err := em.ValidateWithContext(ctx)
			assert.NoError(t, err)

			em = tax.Extensions{
				mx.ExtKeyCFDIPostCode: "123457",
			}
			err = em.ValidateWithContext(ctx)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "mx-cfdi-post-code: does not match pattern")

			kd := mr.ExtensionDef(mx.ExtKeyCFDIPostCode)
			pt := kd.Pattern
			kd.Pattern = "[][" // invalid
			err = em.ValidateWithContext(ctx)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "mx-cfdi-post-code: error parsing regexp: missing closing ]: `[][`")
			kd.Pattern = pt // put back!
		})

		t.Run("test codes", func(t *testing.T) {
			em := tax.Extensions{
				mx.ExtKeyCFDIFiscalRegime: "601",
			}
			err := em.ValidateWithContext(ctx)
			assert.NoError(t, err)

			em = tax.Extensions{
				mx.ExtKeyCFDIFiscalRegime: "000",
			}
			err = em.ValidateWithContext(ctx)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "mx-cfdi-fiscal-regime: code '000' invalid")
		})
	})

	sp := es.New()
	ctx := sp.WithContext(context.Background())
	t.Run("test good key", func(t *testing.T) {
		em := tax.Extensions{
			es.ExtKeyTBAIProduct: "goods",
		}
		err := em.ValidateWithContext(ctx)
		assert.NoError(t, err)
	})
	t.Run("test bad key", func(t *testing.T) {
		em := tax.Extensions{
			es.ExtKeyTBAIProduct: "bads",
		}
		err := em.ValidateWithContext(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "es-tbai-product: key 'bads' invalid")
	})
	t.Run("missing extension", func(t *testing.T) {
		em := tax.Extensions{
			"random-key": "type",
		}
		err := em.ValidateWithContext(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "random-key: undefined")
	})
	t.Run("invalid key", func(t *testing.T) {
		em := tax.Extensions{
			"INVALID": "value",
		}
		err := em.ValidateWithContext(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "INVALID: must be in a valid format")
	})
}

func TestExtensionsHas(t *testing.T) {
	em := tax.Extensions{
		"key": "value",
	}
	assert.True(t, em.Has("key"))
	assert.False(t, em.Has("invalid"))
}

func TestExtensionsEquals(t *testing.T) {
	tests := []struct {
		name string
		em1  tax.Extensions
		em2  tax.Extensions
		want bool
	}{
		{
			name: "empty",
			em1:  tax.Extensions{},
			em2:  tax.Extensions{},
			want: true,
		},
		{
			name: "same",
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key": "value"},
			want: true,
		},
		{
			name: "different",
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key": "value2"},
			want: false,
		},
		{
			name: "different keys",
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key2": "value"},
			want: false,
		},
		{
			name: "different lengths",
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key": "value", "key2": "value"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.em1.Equals(tt.em2))
		})
	}
}

func TestExtensionsContains(t *testing.T) {
	tests := []struct {
		name string
		em1  tax.Extensions
		em2  tax.Extensions
		want bool
	}{
		{
			name: "empty",
			em1:  tax.Extensions{},
			em2:  tax.Extensions{},
			want: false,
		},
		{
			name: "same",
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key": "value"},
			want: true,
		},
		{
			name: "different",
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key": "value2"},
			want: false,
		},
		{
			name: "different keys",
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key2": "value"},
			want: false,
		},
		{
			name: "different lengths",
			em1:  tax.Extensions{"key": "value", "key2": "value"},
			em2:  tax.Extensions{"key": "value"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.em1.Contains(tt.em2))
		})
	}
}

func TestExtensionsMerge(t *testing.T) {
	tests := []struct {
		name string
		em1  tax.Extensions
		em2  tax.Extensions
		want tax.Extensions
	}{
		{
			name: "empty",
			em1:  tax.Extensions{},
			em2:  tax.Extensions{},
			want: tax.Extensions{},
		},
		{
			name: "same",
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key": "value"},
			want: tax.Extensions{"key": "value"},
		},
		{
			name: "nil source",
			em1:  nil,
			em2:  tax.Extensions{"key": "value"},
			want: tax.Extensions{"key": "value"},
		},
		{
			name: "nil destination",
			em1:  tax.Extensions{"key": "value"},
			em2:  nil,
			want: tax.Extensions{"key": "value"},
		},
		{
			name: "different",
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key": "value2"},
			want: tax.Extensions{"key": "value2"},
		},
		{
			name: "different keys",
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key2": "value"},
			want: tax.Extensions{"key": "value", "key2": "value"},
		},
		{
			name: "different lengths",
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key": "value", "key2": "value"},
			want: tax.Extensions{"key": "value", "key2": "value"},
		},
		{
			name: "different lengths 2",
			em1:  tax.Extensions{"key": "value", "key2": "value"},
			em2:  tax.Extensions{"key": "value"},
			want: tax.Extensions{"key": "value", "key2": "value"},
		},
		{
			name: "different lengths 3",
			em1:  tax.Extensions{"key": "value2"},
			em2:  tax.Extensions{"key": "value", "key2": "value"},
			want: tax.Extensions{"key": "value", "key2": "value"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.em1.Merge(tt.em2))
		})
	}

}
