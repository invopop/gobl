package no

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
)

func TestValidateTaxCodeAcceptsMVASuffix(t *testing.T) {
	if err := validateTaxCode(cbc.Code("974760673MVA")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := validateTaxCode(cbc.Code("974 760 673 mva")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateTaxCodeInvalidLength(t *testing.T) {
	if err := validateTaxCode(cbc.Code("123")); err == nil {
		t.Fatalf("expected error for invalid length")
	}
}

func TestValidateTaxCodeAcceptsNOPrefixAndMVASuffix(t *testing.T) {
	if err := validateTaxCode(cbc.Code("NO974760673MVA")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateTaxCodeInvalidChecksum(t *testing.T) {
	if err := validateTaxCode(cbc.Code("974760674MVA")); err == nil {
		t.Fatalf("expected error for invalid checksum")
	}
}

func TestValidateTaxIdentityNilNoop(t *testing.T) {
	if err := validateTaxIdentity(nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateTaxIdentityValid(t *testing.T) {
	id := &tax.Identity{
		Country: l10n.TaxCountryCode(l10n.NO),
		Code:    cbc.Code("974760673MVA"),
	}
	if err := validateTaxIdentity(id); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateTaxIdentityInvalid(t *testing.T) {
	id := &tax.Identity{
		Country: l10n.TaxCountryCode(l10n.NO),
		Code:    cbc.Code("123"),
	}
	if err := validateTaxIdentity(id); err == nil {
		t.Fatalf("expected error")
	}
}

func TestValidateTaxCodeEmptyIsNoop(t *testing.T) {
	if err := validateTaxCode(cbc.Code("")); err != nil {
		t.Fatalf("expected no error for empty code, got %v", err)
	}
}
