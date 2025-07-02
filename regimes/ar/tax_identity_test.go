package ar_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/ar"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityNormalization(t *testing.T) {
	r := ar.New()
	tests := []struct {
		Code     cbc.Code
		Expected cbc.Code
	}{
		{
			Code:     "20-12345678-3",
			Expected: "20123456783",
		},
		{
			Code:     " 20 1234 5678 3 ",
			Expected: "20123456783",
		},
		{
			Code:     "AR20123456783",
			Expected: "20123456783",
		},
		{
			Code:     "30/12345678&9",
			Expected: "30123456789",
		},
		{
			Code:     "20123456783",
			Expected: "20123456783",
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: "AR", Code: ts.Code}
		r.NormalizeObject(tID)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestNormalizeTaxIdentity(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		tID := (*tax.Identity)(nil)
		assert.NotPanics(t, func() {
			ar.NormalizeTaxIdentity(tID)
		})
	})
}

func TestTaxIdentityValidation(t *testing.T) {

	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid 1", code: "20409378472"},
		{name: "valid 2", code: "30625913558"},
		{
			name: "with formatting",
			code: "20-40937847-2",
			err:  "contains invalid characters",
		},
		{
			name: "with prefix and spaces",
			code: " AR 20 4093 7847 2 ",
			err:  "contains invalid characters",
		},
		{
			name: "too short",
			code: "3012345678",
			err:  "invalid length",
		},
		{
			name: "non-digit characters",
			code: "30A2345678B",
			err:  "contains invalid characters",
		},
		{
			name: "bad check digit",
			code: "12345678901",
			err:  "invalid check digit",
		},
		{
			name: "empty",
			code: "",
			err:  "cannot be blank",
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			tID := &tax.Identity{
				Country: "AR",
				Code:    ts.code,
			}
			err := ar.Validate(tID)
			if ts.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), ts.err)
				}
			}
		})
	}
}
