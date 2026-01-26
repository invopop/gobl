package dk_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/dk"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid 1", code: "13585628"},
		{name: "valid 2", code: "88146328"},
		{name: "valid 3", code: "25063864"},
		{
			name: "empty code",
			code: "",
			err:  "",
		},
		{
			name: "too short",
			code: "1234567",
			err:  "invalid format",
		},
		{
			name: "too long",
			code: "123456789",
			err:  "invalid format",
		},
		{
			name: "contains letters",
			code: "1234567A",
			err:  "invalid format",
		},
		{
			name: "bad checksum",
			code: "13585627",
			err:  "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "DK", Code: tt.code}
			err := dk.Validate(tID)
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
			name:     "with DK prefix",
			code:     "DK13585628",
			expected: "13585628",
		},
		{
			name:     "with spaces",
			code:     "13 58 56 28",
			expected: "13585628",
		},
		{
			name:     "already normalized",
			code:     "13585628",
			expected: "13585628",
		},
		{
			name:     "empty",
			code:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "DK", Code: tt.code}
			dk.Normalize(tID)
			assert.Equal(t, tt.expected, tID.Code)
		})
	}
}
