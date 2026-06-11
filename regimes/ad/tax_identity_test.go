package ad_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/ad"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid L", code: "L123456A"},
		{name: "valid U", code: "U765432Z"},
		{name: "valid F", code: "F111222B"},
		{name: "valid A", code: "A654321C"},
		{name: "valid C", code: "C000001Z"},
		{name: "valid D", code: "D999999K"},
		{name: "empty code", code: ""},
		{name: "too short", code: "L12345A", err: "invalid NRT format"},
		{name: "too long", code: "L1234567AB", err: "invalid NRT format"},
		{name: "invalid prefix", code: "Z123456A", err: "invalid NRT format"},
		{name: "no control letter", code: "L1234567", err: "invalid NRT format"},
		{name: "lowercase", code: "l123456a", err: "invalid NRT format"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "AD", Code: tt.code}
			err := ad.Validate(tID)
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

func TestTaxIdentityGeneralCases(t *testing.T) {
	tests := []struct {
		name string
		tID  *tax.Identity
		err  string
	}{
		{name: "just country", tID: &tax.Identity{Country: "AD"}},
		{name: "valid NRT", tID: &tax.Identity{Country: "AD", Code: "L123456A"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ad.Validate(tt.tID)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), tt.err)
			}
		})
	}
}
