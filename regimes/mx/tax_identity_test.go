package mx_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/mx"
	"github.com/invopop/gobl/rules"
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
			Code:     "MXGHI70123123Z",
			Expected: "GHI70123123Z",
		},
		{
			Code:     "MXMXI70123123Z",
			Expected: "MXI70123123Z",
		},
		{
			Code:     "MXI70123123Z",
			Expected: "MXI70123123Z",
		},
		{
			Code:     " GHI 701231 23Z ",
			Expected: "GHI70123123Z",
		},
		{
			Code:     "GHI-701231-23Z",
			Expected: "GHI70123123Z",
		},
		{
			Code:     "K&A010301I16",
			Expected: "K&A010301I16",
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: "MX", Code: ts.Code}
		r.NormalizeObject(tID)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestTaxIdentityRules(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		zone l10n.Code
		err  string
	}{
		{name: "foreign code", code: mx.TaxIdentityCodeForeign, zone: "21000"},
		{name: "valid code 1", code: "MNOP8201019HJ"},
		{name: "valid code 2", code: "UVWX610715JKL"},
		{name: "valid code 3", code: "STU760612MN1"},
		{
			name: "valid with symbol",
			code: "K&A010301I16",
		},
		{
			name: "invalid code 1",
			code: "STU760612MN",
			err:  "IDENTITY-01",
		},
		{
			name: "invalid code 2",
			code: "XXXX",
			err:  "IDENTITY-01",
		},
		{
			name: "empty",
			code: "",
			// empty is allowed
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
			tID := &tax.Identity{Country: "MX", Code: ts.code}
			err := rules.Validate(tID)
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
