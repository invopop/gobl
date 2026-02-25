package no

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
)

func TestValidOrgNrMod11(t *testing.T) {
	// Example from Brønnøysundregistrene: 974 760 673
	if !validOrgNrMod11("974760673") {
		t.Fatalf("expected orgnr to be valid")
	}
	if validOrgNrMod11("974760674") {
		t.Fatalf("expected orgnr to be invalid")
	}
}

func TestNormalizeOrgIdentity(t *testing.T) {
	id := &org.Identity{
		Type: IdentityTypeOrgNr,
		Code: cbc.Code("974 760 673"),
	}
	normalizeOrgIdentity(id)
	if id.Code != "974760673" {
		t.Fatalf("expected normalized code to be 974760673, got %s", id.Code)
	}
}

func TestValidateOrgIdentity(t *testing.T) {
	id := &org.Identity{
		Type: IdentityTypeOrgNr,
		Code: cbc.Code("974760674"),
	}
	if err := validateOrgIdentity(id); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestNormalizeOrgIdentityDoesNothingForOtherTypes(t *testing.T) {
	id := &org.Identity{
		Type: cbc.Code("XX"),
		Code: cbc.Code("974 760 673"),
	}
	normalizeOrgIdentity(id)
	if id.Code != "974 760 673" {
		t.Fatalf("expected code to remain unchanged, got %s", id.Code)
	}
}