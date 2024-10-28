package hu_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/hu"
	"github.com/stretchr/testify/assert"
)

func TestValidateGroupId(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{"Empty code", "", ""},
		{"Invalid length (5)", "12345", "invalid length"},
		{"Invalid length (10)", "1234567890", "invalid length"},
		{"Invalid VAT code", "21114445123", "invalid VAT code"},
		{"Invalid area code", "82713452101", "invalid area code"},
		{"Valid code (11 chars)", "88212131403", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &org.Identity{Country: "HU", Code: tt.code, Key: hu.IdentityKeyGroupNumber}
			err := hu.Validate(tID)
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
