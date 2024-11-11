package br_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/br"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityValidation(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid1", code: "05104582000170"},
		{name: "valid2", code: "10909402000167"},
		{
			name: "too long",
			code: "123456789012345",
			err:  "must have 14 digits",
		},
		{
			name: "too short",
			code: "1234567890123",
			err:  "must have 14 digits",
		},
		{
			name: "non-numeric",
			code: "A2345678901234",
			err:  "must contain only digits",
		},
		{
			name: "non-numeric verification digit",
			code: "123456789012AB",
			err:  "must contain only digits",
		},
		{
			name: "first verification digit wrong",
			code: "05104582000160",
			err:  "verification digit mismatch",
		},
		{
			name: "second verification digit wrong",
			code: "05104582000171",
			err:  "verification digit mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "BR", Code: tt.code}
			err := br.Validate(tID)
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

func TestTaxIdentityNormalization(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		want cbc.Code
	}{
		{name: "valid1", code: "05104582000170", want: "05104582000170"},
		{name: "valid2", code: "10909402000167", want: "10909402000167"},
		{name: "valid3", code: "051.045.820/0017-0", want: "05104582000170"},
		{name: "valid4", code: "109.094.020/0016-7", want: "10909402000167"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "BR", Code: tt.code}
			br.Normalize(tID)
			assert.Equal(t, tt.want, tID.Code)
		})
	}
}
