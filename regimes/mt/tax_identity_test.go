package mt_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/regimes/mt"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	regime := mt.New()
	tests := []struct {
		name string
		code cbc.Code
		want string
	}{
		{
			name: "removes separators",
			code: "1167-9112",
			want: "11679112",
		},
		{
			name: "removes the country prefix and whitespace",
			code: "MT 1167 9112",
			want: "11679112",
		},
		{
			name: "uppercases before stripping the prefix",
			code: "mt11679112",
			want: "11679112",
		},
		{
			name: "leaves an already normalized code alone",
			code: "11679112",
			want: "11679112",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{
				Country: regime.Country,
				Code:    tt.code,
			}
			norm.Normalize(tID)
			assert.Equal(t, tt.want, tID.Code.String())

			// Normalizers must be idempotent.
			norm.Normalize(tID)
			assert.Equal(t, tt.want, tID.Code.String())
		})
	}
}

func TestNormalizeTaxIdentityNil(t *testing.T) {
	var tID *tax.Identity
	assert.NotPanics(t, func() { norm.Normalize(tID) })
}

func TestTaxIdentityRules(t *testing.T) {
	// The check digits are the last two of the eight, computed over the first six
	// with weights [3, 4, 6, 7, 8, 9] as check = 37 - (sum mod 37):
	//
	//   116791 -> 3+4+36+49+72+9    = 173, 173 mod 37 = 25, 37-25 = 12
	//   100000 -> 3                 =   3,   3 mod 37 =  3, 37- 3 = 34
	//   123456 -> 3+8+18+28+40+54   = 151, 151 mod 37 =  3, 37- 3 = 34
	//   987654 -> 27+32+42+42+40+36 = 219, 219 mod 37 = 34, 37-34 =  3
	//   500000 -> 15                =  15,  15 mod 37 = 15, 37-15 = 22
	tests := []struct {
		name    string
		code    cbc.Code
		wantErr bool
	}{
		{
			name:    "valid - 11679112, the published reference example",
			code:    "11679112",
			wantErr: false,
		},
		{
			name:    "valid - 10000034",
			code:    "10000034",
			wantErr: false,
		},
		{
			name:    "valid - 12345634",
			code:    "12345634",
			wantErr: false,
		},
		{
			name:    "valid - 98765403, check digits with a leading zero",
			code:    "98765403",
			wantErr: false,
		},
		{
			name:    "valid - 50000022",
			code:    "50000022",
			wantErr: false,
		},
		{
			// Edge case. Base 111111 sums to exactly 37, so 37 - 0 = 37 and the
			// check digits are "37". See the note in validateTaxCodeChecksum: this
			// branch is a deliberate choice between two undocumented variants.
			name:    "valid - 11111137, weighted sum is an exact multiple of 37",
			code:    "11111137",
			wantErr: false,
		},
		{
			// The other variant of the same edge case, which reduces modulo 37 a
			// second time and expects "00". Rejected here, deliberately and
			// consistently with the case above.
			name:    "invalid - 11111100, the rejected variant of the same edge case",
			code:    "11111100",
			wantErr: true,
		},
		{
			name:    "invalid - 11679113, correct base with a wrong check",
			code:    "11679113",
			wantErr: true,
		},
		{
			name:    "invalid - too short",
			code:    "1167911",
			wantErr: true,
		},
		{
			name:    "invalid - too long",
			code:    "116791123",
			wantErr: true,
		},
		{
			name:    "invalid - contains letters",
			code:    "1167911A",
			wantErr: true,
		},
		{
			name:    "invalid - not eight digits",
			code:    "12345",
			wantErr: true,
		},
		{
			// Base 116791 has check 12, as in the valid case above. This number
			// carries 49 instead, which is 12 + 37 and so congruent to it modulo
			// 37 — enough for python-stdnum, which only tests the congruence.
			// 37 - (sum mod 37) never exceeds 37, so it is rejected here.
			name:    "invalid - 11679149, check value above 37",
			code:    "11679149",
			wantErr: true,
		},
		{
			// Leading zeros are accepted, matching valvat's syntax rule for Malta
			// (/\AMT[0-9]{8}\Z/). python-stdnum rejects them; no official source
			// states a digit count at all, let alone a leading-digit rule.
			name:    "valid - 00000128, leading zeros",
			code:    "00000128",
			wantErr: false,
		},
		{
			name:    "empty code is left to the presence rules",
			code:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{
				Country: "MT",
				Code:    tt.code,
			}
			err := rules.Validate(tID)
			if tt.wantErr {
				if assert.Error(t, err, "expected error for code: %s", tt.code) {
					assert.Contains(t, err.Error(), "IDENTITY-01")
				}
			} else {
				require.NoError(t, err, "unexpected error for code: %s", tt.code)
			}
		})
	}
}
