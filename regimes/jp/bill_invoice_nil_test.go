package jp

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests for nil and empty-list handling in bill_invoice.go validators.

func TestValidateSupplier_NilParty(t *testing.T) {
	err := validateSupplier((*org.Party)(nil))
	assert.NoError(t, err)
}

func TestValidateSupplier_WrongType(t *testing.T) {
	err := validateSupplier("not a party")
	assert.NoError(t, err)
}

func TestValidateCustomer_NilParty(t *testing.T) {
	err := validateCustomer((*org.Party)(nil))
	assert.NoError(t, err)
}

func TestValidateCustomer_WrongType(t *testing.T) {
	err := validateCustomer("not a party")
	assert.NoError(t, err)
}

func TestValidateAddress_Nil(t *testing.T) {
	err := validateAddress(nil)
	assert.NoError(t, err)
}

func TestValidateInvoiceLine_Nil(t *testing.T) {
	err := validateInvoiceLine(nil)
	assert.NoError(t, err)
}

func TestValidateInvoiceLines_NilValue(t *testing.T) {
	// Untyped nil: type assertion fails (!ok), returns nil.
	err := validateInvoiceLines(nil)
	assert.NoError(t, err)
}

func TestValidateInvoiceLines_NilSlice(t *testing.T) {
	// Typed nil slice: ok but lines == nil, returns nil.
	var lines []*bill.Line
	err := validateInvoiceLines(lines)
	assert.NoError(t, err)
}

func TestValidateInvoiceLines_WrongType(t *testing.T) {
	err := validateInvoiceLines("not lines")
	assert.NoError(t, err)
}

func TestValidateInvoiceLines_LineError(t *testing.T) {
	// A line missing its item must propagate a validation error.
	lines := []*bill.Line{
		{}, // no Item, no Quantity
	}
	err := validateInvoiceLines(lines)
	require.Error(t, err)
}

func TestValidateExportLines_NilValue(t *testing.T) {
	// Untyped nil: type assertion fails (!ok), returns nil.
	err := validateExportLines(nil)
	assert.NoError(t, err)
}

func TestValidateExportLines_NilSlice(t *testing.T) {
	// Typed nil slice: ok but lines == nil, returns nil.
	var lines []*bill.Line
	err := validateExportLines(lines)
	assert.NoError(t, err)
}

func TestValidateExportLines_WrongType(t *testing.T) {
	err := validateExportLines("not lines")
	assert.NoError(t, err)
}

func TestValidateExportLines_NilLine(t *testing.T) {
	// Nil line in the slice must be skipped without error.
	err := validateExportLines([]*bill.Line{nil})
	assert.NoError(t, err)
}

func TestValidateExportLines_NilCombo(t *testing.T) {
	// Nil tax combo must be skipped without error.
	err := validateExportLines([]*bill.Line{
		{Taxes: []*tax.Combo{nil}},
	})
	assert.NoError(t, err)
}

func TestValidateExportLines_NonVATCombo(t *testing.T) {
	// Non-VAT combos must be ignored; no error expected.
	err := validateExportLines([]*bill.Line{
		{Taxes: []*tax.Combo{{Category: "OTHER"}}},
	})
	assert.NoError(t, err)
}
