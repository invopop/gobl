package ro_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/ro"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name    string
		code    cbc.Code
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid CUI without RO prefix",
			code:    "18547290",
			wantErr: false,
		},
		{
			name:    "valid CUI with RO prefix",
			code:    "RO18547290",
			wantErr: false,
		},
		{
			name:    "valid short CUI (2 digits)",
			code:    "27",
			wantErr: false,
		},
		{
			name:    "valid max length CUI (10 digits)",
			code:    "1234567897",
			wantErr: false,
		},
		{
			name:    "invalid checksum",
			code:    "18547291",
			wantErr: true,
			errMsg:  "invalid checksum",
		},
		{
			name:    "invalid checksum with RO prefix",
			code:    "RO18547291",
			wantErr: true,
			errMsg:  "invalid checksum",
		},
		{
			name:    "too short - only 1 digit",
			code:    "1",
			wantErr: true,
			errMsg:  "invalid format", // Fails regex
		},
		{
			name:    "too long - more than 10 digits",
			code:    "12345678901",
			wantErr: true,
			errMsg:  "invalid format", // Fails regex
		},
		{
			name:    "contains letters",
			code:    "1854729A",
			wantErr: true,
			errMsg:  "invalid format",
		},
		{
			name:    "invalid format with special characters",
			code:    "1854-7290",
			wantErr: true,
			errMsg:  "invalid format",
		},
		{
			name:    "empty code",
			code:    "",
			wantErr: true,
			errMsg:  "cannot be blank",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{
				Country: "RO",
				Code:    tt.code,
			}
			err := ro.Validate(tID)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
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
			code:     "18547290",
			expected: "18547290",
		},
		{
			name:     "with RO prefix - normalized to no prefix",
			code:     "RO18547290",
			expected: "18547290",
		},
		{
			name:     "with spaces - normalized",
			code:     "1854 7290",
			expected: "18547290",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{
				Country: "RO",
				Code:    tt.code,
			}
			ro.Normalize(tID)
			assert.Equal(t, tt.expected, tID.Code)
		})
	}
}
