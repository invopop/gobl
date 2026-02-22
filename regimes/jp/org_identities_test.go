package jp_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/jp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrgIdentityValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		code    cbc.Code
		wantErr string
	}{
		{
			name: "valid corporate number",
			code: "7000012050002",
		},
		{
			name:    "too short",
			code:    "123456789012",
			wantErr: "must be exactly 13 digits",
		},
		{
			name:    "too long",
			code:    "12345678901234",
			wantErr: "must be exactly 13 digits",
		},
		{
			name:    "non-numeric",
			code:    "123456789012X",
			wantErr: "must be exactly 13 digits",
		},
		{
			name:    "invalid check digit",
			code:    "7000012050003", // correct is 7000012050002
			wantErr: "invalid check digit",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			id := &org.Identity{Type: jp.IdentityTypeCorporateNumber, Code: tt.code}
			err := jp.Validate(id)
			if tt.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}

func TestOrgIdentityValidation_QualifiedInvoiceIssuer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		code    cbc.Code
		wantErr string
	}{
		{
			name: "valid QII number",
			code: "T7000012050002",
		},
		{
			name: "valid QII number with check digit zero",
			code: "T0000000000000",
		},
		{
			name:    "missing T prefix",
			code:    "7000012050002",
			wantErr: "invalid Qualified Invoice Issuer format",
		},
		{
			name:    "too short",
			code:    "T123456789012",
			wantErr: "invalid Qualified Invoice Issuer format",
		},
		{
			name:    "too long",
			code:    "T12345678901234",
			wantErr: "invalid Qualified Invoice Issuer format",
		},
		{
			name:    "non-numeric after T",
			code:    "T123456789012X",
			wantErr: "invalid Qualified Invoice Issuer format",
		},
		{
			name:    "invalid check digit",
			code:    "T7000012050003",
			wantErr: "invalid check digit",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			id := &org.Identity{Type: jp.IdentityTypeQualifiedInvoiceIssuer, Code: tt.code}
			err := jp.Validate(id)
			if tt.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}

func TestOrgIdentityValidation_QII_EmptyCode(t *testing.T) {
	t.Parallel()

	id := &org.Identity{Type: jp.IdentityTypeQualifiedInvoiceIssuer, Code: ""}
	err := jp.Validate(id)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "code")
}

func TestOrgIdentityValidation_MyNumber(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		code    cbc.Code
		wantErr string
	}{
		{
			name: "valid 12-digit my number",
			code: "123456789012",
		},
		{
			name:    "too short",
			code:    "12345678901",
			wantErr: "invalid My Number format",
		},
		{
			name:    "too long",
			code:    "1234567890123",
			wantErr: "invalid My Number format",
		},
		{
			name:    "non-numeric",
			code:    "12345678901X",
			wantErr: "invalid My Number format",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			id := &org.Identity{Type: jp.IdentityTypeMyNumber, Code: tt.code}
			err := jp.Validate(id)
			if tt.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}

func TestOrgIdentityValidation_MyNumber_EmptyCode(t *testing.T) {
	t.Parallel()

	id := &org.Identity{Type: jp.IdentityTypeMyNumber, Code: ""}
	err := jp.Validate(id)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "code")
}

func TestOrgIdentityValidation_UnknownType(t *testing.T) {
	t.Parallel()

	// Unknown type should pass validation (no rules to apply)
	id := &org.Identity{Type: "UNKNOWN", Code: "anything"}
	err := jp.Validate(id)
	assert.NoError(t, err)
}

func TestOrgIdentityValidation_NilIdentity(t *testing.T) {
	t.Parallel()

	err := jp.Validate((*org.Identity)(nil))
	assert.NoError(t, err)
}

func TestOrgIdentityValidation_CorporateNumber_EmptyCode(t *testing.T) {
	t.Parallel()

	id := &org.Identity{Type: jp.IdentityTypeCorporateNumber, Code: ""}
	err := jp.Validate(id)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "code")
}

func TestOrgIdentityNormalization_NilIdentity(t *testing.T) {
	t.Parallel()

	// Should not panic on nil
	jp.Normalize((*org.Identity)(nil))
}

func TestOrgIdentityNormalization_StripCountryPrefix(t *testing.T) {
	t.Parallel()

	id := &org.Identity{Type: jp.IdentityTypeCorporateNumber, Code: "JP7000012050002"}
	jp.Normalize(id)
	assert.Equal(t, cbc.Code("7000012050002"), id.Code)
}

func TestOrgIdentityNormalization_LowercasePrefix(t *testing.T) {
	id := &org.Identity{Type: jp.IdentityTypeCorporateNumber, Code: "jp7000012050002"}
	jp.Normalize(id)
	assert.Equal(t, cbc.Code("7000012050002"), id.Code)
}

func TestOrgIdentityNormalization_TrimSpaces(t *testing.T) {
	id := &org.Identity{Type: jp.IdentityTypeCorporateNumber, Code: "  7000012050002  "}
	jp.Normalize(id)
	assert.Equal(t, cbc.Code("7000012050002"), id.Code)
}

func TestOrgIdentityNormalization_RemovesGarbageChars(t *testing.T) {
	id := &org.Identity{Type: jp.IdentityTypeCorporateNumber, Code: "70.00/01*205_0002"}
	jp.Normalize(id)
	assert.Equal(t, cbc.Code("7000012050002"), id.Code)
}

func TestOrgIdentityNormalization_QII(t *testing.T) {
	id := &org.Identity{Type: jp.IdentityTypeQualifiedInvoiceIssuer, Code: " jp-t7000012050002 "}
	jp.Normalize(id)
	assert.Equal(t, cbc.Code("T7000012050002"), id.Code)
}

func TestOrgIdentityNormalization(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   cbc.Code
		out  cbc.Code
	}{
		{
			name: "strips hyphens",
			in:   "7000-0120-5000-2",
			out:  "7000012050002",
		},
		{
			name: "strips spaces",
			in:   "7000 0120 5000 2",
			out:  "7000012050002",
		},
		{
			name: "already clean",
			in:   "7000012050002",
			out:  "7000012050002",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			id := &org.Identity{Type: jp.IdentityTypeCorporateNumber, Code: tt.in}
			jp.Normalize(id)
			assert.Equal(t, tt.out, id.Code)
		})
	}
}

func TestCorporateNumber_CheckDigitZeroCase(t *testing.T) {
	id := &org.Identity{
		Type: jp.IdentityTypeCorporateNumber,
		Code: "0000000000000", // check digit is 0
	}
	err := jp.Validate(id)
	require.NoError(t, err)
}
