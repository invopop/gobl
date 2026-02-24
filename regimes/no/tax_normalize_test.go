package no

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
)

func TestNormalizeTaxIdentityStripsPrefixAndSuffix(t *testing.T) {
	id := &tax.Identity{
		Country: l10n.TaxCountryCode(l10n.NO),
		Code:    " NO 974 760 673 mva ",
	}

	normalizeTaxIdentity(id)

	if id.Code != "974760673" {
		t.Fatalf("expected normalized code to be 974760673, got %s", id.Code)
	}
}

func TestNormalizeTaxIdentityLeavesUnknownFormats(t *testing.T) {
	id := &tax.Identity{
		Country: l10n.TaxCountryCode(l10n.NO),
		Code:    "SOMETHING",
	}
	normalizeTaxIdentity(id)
	if id.Code != "SOMETHING" {
		t.Fatalf("expected code to remain unchanged, got %s", id.Code)
	}
}
