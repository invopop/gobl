package it

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Inbox keys to universally identify where copies of documents can be sent.
const (
	KeyInboxSDICode cbc.Key = "it-sdi-code"
	KeyInboxSDIPEC  cbc.Key = "it-sdi-pec"
)

var inboxKeyDefinitions = []*cbc.KeyDefinition{
	{
		Key: KeyInboxSDICode,
		Name: i18n.String{
			i18n.EN: "SDI Destination Code",
			i18n.IT: "Codice Destinatario SDI",
		},
	},
	{
		Key: KeyInboxSDIPEC,
		Name: i18n.String{
			i18n.EN: "SDI PEC Destination",
			i18n.IT: "PEC Destinatario SDI",
		},
	},
}
