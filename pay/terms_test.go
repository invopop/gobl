package pay

import (
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/stretchr/testify/assert"
)

func TestTermsValidation(t *testing.T) {
	tm := new(Terms)
	tm.Key = cbc.Key("foo")
	err := tm.Validate()
	assert.Error(t, err, "expected validation error")

	tm.Key = cbc.Key("due_date")
	err = tm.Validate()
	assert.Error(t, err, "expected validation error")
	assert.Contains(t, err.Error(), "key: must be a valid value")

	tm.Key = TermKeyAdvanced
	err = tm.Validate()
	assert.NoError(t, err)

	tm.Key = TermKeyNA
	err = tm.Validate()
	assert.NoError(t, err)
}

func TestTermsCalculateDues(t *testing.T) {
	sum := num.MakeAmount(10000, 2)
	var terms *Terms
	zero := num.MakeAmount(0, 2)
	terms.CalculateDues(zero, sum) // Should not panic
	terms = new(Terms)
	terms.DueDates = []*DueDate{
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

	terms.DueDates = []*DueDate{
		{
			Date:   cal.NewDate(2021, 11, 10),
			Amount: num.MakeAmount(40, 0),
		},
	}
	terms.CalculateDues(zero, sum)
	assert.Equal(t, "40.00", terms.DueDates[0].Amount.String(), "should normalize amounts for currency")
}
