package il_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/il"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestIdentityNormalization(t *testing.T) {
	r := tax.RegimeDefFor("IL")

	t.Run("personal ID strips spaces", func(t *testing.T) {
		id := &org.Identity{
			Key:  il.IdentityKeyPersonalID,
			Code: "123 456 782",
		}
		r.NormalizeObject(id)
		assert.Equal(t, "123456782", id.Code.String())
	})

	t.Run("personal ID strips dashes", func(t *testing.T) {
		id := &org.Identity{
			Key:  il.IdentityKeyPersonalID,
			Code: "12-345-6782",
		}
		r.NormalizeObject(id)
		assert.Equal(t, "123456782", id.Code.String())
	})

	t.Run("corporation number strips spaces", func(t *testing.T) {
		id := &org.Identity{
			Key:  il.IdentityKeyCorporationNumber,
			Code: "510 123 456",
		}
		r.NormalizeObject(id)
		assert.Equal(t, "510123456", id.Code.String())
	})

	t.Run("unknown key not modified", func(t *testing.T) {
		id := &org.Identity{
			Key:  "other",
			Code: "abc 123",
		}
		r.NormalizeObject(id)
		assert.Equal(t, "abc 123", id.Code.String())
	})
}

func TestValidatePersonalID(t *testing.T) {
	t.Run("nil identity", func(t *testing.T) {
		var id *org.Identity
		err := il.Validate(id)
		assert.NoError(t, err)
	})

	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		// Valid: 9-digit numbers passing Luhn
		{name: "valid personal ID", code: "123456782"},
		{name: "valid all zeros with check", code: "000000000"},

		// Invalid: wrong length
		{name: "too short", code: "12345678", err: "code"},
		{name: "too long", code: "1234567890", err: "code"},

		// Invalid: non-numeric
		{name: "contains letters", code: "12345678A", err: "code"},

		// Invalid: bad Luhn checksum
		{name: "bad checksum", code: "123456789", err: "invalid checksum"},

		// Empty: required
		{name: "empty code", code: "", err: "code"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Key: il.IdentityKeyPersonalID, Code: tt.code}
			err := il.Validate(id)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.err)
				}
			}
		})
	}
}

func TestValidateCorporationNumber(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		// Valid: 9 digits (no Luhn check)
		{name: "valid company 51 prefix", code: "510123456"},
		{name: "valid company 58 prefix", code: "580123456"},
		{name: "valid any 9 digits", code: "123456789"},

		// Invalid: wrong length
		{name: "too short", code: "12345678", err: "code"},
		{name: "too long", code: "1234567890", err: "code"},

		// Invalid: non-numeric
		{name: "contains letters", code: "51012345A", err: "code"},

		// Empty: required
		{name: "empty code", code: "", err: "code"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Key: il.IdentityKeyCorporationNumber, Code: tt.code}
			err := il.Validate(id)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.err)
				}
			}
		})
	}
}

func TestValidateUnknownIdentityKey(t *testing.T) {
	id := &org.Identity{Key: "unknown", Code: "anything"}
	err := il.Validate(id)
	assert.NoError(t, err)
}
