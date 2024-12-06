package sdi

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Inbox keys to universally identify where copies of documents can be sent.
const (
	KeyInboxCode cbc.Key = "it-sdi-code"
	KeyInboxPEC  cbc.Key = "it-sdi-pec"
)

var inboxes = []*cbc.Definition{
	{
		Key: KeyInboxCode,
		Name: i18n.String{
			i18n.EN: "SDI Destination Code",
			i18n.IT: "Codice Destinatario SDI",
		},
	},
	{
		Key: KeyInboxPEC,
		Name: i18n.String{
			i18n.EN: "SDI PEC Destination",
			i18n.IT: "PEC Destinatario SDI",
		},
	},
}
