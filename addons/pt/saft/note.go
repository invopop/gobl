package saft

import (
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

func orgNoteRules() *rules.Set {
	return rules.For(new(org.Note),
		rules.When(is.Func("legal exemption note", noteIsLegalExemption),
			rules.Field("text",
				rules.Assert("01", "the length must be between 6 and 60", is.Length(6, 60)),
			),
		),
	)
}

// noteIsLegalExemption checks if the note is a legal exemption note.
func noteIsLegalExemption(val any) bool {
	note, ok := val.(*org.Note)
	return ok && note != nil && note.Key == org.NoteKeyLegal && note.Src == ExtKeyExemption
}
