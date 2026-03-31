package sa_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	_ "github.com/invopop/gobl/regimes/sa"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityRules(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{
			name: "valid 15-digit starting and ending with 3",
			code: "312345678912343",
		},
		{
			name: "valid another pattern",
			code: "399999999900003",
		},
		{
			name: "valid all zeros in middle",
			code: "300000000000003",
		},
		{
			name: "empty code is allowed",
			code: "",
		},
		{
			name: "does not start with 3",
			code: "212345678912343",
			err:  "IDENTITY-01",
		},
		{
			name: "does not end with 3",
			code: "312345678912341",
			err:  "IDENTITY-01",
		},
		{
			name: "too short",
			code: "31234567893",
			err:  "IDENTITY-01",
		},
		{
			name: "too long",
			code: "3123456789123430",
			err:  "IDENTITY-01",
		},
		{
			name: "contains letters",
			code: "31234567891234A",
			err:  "IDENTITY-01",
		},
		{
			name: "only digits but wrong prefix and suffix",
			code: "112345678912341",
			err:  "IDENTITY-01",
		},
		{
			name: "14 digits",
			code: "31234567891233",
			err:  "IDENTITY-01",
		},
		{
			name: "16 digits",
			code: "3123456789123433",
			err:  "IDENTITY-01",
		},
		{
			name: "special characters",
			code: "3-1234567891234-3",
			err:  "IDENTITY-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "SA", Code: tt.code}
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

func TestTaxIdentityNormalization(t *testing.T) {
	t.Run("normalizes identity", func(t *testing.T) {
		tID := &tax.Identity{Country: "SA", Code: "312345678912343"}
		rd := tax.RegimeDefFor("SA")
		assert.NotNil(t, rd)
		rd.NormalizeObject(tID)
		assert.Equal(t, cbc.Code("312345678912343"), tID.Code)
	})
}
