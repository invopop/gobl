package is_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/is"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestOrgIdentityNormalize(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		typeCode cbc.Code
		input    cbc.Code
		expected cbc.Code
	}{
		{name: "already normalized", typeCode: is.IdentityTypeKennitala, input: "5012031220", expected: "5012031220"},
		{name: "with hyphen", typeCode: is.IdentityTypeKennitala, input: "090286-2349", expected: "0902862349"},
		{name: "with spaces", typeCode: is.IdentityTypeKennitala, input: "  5012031220  ", expected: "5012031220"},
		{name: "non-numeric left untouched", typeCode: is.IdentityTypeKennitala, input: "ABCDEFGHIJK", expected: "ABCDEFGHIJK"},
		{name: "empty left untouched", typeCode: is.IdentityTypeKennitala, input: "", expected: ""},
		{name: "other type left untouched", typeCode: "unknown", input: "090286-2349", expected: "090286-2349"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			id := &org.Identity{Type: tt.typeCode, Code: tt.input}
			norm.Normalize(id, tax.RegimeContext("IS"))
			assert.Equal(t, tt.expected, id.Code)
		})
	}
}

func TestOrgIdentityValidate(t *testing.T) {
	tests := []struct {
		name     string
		typeCode cbc.Code
		input    cbc.Code
		err      string
	}{
		// Valid: person, company, and temporary kennitalas are all valid national IDs.
		{name: "valid company kennitala", typeCode: is.IdentityTypeKennitala, input: "5012031220"},
		{name: "valid person kennitala", typeCode: is.IdentityTypeKennitala, input: "0902862349"},
		{name: "valid temporary kennitala", typeCode: is.IdentityTypeKennitala, input: "8101850150"},
		{name: "unknown identity type", typeCode: "unknown", input: "1234567890"},

		// Invalid.
		{name: "too short", typeCode: is.IdentityTypeKennitala, input: "501203122", err: "[GOBL-IS-ORG-IDENTITY-01]"},
		{name: "too long", typeCode: is.IdentityTypeKennitala, input: "50120312201", err: "[GOBL-IS-ORG-IDENTITY-01]"},
		{name: "with letters", typeCode: is.IdentityTypeKennitala, input: "501203122A", err: "[GOBL-IS-ORG-IDENTITY-01]"},
		{name: "invalid checksum", typeCode: is.IdentityTypeKennitala, input: "5012031230", err: "[GOBL-IS-ORG-IDENTITY-02]"},
	}

	opts := []rules.WithContext{
		tax.RegimeContext(is.CountryCode),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			id := &org.Identity{Type: tt.typeCode, Code: tt.input}
			err := rules.Validate(id, opts...)
			if tt.err == "" {
				assert.NoError(t, err)
			} else if assert.Error(t, err) {
				assert.Contains(t, err.Error(), tt.err)
			}
		})
	}
}
