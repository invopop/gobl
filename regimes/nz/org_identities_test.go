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
