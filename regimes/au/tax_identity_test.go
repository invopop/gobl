package au_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/au"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	var tID *tax.Identity
	assert.NotPanics(t, func() {
		au.Normalize(tID)
	}, "nil tax identity")

	tests := []struct {
		Code     cbc.Code
		Expected cbc.Code
	}{
		{
			Code:     "51 824 753 556",
			Expected: "51824753556",
		},
		{
			Code:     "51-824-753-556",
			Expected: "51824753556",
		},
		{
			Code:     " 51824753556 ",
			Expected: "51824753556",
		},
		{
			Code:     "51.824.753.556",
			Expected: "51824753556",
		},
		{
			Code:     "51_824_753_556",
			Expected: "51824753556",
		},
		{
			Code:     "51 824-753.556",
			Expected: "51824753556",
		},
		{
			Code:     "53004085616",
			Expected: "53004085616",
		},
		{
			Code:     "",
			Expected: "",
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: "AU", Code: ts.Code}
		au.Normalize(tID)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid ABN example", code: "51824753556"},
		{name: "valid ABN 2", code: "53004085616"},
		{
			name: "invalid checksum",
			code: "51824753557",
			err:  "invalid ABN checksum",
		},
		{
			name: "too short",
			code: "5182475355",
			err:  "must be 11 digits",
		},
		{
			name: "too long",
			code: "518247535567",
			err:  "must be 11 digits",
		},
		{
			name: "non-numeric",
			code: "5182475355A",
			err:  "must be 11 digits",
		},
		{
			name: "all zeros",
			code: "00000000000",
			err:  "invalid ABN checksum",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "AU", Code: tt.code}
			err := au.Validate(tID)
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
