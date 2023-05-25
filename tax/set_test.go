package tax_test

import (
	"context"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
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
					Rate:     es.TaxRateExempt.With(es.TaxRateArticle20),
				},
			},
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
