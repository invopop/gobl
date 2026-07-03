package cl_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		code     cbc.Code
		expected cbc.Code
	}{
		{
			code:     "76.086.428-5",
			expected: "760864285",
		},
		{
			code:     "CL 12.531.909-2",
			expected: "125319092",
		},
		{
			code:     "1.000.005-k",
			expected: "1000005K",
		},
	}

	for _, tt := range tests {
		t.Run(tt.code.String(), func(t *testing.T) {
			tID := &tax.Identity{Country: "CL", Code: tt.code}
			norm.Normalize(tID)
			assert.Equal(t, tt.expected, tID.Code)
		})
	}
}

func TestTaxIdentityRules(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{
			name: "valid numeric check digit 1",
			code: "760864285",
		},
		{
			name: "valid numeric check digit 2",
			code: "125319092",
		},
		{
			name: "valid zero check digit",
			code: "10000130",
		},
		{
			name: "valid K check digit",
			code: "1000005K",
		},
		{
			name: "invalid checksum",
			code: "125319093",
			err:  "IDENTITY-01",
		},
		{
			name: "invalid check digit character",
			code: "12531909X",
			err:  "IDENTITY-01",
		},
		{
			name: "invalid body characters",
			code: "12A319092",
			err:  "IDENTITY-01",
		},
		{
			name: "invalid length",
			code: "1234567",
			err:  "IDENTITY-01",
		},
		{
			name: "empty code is accepted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "CL", Code: tt.code}
			err := rules.Validate(tID)
			if tt.err == "" {
				assert.NoError(t, err)
			} else if assert.Error(t, err) {
				assert.Contains(t, err.Error(), tt.err)
			}
		})
	}
}
