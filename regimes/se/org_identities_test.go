package se_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/se"
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
		{name: "Organization number already normalized", typeCode: se.IdentityTypeOrgNr, input: "5560360793", expected: "5560360793"},
		{name: "Organization number with hyphens", typeCode: se.IdentityTypeOrgNr, input: "556036-0793", expected: "5560360793"},
		{name: "Organization number with spaces", typeCode: se.IdentityTypeOrgNr, input: "  556036 0793  ", expected: "5560360793"},
		{name: "Organization number with check digit", typeCode: se.IdentityTypeOrgNr, input: "5560360793-01", expected: "5560360793"},

		{name: "Person number already normalized", typeCode: se.IdentityTypePersonNr, input: "800101-0017", expected: "800101-0017"},
		{name: "Person number with hyphen", typeCode: se.IdentityTypePersonNr, input: "800101-0017", expected: "800101-0017"},
		{name: "Person number with plus sign", typeCode: se.IdentityTypePersonNr, input: "800101+0017", expected: "800101+0017"},
		{name: "Person number with spaces", typeCode: se.IdentityTypePersonNr, input: "  800101-0017  ", expected: "800101-0017"},
		{name: "Person number without hyphen or plus sign", typeCode: se.IdentityTypePersonNr, input: "8001010017", expected: "800101-0017"},

		{name: "Coordination number already normalized", typeCode: se.IdentityTypeCoordinationNr, input: "800161-0017", expected: "800161-0017"},
		{name: "Coordination number with hyphen", typeCode: se.IdentityTypeCoordinationNr, input: "800161-0017", expected: "800161-0017"},
		{name: "Coordination number with plus sign", typeCode: se.IdentityTypeCoordinationNr, input: "800161+0017", expected: "800161+0017"},
		{name: "Coordination number with spaces", typeCode: se.IdentityTypeCoordinationNr, input: "  800161-0017  ", expected: "800161-0017"},
		{name: "Coordination number without hyphen or plus sign", typeCode: se.IdentityTypeCoordinationNr, input: "8001610017", expected: "800161-0017"},

		{name: "Unknown key", typeCode: "unknown", input: "1234567890", expected: "1234567890"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			id := &org.Identity{Type: tt.typeCode, Code: tt.input}
			se.Normalize(id)
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
		// Valid cases
		{name: "Valid organization number", typeCode: se.IdentityTypeOrgNr, input: "5560360793", err: ""},
		{name: "Valid person number", typeCode: se.IdentityTypePersonNr, input: "800101-0019", err: ""},
		{name: "Valid coordination number", typeCode: se.IdentityTypeCoordinationNr, input: "800161-0016", err: ""},
		{name: "Valid person number with plus sign", typeCode: se.IdentityTypePersonNr, input: "800101+0019", err: ""},
		{name: "Unknown identity type", typeCode: "unknown", input: "1234567890", err: ""},

		// Invalid cases
		{name: "Organization number too short", typeCode: se.IdentityTypeOrgNr, input: "123456789", err: se.ErrInvalidOrgNumberFormat.Error()},
		{name: "Organization number too long", typeCode: se.IdentityTypeOrgNr, input: "12345678901", err: se.ErrInvalidOrgNumberFormat.Error()},
		{name: "Organization number with letters", typeCode: se.IdentityTypeOrgNr, input: "123456789A", err: se.ErrInvalidOrgNumberFormat.Error()},
		{name: "Organization number with invalid check digit", typeCode: se.IdentityTypeOrgNr, input: "5560360794", err: se.ErrInvalidChecksum.Error()},

		{name: "Person number without separator", typeCode: se.IdentityTypePersonNr, input: "8001010019", err: se.ErrInvalidIndividualNumberFormat.Error()},
		{name: "Person number too short", typeCode: se.IdentityTypePersonNr, input: "800101-001", err: se.ErrInvalidIndividualNumberFormat.Error()},
		{name: "Person number too long", typeCode: se.IdentityTypePersonNr, input: "800101-00177", err: se.ErrInvalidIndividualNumberFormat.Error()},
		{name: "Person number with letters", typeCode: se.IdentityTypePersonNr, input: "800101-001A", err: se.ErrInvalidIndividualNumberFormat.Error()},
		{name: "Person number with invalid check digit", typeCode: se.IdentityTypePersonNr, input: "800101-0018", err: se.ErrInvalidChecksum.Error()},

		{name: "Coordination number without separator", typeCode: se.IdentityTypeCoordinationNr, input: "8001610017", err: se.ErrInvalidIndividualNumberFormat.Error()},
		{name: "Coordination number too short", typeCode: se.IdentityTypeCoordinationNr, input: "800161-001", err: se.ErrInvalidIndividualNumberFormat.Error()},
		{name: "Coordination number too long", typeCode: se.IdentityTypeCoordinationNr, input: "800161-00177", err: se.ErrInvalidIndividualNumberFormat.Error()},
		{name: "Coordination number with letters", typeCode: se.IdentityTypeCoordinationNr, input: "800161-001A", err: se.ErrInvalidIndividualNumberFormat.Error()},
		{name: "Coordination number with invalid check digit", typeCode: se.IdentityTypeCoordinationNr, input: "800161-0018", err: se.ErrInvalidChecksum.Error()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			id := &org.Identity{Type: tt.typeCode, Code: tt.input, Label: tt.name, Country: l10n.SE.ISO()}
			err := se.Validate(id)

			if tt.err != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
