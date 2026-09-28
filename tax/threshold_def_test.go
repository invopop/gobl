package tax_test

import (
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func baseThresholdDef() *tax.ThresholdDef {
	return &tax.ThresholdDef{
		Key:      "buyer-identity",
		Name:     i18n.NewString("Buyer Identity"),
		Currency: currency.AUD,
		Values: []*tax.ThresholdValueDef{
			{
				Since:  cal.NewDate(2027, 7, 1),
				Amount: num.MakeAmount(2000, 0),
			},
			{
				Amount: num.MakeAmount(1000, 0),
			},
		},
	}
}

func TestThresholdDefValue(t *testing.T) {
	td := baseThresholdDef()

	tests := []struct {
		name   string
		date   cal.Date
		amount string
	}{
		{"on since date", cal.MakeDate(2027, 7, 1), "2000"},
		{"day before since", cal.MakeDate(2027, 6, 30), "1000"},
		{"far future", cal.MakeDate(2099, 1, 1), "2000"},
		{"far past", cal.MakeDate(1990, 1, 1), "1000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tv := td.Value(tt.date)
			require.NotNil(t, tv)
			assert.Equal(t, tt.amount, tv.Amount.String())
		})
	}

	t.Run("before any value returns nil", func(t *testing.T) {
		td := baseThresholdDef()
		td.Values = td.Values[:1]
		assert.Nil(t, td.Value(cal.MakeDate(2027, 6, 30)))
	})

	t.Run("nil def", func(t *testing.T) {
		var td *tax.ThresholdDef
		assert.Nil(t, td.Value(cal.MakeDate(2027, 7, 1)))
	})
}

func TestThresholdDefReached(t *testing.T) {
	td := baseThresholdDef()
	before := cal.MakeDate(2027, 6, 30)
	after := cal.MakeDate(2027, 7, 1)

	assert.False(t, td.Reached(before, num.MakeAmount(99999, 2)))
	assert.True(t, td.Reached(before, num.MakeAmount(1000, 0)))
	assert.True(t, td.Reached(before, num.MakeAmount(100000, 2)))
	assert.False(t, td.Reached(after, num.MakeAmount(1000, 0)))
	assert.True(t, td.Reached(after, num.MakeAmount(2000, 0)))

	t.Run("no value", func(t *testing.T) {
		td := baseThresholdDef()
		td.Values = td.Values[:1]
		assert.False(t, td.Reached(before, num.MakeAmount(5000, 0)))
	})

	t.Run("nil def", func(t *testing.T) {
		var td *tax.ThresholdDef
		assert.False(t, td.Reached(after, num.MakeAmount(5000, 0)))
	})
}

func TestThresholdDefExceeded(t *testing.T) {
	td := baseThresholdDef()
	before := cal.MakeDate(2027, 6, 30)
	after := cal.MakeDate(2027, 7, 1)

	assert.False(t, td.Exceeded(before, num.MakeAmount(1000, 0)))
	assert.True(t, td.Exceeded(before, num.MakeAmount(100001, 2)))
	assert.False(t, td.Exceeded(after, num.MakeAmount(2000, 0)))
	assert.True(t, td.Exceeded(after, num.MakeAmount(2001, 0)))

	t.Run("nil def", func(t *testing.T) {
		var td *tax.ThresholdDef
		assert.False(t, td.Exceeded(after, num.MakeAmount(5000, 0)))
	})
}

func TestThresholdDefValidation(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		assert.NoError(t, rules.Validate(baseThresholdDef()))
	})

	t.Run("missing fields", func(t *testing.T) {
		err := rules.Validate(new(tax.ThresholdDef))
		assert.ErrorContains(t, err, "key is required")
		assert.ErrorContains(t, err, "name is required")
		assert.ErrorContains(t, err, "currency is required")
		assert.ErrorContains(t, err, "at least one value is required")
	})

	t.Run("single value without since", func(t *testing.T) {
		td := baseThresholdDef()
		td.Values = td.Values[1:]
		assert.NoError(t, rules.Validate(td))
	})

	t.Run("since missing before last", func(t *testing.T) {
		td := baseThresholdDef()
		td.Values[0], td.Values[1] = td.Values[1], td.Values[0]
		assert.ErrorContains(t, rules.Validate(td), "values must be in descending chronological order")
	})

	t.Run("ascending dates", func(t *testing.T) {
		td := baseThresholdDef()
		td.Values[1].Since = cal.NewDate(2028, 1, 1)
		assert.ErrorContains(t, rules.Validate(td), "values must be in descending chronological order")
	})

	t.Run("repeated dates", func(t *testing.T) {
		td := baseThresholdDef()
		td.Values[1].Since = cal.NewDate(2027, 7, 1)
		assert.ErrorContains(t, rules.Validate(td), "values must be in descending chronological order")
	})

	t.Run("nil value", func(t *testing.T) {
		td := baseThresholdDef()
		td.Values = append(td.Values, nil)
		assert.ErrorContains(t, rules.Validate(td), "values must be in descending chronological order")
	})
}

func TestRegimeDefThresholdDef(t *testing.T) {
	r := &tax.RegimeDef{
		Thresholds: []*tax.ThresholdDef{baseThresholdDef()},
	}
	assert.NotNil(t, r.ThresholdDef("buyer-identity"))
	assert.Nil(t, r.ThresholdDef(cbc.Key("unknown")))

	var nr *tax.RegimeDef
	assert.Nil(t, nr.ThresholdDef("buyer-identity"))
}

func TestAddonDefThresholdDef(t *testing.T) {
	ad := &tax.AddonDef{
		Thresholds: []*tax.ThresholdDef{baseThresholdDef()},
	}
	assert.NotNil(t, ad.ThresholdDef("buyer-identity"))
	assert.Nil(t, ad.ThresholdDef(cbc.Key("unknown")))

	var nad *tax.AddonDef
	assert.Nil(t, nad.ThresholdDef("buyer-identity"))
}
