package tax_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
					Rate:     "standard",
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
					Rate:     "standard",
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
			desc: "duplicate",
			set: tax.Set{
				{
					Category: "VAT",
					Rate:     "standard",
					Percent:  num.NewPercentage(20, 3),
				},
				{
					Category: "VAT",
					Rate:     "reduced",
					Percent:  num.NewPercentage(20, 3),
				},
			},
			err: "duplicated",
		},
		{
			desc: "missing percentage",
			set: tax.Set{
				{
					Category: "VAT",
				},
			},
			err: "percent: cannot be blank",
		},
		{
			desc: "missing percentage with exempt rate",
			set: tax.Set{
				{
					Category: "VAT",
					Rate:     tax.RateExempt,
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
					Rate:     cbc.Key("invalid-tag"),
				},
			},
			err: "rate: must be a valid value.",
		},
		{
			desc: "missing percent with surcharge",
			set: tax.Set{
				{
					Category:  "VAT",
					Surcharge: num.NewPercentage(5, 3),
				},
			},
			err: "percent: cannot be blank; surcharge: required with percent.",
		},
		{
			desc: "exempt rate with reason",
			set: tax.Set{
				{
					Category: "VAT",
					Rate:     tax.RateExempt,
					Ext: tax.Extensions{
						es.ExtKeyTBAIExemption: "E1",
					},
				},
			},
			err: nil,
		},
		{
			desc: "exempt rate with invalid reason",
			set: tax.Set{
				{
					Category: "VAT",
					Rate:     tax.RateExempt,
					Ext: tax.Extensions{
						"foo": "E1",
					},
				},
			},
			err: "0: ext: (foo: invalid.)",
		},
		{
			desc: "category extension",
			set: tax.Set{
				{
					Category: "VAT",
					Rate:     tax.RateExempt,
					Ext: tax.Extensions{
						es.ExtKeyTBAIProduct: "services",
					},
				},
			},
			err: nil,
		},
	}
	es := es.New()
	ctx := es.WithContext(context.Background())
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
					Rate:     "standard",
				},
				{
					Category: "IRPF",
					Rate:     "pro",
				},
			},
			set2: tax.Set{
				{
					Category: "IRPF",
					Rate:     "pro",
				},
				{
					Category: "VAT",
					Rate:     "standard",
				},
			},
			res: true,
		},
		{
			desc: "bad comparison",
			set: tax.Set{
				{
					Category: "VAT",
					Rate:     "standard",
				},
			},
			set2: tax.Set{
				{
					Category: "IRPF",
					Rate:     "pro",
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
			Rate:     "standard",
		},
		{
			Category: "IRPF",
			Rate:     "pro",
		},
	}
	assert.Equal(t, s.Rate("VAT"), cbc.Key("standard"))
	assert.Equal(t, s.Rate("IRPF"), cbc.Key("pro"))
	assert.Empty(t, s.Rate("FOO"))
}

func TestSetGet(t *testing.T) {
	s := tax.Set{
		{
			Category: "VAT",
			Rate:     "standard",
		},
		{
			Category: "IRPF",
			Rate:     "pro",
		},
	}
	assert.NotNil(t, s.Get(cbc.Code("VAT")))
	assert.Nil(t, s.Get(cbc.Code("FOO")))
}

func TestComboUnmarshal(t *testing.T) {
	data := []byte(`{"cat":"VAT","tags":["standard"],"percent":"20%"}`)
	var c tax.Combo
	err := json.Unmarshal(data, &c)
	require.NoError(t, err)
	assert.Equal(t, c.Category, cbc.Code("VAT"))
	assert.Equal(t, c.Rate, cbc.Key("standard"))
}

func TestNormalizeSet(t *testing.T) {
	s := tax.NormalizeSet(nil)
	assert.Nil(t, s)

	s = tax.Set{
		{
			Category: "VAT",
			Rate:     "standard",
		},
		{
			Category: "IRPF",
			Rate:     "pro",
		},
	}
	s = tax.NormalizeSet(s)
	assert.Equal(t, s[0].Category, cbc.Code("VAT"))
	assert.Equal(t, s[1].Category, cbc.Code("IRPF"))

	s = tax.Set{
		{
			Category: "VAT",
			Rate:     "standard",
			Ext: tax.Extensions{
				es.ExtKeyFacturaECorrection: "",
			},
		},
		nil,
	}
	assert.NotNil(t, s[0].Ext)
	assert.Len(t, s, 2)
	s = tax.NormalizeSet(s)
	assert.Nil(t, s[0].Ext)
	assert.Len(t, s, 1)

	s = tax.Set{
		nil,
	}
	assert.Len(t, s, 1)
	s = tax.NormalizeSet(s)
	assert.Nil(t, s)
}
