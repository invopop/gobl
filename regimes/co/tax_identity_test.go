package co_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/co"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	tID := &tax.Identity{Country: l10n.CO, Code: "901.458.652-7"}
	err := co.Calculate(tID)
	require.NoError(t, err)
	assert.Equal(t, co.TaxIdentityTypeTIN, tID.Type, "autoassign type")

	tID = &tax.Identity{Country: l10n.CO, Type: co.TaxIdentityTypeCivil, Code: "XX"}
	err = co.Calculate(tID)
	require.NoError(t, err)
	assert.Equal(t, co.TaxIdentityTypeCivil, tID.Type, "copy type")

	tests := []struct {
		Code     cbc.Code
		Type     cbc.Key
		Expected cbc.Code
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
		{
			Code:     "100 100 100",
			Type:     co.TaxIdentityTypeCivil,
			Expected: "100100100",
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: l10n.CO, Code: ts.Code}
		err := co.Calculate(tID)
		assert.NoError(t, err)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		typ  cbc.Key
		code cbc.Code
		err  string
	}{
		{name: "good 1", typ: "tin", code: "412615332"},
		{name: "good 2", typ: "tin", code: "8110079918"},
		{name: "good 3", typ: "tin", code: "124499654"},
		{name: "good 4", typ: "tin", code: "8300801501"},
		{name: "good 5", typ: "tin", code: "700602703"},
		{name: "good no tin", code: "700602703"},
		{name: "ignore other typ", typ: "passport", code: "1234"},
		{
			name: "too long",
			typ:  "tin",
			code: "123456789100",
			err:  "too long",
		},
		{
			name: "too short",
			typ:  "tin",
			code: "123456",
			err:  "too short",
		},
		{
			name: "not normalized",
			typ:  "tin",
			code: "12.449.965-4",
			err:  "contains invalid characters",
		},
		{
			name: "bad checksum",
			typ:  "tin",
			code: "412615331",
			err:  "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: l10n.CO, Type: tt.typ, Code: tt.code}
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
