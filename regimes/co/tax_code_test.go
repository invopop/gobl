package co_test

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/co"
	"github.com/invopop/gobl/tax"
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
		tID := &tax.Identity{Country: l10n.ES, Code: ts.Code}
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
	assert.Equal(t, p.Addresses[0].Locality, "BOGOTÁ, D.C.")
	assert.Equal(t, p.Addresses[0].Region, "Bogotá")
}

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code string
		zone l10n.Code
		err  string
	}{
		{name: "good 1", code: "412615332", zone: "11001"},
		{name: "good 2", code: "8110079918", zone: "11001"},
		{name: "good 3", code: "124499654", zone: "08638"},
		{name: "good 4", code: "8300801501", zone: "11001"},
		{name: "good 5", code: "700602703", zone: "11001"},
		{
			name: "empty",
			code: "",
			zone: "11001",
			err:  "code: cannot be blank",
		},
		{
			name: "no zone",
			code: "412615332",
			zone: "",
			err:  "zone: cannot be blank",
		},
		{
			name: "invalid zone",
			code: "412615332",
			zone: "99999",
			err:  "zone: must be a valid value",
		},
		{
			name: "too long",
			code: "123456789100",
			zone: "11001",
			err:  "too long",
		},
		{
			name: "too short",
			code: "123456",
			zone: "11001",
			err:  "too short",
		},
		{
			name: "not normalized",
			code: "12.449.965-4",
			zone: "11001",
			err:  "contains invalid characters",
		},
		{
			name: "bad checksum",
			code: "412615331",
			zone: "11001",
			err:  "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: l10n.CO, Code: tt.code, Zone: tt.zone}
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
