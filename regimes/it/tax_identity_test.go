package it_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	_ "github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityRules(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		zone string
		err  string
	}{
		{name: "good 1", code: "12345678903"},
		{name: "good 2", code: "13029381004"},
		{name: "good 3", code: "10182640150"},
		{
			name: "empty",
			code: "",
			err:  "",
		},
		{
			name: "too long",
			code: "123456789001",
			err:  "IDENTITY-01",
		},
		{
			name: "too short",
			code: "1234567890",
			err:  "IDENTITY-01",
		},
		{
			name: "not normalized",
			code: "12.449.965-439",
			err:  "IDENTITY-01",
		},
		{
			name: "includes non-numeric characters",
			code: "A764352056Z",
			err:  "IDENTITY-01",
		},
		{
			name: "invalid check digit",
			code: "12345678901",
			err:  "IDENTITY-01",
		},
		{
			name: "invalid check digit",
			code: "13029381009",
			err:  "IDENTITY-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "IT", Code: tt.code}
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

func TestTaxIdentityValidateGeneralCases(t *testing.T) {
	tests := []struct {
		name string
		tID  *tax.Identity
		err  string
	}{
		{
			name: "just country",
			tID:  &tax.Identity{Country: "IT"},
			err:  "",
		},
		{
			name: "no type, assume biz",
			tID:  &tax.Identity{Country: "IT", Code: "12345678903"},
			err:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rules.Validate(tt.tID)
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
