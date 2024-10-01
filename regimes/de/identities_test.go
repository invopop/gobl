package de_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/de"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestIdentityNormalization(t *testing.T) {
	r := tax.RegimeDefFor("DE")

	id := &org.Identity{
		Key:  de.IdentityKeyTaxNumber,
		Code: "123 / 456.78901 Ab",
	}
	r.NormalizeObject(id)
	assert.Equal(t, "12345678901", id.Code.String())
}

func TestTaxNumberValidation(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid 10 digits", code: "8609574271"},
		{name: "valid 11 digits", code: "65929970489"},
		{name: "valid 12 digits", code: "575492850137"},
		{name: "valid 13 digits", code: "9947036892816"},

		// Invalid formats
		{name: "too short", code: "12345678", err: "invalid length"},
		{name: "too long", code: "12345678901234", err: "invalid length"},
		{name: "non-numeric", code: "123ABC78901", err: "should only contain digits"},
		{name: "starts with zero", code: "01234567890", err: "first digit cannot be 0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Key: de.IdentityKeyTaxNumber, Code: tt.code}
			err := de.Validate(id)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.err)
				}
			}
		})
	}
}

func TestTaxIDValidation(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid 11 digits 1", code: "86095742719"},
		{name: "valid 11 digits 2", code: "65929970489"},
		{name: "valid 11 digits 3", code: "57549285017"},
		{name: "valid 13 digits", code: "9947036892816"},

		// Invalid formats
		{name: "too short", code: "12345678", err: "invalid length"},
		{name: "too long", code: "12345678901234", err: "invalid length"},
		{name: "non-numeric", code: "123ABC78901", err: "should only contain digits"},
		{name: "starts with zero", code: "01234567890", err: "first digit cannot be 0"},

		// Invalid digit occurrence
		{name: "no digit twice", code: "12345678901", err: "digit occurrence rule not satisfied"},
		{name: "two digits twice", code: "11223344556", err: "more than one digit appears twice"},
		{name: "digit four times", code: "11112222333", err: "a digit appears more than three times"},

		// Checksum tests
		{name: "invalid checksum 1", code: "12345678910", err: "checksum mismatch"},
		{name: "invalid checksum 2", code: "98765432130", err: "checksum mismatch"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Key: de.IdentityKeyTaxID, Code: tt.code}
			err := de.Validate(id)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.err)
				}
			}
		})
	}
}
