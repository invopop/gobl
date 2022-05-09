package tax_test

import (
	"testing"

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
				},
				{
					Category: "IRPF",
					Rate:     "pro",
				},
			},
		},
		{
			desc: "duplicate",
			set: tax.Set{
				{
					Category: "VAT",
					Rate:     "standard",
				},
				{
					Category: "VAT",
					Rate:     "reduced",
				},
			},
			err: "duplicated",
		},
		{
			desc: "invalid category",
			set: tax.Set{
				{
					Category: "foo-cat",
					Rate:     "standard",
				},
			},
			err: "cat: the length must be between",
		},
		{
			desc: "invalid rate",
			set: tax.Set{
				{
					Category: "VAT",
					Rate:     "STD",
				},
			},
			err: "rate: must be in a valid format",
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Helper()
			err := test.set.Validate()
			if test.err == nil {
				assert.NoError(t, err)
			} else if e, ok := test.err.(error); ok {
				assert.ErrorIs(t, err, e)
			} else if s, ok := test.err.(string); ok {
				assert.Contains(t, err.Error(), s)
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
	assert.Equal(t, s.Rate("VAT"), tax.Key("standard"))
	assert.Equal(t, s.Rate("IRPF"), tax.Key("pro"))
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
	assert.NotNil(t, s.Get(tax.Code("VAT")))
	assert.Nil(t, s.Get(tax.Code("FOO")))
}
