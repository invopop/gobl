package ro_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	_ "github.com/invopop/gobl/regimes/ro"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

// Unit tests to check the algorithm inside tax_identity.go
func TestTaxIdentityRules(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "good 1 (real ANAF CUI)", code: "16281620"},
		{name: "good 2 (minimal length)", code: "19"},
		{name: "good 3 (max length)", code: "1234567897"},
		{
			name: "bad checksum",
			code: "16281621",
			err:  "IDENTITY-01",
		},
		{
			name: "too short",
			code: "1",
			err:  "IDENTITY-01",
		},
		{
			name: "too long",
			code: "12345678901",
			err:  "IDENTITY-01",
		},
		{
			name: "not numeric",
			code: "RO1628162",
			err:  "IDENTITY-01",
		},
		{
			name: "not normalized",
			code: "16-281-620",
			err:  "IDENTITY-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "RO", Code: tt.code}
			err := rules.Validate(tID)
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
