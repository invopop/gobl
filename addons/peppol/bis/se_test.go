package bis

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func sePct(int64Val int64, exp uint32) num.Percentage {
	return num.MakePercentage(int64Val, exp)
}

func TestSEVATRateAllowed(t *testing.T) {
	assert.True(t, seVATRateAllowed(nil))
	assert.True(t, seVATRateAllowed(&bill.Invoice{}))
	// Build invoice with allowed 25% rate.
	twentyFive := sePct(250, 3)
	good := &bill.Invoice{Totals: &bill.Totals{Taxes: &tax.Total{Categories: []*tax.CategoryTotal{
		{Rates: []*tax.RateTotal{{Percent: &twentyFive}}},
	}}}}
	assert.True(t, seVATRateAllowed(good))
	// 21% — disallowed.
	twentyOne := sePct(210, 3)
	bad := &bill.Invoice{Totals: &bill.Totals{Taxes: &tax.Total{Categories: []*tax.CategoryTotal{
		{Rates: []*tax.RateTotal{{Percent: &twentyOne}}},
	}}}}
	assert.False(t, seVATRateAllowed(bad))
	// Nil rate or category — skipped.
	skip := &bill.Invoice{Totals: &bill.Totals{Taxes: &tax.Total{Categories: []*tax.CategoryTotal{
		nil, {Rates: []*tax.RateTotal{nil, {Percent: nil}}},
	}}}}
	assert.True(t, seVATRateAllowed(skip))
}

func TestSwedishVATLength(t *testing.T) {
	assert.True(t, seVATLength(nil))
	assert.True(t, seVATLength(&org.Party{}))
	// Non-SE — passes.
	assert.True(t, seVATLength(&org.Party{TaxID: &tax.Identity{Country: "DE", Code: "X"}}))
	// SE empty — passes.
	assert.True(t, seVATLength(&org.Party{TaxID: &tax.Identity{Country: l10n.SE.Tax()}}))
	// SE bare 10 — passes.
	assert.True(t, seVATLength(&org.Party{TaxID: &tax.Identity{Country: l10n.SE.Tax(), Code: "5560360793"}}))
	// SE full 14 — passes.
	assert.True(t, seVATLength(&org.Party{TaxID: &tax.Identity{Country: l10n.SE.Tax(), Code: "SE556036079301"}}))
	// SE with weird length — fails.
	assert.False(t, seVATLength(&org.Party{TaxID: &tax.Identity{Country: l10n.SE.Tax(), Code: "12345"}}))
}

func TestSwedishVATTrailingDigits(t *testing.T) {
	assert.True(t, seVATTrailingDigits(nil))
	assert.True(t, seVATTrailingDigits(&org.Party{}))
	// Non-SE passes.
	assert.True(t, seVATTrailingDigits(&org.Party{TaxID: &tax.Identity{Country: "DE"}}))
	// 14-char with digits passes.
	assert.True(t, seVATTrailingDigits(&org.Party{TaxID: &tax.Identity{Country: l10n.SE.Tax(), Code: "SE556036079301"}}))
	// 14-char with letters in trailing — fails.
	assert.False(t, seVATTrailingDigits(&org.Party{TaxID: &tax.Identity{Country: l10n.SE.Tax(), Code: "SE5560360793AB"}}))
	// 10-digit numeric — passes.
	assert.True(t, seVATTrailingDigits(&org.Party{TaxID: &tax.Identity{Country: l10n.SE.Tax(), Code: "5560360793"}}))
	// 10-char with letters — fails.
	assert.False(t, seVATTrailingDigits(&org.Party{TaxID: &tax.Identity{Country: l10n.SE.Tax(), Code: "ABCDEFGHIJ"}}))
	// Other length — passes (handled elsewhere).
	assert.True(t, seVATTrailingDigits(&org.Party{TaxID: &tax.Identity{Country: l10n.SE.Tax(), Code: "12345"}}))
}

func TestSwedishOrgChecks(t *testing.T) {
	// Length
	assert.True(t, seOrgLength(nil))
	assert.True(t, seOrgLength(&org.Party{Identities: []*org.Identity{{Scope: "legal", Code: "5560360793"}}}))
	assert.False(t, seOrgLength(&org.Party{Identities: []*org.Identity{{Scope: "legal", Code: "12345"}}}))

	// Luhn
	assert.True(t, seOrgLuhn(nil))
	assert.True(t, seOrgLuhn(&org.Party{}))
	// Wrong length skipped (handled by R-004).
	assert.True(t, seOrgLuhn(&org.Party{Identities: []*org.Identity{{Scope: "legal", Code: "12345"}}}))
	assert.True(t, seOrgLuhn(&org.Party{Identities: []*org.Identity{{Scope: "legal", Code: "5560360793"}}}))
	assert.False(t, seOrgLuhn(&org.Party{Identities: []*org.Identity{{Scope: "legal", Code: "5560360794"}}}))
	// Nil identity / wrong scope skipped.
	assert.True(t, seOrgLuhn(&org.Party{Identities: []*org.Identity{nil, {Scope: "tax", Code: "X"}}}))
}
