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
		{name: "validCPF", code: "01234567890"},
		{name: "validCPF", code: "35549549506"},

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
		{
			name: "invalid CPF wrong checksum",
			code: "11144477730",
			err:  "verification digit mismatch",
		},
		{
			name: "invalid CPF too short",
			code: "1114447773",
			err:  "must have 11 (CPF) or 14 (CNPJ) digits",
		},
		{
			name: "invalid CPF too long",
			code: "111444777356",
			err:  "must have 11 (CPF) or 14 (CNPJ) digits",
		},
		{
			name: "invalid CPF non-numeric",
			code: "11A44477735",
			err:  "must contain only digits",
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
		{name: "valid5", code: "012.345.678-90", want: "01234567890"},
		{name: "valid6", code: "125.217.532-97", want: "12521753297"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "BR", Code: tt.code}
			br.Normalize(tID)
			assert.Equal(t, tt.want, tID.Code)
		})
	}
}
