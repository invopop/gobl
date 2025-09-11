package tax_test

import (
	"context"
	"testing"

	"github.com/invopop/gobl/addons/es/tbai"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
)

func TestSetValidation(t *testing.T) {
	var tests = []struct {
		desc string
		set  tax.Set
		err  interface{}
	}{
		{
			desc: "simple success",
			set: tax.Set{
				{
					Category: "VAT",
					Key:      "standard",
					Percent:  num.NewPercentage(20, 3),
				},
			},
			err: nil,
		},
		{
			desc: "empty success",
			set:  tax.Set{},
			err:  nil,
		},
		{
			desc: "complex success",
			set: tax.Set{
				{
					Category: "VAT",
					Key:      "standard",
					Percent:  num.NewPercentage(20, 3),
				},
				{
					Category: "IRPF",
					Rate:     "pro",
					Percent:  num.NewPercentage(15, 3),
				},
			},
		},
		{
			desc: "other country",
			set: tax.Set{
				{
					Category: "VAT",
					Country:  "NL",
					Key:      "standard",
					Rate:     "general",
					Percent:  num.NewPercentage(20, 3),
				},
			},
			err: nil,
		},
		{
			desc: "other country no percent",
			set: tax.Set{
				{
					Category: "VAT",
					Country:  "NL",
				},
			},
			err: "0: (percent: cannot be blank.).",
		},
		{
			desc: "exempt rate with percent",
			set: tax.Set{
				{
					Category: "VAT",
					Key:      tax.KeyExempt,
					Percent:  num.NewPercentage(5, 3),
				},
			},
			err: "0: (percent: must be nil for 'exempt' in 'VAT'.)",
		},
		{
			desc: "duplicate",
			set: tax.Set{
				{
					Category: "VAT",
					Key:      "standard",
					Percent:  num.NewPercentage(20, 3),
				},
				{
					Category: "VAT",
					Key:      "reduced",
					Percent:  num.NewPercentage(20, 3),
				},
			},
			err: "duplicated",
		},
		{
			desc: "VAT missing percentage",
			set: tax.Set{
				{
					Category: "VAT",
				},
			},
			err: "0: (percent: cannot be blank.)",
		},
		{
			desc: "VAT missing percentage with key",
			set: tax.Set{
				{
					Category: "VAT",
					Key:      "standard",
				},
			},
			err: "0: (percent: required for 'standard' in 'VAT'.).",
		},
		{
			desc: "IRPF missing percentage",
			set: tax.Set{
				{
					Category: "IRPF",
				},
			},
			err: "0: (percent: cannot be blank.)",
		},
		{
			desc: "missing percentage with exempt rate",
			set: tax.Set{
				{
					Category: "VAT",
					Key:      tax.KeyExempt,
				},
			},
			err: nil, // this is okay
		},
		{
			desc: "undefined category code",
			set: tax.Set{
				{
					Category: "VAT2",
					Percent:  num.NewPercentage(20, 3),
				},
			},
			err: "cat: must be a valid value",
		},
		{
			desc: "undefined category rate",
			set: tax.Set{
				{
					Category: "VAT",
					Percent:  num.NewPercentage(20, 3),
					Key:      cbc.Key("invalid-tag"),
				},
			},
			err: "0: (key: must be a valid value.).",
		},
		{
			desc: "rate with extension",
			set: tax.Set{
				{
					Category: "VAT",
					Key:      "standard",
					Percent:  num.NewPercentage(20, 3),
					Rate:     tax.RateGeneral.With("eqs"),
				},
			},
			err: nil,
		},
		{
			desc: "missing percent with surcharge",
			set: tax.Set{
				{
					Category:  "VAT",
					Surcharge: num.NewPercentage(5, 3),
				},
			},
			err: "surcharge: required with percent.",
		},
		{
			desc: "exempt key with reason",
			set: tax.Set{
				{
					Category: "VAT",
					Key:      tax.KeyExempt,
					Ext: tax.Extensions{
						tbai.ExtKeyExempt: "E1",
					},
				},
			},
			err: nil,
		},
		{
			desc: "exempt, no key, with extension",
			set: tax.Set{
				{
					Category: "VAT",
					// The correct key would be set
					// here automatically in normalization.
					Ext: tax.Extensions{
						tbai.ExtKeyExempt: "E1",
					},
				},
			},
			err: "0: (percent: cannot be blank.).",
		},
		{
			desc: "exempt rate with invalid reason",
			set: tax.Set{
				{
					Category: "VAT",
					Key:      tax.KeyExempt,
					Ext: tax.Extensions{
						"foo": "E1",
					},
				},
			},
			err: "0: (ext: (foo: undefined.).)",
		},
		{
			desc: "category extension",
			set: tax.Set{
				{
					Category: "VAT",
					Key:      tax.KeyExempt,
					Ext: tax.Extensions{
						tbai.ExtKeyProduct: "services",
					},
				},
			},
			err: nil,
		},
	}
	ctx := context.Background()
	ctx = es.New().WithContext(ctx)
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Helper()
			err := test.set.ValidateWithContext(ctx)
			if test.err == nil {
				assert.NoError(t, err)
			} else if e, ok := test.err.(error); ok {
				assert.ErrorIs(t, err, e)
			} else if s, ok := test.err.(string); ok {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), s)
				}
			}
		})
	}
}

