package bis

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestPartyHasContactGroup(t *testing.T) {
	assert.True(t, partyHasContactGroup(nil))
	assert.False(t, partyHasContactGroup(&org.Party{}))
	assert.True(t, partyHasContactGroup(&org.Party{People: []*org.Person{{Name: &org.Name{Given: "X"}}}}))
	assert.True(t, partyHasContactGroup(&org.Party{Telephones: []*org.Telephone{{Number: "1"}}}))
	assert.True(t, partyHasContactGroup(&org.Party{Emails: []*org.Email{{Address: "a@b"}}}))
}

func TestFirstAddressHelpers(t *testing.T) {
	assert.True(t, firstAddressHasLocalityPE(nil))
	assert.True(t, firstAddressHasLocalityPE([]*org.Address{}))
	assert.False(t, firstAddressHasLocalityPE([]*org.Address{nil}))
	assert.False(t, firstAddressHasLocalityPE([]*org.Address{{}}))
	assert.True(t, firstAddressHasLocalityPE([]*org.Address{{Locality: "Berlin"}}))

	assert.True(t, firstAddressHasCodePE([]*org.Address{}))
	assert.False(t, firstAddressHasCodePE([]*org.Address{nil}))
	assert.False(t, firstAddressHasCodePE([]*org.Address{{}}))
	assert.True(t, firstAddressHasCodePE([]*org.Address{{Code: "10115"}}))
}

func TestPartyContactHelpers(t *testing.T) {
	assert.True(t, partyHasContactName(nil))
	assert.False(t, partyHasContactName(&org.Party{}))
	assert.False(t, partyHasContactName(&org.Party{People: []*org.Person{nil}}))
	assert.False(t, partyHasContactName(&org.Party{People: []*org.Person{{}}}))
	assert.True(t, partyHasContactName(&org.Party{People: []*org.Person{{Name: &org.Name{Given: "Anna"}}}}))

	assert.True(t, partyHasContactTelephone(nil))
	assert.False(t, partyHasContactTelephone(&org.Party{}))
	assert.True(t, partyHasContactTelephone(&org.Party{Telephones: []*org.Telephone{{Number: "+49"}}}))
	assert.True(t, partyHasContactTelephone(&org.Party{People: []*org.Person{{Telephones: []*org.Telephone{{Number: "+49"}}}}}))

	assert.True(t, partyHasContactEmail(nil))
	assert.False(t, partyHasContactEmail(&org.Party{}))
	assert.True(t, partyHasContactEmail(&org.Party{Emails: []*org.Email{{Address: "a@b"}}}))
	assert.True(t, partyHasContactEmail(&org.Party{People: []*org.Person{{Emails: []*org.Email{{Address: "a@b"}}}}}))
}

func TestDeSupplierHasTaxIDForCategory(t *testing.T) {
	assert.True(t, deSupplierHasTaxIDForCategory(nil))
	// No qualifying categories — passes regardless of supplier.
	assert.True(t, deSupplierHasTaxIDForCategory(&bill.Invoice{Totals: &bill.Totals{Taxes: &tax.Total{}}}))
	// S category present, supplier has tax id -> passes.
	tot := &bill.Totals{Taxes: &tax.Total{Categories: []*tax.CategoryTotal{
		{Rates: []*tax.RateTotal{{Ext: tax.Extensions{untdid.ExtKeyTaxCategory: "S"}}}},
	}}}
	assert.True(t, deSupplierHasTaxIDForCategory(&bill.Invoice{
		Totals:   tot,
		Supplier: &org.Party{TaxID: &tax.Identity{Code: "DE111"}},
	}))
	// No supplier -> fails.
	assert.False(t, deSupplierHasTaxIDForCategory(&bill.Invoice{Totals: tot}))
	// Supplier with legal-scope identity (tax registration) -> passes.
	assert.True(t, deSupplierHasTaxIDForCategory(&bill.Invoice{
		Totals: tot,
		Supplier: &org.Party{Identities: []*org.Identity{{
			Scope: org.IdentityScopeLegal, Code: "HRB-1234",
		}}},
	}))
	// Non-legal identity (e.g. DUNS, GLN) -> fails.
	assert.False(t, deSupplierHasTaxIDForCategory(&bill.Invoice{
		Totals:   tot,
		Supplier: &org.Party{Identities: []*org.Identity{{Code: "123456789"}}},
	}))
}

