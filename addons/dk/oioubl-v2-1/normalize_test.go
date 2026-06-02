package oioubl_test

import (
	"testing"

	oioubl "github.com/invopop/gobl/addons/dk/oioubl-v2-1"
	en16931 "github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/cef"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func bankPayment() *bill.PaymentDetails {
	return &bill.PaymentDetails{
		Terms: &pay.Terms{Notes: "Net 30 days"},
		Instructions: &pay.Instructions{
			Key:            pay.MeansKeyCreditTransfer,
			CreditTransfer: []*pay.CreditTransfer{{IBAN: "DK5000400440116243", BIC: "DABADKKK"}},
		},
	}
}

// TestNormalizeExemptToZeroRated checks that a VAT-exempt line keeps its GOBL
// exempt category (so EN 16931 still requires the reason) while carrying the
// OIOUBL ZeroRated category in the dk-oioubl-tax-category extension.
func TestNormalizeExemptToZeroRated(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Addons = tax.WithAddons(en16931.V2017, oioubl.V2_1)
	inv.Lines[0].Taxes = tax.Set{{
		Category: "VAT",
		Key:      tax.KeyExempt,
		Ext:      tax.ExtensionsOf(cbc.CodeMap{cef.ExtKeyVATEX: "VATEX-EU-132"}),
	}}
	inv.Payment = bankPayment()
	require.NoError(t, inv.Calculate())

	assert.Equal(t, "E", inv.Lines[0].Taxes[0].Ext.Get(untdid.ExtKeyTaxCategory).String(),
		"the GOBL category stays exempt")
	assert.Equal(t, "ZeroRated", inv.Lines[0].Taxes[0].Ext.Get(oioubl.ExtKeyTaxCategory).String(),
		"OIOUBL reports exempt as ZeroRated")
	require.NoError(t, rules.Validate(inv))
}

// TestNormalizeExemptStillRequiresReason confirms EN 16931's exemption-reason
// requirement is preserved (the GOBL category stays exempt), so the addon needs
// no rule of its own.
func TestNormalizeExemptStillRequiresReason(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Addons = tax.WithAddons(en16931.V2017, oioubl.V2_1)
	inv.Lines[0].Taxes = tax.Set{{Category: "VAT", Key: tax.KeyExempt}}
	inv.Payment = bankPayment()
	require.NoError(t, inv.Calculate())
	assert.ErrorContains(t, rules.Validate(inv), "exempt")
}

// TestNormalizeStandardUnchanged confirms the normalizer only touches exempt.
func TestNormalizeStandardUnchanged(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Addons = tax.WithAddons(en16931.V2017, oioubl.V2_1)
	inv.Payment = bankPayment()
	require.NoError(t, inv.Calculate())
	assert.Equal(t, "S", inv.Lines[0].Taxes[0].Ext.Get(untdid.ExtKeyTaxCategory).String())
	assert.Equal(t, "StandardRated", inv.Lines[0].Taxes[0].Ext.Get(oioubl.ExtKeyTaxCategory).String())
	assert.Equal(t, "IBAN", inv.Payment.Instructions.Ext.Get(oioubl.ExtKeyPaymentChannel).String(),
		"a bank transfer should carry the IBAN payment channel")
	require.NoError(t, rules.Validate(inv))
}
