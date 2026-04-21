package bis

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNLCreditNoteHasPreceding(t *testing.T) {
	assert.True(t, nlCreditNoteHasPreceding(nil))
	assert.True(t, nlCreditNoteHasPreceding(&bill.Invoice{Type: bill.InvoiceTypeStandard}))
	assert.False(t, nlCreditNoteHasPreceding(&bill.Invoice{Type: bill.InvoiceTypeCreditNote}))
	assert.True(t, nlCreditNoteHasPreceding(&bill.Invoice{
		Type:      bill.InvoiceTypeCreditNote,
		Preceding: []*org.DocumentRef{{Code: "ORIG"}},
	}))
	// Nil preceding entry skipped.
	assert.False(t, nlCreditNoteHasPreceding(&bill.Invoice{
		Type:      bill.InvoiceTypeCreditNote,
		Preceding: []*org.DocumentRef{nil, {}},
	}))
}

func TestNLLineOrderRefRequiresOrderingCode(t *testing.T) {
	assert.True(t, nlLineOrderRefRequiresOrderingCode(nil))
	// No line order refs — passes.
	assert.True(t, nlLineOrderRefRequiresOrderingCode(&bill.Invoice{Lines: []*bill.Line{{}}}))
	// Line with order ref but no Ordering.Code — fails.
	assert.False(t, nlLineOrderRefRequiresOrderingCode(&bill.Invoice{
		Lines: []*bill.Line{{Order: "X"}},
	}))
	// Line with order ref AND Ordering.Code — passes.
	assert.True(t, nlLineOrderRefRequiresOrderingCode(&bill.Invoice{
		Lines:    []*bill.Line{{Order: "X"}},
		Ordering: &bill.Ordering{Code: "BR"},
	}))
}

func TestFirstAddressStreetLocalityCode(t *testing.T) {
	assert.True(t, firstAddressStreetLocalityCode(nil))
	assert.True(t, firstAddressStreetLocalityCode([]*org.Address{}))
	// Note: this helper requires non-nil first address, where nil first returns false because (nil != nil) is false.
	assert.False(t, firstAddressStreetLocalityCode([]*org.Address{nil}))
	assert.False(t, firstAddressStreetLocalityCode([]*org.Address{{Street: "X"}}))
	assert.True(t, firstAddressStreetLocalityCode([]*org.Address{{
		Street: "X", Locality: "Amsterdam", Code: "1011",
	}}))
}

func TestNLPartyLegalScheme(t *testing.T) {
	assert.True(t, nlPartyLegalScheme(nil))
	// Non-NL party — passes.
	assert.True(t, nlPartyLegalScheme(&org.Party{TaxID: &tax.Identity{Country: "DE"}}))
	// NL party with KVK scheme — passes.
	assert.True(t, nlPartyLegalScheme(&org.Party{
		TaxID: &tax.Identity{Country: l10n.NL.Tax()},
		Identities: []*org.Identity{{
			Scope: "legal",
			Ext:   tax.Extensions{iso.ExtKeySchemeID: "0106"},
		}},
	}))
	// NL party with OIN scheme — passes.
	assert.True(t, nlPartyLegalScheme(&org.Party{
		TaxID: &tax.Identity{Country: l10n.NL.Tax()},
		Identities: []*org.Identity{{
			Scope: "legal",
			Ext:   tax.Extensions{iso.ExtKeySchemeID: "0190"},
		}},
	}))
	// NL party with wrong scheme — fails.
	assert.False(t, nlPartyLegalScheme(&org.Party{
		TaxID: &tax.Identity{Country: l10n.NL.Tax()},
		Identities: []*org.Identity{{
			Scope: "legal",
			Ext:   tax.Extensions{iso.ExtKeySchemeID: "0184"},
		}},
	}))
	// NL party with no legal-scope identity — fails.
	assert.False(t, nlPartyLegalScheme(&org.Party{
		TaxID:      &tax.Identity{Country: l10n.NL.Tax()},
		Identities: []*org.Identity{{Scope: "tax", Ext: tax.Extensions{iso.ExtKeySchemeID: "0106"}}},
	}))
	// Nil identity entry skipped.
	assert.False(t, nlPartyLegalScheme(&org.Party{
		TaxID:      &tax.Identity{Country: l10n.NL.Tax()},
		Identities: []*org.Identity{nil},
	}))
}

func TestNLPaymentMeansAllowed(t *testing.T) {
	assert.True(t, nlPaymentMeansAllowed(nil))
	assert.True(t, nlPaymentMeansAllowed(&pay.Instructions{}))
	assert.True(t, nlPaymentMeansAllowed(&pay.Instructions{Ext: payExt(cbc.Code("30"))}))
	assert.True(t, nlPaymentMeansAllowed(&pay.Instructions{Ext: payExt(cbc.Code("48"))}))
	assert.False(t, nlPaymentMeansAllowed(&pay.Instructions{Ext: payExt(cbc.Code("99"))}))
}
