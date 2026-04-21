package bis

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNorwegianVATFormat(t *testing.T) {
	assert.True(t, norwegianVATFormat(nil))
	assert.True(t, norwegianVATFormat(&org.Party{}))
	// Non-NO tax id passes.
	assert.True(t, norwegianVATFormat(&org.Party{TaxID: &tax.Identity{Country: "DE"}}))
	// Empty code passes.
	assert.True(t, norwegianVATFormat(&org.Party{TaxID: &tax.Identity{Country: l10n.NO.Tax()}}))
	// Bare 9 digits passes.
	assert.True(t, norwegianVATFormat(&org.Party{TaxID: &tax.Identity{Country: l10n.NO.Tax(), Code: "990983666"}}))
	// Full NOxxxMVA form passes.
	assert.True(t, norwegianVATFormat(&org.Party{TaxID: &tax.Identity{Country: l10n.NO.Tax(), Code: "NO990983666MVA"}}))
	// Non-conforming code fails.
	assert.False(t, norwegianVATFormat(&org.Party{TaxID: &tax.Identity{Country: l10n.NO.Tax(), Code: "ABC"}}))
}
