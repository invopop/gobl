package tax_test

import (
	"testing"
	"time"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegimeTimeLocation(t *testing.T) {
	r := new(tax.RegimeDef)
	r.TimeZone = "Europe/Amsterdam"
	loc, err := time.LoadLocation("Europe/Amsterdam")
	require.NoError(t, err)

	assert.Equal(t, loc, r.TimeLocation())

	r.TimeZone = "INVALID"
	loc = r.TimeLocation()
	assert.Equal(t, loc, time.UTC)
}

func TestRegimeInCategoryRates(t *testing.T) {
	var r *tax.RegimeDef // nil regime
	rate := cbc.Key("standard")
	err := validation.Validate(rate, r.InCategoryRates(tax.CategoryVAT))
	assert.ErrorContains(t, err, "must be blank when regime is undefine")
}
