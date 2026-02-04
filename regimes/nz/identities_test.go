package nz_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/nz"
	"github.com/stretchr/testify/assert"
)

func TestValidNZBNNumbers(t *testing.T) {
	r := nz.New()
	id := &org.Identity{
		Key:  org.IdentityKeyGLN,
		Code: "9429041234563",
	}
	r.Normalizer(id)
	err := r.Validator(id)
	assert.NoError(t, err)
}

func TestInvalidNZBNNumbers(t *testing.T) {
	tests := []struct {
		name string
		nzbn string
	}{
		{"WrongPrefix", "9329041234567"},
		{"TooShort", "942904123456"},
		{"TooLong", "94290412345638"},
		{"BadCheckDigit", "9429041234560"},
		{"Letters", "942904123456A"},
		{"SpecialChars", "942904!123456"},
	}

	r := nz.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{
				Key:  org.IdentityKeyGLN,
				Code: cbc.Code(tt.nzbn),
			}
			r.Normalizer(id)
			err := r.Validator(id)
			assert.Error(t, err, "NZBN %s should be invalid", tt.nzbn)
		})
	}
}

func TestNormalizeNZBN(t *testing.T) {
	r := nz.New()
	id := &org.Identity{
		Key:  org.IdentityKeyGLN,
		Code: "9429041234563",
	}
	r.Normalizer(id)
	assert.Equal(t, "9429041234563", id.Code.String())
}

func TestNormalizeNZBNWithDashes(t *testing.T) {
	r := nz.New()
	id := &org.Identity{
		Key:  org.IdentityKeyGLN,
		Code: "9429-0412-34563",
	}
	r.Normalizer(id)
	assert.Equal(t, "9429041234563", id.Code.String())
}

func TestNormalizeNZBNWithSpaces(t *testing.T) {
	r := nz.New()
	id := &org.Identity{
		Key:  org.IdentityKeyGLN,
		Code: "9429 0412 34563",
	}
	r.Normalizer(id)
	assert.Equal(t, "9429041234563", id.Code.String())
}

func TestValidateNZBNEmptyCode(t *testing.T) {
	r := nz.New()
	id := &org.Identity{
		Key:  org.IdentityKeyGLN,
		Code: "",
	}
	err := r.Validator(id)
	assert.Error(t, err, "GLN identity with empty code should fail")
}

func TestNonGLNIdentitySkipped(t *testing.T) {
	r := nz.New()
	id := &org.Identity{
		Key:  "other",
		Code: "invalid",
	}
	err := r.Validator(id)
	assert.NoError(t, err, "Non-GLN identities should be skipped")
}

// --- IRD as org.Identity tests ---

func TestValidIRDAsOrgIdentity(t *testing.T) {
	tests := []string{
		"49091850",
		"136410132",
	}

	r := nz.New()
	for _, ird := range tests {
		t.Run(ird, func(t *testing.T) {
			id := &org.Identity{
				Key:  nz.IdentityKeyIRD,
				Code: cbc.Code(ird),
			}
			r.Normalizer(id)
			err := r.Validator(id)
			assert.NoError(t, err, "IRD %s should be valid as org.Identity", ird)
		})
	}
}

func TestInvalidIRDAsOrgIdentity(t *testing.T) {
	tests := []struct {
		name string
		ird  string
	}{
		{"TooShort", "1234567"},
		{"TooLong", "1234567890"},
		{"BadCheckDigit", "49091851"},
	}

	r := nz.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{
				Key:  nz.IdentityKeyIRD,
				Code: cbc.Code(tt.ird),
			}
			r.Normalizer(id)
			err := r.Validator(id)
			assert.Error(t, err, "IRD %s should be invalid as org.Identity", tt.ird)
		})
	}
}

func TestNormalizeIRDAsOrgIdentity(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"WithDashes", "49-091-850", "49091850"},
		{"WithSpaces", "49 091 850", "49091850"},
		{"Clean", "49091850", "49091850"},
	}

	r := nz.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{
				Key:  nz.IdentityKeyIRD,
				Code: cbc.Code(tt.input),
			}
			r.Normalizer(id)
			assert.Equal(t, tt.expected, id.Code.String())
		})
	}
}
