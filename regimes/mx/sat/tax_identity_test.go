package sat_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/mx"
	"github.com/invopop/gobl/regimes/mx/sat"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityNormalization(t *testing.T) {
	r := mx.New()
	tests := []struct {
		Code     cbc.Code
		Expected cbc.Code
	}{
		{
			Code:     "GHI70123123Z",
			Expected: "GHI70123123Z",
		},
		{
			Code:     " GHI 701231 23Z ",
			Expected: "GHI70123123Z",
		},
		{
			Code:     "GHI-701231-23Z",
			Expected: "GHI70123123Z",
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: "MX", Code: ts.Code}
		r.NormalizeObject(tID)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestTaxIdentityValidation(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		zone l10n.Code
		err  string
	}{
		{name: "foreign code", code: sat.TaxIdentityCodeForeign, zone: "21000"},
		{name: "valid code 1", code: "MNOP8201019HJ"},
		{name: "valid code 2", code: "UVWX610715JKL"},
		{name: "valid code 3", code: "STU760612MN1"},
		{
			name: "invalid code 1",
			code: "STU760612MN",
			err:  tax.ErrIdentityCodeInvalid.Error(),
		},
		{
			name: "invalid code 2",
			code: "XXXX",
			err:  tax.ErrIdentityCodeInvalid.Error(),
		},
		{
			name: "empty",
			code: "",
			// empty is allowed
			// err:  "code: cannot be blank",
		},
		{
			name: "missing zone",
			code: "MNOP8201019HJ",
			zone: "",
			// deprecated. err:  "zone: cannot be blank",
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "MX", Code: ts.code, Zone: ts.zone}
			err := mx.Validate(tID)
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

func TestTaxIdentityDetermineType(t *testing.T) {
	tests := []struct {
		Code cbc.Code
		Type cbc.Key
	}{
		{
			Code: "",
			Type: cbc.KeyEmpty,
		},
		{
			Code: sat.TaxIdentityCodeForeign,
			Type: sat.TaxIdentityTypePerson,
		},
		{
			Code: "MNOP8201019HJ",
			Type: sat.TaxIdentityTypePerson,
		},
		{
			Code: "ABC830720XYZ",
			Type: sat.TaxIdentityTypeCompany,
		},
		{
			Code: "XXXX",
			Type: cbc.KeyEmpty,
		},
	}
	for _, ts := range tests {
		t.Run(string(ts.Code), func(t *testing.T) {
			res := sat.DetermineTaxCodeType(ts.Code)
			assert.Equal(t, ts.Type, res)
		})
	}
}
