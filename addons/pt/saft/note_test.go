package saft_test

import (
	"testing"

	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNoteValidation(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)

	t.Run("nil note", func(t *testing.T) {
		var note *org.Note
		err := addon.Validator(note)
		assert.NoError(t, err)
	})

	t.Run("invalid exemption note - too long", func(t *testing.T) {
		note := &org.Note{
			Key:  org.NoteKeyLegal,
			Src:  saft.ExtKeyExemption,
			Text: "1234567890123456789012345678901234567890123456789012345678901", // 61 chars
		}
		err := addon.Validator(note)
		assert.ErrorContains(t, err, "text: the length must be between 6 and 60")
	})

	t.Run("invalid exemption note - too short", func(t *testing.T) {
		note := &org.Note{
			Key:  org.NoteKeyLegal,
			Src:  saft.ExtKeyExemption,
			Text: "12345",
		}
		err := addon.Validator(note)
		assert.ErrorContains(t, err, "text: the length must be between 6 and 60")
	})

	t.Run("valid exemption note - min length", func(t *testing.T) {
		note := &org.Note{
			Key:  org.NoteKeyLegal,
			Src:  saft.ExtKeyExemption,
			Text: "123456",
		}

		err := addon.Validator(note)
		assert.NoError(t, err)
	})

	t.Run("valid exemption note - max length", func(t *testing.T) {
		note := &org.Note{
			Key:  org.NoteKeyLegal,
			Src:  saft.ExtKeyExemption,
			Text: "123456789012345678901234567890123456789012345678901234567890", // 60 chars
		}

		err := addon.Validator(note)
		assert.NoError(t, err)
	})

	t.Run("other note", func(t *testing.T) {
		note := &org.Note{
			Text: "1234567890123456789012345678901234567890123456789012345678901", // 61 chars
		}
		err := addon.Validator(note)
		assert.NoError(t, err)
	})
}
