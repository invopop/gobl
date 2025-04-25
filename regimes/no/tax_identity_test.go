// Package no_test provides tests for the Norwegian TRN (Tax Registration Number) validation.
package no_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/no"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid TRN 1", code: "123456785"},
		{name: "valid TRN 2", code: "290883970"},
		{name: "valid TRN 3", code: "974760673"},

		// Invalid formats
		{name: "too short", code: "12345678", err: "must be a 9-digit number"},
		{name: "too long", code: "1234567890", err: "must be a 9-digit number"},
		{name: "non-numeric", code: "12345ABCD", err: "must be a 9-digit number"},
		{name: "invalid checksum", code: "123456789", err: "invalid checksum for TRN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "NO", Code: tt.code}
			err := no.Validate(tID)
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

func TestValidateTRNCode(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		err   string
	}{
		{name: "non-cbc.Code input", value: 12345, err: ""},
		{name: "empty code", value: cbc.Code(""), err: ""},
		{name: "invalid format", value: cbc.Code("12345"), err: "must be a 9-digit number"},
		{name: "invalid checksum", value: cbc.Code("123456789"), err: "invalid checksum for TRN"},
		{name: "valid TRN", value: cbc.Code("290883970"), err: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := no.ValidateTRNCode(tt.value)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.err)
			}
		})
	}
}

func TestValidateChecksum(t *testing.T) {
	tests := []struct {
		name string
		trn  string
		want bool
	}{
		{name: "valid checksum", trn: "290883970", want: true},
		{name: "invalid checksum (10)", trn: "000000060", want: false},
		{name: "invalid checksum (random)", trn: "123456789", want: false},
		{name: "valid checksum (11 treated as 0)", trn: "974760673", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := no.ValidateChecksum(tt.trn)
			assert.Equal(t, tt.want, got)
		})
	}
}
