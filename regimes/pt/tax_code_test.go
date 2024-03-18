package pt_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		Code     cbc.Code
		Expected cbc.Code
	}{
		{
			Code:     "901.458.652",
			Expected: "901458652",
		},
		{
			Code:     "800.134.536",
			Expected: "800134536",
		},
		{
			Code:     "PT 800.134.536",
			Expected: "800134536",
		},
		{
			Code:     "36029785",
			Expected: "36029785",
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: l10n.PT, Code: ts.Code}
		err := pt.Calculate(tID)
		assert.NoError(t, err)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "good 1", code: "999999990"},
		{name: "good 2", code: "287024008"},
		{name: "good 3", code: "501442600"},
		{name: "good 4", code: "501442600"},
		{
			name: "empty",
			code: "",
			err:  "code: cannot be blank",
		},
		{
			name: "invalid zone",
			code: "420000000",
			err:  "invalid prefix",
		},
		{
			name: "too long",
			code: "123456789100",
			err:  "invalid length",
		},
		{
			name: "too short",
			code: "123456",
			err:  "invalid length",
		},
		{
			name: "not normalized",
			code: "12.449.965-4",
			err:  "contains invalid characters",
		},
		{
			name: "bad checksum",
			code: "999999991",
			err:  "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: l10n.PT, Code: tt.code}
			err := pt.Validate(tID)
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
