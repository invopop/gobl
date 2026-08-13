package net

import (
	"errors"
	"testing"

	"github.com/invopop/gobl/cbc"

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
		{
			name:  "U-Label normalised to A-Label",
			input: "münchen.de",
			want:  Address("xn--mnchen-3ya.de"),
		},
		{
			name:  "A-Label accepted verbatim",
			input: "xn--mnchen-3ya.de",
			want:  Address("xn--mnchen-3ya.de"),
		},
		{
			name:  "U-Label with uppercase + trailing dot",
			input: "München.DE.",
			want:  Address("xn--mnchen-3ya.de"),
		},
		{
			name:    "invalid IDN label",
			input:   "-bad.com", // labels MUST NOT start with a hyphen
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

func TestAddressString(t *testing.T) {
	assert.Equal(t, "billing.invopop.com", Address("billing.invopop.com").String())
}

func TestAddressValidate(t *testing.T) {
	assert.NoError(t, Address("billing.invopop.com").Validate())
	assert.ErrorIs(t, Address("").Validate(), ErrAddressEmpty)
	assert.ErrorIs(t, Address("localhost").Validate(), ErrAddressInvalid)
}

func TestAddressKeyURL(t *testing.T) {
	a := Address("billing.invopop.com")
	assert.Equal(t,
		"https://billing.invopop.com/.well-known/gobl/keys/key-1",
		a.KeyURL("key-1"),
	)
}

func TestKeyPath(t *testing.T) {
	assert.Equal(t, "/.well-known/gobl/keys/abc", KeyPath("abc"))
}

func TestAddressWhoURL(t *testing.T) {
	a := Address("billing.invopop.com")
	assert.Equal(t, "https://billing.invopop.com/.well-known/gobl/who", a.WhoURL())
}

func TestAddressInboxURL(t *testing.T) {
	a := Address("billing.invopop.com")
	assert.Equal(t, "https://billing.invopop.com/.well-known/gobl/inbox", a.InboxURL())
}

func TestAddressJWKSURL(t *testing.T) {
	a := Address("billing.invopop.com")
	assert.Equal(t, "https://billing.invopop.com/.well-known/jwks.json", a.JWKSURL())
}

func TestAddressURI(t *testing.T) {
	// The gobl: URI form is for endpoint lists and header routing;
	// signed claims carry the bare address.
	a := Address("billing.invopop.com")
	assert.Equal(t, cbc.URI("gobl:billing.invopop.com"), a.URI())
}
