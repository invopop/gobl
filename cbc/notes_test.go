package cbc_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/stretchr/testify/assert"
)

func TestNotesValidation(t *testing.T) {
	n := new(cbc.Note)
	n.Text = "This is a general note test"

	err := n.Validate()
	assert.NoError(t, err) // empty key ok

	n.Key = cbc.NoteKeyGeneral
	err = n.Validate()
	assert.NoError(t, err)

	n.Key = cbc.NoteKey("fooo")
	err = n.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key: must be a valid value")
}
