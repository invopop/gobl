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
// exempt category while carrying the OIOUBL ZeroRated category in the
// dk-oioubl-tax-category extension. A VATEX reason remains allowed (and is
// carried through), even though OIOUBL no longer requires one.
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

// TestNormalizeExemptNeedsNoReason confirms that, with the OIOUBL addon present,
// EN 16931's exemption-reason requirement is relaxed: OIOUBL 2.1 has no exempt
// category (exempt is reported as ZeroRated, which requires no reason), so a
// VAT-exempt line with neither a VATEX code nor an exemption note validates.
func TestNormalizeExemptNeedsNoReason(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Addons = tax.WithAddons(en16931.V2017, oioubl.V2_1)
	inv.Lines[0].Taxes = tax.Set{{Category: "VAT", Key: tax.KeyExempt}}
	inv.Payment = bankPayment()
	require.NoError(t, inv.Calculate())
	assert.NoError(t, rules.Validate(inv))
}

// TestNormalizeReverseChargeNeedsNoReason confirms the same relaxation for
// reverse-charge: OIOUBL reports it as the ReverseCharge category, which carries
// no exemption reason, so the EN 16931 exemption-note requirement is skipped.
func TestNormalizeReverseChargeNeedsNoReason(t *testing.T) {
	inv := testInvoiceStandard(t)
	inv.Addons = tax.WithAddons(en16931.V2017, oioubl.V2_1)
	inv.Lines[0].Taxes = tax.Set{{Category: "VAT", Key: tax.KeyReverseCharge}}
	inv.Payment = bankPayment()
	require.NoError(t, inv.Calculate())
	assert.NoError(t, rules.Validate(inv))
	assert.Equal(t, "ReverseCharge", inv.Lines[0].Taxes[0].Ext.Get(oioubl.ExtKeyTaxCategory).String())
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
