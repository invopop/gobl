package bis

import (
	"testing"
	"time"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValAsParty(t *testing.T) {
	assert.Nil(t, valAsParty(nil))
	assert.Nil(t, valAsParty("string"))
	p := &org.Party{}
	assert.Equal(t, p, valAsParty(p))
}

func TestPartyHasLegalIdentity(t *testing.T) {
	assert.True(t, partyHasLegalIdentity(nil))
	assert.False(t, partyHasLegalIdentity(&org.Party{}))
	assert.True(t, partyHasLegalIdentity(&org.Party{
		Identities: []*org.Identity{{Scope: "legal", Code: "X"}},
	}))
	assert.True(t, partyHasLegalIdentity(&org.Party{
		TaxID: &tax.Identity{Code: "X"},
	}))
	// Identities without legal scope but with TaxID still passes through TaxID branch.
	assert.True(t, partyHasLegalIdentity(&org.Party{
		Identities: []*org.Identity{{Scope: "tax", Code: "X"}},
		TaxID:      &tax.Identity{Code: "Y"},
	}))
}

func TestFirstAddressStreetAndCode(t *testing.T) {
	assert.True(t, firstAddressStreetAndCode(nil))
	assert.True(t, firstAddressStreetAndCode([]*org.Address{}))
	assert.False(t, firstAddressStreetAndCode([]*org.Address{nil}))
	assert.False(t, firstAddressStreetAndCode([]*org.Address{{Street: "X"}}))
	assert.False(t, firstAddressStreetAndCode([]*org.Address{{Code: "1"}}))
	assert.True(t, firstAddressStreetAndCode([]*org.Address{{Street: "X", Code: "1"}}))
}

func TestValidISAccount(t *testing.T) {
	assert.True(t, validISAccount("123456789012"))                     // 12-digit domestic
	assert.True(t, validISAccount("IS140159260076545510730339"))       // IS IBAN
	assert.True(t, validISAccount("IS14 0159 2600 7654 5510 7303 39")) // IBAN with spaces
	assert.False(t, validISAccount(""))
	assert.False(t, validISAccount("12345"))
	assert.False(t, validISAccount("DE89370400440532013000")) // non-IS IBAN
}

func TestEINDAGIRules(t *testing.T) {
	eindagi := func(text string) *org.Note { return &org.Note{Src: NoteSrcEINDAGI, Text: text} }
	d := func(year int, month time.Month, day int) *cal.Date { return cal.NewDate(year, month, day) }

	t.Run("format: valid and invalid", func(t *testing.T) {
		assert.True(t, isEINDAGIFormatValid(nil))
		assert.True(t, isEINDAGIFormatValid(&bill.Invoice{}))
		assert.True(t, isEINDAGIFormatValid(&bill.Invoice{Notes: []*org.Note{eindagi("2026-06-30")}}))
		assert.False(t, isEINDAGIFormatValid(&bill.Invoice{Notes: []*org.Note{eindagi("30/06/2026")}}))
		// Other notes are ignored.
		assert.True(t, isEINDAGIFormatValid(&bill.Invoice{Notes: []*org.Note{{Src: "other", Text: "junk"}}}))
	})

	t.Run("due-date presence", func(t *testing.T) {
		// No EINDAGI — passes regardless.
		assert.True(t, isEINDAGIDueDatePresent(&bill.Invoice{}))
		// EINDAGI without due date — fails.
		assert.False(t, isEINDAGIDueDatePresent(&bill.Invoice{
			Notes: []*org.Note{eindagi("2026-06-30")},
		}))
		// EINDAGI with due date — passes.
		assert.True(t, isEINDAGIDueDatePresent(&bill.Invoice{
			Notes: []*org.Note{eindagi("2026-06-30")},
			Payment: &bill.PaymentDetails{Terms: &pay.Terms{
				DueDates: []*pay.DueDate{{Date: d(2026, time.June, 30)}},
			}},
		}))
	})

	t.Run("EINDAGI date ≥ first due date", func(t *testing.T) {
		// EINDAGI after due date — passes.
		assert.True(t, isEINDAGIAfterFirstDue(&bill.Invoice{
			Notes: []*org.Note{eindagi("2026-07-15")},
			Payment: &bill.PaymentDetails{Terms: &pay.Terms{
				DueDates: []*pay.DueDate{{Date: d(2026, time.June, 30)}},
			}},
		}))
		// EINDAGI before due date — fails.
		assert.False(t, isEINDAGIAfterFirstDue(&bill.Invoice{
			Notes: []*org.Note{eindagi("2026-06-01")},
			Payment: &bill.PaymentDetails{Terms: &pay.Terms{
				DueDates: []*pay.DueDate{{Date: d(2026, time.June, 30)}},
			}},
		}))
		// Equal dates — passes.
		assert.True(t, isEINDAGIAfterFirstDue(&bill.Invoice{
			Notes: []*org.Note{eindagi("2026-06-30")},
			Payment: &bill.PaymentDetails{Terms: &pay.Terms{
				DueDates: []*pay.DueDate{{Date: d(2026, time.June, 30)}},
			}},
		}))
		// Malformed EINDAGI ignored (format check covers it).
		assert.True(t, isEINDAGIAfterFirstDue(&bill.Invoice{
			Notes: []*org.Note{eindagi("bogus")},
			Payment: &bill.PaymentDetails{Terms: &pay.Terms{
				DueDates: []*pay.DueDate{{Date: d(2026, time.June, 30)}},
			}},
		}))
	})
}

func TestISPaymentCodes(t *testing.T) {
	// Code 9
	assert.True(t, isPaymentCode9Account(nil))
	assert.True(t, isPaymentCode9Account(&pay.Instructions{Ext: payExt("30")}))
	assert.False(t, isPaymentCode9Account(&pay.Instructions{Ext: payExt("9")})) // no transfers
	assert.True(t, isPaymentCode9Account(&pay.Instructions{
		Ext:            payExt("9"),
		CreditTransfer: []*pay.CreditTransfer{{Number: "123456789012"}},
	}))
	assert.True(t, isPaymentCode9Account(&pay.Instructions{
		Ext:            payExt("9"),
		CreditTransfer: []*pay.CreditTransfer{{IBAN: "IS140159260076545510730339"}},
	}))
	assert.False(t, isPaymentCode9Account(&pay.Instructions{
		Ext:            payExt("9"),
		CreditTransfer: []*pay.CreditTransfer{{Number: "12345"}},
	}))

	// Code 42
	assert.True(t, isPaymentCode42Account(nil))
	assert.True(t, isPaymentCode42Account(&pay.Instructions{Ext: payExt("30")}))
	assert.False(t, isPaymentCode42Account(&pay.Instructions{Ext: payExt("42")}))
	assert.True(t, isPaymentCode42Account(&pay.Instructions{
		Ext:            payExt("42"),
		CreditTransfer: []*pay.CreditTransfer{{Number: "123456789012"}},
	}))
	assert.False(t, isPaymentCode42Account(&pay.Instructions{
		Ext:            payExt("42"),
		CreditTransfer: []*pay.CreditTransfer{{Number: "AAA"}},
	}))
}
