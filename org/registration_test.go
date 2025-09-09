package org_test

import (
	"testing"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestRegistrationNormalize_TrimsFieldsAndSetsUUID(t *testing.T) {
	r := &org.Registration{
		Label:   "  My Label  ",
		Office:  "\tOffice Name  ",
		Book:    "  Book ",
		Volume:  "  Volume  ",
		Sheet:   "  Sheet  ",
		Section: "  Section  ",
		Page:    "  Page  ",
		Entry:   "  Entry  ",
		Other:   "  Other  ",
	}

	assert.NotPanics(t, func() {
		r.Normalize()
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
}

func TestRegistrationNormalize_OnNilReceiverDoesNotPanic(t *testing.T) {
	var r *org.Registration
	assert.NotPanics(t, func() {
		r.Normalize()
	})
}

func TestRegistrationValidate_ValidData(t *testing.T) {
	r := &org.Registration{
		Label:    "  Company  ",
		Office:   "  Main Office  ",
		Currency: currency.EUR,
	}
	r.Normalize()

	err := r.Validate()
	assert.NoError(t, err)
}

func TestRegistrationValidate_InvalidCurrency(t *testing.T) {
	r := &org.Registration{
		Currency: currency.Code("ZZZ"),
	}
	err := r.Validate()
	assert.Error(t, err)
}
