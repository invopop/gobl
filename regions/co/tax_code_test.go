package co_test

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regions/co"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		Code     string
		Expected string
	}{
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
	}
	for _, ts := range tests {
		tID := &org.TaxIdentity{Country: l10n.ES, Code: ts.Code}
		err := co.NormalizeTaxIdentity(tID)
		assert.NoError(t, err)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code string
		err  string
	}{
		{name: "empty", code: ""},
		{name: "good 1", code: "412615332"},
		{name: "good 2", code: "8110079918"},
		{name: "good 3", code: "124499654"},
		{name: "good 4", code: "8300801501"},
		{name: "good 5", code: "700602703"},
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
			tID := &org.TaxIdentity{Country: l10n.CO, Code: tt.code}
			err := co.ValidateTaxIdentity(tID)
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
