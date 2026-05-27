package flow6

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// paymentParty returns a French party with a SIREN identity, used as
// supplier or customer on a Flow 6 payment.
func paymentParty(name, siren string) *org.Party {
	return &org.Party{
		Name: name,
		Identities: []*org.Identity{{
			Code: cbc.Code(siren),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				iso.ExtKeySchemeID: identitySchemeIDSIREN,
			}),
		}},
	}
}

func testPaymentReceipt(t *testing.T) *bill.Payment {
	t.Helper()
	issue := cal.MakeDate(2026, 5, 2)
	return &bill.Payment{
		Regime:    tax.WithRegime("FR"),
		Addons:    tax.WithAddons(V1),
		IssueDate: issue,
		Code:      "PMT-2026-0001",
		Currency:  "EUR",
		Type:      bill.PaymentTypeReceipt,
		Supplier:  paymentParty("VENDEUR SARL", "732829320"),
		Customer:  paymentParty("ACHETEUR SARL", "200000008"),
		Methods:   []*pay.Record{{Key: pay.MeansKeyCreditTransfer}},
		Lines: []*bill.PaymentLine{{
			Amount: num.MakeAmount(120000, 2),
			Document: &org.DocumentRef{
				Code:      "2026-00042",
				IssueDate: cal.NewDate(2026, 4, 15),
			},
		}},
	}
}

func TestPaymentReceiptHappyPath(t *testing.T) {
	pmt := testPaymentReceipt(t)
	runNormalize(t, pmt)
	require.NoError(t, rules.Validate(pmt))
}

func TestPaymentReceiptSetsCDARStatusCode212(t *testing.T) {
	pmt := testPaymentReceipt(t)
	runNormalize(t, pmt)
	assert.Equal(t, cbc.Code("212"), pmt.Ext.Get(ExtKeyStatus))
}

// Default Condition extension on a receipt is MEN (Amount received).
func TestPaymentReceiptDefaultsConditionToMEN(t *testing.T) {
	pmt := testPaymentReceipt(t)
	runNormalize(t, pmt)
	assert.Equal(t, ConditionAmountReceived, pmt.Ext.Get(ExtKeyCondition))
}

// Default Condition extension on an advice is MPA (Amount paid).
func TestPaymentAdviceDefaultsConditionToMPA(t *testing.T) {
	pmt := testPaymentReceipt(t)
	pmt.Type = bill.PaymentTypeAdvice
	runNormalize(t, pmt)
	assert.Equal(t, ConditionAmountPaid, pmt.Ext.Get(ExtKeyCondition))
}

// Partial-payment scenario: caller pins RAP (Amount remaining); the
// SetOneOf chain keeps the explicit override.
func TestPaymentAcceptsRAPOverride(t *testing.T) {
	pmt := testPaymentReceipt(t)
	pmt.Ext = pmt.Ext.Set(ExtKeyCondition, ConditionAmountRemaining)
	runNormalize(t, pmt)
	require.NoError(t, rules.Validate(pmt))
	assert.Equal(t, ConditionAmountRemaining, pmt.Ext.Get(ExtKeyCondition))
}

// Status-only Condition codes are rejected on a Payment.
func TestPaymentRejectsStatusOnlyConditionCodes(t *testing.T) {
	for _, code := range []cbc.Code{
		ConditionBankDetailsUpdate, ConditionInvalidData,
		ConditionExpectedData, ConditionReplacementData,
		ConditionAmountApprovedHT, ConditionDiscount,
	} {
		pmt := testPaymentReceipt(t)
		runNormalize(t, pmt)
		// Replace the normalized MEN default with a Status-only code.
		pmt.Ext = pmt.Ext.Set(ExtKeyCondition, code)
		err := rules.Validate(pmt)
		assert.ErrorContains(t, err, "Payment-applicable", "code %s", code)
	}
}

// Status-only ProcessConditionCodes are rejected on a Payment.
func TestPaymentRejectsStatusProcessCodes(t *testing.T) {
	pmt := testPaymentReceipt(t)
	runNormalize(t, pmt)
	pmt.Ext = pmt.Ext.Set(ExtKeyStatus, "205") // Approved — Status-only
	err := rules.Validate(pmt)
	assert.ErrorContains(t, err, "Payment-applicable")
}

func TestPaymentAdviceSetsCDARStatusCode211(t *testing.T) {
	pmt := testPaymentReceipt(t)
	pmt.Type = bill.PaymentTypeAdvice
	runNormalize(t, pmt)
	assert.Equal(t, cbc.Code("211"), pmt.Ext.Get(ExtKeyStatus))
}

