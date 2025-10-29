package ie_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/ie"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	regime := ie.New()
	tID := &tax.Identity{
		Country: regime.Country,
		Code:    "36.28-739L",
	}
	regime.Normalizer(tID)
	// Normalization removes separator characters like dots and dashes from the tax identity code,
	// as implemented by tax.NormalizeIdentity().
	assert.Equal(t, "3628739L", tID.Code.String())
}

func TestValidateTaxIdentity(t *testing.T) {
	regime := ie.New()

	tests := []struct {
		name    string
		code    cbc.Code
		wantErr bool
	}{
		{
			name:    "valid new format - 3628739L",
			code:    "3628739L",
			wantErr: false,
		},
		{
			name:    "valid new format - 5343381W",
			code:    "5343381W",
			wantErr: false,
		},
		{
			name:    "valid new format with two letters ending in A - 6433435OA",
			code:    "6433435OA",
			wantErr: false,
		},
		{
			name:    "valid new format with two letters ending in H - 3336483DH",
			code:    "3336483DH",
			wantErr: false,
		},
		{
			name:    "valid old format - 8Z49289F",
			code:    "8Z49289F",
			wantErr: false,
		},
		{
			name:    "invalid checksum new format",
			code:    "3628739A",
			wantErr: true,
		},
		{
			name:    "invalid checksum old format",
			code:    "8Z49289A",
			wantErr: true,
		},
		{
			name:    "invalid format - too short",
			code:    "123456",
			wantErr: true,
		},
		{
			name:    "invalid format - no letters",
			code:    "12345678",
			wantErr: true,
		},
		{
			name:    "empty code",
			code:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{
				Country: regime.Country,
				Code:    tt.code,
			}
			err := regime.Validator(tID)
			if tt.wantErr {
				assert.Error(t, err, "expected error for code: %s", tt.code)
			} else {
				require.NoError(t, err, "unexpected error for code: %s", tt.code)
			}
		})
	}
}
