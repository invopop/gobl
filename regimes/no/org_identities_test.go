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

func TestValidateOrgIdentityNilIsNoop(t *testing.T) {
	if err := validateOrgIdentity(nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateOrgIdentityValidPasses(t *testing.T) {
	id := &org.Identity{
		Type: IdentityTypeOrgNr,
		Code: cbc.Code("974 760 673"),
	}
	if err := validateOrgIdentity(id); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateOrgIdentityInvalidFormat(t *testing.T) {
	id := &org.Identity{
		Type: IdentityTypeOrgNr,
		Code: cbc.Code("123"), // not 9 digits
	}
	if err := validateOrgIdentity(id); err == nil {
		t.Fatalf("expected validation error for invalid format")
	}
}

func TestValidOrgNrMod11InvalidLength(t *testing.T) {
	if validOrgNrMod11("123") {
		t.Fatalf("expected false for invalid length")
	}
}

func TestValidOrgNrMod11NonDigit(t *testing.T) {
	// 9 chars but includes a non-digit -> should be false
	if validOrgNrMod11("97476067A") {
		t.Fatalf("expected false for non-digit char")
	}
}
func TestValidOrgNrMod11CheckDigitTenIsInvalid(t *testing.T) {
	// Crafted so sum%11==1 -> check digit would be 10 -> invalid by definition
	if validOrgNrMod11("370000000") {
		t.Fatalf("expected false when computed check digit is 10")
	}
}
