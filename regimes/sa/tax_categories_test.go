package sa_test

import (
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/sa"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaxCategories(t *testing.T) {
	r := sa.New()
	require.NotNil(t, r)
	require.Len(t, r.Categories, 1)

	vat := r.Categories[0]
	assert.Equal(t, tax.CategoryVAT, vat.Code)

	require.Len(t, vat.Rates, 1)
	rate := vat.Rates[0]
	require.Len(t, rate.Values, 2)

	// 15% since July 2020
	assert.Equal(t, num.MakePercentage(150, 3), rate.Values[0].Percent)
	assert.Equal(t, cal.NewDate(2020, 7, 1), rate.Values[0].Since)

	// 5% since January 2018
	assert.Equal(t, num.MakePercentage(50, 3), rate.Values[1].Percent)
	assert.Equal(t, cal.NewDate(2018, 1, 1), rate.Values[1].Since)
}
