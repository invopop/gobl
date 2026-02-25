package no

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func TestNormalizeSwitchTaxIdentity(t *testing.T) {
	id := &tax.Identity{
		Country: l10n.TaxCountryCode(l10n.NO),
		Code:    cbc.Code(" NO 974 760 673 mva "),
	}
	Normalize(id)

	if id.Code != "974760673" {
		t.Fatalf("expected normalized code to be 974760673, got %s", id.Code)
	}
}

func TestNormalizeSwitchOrgIdentity(t *testing.T) {
	id := &org.Identity{
		Type: IdentityTypeOrgNr,
		Code: cbc.Code("974 760 673"),
	}
	Normalize(id)

	if id.Code != "974760673" {
		t.Fatalf("expected normalized code to be 974760673, got %s", id.Code)
	}
}

func TestNormalizeSwitchUnknownTypeNoop(_ *testing.T) {
	Normalize(struct{}{})
}

func TestValidateSwitchInvoice(t *testing.T) {
	inv := &bill.Invoice{
		Supplier: &org.Party{
			Name: "Supplier",
			TaxID: &tax.Identity{
				Country: l10n.TaxCountryCode(l10n.NO),
				Code:    "974760673MVA",
			},
		},
	}
	inv.Tags.SetTags(tax.TagSimplified)

	if err := Validate(inv); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateSwitchTaxIdentity(t *testing.T) {
	id := &tax.Identity{
		Country: l10n.TaxCountryCode(l10n.NO),
		Code:    cbc.Code("974760673MVA"),
	}
	if err := Validate(id); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateSwitchOrgIdentity(t *testing.T) {
	id := &org.Identity{
		Type: IdentityTypeOrgNr,
		Code: cbc.Code("974760673"),
	}
	if err := Validate(id); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateSwitchUnknownTypeNoop(t *testing.T) {
	if err := Validate(struct{}{}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