func TestSkontoFormatValid(t *testing.T) {
	assert.True(t, skontoFormatValid(nil))
	assert.True(t, skontoFormatValid(&bill.Invoice{}))
	assert.True(t, skontoFormatValid(&bill.Invoice{Payment: &bill.PaymentDetails{}}))
	// Intro paragraph + well-formed Skonto line passes.
	good := &bill.Invoice{Payment: &bill.PaymentDetails{Terms: &pay.Terms{
		Notes: "Payment within 30 days net.\n#SKONTO#TAGE=14#PROZENT=3.00#",
	}}}
	assert.True(t, skontoFormatValid(good))
	// Malformed Skonto line fails.
	bad := &bill.Invoice{Payment: &bill.PaymentDetails{Terms: &pay.Terms{
		Notes: "#SKONTO#TAGE=14#PROZENT#",
	}}}
	assert.False(t, skontoFormatValid(bad))
	// BASISBETRAG variant passes.
	withBasis := &bill.Invoice{Payment: &bill.PaymentDetails{Terms: &pay.Terms{
		Notes: "#SKONTO#TAGE=7#PROZENT=5.00#BASISBETRAG=1000.00#",
	}}}
	assert.True(t, skontoFormatValid(withBasis))
	// Per-due-date notes are also checked.
	duePerDate := &bill.Invoice{Payment: &bill.PaymentDetails{Terms: &pay.Terms{
		DueDates: []*pay.DueDate{{Notes: "#SKONTO#oops"}},
	}}}
	assert.False(t, skontoFormatValid(duePerDate))
	// Non-# lines pass regardless.
	plain := &bill.Invoice{Payment: &bill.PaymentDetails{Terms: &pay.Terms{
		Notes: "Anything goes here.\nNo Skonto.",
	}}}
	assert.True(t, skontoFormatValid(plain))
}

func TestDePaymentExclusivity(t *testing.T) {
	// deCreditTransferExclusive
	assert.True(t, deCreditTransferExclusive(nil))
	assert.True(t, deCreditTransferExclusive(&pay.Instructions{Ext: payExt("99")}))
	assert.False(t, deCreditTransferExclusive(&pay.Instructions{Ext: payExt("30")})) // no transfer
	assert.True(t, deCreditTransferExclusive(&pay.Instructions{
		Ext:            payExt("30"),
		CreditTransfer: []*pay.CreditTransfer{{IBAN: "DE89"}},
	}))
	assert.False(t, deCreditTransferExclusive(&pay.Instructions{
		Ext:            payExt("30"),
		CreditTransfer: []*pay.CreditTransfer{{IBAN: "DE89"}},
		Card:           &pay.Card{},
	}))

	// deCardExclusive
	assert.True(t, deCardExclusive(nil))
	assert.True(t, deCardExclusive(&pay.Instructions{Ext: payExt("30")}))
	assert.False(t, deCardExclusive(&pay.Instructions{Ext: payExt("48")})) // missing card
	assert.True(t, deCardExclusive(&pay.Instructions{Ext: payExt("48"), Card: &pay.Card{}}))
	assert.False(t, deCardExclusive(&pay.Instructions{
		Ext:            payExt("48"),
		Card:           &pay.Card{},
		CreditTransfer: []*pay.CreditTransfer{{}},
	}))

	// deDirectDebitExclusive
	assert.True(t, deDirectDebitExclusive(nil))
	assert.True(t, deDirectDebitExclusive(&pay.Instructions{Ext: payExt("30")}))
	assert.False(t, deDirectDebitExclusive(&pay.Instructions{Ext: payExt("59")}))
	assert.True(t, deDirectDebitExclusive(&pay.Instructions{Ext: payExt("59"), DirectDebit: &pay.DirectDebit{}}))
	assert.False(t, deDirectDebitExclusive(&pay.Instructions{
		Ext:            payExt("59"),
		DirectDebit:    &pay.DirectDebit{},
		CreditTransfer: []*pay.CreditTransfer{{}},
	}))

	// deDirectDebitFieldsComplete
	assert.True(t, deDirectDebitFieldsComplete(nil))
	assert.True(t, deDirectDebitFieldsComplete(&pay.Instructions{}))
	assert.False(t, deDirectDebitFieldsComplete(&pay.Instructions{DirectDebit: &pay.DirectDebit{}}))
	assert.True(t, deDirectDebitFieldsComplete(&pay.Instructions{
		DirectDebit: &pay.DirectDebit{Creditor: "C", Account: "A"},
	}))
}
