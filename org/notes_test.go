package org_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestNotesValidation(t *testing.T) {
	n := new(org.Note)
	n.Text = "This is a general note test"

	err := n.Validate()
	assert.NoError(t, err) // empty key ok

	n.Key = org.NoteKeyGeneral
	err = n.Validate()
	assert.NoError(t, err)

	n.Key = org.NoteKey("fooo")
	err = n.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key: must be a valid value")
}
