package cal_test

import (
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/stretchr/testify/assert"
)

func TestPeriodValidation(t *testing.T) {
	p := &cal.Period{}
	assert.Error(t, p.Validate())

	p = &cal.Period{
		Start: cal.MakeDate(2022, 1, 25),
		End:   cal.MakeDate(2022, 2, 28),
	}
	assert.NoError(t, p.Validate())
	p = &cal.Period{
		Start: cal.MakeDate(2022, 1, 25),
		End:   cal.MakeDate(2022, 1, 25),
	}
	assert.NoError(t, p.Validate(), "allow same day")

	p = &cal.Period{
		Start: cal.MakeDate(2022, 1, 25),
	}
	assert.Error(t, p.Validate())

	p = &cal.Period{
		Start: cal.MakeDate(2022, 1, 25),
		End:   cal.MakeDate(2022, 1, 20),
	}
	assert.Error(t, p.Validate())
}
