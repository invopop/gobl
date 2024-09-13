package co_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/co"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	var tID *tax.Identity
	assert.NotPanics(t, func() {
		co.Normalize(tID)
	})

	tests := []struct {
		Code     cbc.Code
		Expected cbc.Code
	}{
		{
			Code:     "XX",
			Expected: "XX",
		},
		{
			Code:     "901.458.652-7",
			Expected: "9014586527",
		},
		{
			Code:     "800.134.536-3",
			Expected: "8001345363",
		},
		{
			Code:     "36029785",
			Expected: "36029785",
		},
		{
			Code:     "100 100 100",
			Expected: "100100100",
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: "CO", Code: ts.Code}
		co.Normalize(tID)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "good 1", code: "412615332"},
		{name: "good 2", code: "8110079918"},
		{name: "good 3", code: "124499654"},
		{name: "good 4", code: "8300801501"},
		{name: "good 5", code: "700602703"},
		{name: "good no tin", code: "700602703"},
		{
			name: "too long",
			code: "123456789100",
			err:  "too long",
		},
		{
			name: "too short",
			code: "123456",
			err:  "too short",
		},
		{
			name: "not normalized",
			code: "12.449.965-4",
			err:  "contains invalid characters",
		},
		{
			name: "bad checksum",
			code: "412615331",
			err:  "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "CO", Code: tt.code}
			err := co.Validate(tID)
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
