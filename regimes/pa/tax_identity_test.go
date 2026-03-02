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
		// Verified against juancorradine/Panama-RUC-DV-Calculator (tested with DGI ETax 2.0)
		{name: "valid natural 8-769-1080", code: "8-769-1080-56"},
		{name: "valid natural 5-257-218", code: "5-257-218-09"},
		{name: "valid natural 6-108-289", code: "6-108-289-79"},
		{name: "valid natural 8-28-1284", code: "8-28-1284-33"},
		{name: "valid natural 2-7-89", code: "2-7-89-20"},
		{name: "valid natural 1-184-921", code: "1-184-921-49"},
		{name: "valid natural max province", code: "10-1234-567890-22"},

		// Valid: foreigner (E)
		{name: "valid foreigner", code: "E-12-342-10"},
		{name: "valid foreigner short", code: "E-8-1234-97"},

		// Valid: naturalized (N)
		{name: "valid naturalized", code: "N-45-832-58"},

		// Valid: PE (Panameño en el Exterior)
		{name: "valid PE", code: "PE-10-442-50"},

		// Valid: AV (Antes de la Vigencia)
		{name: "valid AV", code: "1AV-432-658-31"},
		{name: "valid AV max province", code: "12AV-1-1-05"},

		// Valid: PI (Población Indígena)
		{name: "valid PI", code: "4PI-234-123-31"},

		// Valid: legal entity
		// Verified against juancorradine/Panama-RUC-DV-Calculator (tested with DGI ETax 2.0)
		{name: "valid legal 155720753-2-2022", code: "155720753-2-2022-39"},
		{name: "valid legal 2588017-1-831938", code: "2588017-1-831938-20"},
		{name: "valid legal 1489806-1-645353", code: "1489806-1-645353-68"},
		{name: "valid legal 1956569-1-732877", code: "1956569-1-732877-00"},
		{name: "valid legal 797609-1-493865", code: "797609-1-493865-12"},
		{name: "valid legal 15565624-2-2017", code: "15565624-2-2017-63"},
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
		// Verified against juancorradine/Panama-RUC-DV-Calculator (tested with DGI ETax 2.0)
		{name: "valid NT 8NT-0-1-24", code: "8NT-0-1-24-33"},
		{name: "valid NT 3NT-0-465-45624", code: "3NT-0-465-45624-03"},
		{name: "valid NT 9NT-0-2-421578", code: "9NT-0-2-421578-50"},
		{name: "valid NT 6NT-0-227-888555", code: "6NT-0-227-888555-09"},
		{name: "valid NT 12NT-0-45-2154", code: "12NT-0-45-2154-17"},

		// Valid: mod-11 remainder edge cases
		{name: "valid remainder 10", code: "8-1-10-07"},
		{name: "valid remainder 1", code: "E-1-1-91"},
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
			name: "invalid first segment too long",
			code: "12345678901-442-445-90",
			err:  "code",
		},
		{
			name: "invalid format does not leak internal error",
			code: "ABC-DEF-GHI-99",
			err:  "code",
		},

		// Invalid: DV checksum errors
		{
			name: "invalid DV natural",
			code: "8-769-1080-99",
			err:  "dv checksum failed",
		},
		{
			name: "invalid DV foreigner",
			code: "E-12-342-99",
			err:  "dv checksum failed",
		},
		{
			name: "invalid DV legal",
			code: "2588017-1-831938-99",
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
			code: "PE-10-442-99",
			err:  "dv checksum failed",
		},
		{
			name: "invalid DV AV",
			code: "1AV-432-658-99",
			err:  "dv checksum failed",
		},
		{
			name: "invalid DV PI",
			code: "4PI-234-123-99",
			err:  "dv checksum failed",
		},
		{
			name: "invalid DV NT",
			code: "8NT-0-1-24-99",
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
			code:     "8-769-1080-56",
			expected: "8-769-1080-56",
		},
		{
			name:     "with PA prefix",
			code:     "PA8-769-1080-56",
			expected: "8-769-1080-56",
		},
		{
			name:     "with PA prefix and hyphen",
			code:     "PA-8-769-1080-56",
			expected: "8-769-1080-56",
		},
		{
			name:     "lowercase",
			code:     "e-12-342-10",
			expected: "E-12-342-10",
		},
		{
			name:     "lowercase with PA prefix",
			code:     "pae-12-342-10",
			expected: "E-12-342-10",
		},
		{
			name:     "with spaces",
			code:     "8 -769- 1080-56",
			expected: "8-769-1080-56",
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
