package nl_test

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regions/nl"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		Code     string
		Expected string
	}{
		{
			Code:     "000099995b57",
			Expected: "000099995B57",
		},
		{
			Code:     "NL000099995b57",
			Expected: "000099995B57",
		},
		{
			Code:     " 4359 6386 R ",
			Expected: "43596386R",
		},
	}
	for _, ts := range tests {
		tID := &org.TaxIdentity{Country: l10n.ES, Code: ts.Code}
		err := nl.NormalizeTaxIdentity(tID)
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
		{
			name: "empty",
			code: "",
		},
		{
			name: "too long",
			code: "a really really long string that's way too long",
			err:  "invalid length",
		},
		{
			name: "too short",
			code: "shorty",
			err:  "invalid length",
		},
		{
			name: "valid",
			code: "000099995B57",
		},
		{
			name: "not normalized",
			code: "000099995b57",
			err:  "invalid company code",
		},
		{
			name: "no B",
			code: "000099998X57",
			err:  "invalid company code",
		},
		{
			name: "non numbers",
			code: "000099998B5a",
			err:  "invalid VAT number",
		},
		{
			name: "invalid checksum",
			code: "123456789B12",
			err:  "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &org.TaxIdentity{Country: l10n.NL, Code: tt.code}
			err := nl.ValidateTaxIdentity(tID)
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
