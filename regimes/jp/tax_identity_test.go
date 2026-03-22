package jp_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/jp"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		name     string
		code     cbc.Code
		expected cbc.Code
	}{
		{name: "already clean", code: "9700150098417", expected: "9700150098417"},
		{name: "with spaces", code: "9700 1500 9841 7", expected: "9700150098417"},
		{name: "with hyphens", code: "9700-1500-9841-7", expected: "9700150098417"},
		{name: "with T prefix", code: "T9700150098417", expected: "9700150098417"},
		{name: "with T prefix and spaces", code: "T 9700 1500 9841 7", expected: "9700150098417"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "JP", Code: tt.code}
			jp.Normalize(tID)
			assert.Equal(t, tt.expected, tID.Code)
		})
	}
}

func TestNormalizeTaxIdentityNil(_ *testing.T) {
	// Should not panic on nil.
	jp.Normalize(nil)
}

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		// Valid Corporate Numbers
		{name: "valid checksum (1st digit) for real NTA example", code: "9700150098417"},
		{name: "valid checksum (1st digit)", code: "5050005005266"},
		{name: "empty code", code: ""},

		// Invalid formats
		{name: "too short", code: "970015009841", err: "must be a 13-digit number"},
		{name: "too long", code: "97001500984170", err: "must be a 13-digit number"},
		{name: "non-numeric", code: "970015009841A", err: "must be a 13-digit number"},
		{name: "with hyphens", code: "9700-150-09841", err: "must be a 13-digit number"},

		// Invalid checksum
		{name: "bad checksum (1st digit) for real NTA example", code: "1700150098417", err: "invalid checksum"},
		{name: "bad checksum (1st digit)", code: "9050005005266", err: "invalid checksum"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "JP", Code: tt.code}
			err := jp.Validate(tID)
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
