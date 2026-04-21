package bis

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestDEInvoiceDocumentTypeValid(t *testing.T) {
	assert.True(t, deInvoiceDocumentTypeValid(nil))
	assert.True(t, deInvoiceDocumentTypeValid(tax.Extensions{}))
	assert.True(t, deInvoiceDocumentTypeValid(tax.Extensions{untdid.ExtKeyDocumentType: "380"}))
	assert.False(t, deInvoiceDocumentTypeValid(tax.Extensions{untdid.ExtKeyDocumentType: "999"}))
}

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

func TestPartyTelephoneMinLength(t *testing.T) {
	assert.True(t, partyTelephoneMinLength(nil))
	// Three digits passes.
	assert.True(t, partyTelephoneMinLength(&org.Party{Telephones: []*org.Telephone{{Number: "123"}}}))
	// Two digits fails.
	assert.False(t, partyTelephoneMinLength(&org.Party{Telephones: []*org.Telephone{{Number: "+1"}}}))
	// Looks at people fallback.
	assert.True(t, partyTelephoneMinLength(&org.Party{People: []*org.Person{{Telephones: []*org.Telephone{{Number: "+49 123 4567"}}}}}))
	// Nil entry skipped.
	assert.True(t, partyTelephoneMinLength(&org.Party{Telephones: []*org.Telephone{nil, {Number: "12345"}}}))
}

func TestPartyEmailFormat(t *testing.T) {
	assert.True(t, partyEmailFormat(nil))
	assert.True(t, partyEmailFormat(&org.Party{Emails: []*org.Email{{Address: "a@b"}}}))
	assert.False(t, partyEmailFormat(&org.Party{Emails: []*org.Email{{Address: "no-at-sign"}}}))
	assert.False(t, partyEmailFormat(&org.Party{Emails: []*org.Email{{Address: "a@b@c"}}}))
	assert.False(t, partyEmailFormat(&org.Party{Emails: []*org.Email{{Address: "@b"}}}))
	assert.False(t, partyEmailFormat(&org.Party{Emails: []*org.Email{{Address: "a@"}}}))
	// Nil entry skipped.
	assert.True(t, partyEmailFormat(&org.Party{Emails: []*org.Email{nil}}))
}

func TestCorrectivePrecedingPresent(t *testing.T) {
	assert.True(t, correctivePrecedingPresent(nil))
	assert.True(t, correctivePrecedingPresent(&bill.Invoice{}))
	assert.True(t, correctivePrecedingPresent(&bill.Invoice{Tax: taxExt("380")}))
	assert.False(t, correctivePrecedingPresent(&bill.Invoice{Tax: taxExt("384")}))
	assert.True(t, correctivePrecedingPresent(&bill.Invoice{
		Tax:       taxExt("384"),
		Preceding: []*org.DocumentRef{{Code: "ORIG"}},
	}))
}

func TestDeVATRatePercentSet(t *testing.T) {
	assert.True(t, deVATRatePercentSet(nil))
	assert.True(t, deVATRatePercentSet(&bill.Invoice{}))
	pct := num.MakePercentage(190, 3)
	// Standard-rated with percent → passes.
	good := &bill.Invoice{Totals: &bill.Totals{Taxes: &tax.Total{Categories: []*tax.CategoryTotal{
		{Rates: []*tax.RateTotal{{Percent: &pct, Ext: tax.Extensions{untdid.ExtKeyTaxCategory: "S"}}}},
	}}}}
	assert.True(t, deVATRatePercentSet(good))
	// Standard-rated without percent → fails.
	bad := &bill.Invoice{Totals: &bill.Totals{Taxes: &tax.Total{Categories: []*tax.CategoryTotal{
		{Rates: []*tax.RateTotal{{Percent: nil, Ext: tax.Extensions{untdid.ExtKeyTaxCategory: "S"}}}},
	}}}}
	assert.False(t, deVATRatePercentSet(bad))
	// Exempt (E) rate without percent → passes (rule scoped to standard-rated).
	exempt := &bill.Invoice{Totals: &bill.Totals{Taxes: &tax.Total{Categories: []*tax.CategoryTotal{
		{Rates: []*tax.RateTotal{{Percent: nil, Ext: tax.Extensions{untdid.ExtKeyTaxCategory: "E"}}}},
	}}}}
	assert.True(t, deVATRatePercentSet(exempt))
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

func TestDeSEPAIBANValid(t *testing.T) {
	assert.True(t, deSEPAIBANValid(nil))
	assert.True(t, deSEPAIBANValid(&pay.Instructions{Ext: payExt("30")}))
	// Code 58 with valid IBAN passes.
	assert.True(t, deSEPAIBANValid(&pay.Instructions{
		Ext:            payExt("58"),
		CreditTransfer: []*pay.CreditTransfer{{IBAN: "DE89370400440532013000"}},
	}))
	// Code 58 with junk fails.
	assert.False(t, deSEPAIBANValid(&pay.Instructions{
		Ext:            payExt("58"),
		CreditTransfer: []*pay.CreditTransfer{{IBAN: "not-an-iban!"}},
	}))
	// Empty account passes (handled elsewhere).
	assert.True(t, deSEPAIBANValid(&pay.Instructions{
		Ext:            payExt("58"),
		CreditTransfer: []*pay.CreditTransfer{{}},
	}))
}

func TestDeSEPADebitIBANValid(t *testing.T) {
	assert.True(t, deSEPADebitIBANValid(nil))
	assert.True(t, deSEPADebitIBANValid(&pay.Instructions{Ext: payExt("30")}))
	assert.True(t, deSEPADebitIBANValid(&pay.Instructions{Ext: payExt("59")}))
	assert.True(t, deSEPADebitIBANValid(&pay.Instructions{
		Ext:         payExt("59"),
		DirectDebit: &pay.DirectDebit{Account: "DE89370400440532013000"},
	}))
	assert.False(t, deSEPADebitIBANValid(&pay.Instructions{
		Ext:         payExt("59"),
		DirectDebit: &pay.DirectDebit{Account: "junk!"},
	}))
}
