package pa_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/pa"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/require"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		// Valid: natural person (cédula)
		{name: "valid natural", code: "8-442-445-90"},
		{name: "valid natural short", code: "1-12-345-72"},
		{name: "valid natural max province", code: "10-1234-567890-41"},

		// Valid: foreigner (E)
		{name: "valid foreigner", code: "E-12-342-35"},
		{name: "valid foreigner short", code: "E-8-1234-62"},

		// Valid: naturalized (N)
		{name: "valid naturalized", code: "N-45-832-03"},

		// Valid: PE (Panameño en el Exterior)
		{name: "valid PE", code: "PE-10-442-98"},

		// Valid: AV (Antes de la Vigencia)
		{name: "valid AV", code: "1AV-432-658-15"},
		{name: "valid AV max province", code: "12AV-1-1-05"},

		// Valid: PI (Población Indígena)
		{name: "valid PI", code: "4PI-234-123-51"},

		// Valid: legal entity
		{name: "valid legal", code: "2486589-1-816994-62"},
		{name: "valid legal 2", code: "155596713-2-2015-59"},
		{name: "valid legal short", code: "1000-1-1-18"},

		// Valid: old-format legal entity (triggers DV cross-reference substitution)
		{name: "valid legal old format crossref 00", code: "100-2-3-06"},
		{name: "valid legal old format crossref 10", code: "10000-1-1-18"},
		{name: "valid legal old format crossref 11", code: "11000-3-5-02"},
		{name: "valid legal old format crossref 12", code: "12000-2-7-30"},
		{name: "valid legal old format crossref 13", code: "13000-4-9-54"},
		{name: "valid legal old format crossref 14", code: "14000-3-11-42"},
		{name: "valid legal old format crossref 16", code: "16000-2-13-88"},
		{name: "valid legal old format crossref 17", code: "17000-4-15-00"},
		{name: "valid legal old format crossref 18", code: "18000-3-17-38"},
		{name: "valid legal old format crossref 30", code: "30000-5-100-05"},
		{name: "valid legal old format default", code: "1000-2-19-73"},

		// Valid: NT (Número Tributario)
		{name: "valid NT", code: "1NT-100-200-3000-30"},
		{name: "valid NT min", code: "12NT-1-1-1-96"},

		// Valid: mod-11 remainder edge cases
		{name: "valid remainder 10", code: "8-1-10-15"},
		{name: "valid remainder 1", code: "E-1-1-08"},
		{name: "valid remainder 0", code: "8-1-100-06"},

		// Valid: special cases
		{name: "valid final consumer", code: "CIP-000-000-0000"},
		{name: "empty code", code: ""},

		// Invalid: format errors
		{
			name: "invalid no hyphens",
			code: "12345678",
			err:  "code",
		},
		{
			name: "invalid too few segments",
			code: "8-442",
			err:  "code",
		},
		{
			name: "invalid letters in entry",
			code: "8-442-ABC-90",
			err:  "code",
		},
		{
			name: "invalid unknown prefix",
			code: "X-12-342-35",
			err:  "code",
		},
		{
			name: "invalid province too long for natural",
			code: "123-442-445-90",
			err:  "code",
		},

		// Invalid: DV checksum errors
		{
			name: "invalid DV natural",
			code: "8-442-445-91",
			err:  "dv checksum failed",
		},
		{
			name: "invalid DV foreigner",
			code: "E-12-342-00",
			err:  "dv checksum failed",
		},
		{
			name: "invalid DV legal",
			code: "2486589-1-816994-99",
			err:  "dv checksum failed",
		},
		{
			name: "invalid DV legal old format",
			code: "10000-1-1-99",
			err:  "dv checksum failed",
		},
		{
			name: "invalid DV naturalized",
			code: "N-45-832-99",
			err:  "dv checksum failed",
		},
		{
			name: "invalid DV PE",
			code: "PE-10-442-00",
			err:  "dv checksum failed",
		},
		{
			name: "invalid DV AV",
			code: "1AV-432-658-00",
			err:  "dv checksum failed",
		},
		{
			name: "invalid DV PI",
			code: "4PI-234-123-00",
			err:  "dv checksum failed",
		},
		{
			name: "invalid DV NT",
			code: "1NT-100-200-3000-00",
			err:  "dv checksum failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "PA", Code: tt.code}
			err := pa.Validate(tID)
			if tt.err == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tt.err)
		})
	}
}

func TestValidateTaxIdentityNil(t *testing.T) {
	require.NotPanics(t, func() {
		err := pa.Validate((*tax.Identity)(nil))
		require.NoError(t, err)
	})
}

func TestValidateUnsupportedDocType(t *testing.T) {
	err := pa.Validate("unsupported")
	require.NoError(t, err)
}

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		name     string
		code     cbc.Code
		expected cbc.Code
	}{
		{
			name:     "already normalized",
			code:     "8-442-445-90",
			expected: "8-442-445-90",
		},
		{
			name:     "with PA prefix",
			code:     "PA8-442-445-90",
			expected: "8-442-445-90",
		},
		{
			name:     "lowercase",
			code:     "e-12-342-35",
			expected: "E-12-342-35",
		},
		{
			name:     "lowercase with PA prefix",
			code:     "pae-12-342-35",
			expected: "E-12-342-35",
		},
		{
			name:     "with spaces",
			code:     "8 -442- 445-90",
			expected: "8-442-445-90",
		},
		{
			name:     "with dots and special chars",
			code:     "8.442.445.90",
			expected: "844244590",
		},
		{
			name:     "preserves hyphens",
			code:     "2486589-1-816994-62",
			expected: "2486589-1-816994-62",
		},
		{
			name:     "final consumer",
			code:     "cip-000-000-0000",
			expected: "CIP-000-000-0000",
		},
		{
			name:     "empty",
			code:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "PA", Code: tt.code}
			pa.Normalize(tID)
			require.Equal(t, tt.expected, tID.Code)
		})
	}
}

func TestNormalizeTaxIdentityNil(t *testing.T) {
	require.NotPanics(t, func() {
		pa.Normalize((*tax.Identity)(nil))
	})
}
