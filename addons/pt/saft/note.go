package saft

import (
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

func validateNote(note *org.Note) error {
	if note == nil {
		return nil
	}

	return validation.ValidateStruct(note,
		validation.Field(&note.Text,
			validation.When(
				note.Key == org.NoteKeyLegal && note.Src == ExtKeyExemption,
				validation.Length(0, 60),
				validation.Skip,
			),
			validation.Skip,
		),
	)
}
