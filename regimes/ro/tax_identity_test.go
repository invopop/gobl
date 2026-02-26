package ro_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/ro"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "good 1", code: "4221306"},
		{name: "good 2", code: "22891860"},
		{name: "good 3", code: "16109528"},
		{name: "good 4", code: "37616299"},
		{name: "good 5", code: "14399840"},
		{name: "good 6", code: "18547290"},
		{name: "good 7", code: "38144933"},
		{name: "good 8", code: "22252394"},
		{name: "good 9", code: "40188877"},
		{name: "good 10", code: "40425604"},
		{name: "good 11", code: "32988399"},
		{name: "good 12", code: "20959993"},
		{name: "good 13", code: "18834610"},
		{name: "good 14", code: "14752305"},
		{name: "good 15", code: "21237710"},
		{name: "good 16", code: "17732264"},
		{name: "empty", code: ""},
		{
			name: "too short",
			code: "1",
			err:  "invalid format",
		},
		{
			name: "too long",
			code: "12345678901",
			err:  "invalid format",
		},
		{
			name: "contains letters",
			code: "1234567A",
			err:  "invalid format",
		},
		{
			name: "cnp rejected",
			code: "1931113013515",
			err:  "invalid format",
		},
		{
			name: "bad checksum",
			code: "14399841",
			err:  "checksum mismatch",
		},
		{
			name: "bad checksum short",
			code: "13181",
			err:  "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "RO", Code: tt.code}
			err := ro.Validate(tID)
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

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		exp  cbc.Code
	}{
		{name: "already clean", code: "14399840", exp: "14399840"},
		{name: "strip uppercase prefix", code: "RO4221306", exp: "4221306"},
		{name: "strip lowercase prefix", code: "ro4221306", exp: "4221306"},
		{name: "strip mixed case prefix", code: "Ro14399840", exp: "14399840"},
		{name: "leading whitespace", code: " 14399840", exp: "14399840"},
		{name: "trailing whitespace", code: "14399840 ", exp: "14399840"},
		{name: "whitespace within", code: "143 998 40", exp: "14399840"},
		{name: "hyphens", code: "143-998-40", exp: "14399840"},
		{name: "dots", code: "143.998.40", exp: "14399840"},
		{name: "slashes", code: "143/998/40", exp: "14399840"},
		{name: "mixed separators", code: "RO 143-998.40", exp: "14399840"},
		{name: "empty code", code: "", exp: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "RO", Code: tt.code}
			ro.Normalize(tID)
			assert.Equal(t, tt.exp, tID.Code)
		})
	}
}

func TestNormalizeThenValidate(t *testing.T) {
	t.Run("dirty code normalizes to valid", func(t *testing.T) {
		tID := &tax.Identity{Country: "RO", Code: "RO 143-998.40"}
		ro.Normalize(tID)
		assert.Equal(t, cbc.Code("14399840"), tID.Code)
		assert.NoError(t, ro.Validate(tID))
	})

	t.Run("lowercase prefix normalizes to valid", func(t *testing.T) {
		tID := &tax.Identity{Country: "RO", Code: "ro18547290"}
		ro.Normalize(tID)
		assert.Equal(t, cbc.Code("18547290"), tID.Code)
		assert.NoError(t, ro.Validate(tID))
	})

	t.Run("whitespace-padded normalizes to valid", func(t *testing.T) {
		tID := &tax.Identity{Country: "RO", Code: " 4221306 "}
		ro.Normalize(tID)
		assert.Equal(t, cbc.Code("4221306"), tID.Code)
		assert.NoError(t, ro.Validate(tID))
	})
}
