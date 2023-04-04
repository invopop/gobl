package co_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
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

func TestNormalizeParty(t *testing.T) {
	p := &org.Party{
		Name: "Test Party",
		TaxID: &tax.Identity{
			Country: l10n.CO,
			Code:    "412615332",
			Zone:    "11001",
		},
		Addresses: []*org.Address{
			{
				Locality: "Foo",
			},
		},
	}
	err := co.Calculate(p)
	assert.NoError(t, err)
	assert.Equal(t, p.Addresses[0].Locality, "Bogotá, D.C.")
	assert.Equal(t, p.Addresses[0].Region, "Bogotá")

	p = &org.Party{
		Name: "Test Party No Zone",
		TaxID: &tax.Identity{
			Country: l10n.CO,
			Type:    co.TaxIdentityTypeCitizen,
			Code:    "100100100",
		},
		Addresses: []*org.Address{
			{
				Locality: "Foo",
			},
		},
	}
	err = co.Calculate(p)
	require.NoError(t, err)
	err = co.Validate(p)
	assert.NoError(t, err)
}

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		typ  cbc.Key
		code cbc.Code
		zone l10n.Code
		err  string
	}{
		{name: "good 1", typ: "tin", code: "412615332", zone: "11001"},
		{name: "good 2", typ: "tin", code: "8110079918", zone: "11001"},
		{name: "good 3", typ: "tin", code: "124499654", zone: "08638"},
		{name: "good 4", typ: "tin", code: "8300801501", zone: "11001"},
		{name: "good 5", typ: "tin", code: "700602703", zone: "11001"},
		{
			name: "missing code",
			typ:  "tin",
			code: "",
			zone: "11001",
			err:  "code: cannot be blank",
		},
		{
			name: "missing zone for citizen",
			typ:  co.TaxIdentityTypeCitizen,
			code: "100100100",
			zone: "",
			err:  "",
		},
		{
			name: "missing type",
			typ:  "",
			code: "100100100",
			zone: "11001",
			err:  "type: cannot be blank",
		},
		{
			name: "no zone",
			typ:  "tin",
			code: "412615332",
			zone: "",
			err:  "zone: cannot be blank",
		},
		{
			name: "invalid zone",
			typ:  "tin",
			code: "412615332",
			zone: "99999",
			err:  "zone: must be a valid value",
		},
		{
			name: "too long",
			typ:  "tin",
			code: "123456789100",
			zone: "11001",
			err:  "too long",
		},
		{
			name: "too short",
			typ:  "tin",
			code: "123456",
			zone: "11001",
			err:  "too short",
		},
		{
			name: "not normalized",
			typ:  "tin",
			code: "12.449.965-4",
			zone: "11001",
			err:  "contains invalid characters",
		},
		{
			name: "bad checksum",
			typ:  "tin",
			code: "412615331",
			zone: "11001",
			err:  "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: l10n.CO, Type: tt.typ, Code: tt.code, Zone: tt.zone}
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
