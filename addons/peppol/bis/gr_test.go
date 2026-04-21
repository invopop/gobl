package bis

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestGRFullInvoiceID(t *testing.T) {
	assert.Equal(t, "", grFullInvoiceID(nil))
	assert.Equal(t, "001", grFullInvoiceID(&bill.Invoice{Code: "001"}))
	assert.Equal(t, "S_001", grFullInvoiceID(&bill.Invoice{Series: "S", Code: "001"}))
}

func TestGRIDSixSegments(t *testing.T) {
	assert.True(t, grIDSixSegments(nil))
	assert.True(t, grIDSixSegments(&bill.Invoice{}))
	// 6 segments via Series + Code joined.
	good := &bill.Invoice{Series: "100000310_20240413_3_1.1_REPORT", Code: "AAB123"}
	assert.True(t, grIDSixSegments(good))
	// 5 segments fails.
	bad := &bill.Invoice{Code: "a_b_c_d_e"}
	assert.False(t, grIDSixSegments(bad))
}

func TestGRMARK(t *testing.T) {
	assert.True(t, grMARKExactlyOne(nil))
	assert.False(t, grMARKExactlyOne(&bill.Invoice{}))
	assert.False(t, grMARKExactlyOne(&bill.Invoice{Ordering: &bill.Ordering{}}))
	one := &bill.Invoice{Ordering: &bill.Ordering{Identities: []*org.Identity{{Key: IdentityKeyGreekMARK, Code: "1"}}}}
	assert.True(t, grMARKExactlyOne(one))
	two := &bill.Invoice{Ordering: &bill.Ordering{Identities: []*org.Identity{
		{Key: IdentityKeyGreekMARK, Code: "1"},
		{Key: IdentityKeyGreekMARK, Code: "2"},
	}}}
	assert.False(t, grMARKExactlyOne(two))

	// Positive integer check.
	assert.True(t, grMARKPositive(nil))
	assert.True(t, grMARKPositive(&bill.Invoice{}))
	assert.True(t, grMARKPositive(one))
	bad := &bill.Invoice{Ordering: &bill.Ordering{Identities: []*org.Identity{{Key: IdentityKeyGreekMARK, Code: "abc"}}}}
	assert.False(t, grMARKPositive(bad))
	zero := &bill.Invoice{Ordering: &bill.Ordering{Identities: []*org.Identity{{Key: IdentityKeyGreekMARK, Code: "0"}}}}
	assert.False(t, grMARKPositive(zero))
}

func TestGRInvoiceURLCardinality(t *testing.T) {
	assert.True(t, grInvoiceURLCardinality(nil))
	assert.True(t, grInvoiceURLCardinality(&bill.Invoice{}))
	assert.True(t, grInvoiceURLCardinality(&bill.Invoice{
		Preceding: []*org.DocumentRef{{URL: "https://x"}},
	}))
	assert.False(t, grInvoiceURLCardinality(&bill.Invoice{
		Preceding: []*org.DocumentRef{{URL: "https://x"}, {URL: "https://y"}},
	}))
}

func TestGRSupplierVATValid(t *testing.T) {
	assert.True(t, grSupplierVATValid(nil))
	assert.True(t, grSupplierVATValid(&org.Party{}))
	// Non-Greek tax id — passes.
	assert.True(t, grSupplierVATValid(&org.Party{TaxID: &tax.Identity{Country: "DE", Code: "123"}}))
	// Greek 9-digit code passes.
	assert.True(t, grSupplierVATValid(&org.Party{TaxID: &tax.Identity{Country: l10n.GR.Tax(), Code: "123456789"}}))
	// Greek non-9-digit fails.
	assert.False(t, grSupplierVATValid(&org.Party{TaxID: &tax.Identity{Country: l10n.GR.Tax(), Code: "ABC"}}))
}

func TestGRSupplierInboxValid(t *testing.T) {
	assert.True(t, grSupplierInboxValid(nil))
	assert.True(t, grSupplierInboxValid(&org.Party{})) // no inboxes — handled elsewhere
	good := &org.Party{
		TaxID:   &tax.Identity{Country: l10n.GR.Tax(), Code: "123456789"},
		Inboxes: []*org.Inbox{{Scheme: "9933", Code: "123456789"}},
	}
	assert.True(t, grSupplierInboxValid(good))
	wrongScheme := &org.Party{
		TaxID:   &tax.Identity{Country: l10n.GR.Tax(), Code: "123456789"},
		Inboxes: []*org.Inbox{{Scheme: "0184", Code: "123456789"}},
	}
	assert.False(t, grSupplierInboxValid(wrongScheme))
	mismatch := &org.Party{
		TaxID:   &tax.Identity{Country: l10n.GR.Tax(), Code: "123456789"},
		Inboxes: []*org.Inbox{{Scheme: "9933", Code: "987654321"}},
	}
	assert.False(t, grSupplierInboxValid(mismatch))
	badFormat := &org.Party{
		Inboxes: []*org.Inbox{{Scheme: "9933", Code: "abc"}},
	}
	assert.False(t, grSupplierInboxValid(badFormat))
	nilInbox := &org.Party{Inboxes: []*org.Inbox{nil}}
	assert.True(t, grSupplierInboxValid(nilInbox))
}

func TestGRCustomerVATWhenGreek(t *testing.T) {
	assert.True(t, grCustomerVATWhenGreek(nil))
	// Non-Greek customer — passes.
	assert.True(t, grCustomerVATWhenGreek(&org.Party{TaxID: &tax.Identity{Country: "DE"}}))
	// Greek customer without VAT — fails.
	assert.False(t, grCustomerVATWhenGreek(&org.Party{TaxID: &tax.Identity{Country: l10n.GR.Tax()}}))
	// Greek customer with VAT — passes.
	assert.True(t, grCustomerVATWhenGreek(&org.Party{TaxID: &tax.Identity{Country: l10n.GR.Tax(), Code: "123"}}))
}

func TestGRCustomerInboxWhenGreek(t *testing.T) {
	assert.True(t, grCustomerInboxWhenGreek(nil))
	// Non-Greek customer — passes.
	assert.True(t, grCustomerInboxWhenGreek(&org.Party{TaxID: &tax.Identity{Country: "DE"}}))
	// Greek customer with valid inbox — passes.
	assert.True(t, grCustomerInboxWhenGreek(&org.Party{
		TaxID:   &tax.Identity{Country: l10n.GR.Tax(), Code: "123456789"},
		Inboxes: []*org.Inbox{{Scheme: "9933", Code: "123456789"}},
	}))
}
