package oioubl

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/stretchr/testify/assert"
)

func TestOioublTaxCategory(t *testing.T) {
	assert.Equal(t, ExtValueTaxCategoryStandardRated, oioublTaxCategory("S"))
	assert.Equal(t, ExtValueTaxCategoryZeroRated, oioublTaxCategory("Z"))
	assert.Equal(t, ExtValueTaxCategoryZeroRated, oioublTaxCategory("E"), "exempt reports as ZeroRated")
	assert.Equal(t, ExtValueTaxCategoryReverseCharge, oioublTaxCategory("AE"))
	assert.Equal(t, cbc.Code(""), oioublTaxCategory("X"))
}

func TestOioublPaymentChannel(t *testing.T) {
	assert.Equal(t, ExtValuePaymentChannelGiro, oioublPaymentChannel("50"))
	assert.Equal(t, ExtValuePaymentChannelFIK, oioublPaymentChannel("93"))
	assert.Equal(t, ExtValuePaymentChannelIBAN, oioublPaymentChannel("30"), "settled means default to IBAN")
	assert.Equal(t, cbc.Code(""), oioublPaymentChannel("49"), "direct debit carries no channel")
	assert.Equal(t, cbc.Code(""), oioublPaymentChannel(""))
}
