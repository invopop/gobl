package bis

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestDKCreditNoteNonNegative(t *testing.T) {
	assert.True(t, dkCreditNoteNonNegative(nil))
	assert.True(t, dkCreditNoteNonNegative(&bill.Invoice{Type: bill.InvoiceTypeStandard}))
	assert.True(t, dkCreditNoteNonNegative(&bill.Invoice{Type: bill.InvoiceTypeCreditNote}))
	pos := num.MakeAmount(100, 2)
	neg := num.MakeAmount(-100, 2)
	assert.True(t, dkCreditNoteNonNegative(&bill.Invoice{
		Type:   bill.InvoiceTypeCreditNote,
		Totals: &bill.Totals{Payable: pos},
	}))
	assert.False(t, dkCreditNoteNonNegative(&bill.Invoice{
		Type:   bill.InvoiceTypeCreditNote,
		Totals: &bill.Totals{Payable: neg},
	}))
}

func TestPartyHasCVRIdentity(t *testing.T) {
	assert.True(t, partyHasCVRIdentity(nil))
	assert.False(t, partyHasCVRIdentity(&org.Party{}))
	assert.True(t, partyHasCVRIdentity(&org.Party{
		Identities: []*org.Identity{{Ext: tax.Extensions{iso.ExtKeySchemeID: "0184"}}},
	}))
	assert.True(t, partyHasCVRIdentity(&org.Party{
		TaxID: &tax.Identity{Country: l10n.DK.Tax(), Code: "13585628"},
	}))
	assert.False(t, partyHasCVRIdentity(&org.Party{
		Identities: []*org.Identity{nil, {Ext: tax.Extensions{iso.ExtKeySchemeID: "0007"}}},
	}))
}

func TestIdentityHasSchemeID(t *testing.T) {
	assert.True(t, identityHasSchemeID(nil))
	assert.False(t, identityHasSchemeID(&org.Identity{}))
	assert.True(t, identityHasSchemeID(&org.Identity{Ext: tax.Extensions{iso.ExtKeySchemeID: "0184"}}))
}

func TestCustomerIsDK(t *testing.T) {
	assert.False(t, customerIsDK(nil))
	assert.False(t, customerIsDK(&org.Party{TaxID: &tax.Identity{Country: "FR"}}))
	assert.True(t, customerIsDK(&org.Party{TaxID: &tax.Identity{Country: "DK"}}))
}

func TestDKItemClassificationsValid(t *testing.T) {
	assert.True(t, dkItemClassificationsValid(nil))
	// Item without classification description passes.
	inv := &bill.Invoice{
		Lines: []*bill.Line{
			{Item: &org.Item{Identities: []*org.Identity{{Code: "ABC"}}}},
		},
	}
	assert.True(t, dkItemClassificationsValid(inv))
	// Allowed UNSPSC version.
	inv2 := &bill.Invoice{
		Lines: []*bill.Line{
			{Item: &org.Item{Identities: []*org.Identity{{Code: "X", Description: "19.05.01"}}}},
		},
	}
	assert.True(t, dkItemClassificationsValid(inv2))
	// Disallowed version.
	inv3 := &bill.Invoice{
		Lines: []*bill.Line{
			{Item: &org.Item{Identities: []*org.Identity{{Code: "X", Description: "1.0"}}}},
		},
	}
	assert.False(t, dkItemClassificationsValid(inv3))
	// Nil line / item / identity skipped.
	inv4 := &bill.Invoice{Lines: []*bill.Line{nil, {Item: nil}, {Item: &org.Item{Identities: []*org.Identity{nil}}}}}
	assert.True(t, dkItemClassificationsValid(inv4))
}

func TestDKPaymentMeansAllowed(t *testing.T) {
	assert.True(t, dkPaymentMeansAllowed(nil))
	assert.True(t, dkPaymentMeansAllowed(&pay.Instructions{}))
	assert.True(t, dkPaymentMeansAllowed(&pay.Instructions{Ext: payExt(cbc.Code("31"))}))
	assert.False(t, dkPaymentMeansAllowed(&pay.Instructions{Ext: payExt(cbc.Code("99"))}))
}

func TestDKCreditTransferComplete(t *testing.T) {
	assert.True(t, dkCreditTransferComplete(nil))
	// Other payment means -> passes regardless.
	assert.True(t, dkCreditTransferComplete(&pay.Instructions{Ext: payExt(cbc.Code("30"))}))
	// 31 with no transfer entries -> fails.
	assert.False(t, dkCreditTransferComplete(&pay.Instructions{Ext: payExt(cbc.Code("31"))}))
	// 31 with a transfer missing Number -> fails.
	assert.False(t, dkCreditTransferComplete(&pay.Instructions{
		Ext:            payExt(cbc.Code("31")),
		CreditTransfer: []*pay.CreditTransfer{{}},
	}))
	// 31 with a transfer that has Number -> passes.
	assert.True(t, dkCreditTransferComplete(&pay.Instructions{
		Ext:            payExt(cbc.Code("31")),
		CreditTransfer: []*pay.CreditTransfer{{Number: "1234567890"}},
	}))
}

func TestDKDirectDebit49Complete(t *testing.T) {
	assert.True(t, dkDirectDebit49Complete(nil))
	assert.True(t, dkDirectDebit49Complete(&pay.Instructions{Ext: payExt(cbc.Code("30"))}))
	// 49 without DirectDebit -> fails.
	assert.False(t, dkDirectDebit49Complete(&pay.Instructions{Ext: payExt(cbc.Code("49"))}))
	// 49 with empty DirectDebit -> fails.
	assert.False(t, dkDirectDebit49Complete(&pay.Instructions{
		Ext:         payExt(cbc.Code("49")),
		DirectDebit: &pay.DirectDebit{},
	}))
	// 49 with both Ref and Account -> passes.
	assert.True(t, dkDirectDebit49Complete(&pay.Instructions{
		Ext:         payExt(cbc.Code("49")),
		DirectDebit: &pay.DirectDebit{Ref: "M", Account: "A"},
	}))
}
