package no

import (
	"testing"

	"github.com/invopop/gobl/cbc"
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