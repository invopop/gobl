package flow10

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testPaymentB2B(t *testing.T) *bill.Payment {
	t.Helper()
	issued := cal.MakeDate(2026, 1, 10)
	paid := cal.MakeDate(2026, 2, 1)
	return &bill.Payment{
		Regime:    tax.WithRegime("FR"),
		Addons:    tax.WithAddons(V1),
		Type:      bill.PaymentTypeReceipt,
		Code:      "PAY-2026-001",
		Currency:  "EUR",
		IssueDate: cal.MakeDate(2026, 2, 1),
		ValueDate: &paid,
		Method:    &pay.Instructions{Key: pay.MeansKeyCreditTransfer},
		Supplier:  frPartyWithSIREN(),
		Customer:  frCustomerWithSIREN(),
		Lines: []*bill.PaymentLine{
			{
				Document: &org.DocumentRef{
					Code:      "INV-2026-001",
					IssueDate: &issued,
				},
				Amount: num.MakeAmount(12000, 2),
			},
		},
	}
}

func testPaymentB2C(t *testing.T) *bill.Payment {
	t.Helper()
	paid := cal.MakeDate(2026, 2, 1)
	return &bill.Payment{
		Regime:    tax.WithRegime("FR"),
		Addons:    tax.WithAddons(V1),
		Tags:      tax.WithTags(TagB2C),
		Type:      bill.PaymentTypeReceipt,
		Code:      "PAY-2026-B2C-001",
		Currency:  "EUR",
		IssueDate: cal.MakeDate(2026, 2, 1),
		ValueDate: &paid,
		Method:    &pay.Instructions{Key: pay.MeansKeyCreditTransfer},
		Supplier:  frPartyWithSIREN(),
		Lines: []*bill.PaymentLine{
			{
				Amount: num.MakeAmount(12000, 2),
			},
		},
	}
}

func TestPaymentB2BHappyPath(t *testing.T) {
	p := testPaymentB2B(t)
	require.NoError(t, p.Calculate())
	require.NoError(t, rules.Validate(p))
}

func TestPaymentB2CHappyPath(t *testing.T) {
	p := testPaymentB2C(t)
	require.NoError(t, p.Calculate())
	require.NoError(t, rules.Validate(p))
}

func TestPaymentMissingValueDate(t *testing.T) {
	p := testPaymentB2B(t)
	p.ValueDate = nil
	require.NoError(t, p.Calculate())
	err := rules.Validate(p)
	assert.ErrorContains(t, err, "value_date")
}

func TestPaymentB2BRequiresDocumentRef(t *testing.T) {
	p := testPaymentB2B(t)
	p.Lines[0].Document = nil
	require.NoError(t, p.Calculate())
	err := rules.Validate(p)
	assert.ErrorContains(t, err, "document")
}

func TestPaymentB2BRequiresDocumentCode(t *testing.T) {
	p := testPaymentB2B(t)
	p.Lines[0].Document.Code = ""
	require.NoError(t, p.Calculate())
	err := rules.Validate(p)
	assert.ErrorContains(t, err, "invoice ID")
}

func TestPaymentB2BRequiresDocumentIssueDate(t *testing.T) {
	p := testPaymentB2B(t)
	p.Lines[0].Document.IssueDate = nil
	require.NoError(t, p.Calculate())
	err := rules.Validate(p)
	assert.ErrorContains(t, err, "invoice issue date")
}

func TestPaymentB2CDoesNotRequireDocumentRef(t *testing.T) {
	p := testPaymentB2C(t)
	// A B2C payment line has no Document at all — should still pass.
	require.NoError(t, p.Calculate())
	require.NoError(t, rules.Validate(p))
}

func TestPaymentSupplierSIRENRequired(t *testing.T) {
	p := testPaymentB2B(t)
	p.Supplier.TaxID = nil
	p.Supplier.Identities = nil
	require.NoError(t, p.Calculate())
	err := rules.Validate(p)
	assert.ErrorContains(t, err, "SIREN")
}

func TestPaymentVATRateNotInWhitelist(t *testing.T) {
	p := testPaymentB2B(t)
	pct := num.MakePercentage(17, 2) // 17%, not allowed
	p.Lines[0].Tax = &tax.Total{
		Categories: []*tax.CategoryTotal{
			{
				Code: tax.CategoryVAT,
				Rates: []*tax.RateTotal{
					{Percent: &pct, Base: num.MakeAmount(10000, 2), Amount: num.MakeAmount(1700, 2)},
				},
			},
		},
	}
	require.NoError(t, p.Calculate())
	err := rules.Validate(p)
	assert.ErrorContains(t, err, "G1.24")
}

func TestPaymentRejectsNonReceiptType(t *testing.T) {
	p := testPaymentB2B(t)
	p.Type = bill.PaymentTypeRequest
	require.NoError(t, p.Calculate())
	err := rules.Validate(p)
	assert.ErrorContains(t, err, "payment type must be 'receipt'")
}

// --- Internal helper coverage (bill.go) ---------------------------------

func TestPaymentIsB2BWrongType(t *testing.T) {
	assert.False(t, paymentIsB2BAny("x"))
}

func TestPaymentVATRatesAllowedWrongType(t *testing.T) {
	assert.True(t, paymentVATRatesAllowed("x"))
}

func TestPaymentVATRatesAllowedNilLine(t *testing.T) {
	p := &bill.Payment{Lines: []*bill.PaymentLine{nil}}
	assert.True(t, paymentVATRatesAllowed(p))
}
