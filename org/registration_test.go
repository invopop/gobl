package org_test

import (
	"testing"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

func TestRegistrationNormalize_TrimsFieldsAndSetsUUID(t *testing.T) {
	r := &org.Registration{
		Label:            "  My Label  ",
		Office:           "\tOffice Name  ",
		Book:             "  Book ",
		Volume:           "  Volume  ",
		Sheet:            "  Sheet  ",
		Section:          "  Section  ",
		Page:             "  Page  ",
		Entry:            "  Entry  ",
		Other:            "  Other  ",
		LiquidationState: "  LS  ",
		SoleShareholder:  "\tSU  ",
	}

	assert.NotPanics(t, func() {
		norm.Normalize(r)
	})

	assert.Equal(t, "My Label", r.Label)
	assert.Equal(t, "Office Name", r.Office)
	assert.Equal(t, "Book", r.Book)
	assert.Equal(t, "Volume", r.Volume)
	assert.Equal(t, "Sheet", r.Sheet)
	assert.Equal(t, "Section", r.Section)
	assert.Equal(t, "Page", r.Page)
	assert.Equal(t, "Entry", r.Entry)
	assert.Equal(t, "Other", r.Other)
	assert.Equal(t, "LS", r.LiquidationState)
	assert.Equal(t, "SU", r.SoleShareholder)
}

func TestRegistrationNormalize_OnNilReceiverDoesNotPanic(t *testing.T) {
	var r *org.Registration
	assert.NotPanics(t, func() {
		norm.Normalize(r)
	})
}

func TestRegistrationValidate_ValidData(t *testing.T) {
	r := &org.Registration{
		Label:    "  Company  ",
		Office:   "  Main Office  ",
		Currency: currency.EUR,
	}
	norm.Normalize(r)

	assert.NoError(t, rules.Validate(r))
}

func TestRegistrationValidate_InvalidCurrency(t *testing.T) {
	r := &org.Registration{
		Currency: currency.Code("ZZZ"),
	}
	assert.Error(t, rules.Validate(r))
}
