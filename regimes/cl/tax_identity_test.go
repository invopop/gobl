package cl_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	_ "github.com/invopop/gobl/regimes/cl"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	var tID *tax.Identity
	assert.NotPanics(t, func() {
		norm.Normalize(tID)
	}, "nil tax identity")

	tests := []struct {
		Code     cbc.Code
		Expected cbc.Code
	}{
		{
			Code:     "76.086.428-5",
			Expected: "760864285",
		},
		{
			Code:     "12.345.670-k",
			Expected: "12345670K",
		},
		{
			Code:     " 18.765.432-7 ",
			Expected: "187654327",
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: "CL", Code: ts.Code}
		norm.Normalize(tID)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestTaxIdentityRules(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		// The historic "1-9" RUT (Bernardo O'Higgins), frequently used in
		// Chile as a test/example identifier — single-digit body.
		{name: "good 1", code: "19"},
		{name: "good 2 (company)", code: "760864285"},
		{name: "good 3 (check digit K)", code: "12345670K"},
		{name: "good 4 (check digit 0)", code: "761234560"},
		{name: "good 5", code: "187654327"},
		{
			name: "too short",
			code: "1",
			err:  "IDENTITY-01",
		},
		{
			name: "body too long",
			code: "123456789K",
			err:  "IDENTITY-01",
		},
		{
			name: "invalid check character",
			code: "76086428X",
			err:  "IDENTITY-01",
		},
		{
			name: "not normalized",
			code: "76.086.428-5",
			err:  "IDENTITY-01",
		},
		{
			name: "bad checksum",
			code: "760864281",
			err:  "IDENTITY-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "CL", Code: tt.code}
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
