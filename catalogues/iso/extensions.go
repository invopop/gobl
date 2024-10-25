package iso

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

const (
	// ExtKeySchemeID is used by the ISO 6523 scheme identifier.
	ExtKeySchemeID cbc.Key = "iso-scheme-id"
)

var extensions = []*cbc.KeyDefinition{
	{
		Key:  ExtKeySchemeID,
		Name: i18n.NewString("ISO/IEC 6523 Identifier scheme code"),
		Desc: i18n.NewString(here.Doc(`
			Defines a global structure for uniquely identifying organizations or entities.
			This standard is essential in environments where electronic communications require
			unambiguous identification of organizations, especially in automated systems or
			electronic data interchange (EDI).

			The ISO 6523 set of identifies is used by the EN16931 standard for electronic invoicing.
		`)),
		Pattern: `^\d{4}$`,
	},
}
