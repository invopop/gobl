package ro_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/ro"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		name     string
		code     cbc.Code
		expected cbc.Code
	}{
		{name: "with RO prefix", code: "RO13547272", expected: "13547272"},
		{name: "with lowercase ro prefix", code: "ro13547272", expected: "13547272"},
		{name: "with mixed case ro prefix", code: "Ro13547272", expected: "13547272"},
		{name: "with special characters", code: "RO-13547272", expected: "13547272"},
		{name: "with spaces", code: "135 472 72", expected: "13547272"},
		{name: "already normalized", code: "13547272", expected: "13547272"},
		{name: "empty", code: "", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "RO", Code: tt.code}
			ro.Normalize(tID)
			assert.Equal(t, tt.expected, tID.Code)
		})
	}
}

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "Valid 1 - Romanian Tech SRL", code: "13547272"},
		{name: "Valid 2 - Global Logistics SRL", code: "10864098"},
		{name: "Valid 3 - Storage Shop SRL", code: "4376262"},
		{name: "empty code", code: ""},
		{name: "too short", code: "1", err: "invalid format"},
		{name: "too long", code: "12345678901", err: "invalid format"},
		{name: "contains letters", code: "1354727A", err: "invalid format"},
		{name: "bad checksum", code: "13547271", err: "checksum mismatch"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "RO", Code: tt.code}
			err := ro.Validate(tID)
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
