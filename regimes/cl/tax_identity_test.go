package cl_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/cl"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	t.Parallel()

	var tID *tax.Identity
	assert.NotPanics(t, func() {
		cl.Normalize(tID)
	})

	tests := []struct {
		Code     cbc.Code
		Expected cbc.Code
	}{
		{
			Code:     "12.345.678-5",
			Expected: "123456785",
		},
		{
			Code:     "11.111.111-1",
			Expected: "111111111",
		},
		{
			Code:     "12345678-5",
			Expected: "123456785",
		},
		{
			Code:     "76.123.456-K",
			Expected: "76123456K",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(string(tt.Code), func(t *testing.T) {
			t.Parallel()
			tID := &tax.Identity{Country: "CL", Code: tt.Code}
			cl.Normalize(tID)
			assert.Equal(t, tt.Expected, tID.Code)
		})
	}
}

func TestValidateTaxIdentity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		country l10n.TaxCountryCode
		code    cbc.Code
		err     string
	}{
		{
			name:    "valid RUT 1",
			country: "CL",
			code:    "713254975",
		},
		{
			name:    "valid RUT 2",
			country: "CL",
			code:    "111111111",
		},
		{
			name:    "valid RUT 3",
			country: "CL",
			code:    "100000008",
		},
		{
			name:    "valid RUT 4 with K",
			country: "CL",
			code:    "77668208K",
		},
		{
			name:    "valid RUT 5 with lowercase k",
			country: "CL",
			code:    "77668208k",
		},
		{
			name:    "valid short RUT",
			country: "CL",
			code:    "10000009",
		},
		{
			name:    "valid RUT with check digit 0",
			country: "CL",
			code:    "10000130",
		},
		{
			name:    "empty code",
			country: "CL",
			code:    "",
		},
		{
			name:    "too short - less than 7 digits",
			country: "CL",
			code:    "12345",
			err:     "invalid RUT format",
		},
		{
			name:    "invalid check digit",
			country: "CL",
			code:    "123456780",
			err:     "invalid RUT check digit",
		},
		{
			name:    "too long",
			country: "CL",
			code:    "1234567890123",
			err:     "invalid RUT format",
		},
		{
			name:    "invalid characters",
			country: "CL",
			code:    "1234567A5",
			err:     "invalid RUT format",
		},
		{
			name:    "wrong check digit K when should be number",
			country: "CL",
			code:    "12345678K",
			err:     "invalid RUT check digit",
		},
		{
			name:    "non-CL country - should not validate",
			country: "AR",
			code:    "invalid",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tID := &tax.Identity{Country: tt.country, Code: tt.code}
			err := cl.Validate(tID)
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

func TestValidateOtherDocuments(t *testing.T) {
	t.Parallel()

	type unsupportedDoc struct{}
	err := cl.Validate(&unsupportedDoc{})
	assert.NoError(t, err)

	err = cl.Validate(nil)
	assert.NoError(t, err)

	// Test with nil tax.Identity specifically
	var nilIdentity *tax.Identity
	err = cl.Validate(nilIdentity)
	assert.NoError(t, err)
}
