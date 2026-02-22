package jp

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests for nil/empty-value handling in the internal validator functions of
// org_identities.go and tax_identities.go.

func TestValidateCorporateNumber_Empty(t *testing.T) {
	// Empty cbc.Code: should return nil (guard clause).
	err := validateCorporateNumber(cbc.Code(""))
	assert.NoError(t, err)
}

func TestValidateCorporateNumber_WrongType(t *testing.T) {
	// Non-cbc.Code value: type assertion fails, returns nil.
	err := validateCorporateNumber("not-a-cbc-code")
	assert.NoError(t, err)
}

func TestValidateCorpNumberCheckDigit_WrongLength(t *testing.T) {
	// String shorter than 13 digits must return an error.
	err := validateCorpNumberCheckDigit("12345")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid length")
}

func TestValidateQualifiedInvoiceIssuer_Empty(t *testing.T) {
	// Empty cbc.Code: should return nil (guard clause).
	err := validateQualifiedInvoiceIssuer(cbc.Code(""))
	assert.NoError(t, err)
}

func TestValidateQualifiedInvoiceIssuer_WrongType(t *testing.T) {
	// Non-cbc.Code value: type assertion fails, returns nil.
	err := validateQualifiedInvoiceIssuer("not-a-cbc-code")
	assert.NoError(t, err)
}

func TestValidateMyNumber_Empty(t *testing.T) {
	// Empty cbc.Code: should return nil (guard clause).
	err := validateMyNumber(cbc.Code(""))
	assert.NoError(t, err)
}

func TestValidateMyNumber_WrongType(t *testing.T) {
	// Non-cbc.Code value: type assertion fails, returns nil.
	err := validateMyNumber("not-a-cbc-code")
	assert.NoError(t, err)
}

func TestValidateTNumberCheckDigit_WrongLength(t *testing.T) {
	// String shorter than 14 chars must return an error.
	err := validateTNumberCheckDigit("T123")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid format")
}

func TestValidateTNumberCheckDigit_WrongPrefix(t *testing.T) {
	// 14-char string not starting with 'T' must return an error.
	err := validateTNumberCheckDigit("X1234567890123")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid format")
}
