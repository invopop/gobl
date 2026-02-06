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
		name   string
		nzbn   string
		errMsg string
	}{
		{"TooShort", "942904123456", "invalid NZBN"},
		{"TooLong", "94290412345638", "invalid NZBN"},
		{"BadCheckDigit", "9429041234560", "invalid NZBN"},
		{"Letters", "942904123456A", "invalid NZBN"},
		{"SpecialChars", "942904!123456", "invalid NZBN"},
		// Valid GLN but wrong GS1 prefix (UK "50" instead of NZ "94")
		{"WrongPrefix", "5012345678900", "must start with '94'"},
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
			assert.Contains(t, err.Error(), tt.errMsg)
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