func TestPaymentReceiptDefaultsSupplierRoleSeller(t *testing.T) {
	pmt := testPaymentReceipt(t)
	runNormalize(t, pmt)
	assert.Equal(t, RoleSeller, pmt.Supplier.Ext.Get(ExtKeyRole))
	assert.Equal(t, RoleBuyer, pmt.Customer.Ext.Get(ExtKeyRole))
}

func TestPaymentAdviceFlipsRoles(t *testing.T) {
	pmt := testPaymentReceipt(t)
	pmt.Type = bill.PaymentTypeAdvice
	runNormalize(t, pmt)
	// Advice = payer-issued: customer (payee) becomes SE, supplier
	// (payer in the payment-doc sense) becomes BY.
	assert.Equal(t, RoleSeller, pmt.Customer.Ext.Get(ExtKeyRole))
	assert.Equal(t, RoleBuyer, pmt.Supplier.Ext.Get(ExtKeyRole))
}

func TestPaymentRejectsRequestType(t *testing.T) {
	pmt := testPaymentReceipt(t)
	pmt.Type = bill.PaymentTypeRequest
	runNormalize(t, pmt)
	err := rules.Validate(pmt)
	assert.ErrorContains(t, err, "advice")
}

func TestPaymentRequiresSupplierSIREN(t *testing.T) {
	pmt := testPaymentReceipt(t)
	pmt.Supplier.Identities = nil
	runNormalize(t, pmt)
	err := rules.Validate(pmt)
	assert.ErrorContains(t, err, "SIREN")
}

func TestPaymentRequiresCustomerSIREN(t *testing.T) {
	pmt := testPaymentReceipt(t)
	pmt.Customer.Identities = nil
	runNormalize(t, pmt)
	err := rules.Validate(pmt)
	assert.ErrorContains(t, err, "SIREN")
}

func TestPaymentRequiresExactlyOneLine(t *testing.T) {
	pmt := testPaymentReceipt(t)
	pmt.Lines = append(pmt.Lines, &bill.PaymentLine{
		Amount: num.MakeAmount(5000, 2),
		Document: &org.DocumentRef{
			Code:      "2026-00043",
			IssueDate: cal.NewDate(2026, 4, 15),
		},
	})
	runNormalize(t, pmt)
	err := rules.Validate(pmt)
	assert.ErrorContains(t, err, "exactly one")
}

func TestPaymentRequiresDocumentReference(t *testing.T) {
	pmt := testPaymentReceipt(t)
	pmt.Lines[0].Document = nil
	runNormalize(t, pmt)
	err := rules.Validate(pmt)
	assert.ErrorContains(t, err, "payment line document is required")
}

func TestPaymentRequiresDocumentCode(t *testing.T) {
	pmt := testPaymentReceipt(t)
	pmt.Lines[0].Document.Code = ""
	runNormalize(t, pmt)
	err := rules.Validate(pmt)
	assert.ErrorContains(t, err, "payment line document code")
}

func TestPaymentRequiresDocumentIssueDate(t *testing.T) {
	pmt := testPaymentReceipt(t)
	pmt.Lines[0].Document.IssueDate = nil
	runNormalize(t, pmt)
	err := rules.Validate(pmt)
	assert.ErrorContains(t, err, "payment line document issue_date")
}

func TestPaymentRejectsSTCIdentityScheme(t *testing.T) {
	pmt := testPaymentReceipt(t)
	pmt.Supplier.Identities = append(pmt.Supplier.Identities, &org.Identity{
		Code: "12345678",
		Ext: tax.ExtensionsOf(cbc.CodeMap{
			iso.ExtKeySchemeID: "0231",
		}),
	})
	runNormalize(t, pmt)
	err := rules.Validate(pmt)
	assert.ErrorContains(t, err, "Flow 6 allow-list")
}

func TestPaymentStatusCodeMismatchRejected(t *testing.T) {
	pmt := testPaymentReceipt(t)
	runNormalize(t, pmt)
	pmt.Ext = pmt.Ext.Set(ExtKeyStatus, "211") // wrong code for receipt
	err := rules.Validate(pmt)
	assert.ErrorContains(t, err, "ProcessConditionCode")
}

// Document the assumption that the payment-line currency is not
// inspected at the Flow 6 layer — it is taken from bill.Payment.Currency
// at the top level.
func TestPaymentTotalCurrencyEUR(t *testing.T) {
	pmt := testPaymentReceipt(t)
	assert.Equal(t, currency.Code("EUR"), pmt.Currency)
}