func TestSetEquals(t *testing.T) {
	var tests = []struct {
		desc string
		set  tax.Set
		set2 tax.Set
		res  bool
	}{
		{
			desc: "simple comparison",
			set: tax.Set{
				{
					Category: "VAT",
					Key:      "standard",
				},
				{
					Category: "IRPF",
					Key:      "pro",
				},
			},
			set2: tax.Set{
				{
					Category: "IRPF",
					Key:      "pro",
				},
				{
					Category: "VAT",
					Key:      "standard",
				},
			},
			res: true,
		},
		{
			desc: "bad comparison",
			set: tax.Set{
				{
					Category: "VAT",
					Key:      "standard",
				},
			},
			set2: tax.Set{
				{
					Category: "IRPF",
					Key:      "pro",
				},
			},
			res: false,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Helper()
			res := test.set.Equals(test.set2)
			assert.Equal(t, test.res, res)
		})
	}
}

func TestSetRate(t *testing.T) {
	s := tax.Set{
		{
			Category: "VAT",
			Key:      "standard",
		},
		{
			Category: "IRPF",
			Key:      "pro",
		},
	}
	assert.Equal(t, s.Key("VAT"), cbc.Key("standard"))
	assert.Equal(t, s.Key("IRPF"), cbc.Key("pro"))
	assert.Empty(t, s.Key("FOO"))
}

func TestSetGet(t *testing.T) {
	s := tax.Set{
		{
			Category: "VAT",
			Key:      "standard",
		},
		{
			Category: "IRPF",
			Key:      "pro",
		},
	}
	assert.NotNil(t, s.Get(cbc.Code("VAT")))
	assert.Nil(t, s.Get(cbc.Code("FOO")))
}

func TestCleanSet(t *testing.T) {
	s := tax.CleanSet(nil)
	assert.Nil(t, s)

	s = tax.Set{
		{
			Category: "VAT",
			Key:      "standard",
		},
		{
			Category: "IRPF",
			Key:      "pro",
		},
	}
	s = tax.CleanSet(s)
	assert.Equal(t, s[0].Category, cbc.Code("VAT"))
	assert.Equal(t, s[1].Category, cbc.Code("IRPF"))

	s = tax.Set{
		{
			Category: "VAT",
			Key:      "standard",
		},
		nil,
	}
	assert.Len(t, s, 2)
	s = tax.CleanSet(s)
	assert.Len(t, s, 1)

	s = tax.Set{
		nil,
	}
	assert.Len(t, s, 1)
	s = tax.CleanSet(s)
	assert.Nil(t, s)
}

func TestSetHasCategory(t *testing.T) {
	s := tax.Set{
		{
			Category: "VAT",
			Key:      "standard",
		},
		{
			Category: "IRPF",
			Key:      "pro",
		},
	}
	t.Run("has VAT", func(t *testing.T) {
		err := validation.Validate(s, tax.SetHasCategory("VAT"))
		assert.NoError(t, err)
	})
	t.Run("has multiple", func(t *testing.T) {
		err := validation.Validate(s, tax.SetHasCategory(tax.CategoryVAT, es.TaxCategoryIRPF))
		assert.NoError(t, err)
	})
	t.Run("missing category", func(t *testing.T) {
		err := validation.Validate(s, tax.SetHasCategory("FOO"))
		assert.Error(t, err)
		assert.Equal(t, "missing category FOO", err.Error())
	})
	t.Run("with different type", func(t *testing.T) {
		var s2 string
		assert.NotPanics(t, func() {
			err := validation.Validate(s2, tax.SetHasCategory("FOO"))
			assert.NoError(t, err)
		})
	})
}
