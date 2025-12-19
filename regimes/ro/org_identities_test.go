package ro_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/ro"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeOrgIdentity(t *testing.T) {
	tests := []struct {
		name     string
		identity *org.Identity
		expected cbc.Code
	}{
		{
			name: "CUI with RO prefix",
			identity: &org.Identity{
				Type: ro.IdentityTypeCUI,
				Code: "RO18547290",
			},
			expected: "18547290",
		},
		{
			name: "CUI with spaces",
			identity: &org.Identity{
				Type: ro.IdentityTypeCUI,
				Code: "1854 7290",
			},
			expected: "18547290",
		},
		{
			name: "CUI with dashes",
			identity: &org.Identity{
				Type: ro.IdentityTypeCUI,
				Code: "1854-7290",
			},
			expected: "18547290",
		},
		{
			name: "CUI lowercase with RO",
			identity: &org.Identity{
				Type: ro.IdentityTypeCUI,
				Code: "ro18547290",
			},
			expected: "18547290",
		},
		{
			name: "CNP with spaces",
			identity: &org.Identity{
				Type: ro.IdentityTypeCNP,
				Code: "5000101 010003",
			},
			expected: "5000101010003",
		},
		{
			name: "CNP with dashes",
			identity: &org.Identity{
				Type: ro.IdentityTypeCNP,
				Code: "5000101-010003",
			},
			expected: "5000101010003",
		},
		{
			name:     "nil identity",
			identity: nil,
			expected: "",
		},
		{
			name: "empty code",
			identity: &org.Identity{
				Type: ro.IdentityTypeCUI,
				Code: "",
			},
			expected: "",
		},
		{
			name: "unknown identity type",
			identity: &org.Identity{
				Type: "UNKNOWN",
				Code: "12345",
			},
			expected: "12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.identity == nil {
				ro.Normalize(tt.identity)
				// Should not panic
				return
			}
			ro.Normalize(tt.identity)
			assert.Equal(t, tt.expected, tt.identity.Code)
		})
	}
}

func TestValidateOrgIdentity(t *testing.T) {
	tests := []struct {
		name     string
		identity *org.Identity
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid CUI",
			identity: &org.Identity{
				Type: ro.IdentityTypeCUI,
				Code: "18547290",
			},
			wantErr: false,
		},
		{
			name: "invalid CUI checksum",
			identity: &org.Identity{
				Type: ro.IdentityTypeCUI,
				Code: "18547291",
			},
			wantErr: true,
			errMsg:  "invalid checksum",
		},
		{
			name: "valid CNP",
			// A valid dummy CNP: Male (5), Born 2000-01-01, County 01, Seq 000, Control 3
			identity: &org.Identity{
				Type: ro.IdentityTypeCNP,
				Code: "5000101010003",
			},
			wantErr: false,
		},
		{
			name: "invalid CNP checksum",
			identity: &org.Identity{
				Type: ro.IdentityTypeCNP,
				Code: "5000101010002", // Wrong control digit
			},
			wantErr: true,
			errMsg:  "invalid CNP checksum",
		},
		{
			name: "invalid CNP - too short",
			identity: &org.Identity{
				Type: ro.IdentityTypeCNP,
				Code: "500010101000",
			},
			wantErr: true,
			errMsg:  "13 digits",
		},
		{
			name: "invalid CNP - too long",
			identity: &org.Identity{
				Type: ro.IdentityTypeCNP,
				Code: "50001010100012",
			},
			wantErr: true,
			errMsg:  "13 digits",
		},
		{
			name: "invalid CNP - contains letters",
			identity: &org.Identity{
				Type: ro.IdentityTypeCNP,
				Code: "500010101000A",
			},
			wantErr: true,
			errMsg:  "only digits",
		},
		{
			name: "invalid CNP - invalid first digit",
			identity: &org.Identity{
				Type: ro.IdentityTypeCNP,
				Code: "0000101010001", // 0 is invalid
			},
			wantErr: true,
			errMsg:  "invalid CNP first digit",
		},
		{
			name: "valid CNP - foreign resident (starts with 9)",
			identity: &org.Identity{
				Type: ro.IdentityTypeCNP,
				// 9 + 000101 + 01 + 000 + checksum (0)
				// Foreign resident CNP starting with 9
				Code: "9000101010000",
			},
			wantErr: false,
		},
		{
			name: "empty code",
			identity: &org.Identity{
				Type: ro.IdentityTypeCUI,
				Code: "",
			},
			wantErr: true,
			errMsg:  "cannot be blank",
		},
		{
			name: "unknown identity type",
			identity: &org.Identity{
				Type: "UNKNOWN",
				Code: "12345",
			},
			// Unknown types pass through without validation
			// NOTE: Should this be the case, or do we want to panic or error here to identify when a new type has been added
			// and hasn't been implemented?
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ro.Validate(tt.identity)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOrgIdentityTypes(t *testing.T) {
	// Test that identity type definitions are properly set
	regime := ro.New()
	require.NotNil(t, regime)
	require.NotEmpty(t, regime.Identities)

	// Check CUI definition
	var foundCUI bool
	for _, def := range regime.Identities {
		if def.Code == ro.IdentityTypeCUI {
			foundCUI = true
			assert.NotEmpty(t, def.Name[i18n.EN])
			assert.NotEmpty(t, def.Desc[i18n.EN])
		}
	}
	assert.True(t, foundCUI, "CUI identity type should be defined")

	// Check CNP definition
	var foundCNP bool
	for _, def := range regime.Identities {
		if def.Code == ro.IdentityTypeCNP {
			foundCNP = true
			assert.NotEmpty(t, def.Name[i18n.EN])
			assert.NotEmpty(t, def.Desc[i18n.EN])
		}
	}
	assert.True(t, foundCNP, "CNP identity type should be defined")
}
