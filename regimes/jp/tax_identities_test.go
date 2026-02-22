package jp_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/jp"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaxIdentityNormalization(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   cbc.Code
		out  cbc.Code
	}{
		{
			name: "strips hyphens",
			in:   "T7000-0120-5000-2",
			out:  "T7000012050002",
		},
		{
			name: "strips spaces and uppercases",
			in:   " t7000 0120 5000 2 ",
			out:  "T7000012050002",
		},
		{
			name: "strips country prefix JP",
			in:   "JPT7000012050002",
			out:  "T7000012050002",
		},
		{
			name: "already clean T-number",
			in:   "T7000012050002",
			out:  "T7000012050002",
		},
		{
			name: "13-digit corporate number without T prefix",
			in:   "7000012050002",
			out:  "7000012050002",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			id := &tax.Identity{Country: "JP", Code: tt.in}
			jp.Normalize(id)
			assert.Equal(t, tt.out, id.Code)
		})
	}
}

func TestTaxIdentityNormalization_Nil(t *testing.T) {
	t.Parallel()

	// Should not panic on nil
	jp.Normalize((*tax.Identity)(nil))
}

func TestTaxIdentityValidation_Nil(t *testing.T) {
	t.Parallel()

	err := jp.Validate((*tax.Identity)(nil))
	assert.NoError(t, err)
}

func TestTaxIdentityValidation_EmptyCode(t *testing.T) {
	t.Parallel()

	id := &tax.Identity{Country: "JP", Code: ""}
	err := jp.Validate(id)
	assert.NoError(t, err, "empty code should pass tax identity validation")
}

func TestTaxIdentityValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		code    cbc.Code
		wantErr string
	}{
		{
			name: "valid T-number",
			code: "T7000012050002", // check digit verified
		},
		{
			name:    "missing T prefix",
			code:    "7000012050002",
			wantErr: "must be 'T' followed by 13 digits",
		},
		{
			name:    "too short",
			code:    "T123456789012",
			wantErr: "must be 'T' followed by 13 digits",
		},
		{
			name:    "first digit zero with invalid check digit",
			code:    "T0123456789012",
			wantErr: "invalid check digit",
		},
		{
			name:    "invalid check digit",
			code:    "T7000012050003",
			wantErr: "invalid check digit",
		},
		{
			name:    "lowercase t",
			code:    "t7000012050002",
			wantErr: "must be 'T' followed by 13 digits",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			id := &tax.Identity{Country: "JP", Code: tt.code}
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
