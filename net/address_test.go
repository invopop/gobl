package net

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAddress(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Address
		wantErr error
	}{
		{
			name:  "valid FQDN",
			input: "billing.invopop.com",
			want:  Address("billing.invopop.com"),
		},
		{
			name:  "valid subdomain",
			input: "sub.domain.example.org",
			want:  Address("sub.domain.example.org"),
		},
		{
			name:  "uppercase normalized",
			input: "Billing.Invopop.COM",
			want:  Address("billing.invopop.com"),
		},
		{
			name:  "trailing dot stripped",
			input: "billing.invopop.com.",
			want:  Address("billing.invopop.com"),
		},
		{
			name:  "whitespace trimmed",
			input: "  billing.invopop.com  ",
			want:  Address("billing.invopop.com"),
		},
		{
			name:    "empty",
			input:   "",
			wantErr: ErrAddressEmpty,
		},
		{
			name:    "single label",
			input:   "localhost",
			wantErr: ErrAddressInvalid,
		},
		{
			name:    "has scheme",
			input:   "http://example.com",
			wantErr: ErrAddressInvalid,
		},
		{
			name:    "has path",
			input:   "example.com/path",
			wantErr: ErrAddressInvalid,
		},
		{
			name:    "has port",
			input:   "example.com:8080",
			wantErr: ErrAddressInvalid,
		},
		{
			name:    "invalid characters",
			input:   "not valid!.com",
			wantErr: ErrAddressInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAddress(tt.input)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr), "expected %v, got %v", tt.wantErr, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAddressJWKSURL(t *testing.T) {
	a := Address("billing.invopop.com")
	assert.Equal(t, "https://billing.invopop.com/.well-known/gobl/jwks.json", a.JWKSURL())
}

func TestAddressTopic(t *testing.T) {
	tests := []struct {
		addr Address
		want string
	}{
		{Address("billing.invopop.com"), "com.invopop.billing"},
		{Address("sub.domain.example.org"), "org.example.domain.sub"},
		{Address("example.com"), "com.example"},
	}
	for _, tt := range tests {
		t.Run(string(tt.addr), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.addr.Topic())
		})
	}
}
