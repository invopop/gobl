package saft_test

import (
	"testing"

	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

func TestNoteValidation(t *testing.T) {
	t.Run("nil note", func(t *testing.T) {
		var note *org.Note
		err := rules.Validate(note, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("invalid exemption note - too long", func(t *testing.T) {
		note := &org.Note{
			Key:  org.NoteKeyLegal,
			Src:  saft.ExtKeyExemption,
			Text: "1234567890123456789012345678901234567890123456789012345678901", // 61 chars
		}
		err := rules.Validate(note, withAddonContext())
		assert.ErrorContains(t, err, "the length must be between 6 and 60")
	})

	t.Run("invalid exemption note - too short", func(t *testing.T) {
		note := &org.Note{
			Key:  org.NoteKeyLegal,
			Src:  saft.ExtKeyExemption,
			Text: "12345",
		}
		err := rules.Validate(note, withAddonContext())
		assert.ErrorContains(t, err, "the length must be between 6 and 60")
	})

	t.Run("valid exemption note - min length", func(t *testing.T) {
		note := &org.Note{
			Key:  org.NoteKeyLegal,
			Src:  saft.ExtKeyExemption,
			Text: "123456",
		}

		err := rules.Validate(note, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("valid exemption note - max length", func(t *testing.T) {
		note := &org.Note{
			Key:  org.NoteKeyLegal,
			Src:  saft.ExtKeyExemption,
			Text: "123456789012345678901234567890123456789012345678901234567890", // 60 chars
		}

		err := rules.Validate(note, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("other note", func(t *testing.T) {
		note := &org.Note{
			Text: "1234567890123456789012345678901234567890123456789012345678901", // 61 chars
		}
		err := rules.Validate(note, withAddonContext())
		assert.NoError(t, err)
	})
}
