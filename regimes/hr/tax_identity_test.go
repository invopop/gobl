package hr_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/hr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{
			name: "empty code",
			code: "",
		},
		{
			name: "valid OIB 1",
			code: "12345678903",
		},
		{
			name: "valid OIB 2",
			code: "11111111119",
		},
		{
			name: "too short",
			code: "1234567890",
			err:  "invalid format",
		},
		{
			name: "too long",
			code: "123456789012",
			err:  "invalid format",
		},
		{
			name: "contains letters",
			code: "1234567890A",
			err:  "invalid format",
		},
		{
			name: "invalid checksum",
			code: "12345678900",
			err:  "invalid checksum",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "HR", Code: tt.code}
			err := hr.Validate(tID)
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

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		name     string
		code     cbc.Code
		expected cbc.Code
	}{
		{
			name:     "already normalized",
			code:     "12345678903",
			expected: "12345678903",
		},
		{
			name:     "with spaces",
			code:     "1234 5678 903",
			expected: "12345678903",
		},
		{
			name:     "empty",
			code:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "HR", Code: tt.code}
			hr.Normalize(tID)
			assert.Equal(t, tt.expected, tID.Code)
		})
	}
}
