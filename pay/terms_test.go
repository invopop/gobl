package pay_test

import (
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTermsValidation(t *testing.T) {
	tm := new(pay.Terms)
	tm.Key = cbc.Key("foo")
	err := tm.Validate()
	assert.Error(t, err, "expected validation error")

	tm.Key = cbc.Key("due_date")
	err = tm.Validate()
	assert.Error(t, err, "expected validation error")
	assert.Contains(t, err.Error(), "key: must be a valid value")

	tm.Key = pay.TermKeyAdvanced
	err = tm.Validate()
	assert.NoError(t, err)

	tm.Key = ""
	err = tm.Validate()
	assert.NoError(t, err)

	t.Run("with due dates and missing amount", func(t *testing.T) {
		tm := new(pay.Terms)
		tm.Key = pay.TermKeyDueDate
		tm.DueDates = []*pay.DueDate{
			{
				Date: cal.NewDate(2021, 11, 10),
			},
		}
		err := tm.Validate()
		assert.ErrorContains(t, err, "due_dates: (0: (amount: must not be zero.).)")
	})
}

func TestTermsUNTDID4279(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		tm := new(pay.Terms)
		tm.Key = pay.TermKeyEndOfMonth
		assert.Equal(t, "2", tm.UNTDID4279().String())
	})
	t.Run("non-existing", func(t *testing.T) {
		tm := new(pay.Terms)
		tm.Key = cbc.Key("non-existing")
		assert.Equal(t, cbc.CodeEmpty, tm.UNTDID4279())
	})
}

func TestTermsNormalize(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		pt := &pay.Terms{
			Key: pay.TermKeyUndefined,
			Ext: tax.Extensions{
				"random": "",
			},
		}
		pt.Normalize(nil)
		assert.Empty(t, pt.Ext)
		assert.Equal(t, "undefined", pt.Key.String())
	})
	t.Run("nil", func(t *testing.T) {
		var pt *pay.Terms
		assert.NotPanics(t, func() {
			pt.Normalize(nil)
		})
	})
}

func TestTermsCalculateDues(t *testing.T) {
	sum := num.MakeAmount(10000, 2)
	var terms *pay.Terms
	zero := num.MakeAmount(0, 2)
	terms.CalculateDues(zero, sum) // Should not panic
	terms = new(pay.Terms)
	terms.DueDates = []*pay.DueDate{
		{
			Date:    cal.NewDate(2021, 11, 10),
			Percent: num.NewPercentage(40, 2),
		},
		{
			Date:    cal.NewDate(2021, 12, 10),
			Percent: num.NewPercentage(60, 2),
		},
	}
	terms.CalculateDues(zero, sum)

	assert.Equal(t, num.MakeAmount(4000, 2), terms.DueDates[0].Amount)
	assert.Equal(t, num.MakeAmount(6000, 2), terms.DueDates[1].Amount)

	terms.DueDates = []*pay.DueDate{
		{
			Date:   cal.NewDate(2021, 11, 10),
			Amount: num.MakeAmount(40, 0),
		},
	}
	terms.CalculateDues(zero, sum)
	assert.Equal(t, "40.00", terms.DueDates[0].Amount.String(), "should normalize amounts for currency")
}

func TestTermsJSONSchemaExtend(t *testing.T) {
	schema := &jsonschema.Schema{
		Properties: jsonschema.NewProperties(),
	}
	schema.Properties.Set("key", &jsonschema.Schema{
		Type: "string",
	})
	terms := &pay.Terms{}
	terms.JSONSchemaExtend(schema)
	prop, ok := schema.Properties.Get("key")
	require.True(t, ok)
	assert.Len(t, prop.OneOf, 10)
	assert.Equal(t, cbc.Key("end-of-month"), prop.OneOf[0].Const)
}
