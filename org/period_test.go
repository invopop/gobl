package org_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestPeriodValidation(t *testing.T) {
	p := &org.Period{}
	assert.Error(t, p.Validate())

	p = &org.Period{
		Start: org.MakeDate(2022, 1, 25),
		End:   org.MakeDate(2022, 2, 28),
	}
	assert.NoError(t, p.Validate())
	p = &org.Period{
		Start: org.MakeDate(2022, 1, 25),
		End:   org.MakeDate(2022, 1, 25),
	}
	assert.NoError(t, p.Validate(), "allow same day")

	p = &org.Period{
		Start: org.MakeDate(2022, 1, 25),
	}
	assert.Error(t, p.Validate())

	p = &org.Period{
		Start: org.MakeDate(2022, 1, 25),
		End:   org.MakeDate(2022, 1, 20),
	}
	assert.Error(t, p.Validate())
}
