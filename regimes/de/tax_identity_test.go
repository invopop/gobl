package de_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/de"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "good 1", code: "111111125"},
		{name: "good 2", code: "160459932"},
		{name: "good 3", code: "282741168"},
		{name: "good 4", code: "813495425"},
		{
			name: "zeros",
			code: "000000000",
			err:  "invalid format",
		},
		{
			name: "start with zero",
			code: "011111112",
			err:  "invalid format",
		},
		{
			name: "bad mid length",
			code: "12345678910",
			err:  "invalid format",
		},
		{
			name: "too long",
			code: "1234567890123",
			err:  "invalid format",
		},
		{
			name: "too short",
			code: "123456",
			err:  "invalid format",
		},
		{
			name: "not normalized",
			code: "12.449.965-4",
			err:  "invalid format",
		},
		{
			name: "bad checksum",
			code: "999999991",
			err:  "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "DE", Code: tt.code}
			err := de.Validate(tID)
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
