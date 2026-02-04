package nz_test

import (
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/regimes/nz"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaxCategories(t *testing.T) {
	r := nz.New()

	gst := r.CategoryDef(tax.CategoryGST)
	require.NotNil(t, gst, "GST category should exist")

	assert.NotEmpty(t, gst.Keys, "GST category should have keys defined")

	standardRate := gst.RateDef(tax.KeyStandard, tax.RateGeneral)
	require.NotNil(t, standardRate, "Standard rate should exist")
	assert.Equal(t, "15.0%", standardRate.Values[0].Percent.String())

	accommodationRate := gst.RateDef(tax.KeyStandard, nz.TaxRateAccommodation)
	require.NotNil(t, accommodationRate, "Accommodation rate should exist")
	assert.Equal(t, "9.0%", accommodationRate.Values[0].Percent.String())

	zeroRate := gst.RateDef(tax.KeyZero, tax.RateZero)
	require.NotNil(t, zeroRate, "Zero rate should exist")
	assert.Equal(t, "0.0%", zeroRate.Values[0].Percent.String())
}

func TestAccommodationRateEffectiveDate(t *testing.T) {
	r := nz.New()
	gst := r.CategoryDef(tax.CategoryGST)
	require.NotNil(t, gst)

	accommodationRate := gst.RateDef(tax.KeyStandard, nz.TaxRateAccommodation)
	require.NotNil(t, accommodationRate)

	t.Run("After2024April", func(t *testing.T) {
		d := cal.MakeDate(2024, 6, 1)
		v := accommodationRate.Value(d, nil)
		require.NotNil(t, v)
		assert.Equal(t, "9.0%", v.Percent.String())
	})

	t.Run("Before2024April", func(t *testing.T) {
		d := cal.MakeDate(2024, 1, 1)
		v := accommodationRate.Value(d, nil)
		assert.Nil(t, v, "Accommodation rate should not exist before 2024-04-01")
	})
}
